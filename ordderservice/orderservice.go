package orderservice

import (
	"errors"

	"vendingmachine/cartservice"
)

type OrderService struct {
	order *Order
}

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotFound      = errors.New("order not found")
	ErrEmptyCart          = errors.New("cart is empty")
)

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (os *OrderService) CreateOrder(cartItems []*cartservice.CartItem) (*Order, error) {
	if os.order != nil {
		return nil, ErrOrderAlreadyExists
	}

	if len(cartItems) == 0 {
		return nil, ErrEmptyCart
	}

	orderItems := make([]*OrderItem, 0, len(cartItems))
	var totalAmount int64

	for _, item := range cartItems {
		if item == nil || item.GetItemQuantity() <= 0 {
			return nil, errors.New("invalid cart item")
		}

		orderItem := NewOrderItem(
			item.GetSlotID(),
			item.GetFloorID(),
			item.GetItemName(),
			item.GetUnitPrice(),
			item.GetItemQuantity(),
		)

		orderItems = append(orderItems, orderItem)
		totalAmount += item.GetTotalPrice()
	}

	os.order = NewOrder(orderItems, totalAmount)

	return os.order, nil
}

func (os *OrderService) GetOrder() (*Order, error) {
	if os.order == nil {
		return nil, ErrOrderNotFound
	}

	return os.order, nil
}

func (os *OrderService) UpdateOrderState(state OrderState) error {
	if os.order == nil {
		return ErrOrderNotFound
	}

	return os.order.SetOrderState(state)
}

func (os *OrderService) ClearOrder() {
	os.order = nil
}
