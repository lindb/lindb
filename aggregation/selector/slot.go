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

package selector

import (
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source ./slot.go -destination=./slot_mock.go -package selector

// SlotSelector represents a slot selector for aggregator value's index
type SlotSelector interface {
	// IndexOf returns the index of the specified element in aggregator values
	IndexOf(timeSlot int) (idx int, completed bool)
	// Range returns the slot time range
	Range() (start int, end int)
	// PointCount returns the data point count in time range
	PointCount() int
}

// indexSlotSelector represents an index slot selector based on start/ratio
type indexSlotSelector struct {
	start, end    int
	intervalRatio int
	pointCount    int
}

// NewIndexSlotSelector creates an index slot selector using given ratio
func NewIndexSlotSelector(start, end, intervalRatio int) SlotSelector {
	return &indexSlotSelector{
		start:         start,
		end:           end,
		intervalRatio: intervalRatio,
		pointCount:    timeutil.CalPointCount(int64(start), int64(end), int64(intervalRatio)),
	}
}

// Range returns the slot time range
func (s *indexSlotSelector) Range() (start int, end int) {
	return s.start, s.end
}

// PointCount returns the data point count in time range
func (s *indexSlotSelector) PointCount() int {
	return s.pointCount
}

// IndexOf returns the index of the specified element in aggregator values,
// index = (timeSlot - start)/ratio, if timeSlot < start return -1
func (s *indexSlotSelector) IndexOf(timeSlot int) (idx int, completed bool) {
	switch {
	case timeSlot < s.start:
		return -1, false
	case timeSlot > s.end:
		return -1, true
	default:
		return (timeSlot - s.start) / s.intervalRatio, false
	}
}
