package invetoryservice

type Inventory struct {
	floors []*Floor
}

func NewInventory(floors []*Floor) *Inventory {
	copiedFloors := make([]*Floor, len(floors))
	copy(copiedFloors, floors)

	return &Inventory{
		floors: copiedFloors,
	}
}

func (i *Inventory) GetFloors() []*Floor {
	copiedFloors := make([]*Floor, len(i.floors))
	copy(copiedFloors, i.floors)

	return copiedFloors
}

func (i *Inventory) AddFloor(floor *Floor) {
	if floor == nil {
		return
	}

	i.floors = append(i.floors, floor)
}

func (i *Inventory) RemoveFloor(floorID string) bool {
	for j, floor := range i.floors {
		if floor.floorID == floorID {
			i.floors = append(i.floors[:j], i.floors[j+1:]...)
			return true
		}
	}

	return false
}

func (i *Inventory) FindFloor(floorID string) *Floor {
	for _, floor := range i.floors {
		if floor.floorID == floorID {
			return floor
		}
	}

	return nil
}

func (i *Inventory) FindSlot(slotID string) *Slot {
	for _, floor := range i.floors {
		slot := floor.FindSlot(slotID)

		if slot != nil {
			return slot
		}
	}

	return nil
}
