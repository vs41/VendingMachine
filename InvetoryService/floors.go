package invetoryservice

type Floor struct {
	floorID string
	slots   []*Slot
}

func NewFloor(floorID string, slots []*Slot) *Floor {
	copiedSlots := make([]*Slot, len(slots))
	copy(copiedSlots, slots)

	return &Floor{
		floorID: floorID,
		slots:   copiedSlots,
	}
}

func (f *Floor) GetFloorID() string {
	return f.floorID
}

func (f *Floor) GetSlots() []*Slot {
	copiedSlots := make([]*Slot, len(f.slots))
	copy(copiedSlots, f.slots)

	return copiedSlots
}

func (f *Floor) GetNumberOfSlots() int {
	return len(f.slots)
}

func (f *Floor) AddSlot(slot *Slot) {
	if slot == nil {
		return
	}

	f.slots = append(f.slots, slot)
}

func (f *Floor) RemoveSlot(slotID string) bool {
	for i, slot := range f.slots {
		if slot.slotID == slotID {
			f.slots = append(f.slots[:i], f.slots[i+1:]...)
			return true
		}
	}

	return false
}

func (f *Floor) FindSlot(slotID string) *Slot {
	for _, slot := range f.slots {
		if slot.slotID == slotID {
			return slot
		}
	}

	return nil
}
