package memdb

// familySlotRange represents the time slot range based on family time
type familySlotRange struct {
	start, end uint16 // time slot range
}

// newFamilySlotRange creates a new family slot range with start/end
func newFamilySlotRange(start, end uint16) *familySlotRange {
	return &familySlotRange{
		start: start,
		end:   end,
	}
}

// setSlot sets the time slot range
func (sr *familySlotRange) setSlot(slot uint16) {
	if slot < sr.start {
		sr.start = slot
	}
	if slot > sr.end {
		sr.end = slot
	}
}

// getSlotRange returns return the time slot range
func (sr *familySlotRange) getSlotRange() (start, end uint16) {
	return sr.start, sr.end
}
