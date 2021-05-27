// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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

// Intersect returns the intersection of two slot range
func (sr *SlotRange) Intersect(o *SlotRange) *SlotRange {
	result := &SlotRange{}
	result.Start = sr.Start
	if o.Start > sr.Start {
		result.Start = o.Start
	}
	result.End = sr.End
	if o.End < sr.End {
		result.End = o.End
	}
	return result
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
