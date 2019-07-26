package selector

// SlotSelector represents a slot selector for aggregator value's index
type SlotSelector interface {
	// IndexOf returns the index of the specified element in aggregator values
	IndexOf(timeSlot int) int
}

// indexSlotSelector represents an index slot selector based on start/ratio
type indexSlotSelector struct {
	startSlot     int
	intervalRatio int
}

// NewIndexSlotSelector creates an index slot selector using given start and ratio
func NewIndexSlotSelector(startSlot, intervalRatio int) SlotSelector {
	return &indexSlotSelector{
		startSlot:     startSlot,
		intervalRatio: intervalRatio,
	}
}

// IndexOf returns the index of the specified element in aggregator values
// index = (timeSlot - start)/ratio, if timeSlot < start return -1
func (s *indexSlotSelector) IndexOf(timeSlot int) int {
	if timeSlot < s.startSlot {
		return -1
	}
	return (timeSlot - s.startSlot) / s.intervalRatio
}
