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

package aggregation

import (
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./field_agg.go -destination=./field_agg_mock.go -package=aggregation

// FieldAggregator represents a field aggregator, aggregator the field series which with same field id.
type FieldAggregator interface {
	// Aggregate aggregates the field series into current aggregator.
	Aggregate(it series.FieldIterator)
	// AggregateBySlot aggregates the field series into current aggregator.
	AggregateBySlot(slot int, value float64)
	// ResultSet returns the result set of field aggregator.
	ResultSet() (startTime int64, it series.FieldIterator)
	SlotRange() (start, end int)
	reset()
}

// fieldAggregator implements field aggregator interface, aggregator field series based on aggregator spec.
type fieldAggregator struct {
	aggTypes         []field.AggType
	segmentStartTime int64
	start, end       int

	fieldSeriesList []collections.FloatArray
}

// NewFieldAggregator creates a field aggregator,
// time range 's start and end is index based on segment start time and interval.
// e.g. segment start time = 20190905 10:00:00, start = 10, end = 50, interval = 10 seconds,
// real query time range {20190905 10:01:40 ~ 20190905 10:08:20}
func NewFieldAggregator(aggSpec AggregatorSpec, segmentStartTime int64, start, end int) FieldAggregator {
	//TODO maybe agg type has duplicate?
	var aggTypes []field.AggType
	for f := range aggSpec.Functions() {
		aggTypes = append(aggTypes, aggSpec.GetFieldType().GetFuncFieldParams(f)...)
	}

	agg := &fieldAggregator{
		aggTypes:         aggTypes,
		segmentStartTime: segmentStartTime,
		start:            start,
		end:              end,
		fieldSeriesList:  make([]collections.FloatArray, len(aggTypes)),
	}
	return agg
}

func (a *fieldAggregator) SlotRange() (start, end int) {
	return a.start, a.end
}

// ResultSet returns the result set of field aggregator
func (a *fieldAggregator) ResultSet() (startTime int64, it series.FieldIterator) {
	return a.segmentStartTime, newFieldIterator(a.start, a.aggTypes, a.fieldSeriesList)
}

// Aggregate aggregates the field series into current aggregator
func (a *fieldAggregator) Aggregate(it series.FieldIterator) {
	for it.HasNext() {
		pIt := it.Next()
		for pIt.HasNext() {
			slot, value := pIt.Next()
			a.AggregateBySlot(slot, value)
		}
	}
}

// Aggregate aggregates the field series into current aggregator
func (a *fieldAggregator) AggregateBySlot(slot int, value float64) {
	for idx, aggType := range a.aggTypes {
		values := a.fieldSeriesList[idx]
		if values == nil {
			values = collections.NewFloatArray(a.end - a.start + 1)
			values.SetValue(slot, value)
			a.fieldSeriesList[idx] = values
		} else {
			if values.HasValue(slot) {
				values.SetValue(slot, aggType.AggFunc().Aggregate(values.GetValue(slot), value))
			} else {
				values.SetValue(slot, value)
			}
		}
	}
}

func (a *fieldAggregator) reset() {
	for idx := range a.fieldSeriesList {
		a.fieldSeriesList[idx].Reset()
	}
}
