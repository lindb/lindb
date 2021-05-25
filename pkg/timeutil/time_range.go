package timeutil

// SlotRange represents time range with start/end timestamp using low value.
type SlotRange struct {
	Start, End uint16
}

// NewSlotRange creates a new slot range with start/end
func NewSlotRange(start, end uint16) SlotRange {
	return SlotRange{
		Start: start,
		End:   end,
	}
}

// setSlot sets the time slot range
func (sr *SlotRange) SetSlot(slot uint16) {
	if slot < sr.Start {
		sr.Start = slot
	}
	if slot > sr.End {
		sr.End = slot
	}
}

// getSlotRange returns return the time slot range
func (sr *SlotRange) GetRange() (start, end uint16) {
	return sr.Start, sr.End
}

// TimeRange represents time range with start/end timestamp.
type TimeRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// IsEmpty tests if empty, start>=end => empty
func (r *TimeRange) IsEmpty() bool {
	return r.Start >= r.End
}

// Contains tests if timestamp in current time range
func (r *TimeRange) Contains(timestamp int64) bool {
	return timestamp >= r.Start && timestamp <= r.End
}

// Overlap tests if overlap with current time range
func (r *TimeRange) Overlap(o *TimeRange) bool {
	return r.Contains(o.Start) || o.Contains(r.Start)
}

// Intersect returns the intersection of two time range
func (r *TimeRange) Intersect(o *TimeRange) *TimeRange {
	result := &TimeRange{}
	result.Start = r.Start
	if o.Start > r.Start {
		result.Start = o.Start
	}
	result.End = r.End
	if o.End < r.End {
		result.End = o.End
	}
	return result
}
