package invetoryservice

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrSlotNotFound = errors.New("slot not found")
)

type InventoryService struct {
	inv *Inventory
	mu  sync.Mutex
}

func NewInventoryService(inv *Inventory) *InventoryService {
	return &InventoryService{
		inv: inv,
	}
}

// Find a physical slot across all floors.
func (is *InventoryService) FindSlot(slotID string) *Slot {
	is.mu.Lock()
	defer is.mu.Unlock()

	return is.inv.FindSlot(slotID)
}

// Represents a requested reservation.
type ReservationItem struct {
	SlotID   string
	Quantity int
}

// Represents the final dispensing result for one slot.
type DispenseSettlement struct {
	SlotID            string
	ReservedQuantity  int
	DispensedQuantity int
}

// Reserve multiple slots atomically.
func (is *InventoryService) Reserve(items []ReservationItem) error {
	is.mu.Lock()
	defer is.mu.Unlock()

	if len(items) == 0 {
		return errors.New("no items to reserve")
	}

	// Aggregate quantities by slot.
	requested := make(map[string]int)

	for _, item := range items {
		if item.Quantity <= 0 {
			return ErrInvalidQuantity
		}

		requested[item.SlotID] += item.Quantity
	}

	// First validate every slot.
	slots := make(map[string]*Slot)

	for slotID, quantity := range requested {
		slot := is.inv.FindSlot(slotID)

		if slot == nil {
			return fmt.Errorf("%w: %s", ErrSlotNotFound, slotID)
		}

		if slot.product == nil {
			return fmt.Errorf("%w: %s", ErrSlotEmpty, slotID)
		}

		if slot.GetAvailableQuantity() < quantity {
			return fmt.Errorf(
				"%w: slot %s",
				ErrInsufficientStock,
				slotID,
			)
		}

		slots[slotID] = slot
	}

	// All validations passed. Now update the reservations.
	for slotID, quantity := range requested {
		slots[slotID].reservedQuantity += quantity
	}

	return nil
}

// Release reservations when payment fails or an order is cancelled.
func (is *InventoryService) Release(items []ReservationItem) error {
	is.mu.Lock()
	defer is.mu.Unlock()

	requested := make(map[string]int)

	for _, item := range items {
		if item.Quantity <= 0 {
			return ErrInvalidQuantity
		}

		requested[item.SlotID] += item.Quantity
	}

	// Validate everything before changing anything.
	slots := make(map[string]*Slot)

	for slotID, quantity := range requested {
		slot := is.inv.FindSlot(slotID)

		if slot == nil {
			return fmt.Errorf("%w: %s", ErrSlotNotFound, slotID)
		}

		if slot.reservedQuantity < quantity {
			return fmt.Errorf(
				"%w: slot %s",
				ErrInsufficientReserved,
				slotID,
			)
		}

		slots[slotID] = slot
	}

	// Apply all releases.
	for slotID, quantity := range requested {
		slots[slotID].reservedQuantity -= quantity
	}

	return nil
}

// Settle reservations after dispensing.
//
// Example:
// ReservedQuantity = 5
// DispensedQuantity = 3
//
// Physical stock decreases by 3.
// Reserved stock decreases by 5.
// The remaining 2 reserved units are released.
func (is *InventoryService) SettleDispensed(
	settlements []DispenseSettlement,
) error {
	is.mu.Lock()
	defer is.mu.Unlock()

	if len(settlements) == 0 {
		return errors.New("no dispensing results provided")
	}

	slots := make(map[string]*Slot)

	// Validate all results first.
	for _, result := range settlements {
		if result.ReservedQuantity <= 0 {
			return ErrInvalidQuantity
		}

		if result.DispensedQuantity < 0 ||
			result.DispensedQuantity > result.ReservedQuantity {
			return errors.New("invalid dispensed quantity")
		}

		if _, exists := slots[result.SlotID]; exists {
			return fmt.Errorf("duplicate settlement for slot %s", result.SlotID)
		}

		slot := is.inv.FindSlot(result.SlotID)

		if slot == nil {
			return fmt.Errorf("%w: %s", ErrSlotNotFound, result.SlotID)
		}

		if slot.reservedQuantity < result.ReservedQuantity {
			return fmt.Errorf(
				"%w: slot %s",
				ErrInsufficientReserved,
				result.SlotID,
			)
		}

		if slot.productQuantity < result.DispensedQuantity {
			return fmt.Errorf(
				"%w: slot %s",
				ErrInsufficientStock,
				result.SlotID,
			)
		}

		slots[result.SlotID] = slot
	}

	// Apply all settlements after validation succeeds.
	for _, result := range settlements {
		slot := slots[result.SlotID]

		slot.productQuantity -= result.DispensedQuantity
		slot.reservedQuantity -= result.ReservedQuantity

		if slot.productQuantity == 0 {
			slot.product = nil
		}
	}

	return nil
}

// Restock a slot.
func (is *InventoryService) Restock(
	slotID string,
	product *Product,
	quantity int,
) error {
	is.mu.Lock()
	defer is.mu.Unlock()

	slot := is.inv.FindSlot(slotID)

	if slot == nil {
		return ErrSlotNotFound
	}

	return slot.restock(product, quantity)
}

// Add a new floor.
func (is *InventoryService) AddFloor(floor *Floor) {
	is.mu.Lock()
	defer is.mu.Unlock()

	is.inv.AddFloor(floor)
}

// Remove an existing floor.
func (is *InventoryService) RemoveFloor(floorID string) bool {
	is.mu.Lock()
	defer is.mu.Unlock()

	return is.inv.RemoveFloor(floorID)
}
