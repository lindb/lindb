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

package indexdb

const full = 10000

// seriesEvent represents the series data(tags hash=>series id)
type seriesEvent struct {
	tagsHash uint64
	seriesID uint32
}

// metricEvent represents metric id mapping include series/metric id sequence
type metricEvent struct {
	metricIDSeq uint32
	events      []seriesEvent
}

// mappingEvent represents the pending persist id mapping events
type mappingEvent struct {
	events map[uint32]*metricEvent

	pending int
}

// newMappingEvent creates a id mapping event
func newMappingEvent() *mappingEvent {
	return &mappingEvent{
		events: make(map[uint32]*metricEvent),
	}
}

// addSeriesID adds series data for metric
func (event *mappingEvent) addSeriesID(metricID uint32, tagsHash uint64, seriesID uint32) {
	e, ok := event.events[metricID]
	if !ok {
		e = &metricEvent{}
		event.events[metricID] = e
	}
	e.events = append(e.events, seriesEvent{
		tagsHash: tagsHash,
		seriesID: seriesID,
	})
	// set id sequence directly, because gen series id in order
	e.metricIDSeq = seriesID
	event.pending++
}

// isFull returns if events is full
func (event *mappingEvent) isFull() bool {
	return event.pending > full
}

// isEmpty returns if evens is empty
func (event *mappingEvent) isEmpty() bool {
	return event.pending == 0
}
