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
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./series_agg.go -destination=./series_agg_mock.go -package=aggregation

// newBlockFunc represents create series block function by query time range.
//type newBlockFunc func() series.Block
//
// FieldAggregates represents aggregator which aggregates fields of a time series
type FieldAggregates []SeriesAggregator

// ResultSet returns the result set of aggregator
func (agg FieldAggregates) ResultSet(tags string) series.GroupedIterator {
	return newGroupedIterator(tags, agg)
}

// Reset resets the aggregator's context for reusing
func (agg FieldAggregates) Reset() {
	for _, aggregator := range agg {
		aggregator.Reset()
	}
}

// NewFieldAggregates creates the field aggregates based on aggregator specs and query time range.
// NOTICE: if do down sampling aggregator, aggregator specs must be in order by field id.
func NewFieldAggregates(
	queryInterval timeutil.Interval,
	intervalRatio int,
	queryTimeRange timeutil.TimeRange,
	aggSpecs AggregatorSpecs,
) FieldAggregates {
	aggregates := make(FieldAggregates, len(aggSpecs))
	for idx, aggSpec := range aggSpecs {
		aggregates[idx] = NewSeriesAggregator(queryInterval, intervalRatio, queryTimeRange, aggSpec)
	}
	return aggregates
}

// SeriesAggregator represents a series aggregator which aggregates one field of a time series
type SeriesAggregator interface {
	// FieldName returns field name
	FieldName() field.Name
	// GetFieldType returns field type
	GetFieldType() field.Type
	GetAggregator(segmentStartTime int64) (agg FieldAggregator, ok bool)
	GetAggregates() []FieldAggregator

	ResultSet() series.Iterator
	Reset()
}

type seriesAggregator struct {
	fieldName field.Name
	fieldType field.Type

	queryInterval  timeutil.Interval
	queryTimeRange timeutil.TimeRange
	intervalRatio  int

	aggregates []FieldAggregator
	aggSpec    AggregatorSpec
	calc       timeutil.Calculator

	startTime int64
}

// NewSeriesAggregator creates a series aggregator.
func NewSeriesAggregator(
	queryInterval timeutil.Interval,
	intervalRatio int,
	queryTimeRange timeutil.TimeRange,
	aggSpec AggregatorSpec,
) SeriesAggregator {
	calc := queryInterval.Calculator()
	segmentTime := calc.CalcSegmentTime(queryTimeRange.Start)
	startTime := calc.CalcFamilyStartTime(segmentTime, calc.CalcFamily(queryTimeRange.Start, segmentTime))

	length := calc.CalcTimeWindows(queryTimeRange.Start, queryTimeRange.End)

	agg := &seriesAggregator{
		fieldName:      aggSpec.FieldName(),
		fieldType:      aggSpec.GetFieldType(),
		startTime:      startTime,
		calc:           calc,
		intervalRatio:  intervalRatio,
		queryInterval:  queryInterval,
		queryTimeRange: queryTimeRange,
		aggSpec:        aggSpec,
	}
	if length > 0 {
		agg.aggregates = make([]FieldAggregator, length)
	}
	return agg
}

// FieldName returns field name.
func (a *seriesAggregator) FieldName() field.Name {
	return a.fieldName
}

// GetFieldType returns the field type.
func (a *seriesAggregator) GetFieldType() field.Type {
	return a.fieldType
}

// ResultSet returns the result set of series aggregator.
func (a *seriesAggregator) ResultSet() series.Iterator {
	return newSeriesIterator(a)
}

// Reset resets the aggregator's context for reusing.
func (a *seriesAggregator) Reset() {
	for _, agg := range a.aggregates {
		if agg != nil {
			agg.reset()
		}
	}
}

func (a *seriesAggregator) GetAggregates() []FieldAggregator {
	return a.aggregates
}

// GetAggregator gets field aggregator by segment start time, if not exist return (nil,false).
func (a *seriesAggregator) GetAggregator(segmentStartTime int64) (agg FieldAggregator, ok bool) {
	if segmentStartTime < a.startTime {
		return
	}
	idx := a.calc.CalcTimeWindows(a.startTime, segmentStartTime) - 1
	if idx < 0 || idx >= len(a.aggregates) {
		return
	}
	agg = a.aggregates[idx]
	if agg == nil {
		storageTimeRange := &timeutil.TimeRange{
			Start: segmentStartTime,
			End:   a.calc.CalcFamilyEndTime(segmentStartTime),
		}
		timeRange := a.queryTimeRange.Intersect(storageTimeRange)
		storageInterval := a.queryInterval.Int64() / int64(a.intervalRatio)
		startIdx := a.calc.CalcSlot(timeRange.Start, segmentStartTime, storageInterval)
		endIdx := a.calc.CalcSlot(timeRange.End, segmentStartTime, storageInterval)

		agg = NewFieldAggregator(a.aggSpec, segmentStartTime, startIdx, endIdx)
		a.aggregates[idx] = agg
	}
	ok = true
	return
}
