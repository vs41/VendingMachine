package invetoryservice

import "errors"

var (
	ErrInvalidQuantity       = errors.New("quantity must be greater than zero")
	ErrInsufficientStock     = errors.New("insufficient available quantity")
	ErrInsufficientReserved  = errors.New("insufficient reserved quantity")
	ErrSlotEmpty             = errors.New("slot does not contain a product")
	ErrProductAlreadyPresent = errors.New("slot already contains a product")
)

type Slot struct {
	slotID           string
	floorID          string
	product          *Product
	productQuantity  int
	reservedQuantity int
}

func NewSlot(
	slotID string,
	floorID string,
	product *Product,
	productQuantity int,
) *Slot {
	return &Slot{
		slotID:           slotID,
		floorID:          floorID,
		product:          product,
		productQuantity:  productQuantity,
		reservedQuantity: 0,
	}
}

// Getters

func (s *Slot) GetSlotID() string {
	return s.slotID
}

func (s *Slot) GetFloorID() string {
	return s.floorID
}

func (s *Slot) GetProduct() *Product {
	return s.product
}

func (s *Slot) GetProductQuantity() int {
	return s.productQuantity
}

func (s *Slot) GetReservedQuantity() int {
	return s.reservedQuantity
}

func (s *Slot) GetAvailableQuantity() int {
	return s.productQuantity - s.reservedQuantity
}

// Inventory operations

func (s *Slot) reserve(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	if s.product == nil {
		return ErrSlotEmpty
	}

	if s.GetAvailableQuantity() < quantity {
		return ErrInsufficientStock
	}

	s.reservedQuantity += quantity
	return nil
}

func (s *Slot) release(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	if s.reservedQuantity < quantity {
		return ErrInsufficientReserved
	}

	s.reservedQuantity -= quantity
	return nil
}

// Dispense a quantity that has already been reserved.
func (s *Slot) dispense(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	if s.reservedQuantity < quantity {
		return ErrInsufficientReserved
	}

	if s.productQuantity < quantity {
		return ErrInsufficientStock
	}

	s.productQuantity -= quantity
	s.reservedQuantity -= quantity

	if s.productQuantity == 0 {
		s.product = nil
	}

	return nil
}

// Add stock to an existing slot.
func (s *Slot) restock(product *Product, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	if product == nil {
		return ErrSlotEmpty
	}

	if s.product != nil && s.product.productID != product.productID {
		return ErrProductAlreadyPresent
	}

	s.product = product
	s.productQuantity += quantity

	return nil
}

// Remove a product only when no stock is reserved.
func (s *Slot) clear() error {
	if s.reservedQuantity > 0 {
		return ErrInsufficientReserved
	}

	s.product = nil
	s.productQuantity = 0

	return nil
}
