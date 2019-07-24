package timeutil

// TimeRange represents time range with start/end timestamp
type TimeRange struct {
	Start, End int64
}

func (r *TimeRange) Contains(t int64) bool {
	return t >= r.Start && t <= r.End
}

func (r *TimeRange) Overlaps(o *TimeRange) bool {
	return r.Contains(o.Start) || o.Contains(r.Start)
}
