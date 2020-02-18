package memdb

// slotRange represents the time slot range based on family time
type slotRange struct {
	start, end uint16 // time slot range
}

// newSlotRange creates a new slot range with start/end
func newSlotRange(start, end uint16) slotRange {
	return slotRange{
		start: start,
		end:   end,
	}
}

// setSlot sets the time slot range
func (sr *slotRange) setSlot(slot uint16) {
	if slot < sr.start {
		sr.start = slot
	}
	if slot > sr.end {
		sr.end = slot
	}
}

// getSlotRange returns return the time slot range
func (sr *slotRange) getRange() (start, end uint16) {
	return sr.start, sr.end
}
