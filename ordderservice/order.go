package orderservice

import (
	"errors"
	"vendingmachine/utiles"
)

type OrderState string

const (
	OrderPending            OrderState = "PENDING"
	OrderCreated            OrderState = "CREATED"
	OrderCancelled          OrderState = "CANCELLED"
	OrderFailed             OrderState = "FAILED"
	OrderCompleted          OrderState = "COMPLETED"
	OrderPartiallyCompleted OrderState = "PARTIALLY_COMPLETED"
)

var (
	ErrStateAlreadySet        = errors.New("order state is already set")
	ErrInvalidStateTransition = errors.New("invalid order state transition")
)

type Order struct {
	orderID     string
	state       OrderState
	items       []*OrderItem
	totalAmount int64
}

type OrderItem struct {
	slotID    string
	floorID   string
	name      string
	unitPrice int64
	quantity  int
}

func NewOrderItem(
	slotID string,
	floorID string,
	name string,
	unitPrice int64,
	quantity int,
) *OrderItem {
	return &OrderItem{
		slotID:    slotID,
		floorID:   floorID,
		name:      name,
		unitPrice: unitPrice,
		quantity:  quantity,
	}
}

func NewOrder(orderItems []*OrderItem, total int64) *Order {
	items := make([]*OrderItem, 0, len(orderItems))

	for _, item := range orderItems {
		if item == nil {
			continue
		}

		itemCopy := *item
		items = append(items, &itemCopy)
	}

	return &Order{
		orderID:     utiles.UniqueIdGeneretor(),
		state:       OrderPending,
		items:       items,
		totalAmount: total,
	}
}

func (o *Order) GetOrderID() string {
	return o.orderID
}

func (o *Order) GetTotalAmount() int64 {
	return o.totalAmount
}

func (o *Order) GetItems() []*OrderItem {
	items := make([]*OrderItem, 0, len(o.items))

	for _, item := range o.items {
		itemCopy := *item
		items = append(items, &itemCopy)
	}

	return items
}

func (o *Order) GetOrderState() OrderState {
	return o.state
}

func (o *Order) SetOrderState(newState OrderState) error {
	if o.state == newState {
		return ErrStateAlreadySet
	}

	valid := false

	switch o.state {
	case OrderPending:
		valid = newState == OrderCreated ||
			newState == OrderCancelled ||
			newState == OrderFailed

	case OrderCreated:
		valid = newState == OrderCompleted ||
			newState == OrderPartiallyCompleted ||
			newState == OrderFailed
	}

	if !valid {
		return ErrInvalidStateTransition
	}

	o.state = newState
	return nil
}

// OrderItem getters

func (oi *OrderItem) GetSlotID() string {
	return oi.slotID
}

func (oi *OrderItem) GetName() string {
	return oi.name
}

func (oi *OrderItem) GetUnitPrice() int64 {
	return oi.unitPrice
}

func (oi *OrderItem) GetQuantity() int {
	return oi.quantity
}
