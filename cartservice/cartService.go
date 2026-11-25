package cartservice

import (
	"errors"

	inventoryservice "vendingmachine/InvetoryService"
	"vendingmachine/utiles"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrSlotNotFound      = errors.New("slot not found")
	ErrCartNotStarted    = errors.New("cart not started")
	ErrCartAlreadyExists = errors.New("cart already exists")
)

type CartService struct {
	cart *Cart
	inv  *inventoryservice.InventoryService
}

func NewCartService(cart *Cart, inv *inventoryservice.InventoryService) *CartService {
	return &CartService{cart: cart, inv: inv}
}

func (cs *CartService) StartCart(userID string) error {
	if cs.cart != nil {
		return ErrCartAlreadyExists
	}
	cs.cart = NewCart(utiles.UniqueIdGeneretor(), userID)
	return nil
}

// EndCart clears and releases the active cart so the next customer can start one.
func (cs *CartService) EndCart() {
	if cs.cart != nil {
		cs.cart.Clear()
		cs.cart = nil
	}
}

func (cs *CartService) AddItemToCart(slotID string, quantity int) error {
	if cs.cart == nil {
		return ErrCartNotStarted
	}
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	slot := cs.inv.FindSlot(slotID)
	if slot == nil {
		return ErrSlotNotFound
	}
	product := slot.GetProduct()
	if product == nil {
		return inventoryservice.ErrSlotEmpty
	}

	_, existingItem := cs.findItemBySlotID(slotID)
	requestedQuantity := quantity
	if existingItem != nil {
		requestedQuantity += existingItem.GetItemQuantity()
	}
	if err := cs.checkAvailability(slotID, requestedQuantity); err != nil {
		return err
	}

	if existingItem != nil {
		for i := 0; i < quantity; i++ {
			existingItem.increaseQuantityByOne()
		}
		return nil
	}

	item := NewCartItem(
		utiles.UniqueIdGeneretor(), quantity,
		product.GetProductPrice(), slot.GetSlotID(), slot.GetFloorID(), product.GetProductName(),
	)
	return cs.cart.AddItem(item)
}

func (cs *CartService) RemoveItemFromCart(itemID string) error {
	if cs.cart == nil {
		return ErrCartNotStarted
	}
	return cs.cart.RemoveItem(itemID)
}

func (cs *CartService) IncreaseItemQuantity(itemID string) error {
	if cs.cart == nil {
		return ErrCartNotStarted
	}
	_, item := cs.cart.FindItem(itemID)
	if item == nil {
		return ErrItemNotFound
	}
	if err := cs.checkAvailability(item.GetSlotID(), item.GetItemQuantity()+1); err != nil {
		return err
	}
	return cs.cart.IncreaseQuantity(itemID)
}

func (cs *CartService) DecreaseItemQuantity(itemID string) error {
	if cs.cart == nil {
		return ErrCartNotStarted
	}
	return cs.cart.DecreaseQuantity(itemID)
}

func (cs *CartService) GetCartTotal() int64 {
	if cs.cart == nil {
		return 0
	}
	return cs.cart.GetTotal()
}

func (cs *CartService) GetCartItems() ([]*CartItem, error) {
	if cs.cart == nil {
		return nil, ErrCartNotStarted
	}
	return cs.cart.GetCartItems(), nil
}

func (cs *CartService) ClearCart() error {
	if cs.cart == nil {
		return ErrCartNotStarted
	}
	cs.cart.Clear()
	return nil
}

func (cs *CartService) findItemBySlotID(slotID string) (int, *CartItem) {
	for index, item := range cs.cart.GetCartItems() {
		if item.GetSlotID() == slotID {
			return index, item
		}
	}
	return -1, nil
}

func (cs *CartService) checkAvailability(slotID string, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	slot := cs.inv.FindSlot(slotID)
	if slot == nil {
		return ErrSlotNotFound
	}
	if slot.GetAvailableQuantity() < quantity {
		return ErrInsufficientStock
	}
	return nil
}
