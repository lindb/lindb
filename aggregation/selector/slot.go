package selector

// SlotSelector represents a slot selector for aggregator value's index
type SlotSelector interface {
	// IndexOf returns the index of the specified element in aggregator values
	IndexOf(startSlot, timeSlot int) int
}

// indexSlotSelector represents an index slot selector based on start/ratio
type indexSlotSelector struct {
	intervalRatio int
}

// NewIndexSlotSelector creates an index slot selector using given ratio
func NewIndexSlotSelector(intervalRatio int) SlotSelector {
	return &indexSlotSelector{
		intervalRatio: intervalRatio,
	}
}

// IndexOf returns the index of the specified element in aggregator values
// index = (timeSlot - start)/ratio, if timeSlot < start return -1
func (s *indexSlotSelector) IndexOf(startSlot, timeSlot int) int {
	if timeSlot < startSlot {
		return -1
	}
	return (timeSlot - startSlot) / s.intervalRatio
}
