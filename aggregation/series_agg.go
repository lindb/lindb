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
	"github.com/lindb/lindb/aggregation/selector"
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

//// Reset resets the aggregator's context for reusing
//func (agg FieldAggregates) Reset() {
//	for _, aggregator := range agg {
//		aggregator.Reset()
//	}
//}
//
// NewFieldAggregates creates the field aggregates based on aggregator specs and query time range.
// NOTICE: if do down sampling aggregator, aggregator specs must be in order by field id.
func NewFieldAggregates(
	queryInterval timeutil.Interval,
	queryTimeRange timeutil.TimeRange,
	aggSpecs AggregatorSpecs,
) FieldAggregates {
	aggregates := make(FieldAggregates, len(aggSpecs))
	for idx, aggSpec := range aggSpecs {
		aggregates[idx] = NewSeriesAggregator(queryInterval, queryTimeRange, aggSpec)
	}
	return aggregates
}

// SeriesAggregator represents a series aggregator which aggregates one field of a time series
type SeriesAggregator interface {
	// FieldName returns field name
	FieldName() field.Name
	// GetFieldType returns field type
	GetFieldType() field.Type
	SetFieldType(fieldType field.Type)

	GetFiledAggregator() FieldAggregator
	ResultSet() series.Iterator
}

type seriesAggregator struct {
	fieldName field.Name
	fieldType field.Type
	//ratio          int
	aggregator     FieldAggregator
	queryInterval  timeutil.Interval
	queryTimeRange timeutil.TimeRange
	aggSpec        AggregatorSpec
	calc           timeutil.Calculator

	startTime int64
}

//
// NewSeriesAggregator creates a series aggregator
func NewSeriesAggregator(
	queryInterval timeutil.Interval,
	queryTimeRange timeutil.TimeRange,
	aggSpec AggregatorSpec,
) SeriesAggregator {
	calc := queryInterval.Calculator()
	segmentTime := calc.CalcSegmentTime(queryTimeRange.Start)
	startTime := calc.CalcFamilyStartTime(segmentTime, calc.CalcFamily(queryTimeRange.Start, segmentTime))

	agg := &seriesAggregator{
		fieldName:      aggSpec.FieldName(),
		fieldType:      aggSpec.GetFieldType(),
		startTime:      startTime,
		calc:           calc,
		queryInterval:  queryInterval,
		queryTimeRange: queryTimeRange,
		aggSpec:        aggSpec,
	}
	//TODO maybe agg type has duplicate?
	var aggTypes []field.AggType
	for f := range aggSpec.Functions() {
		aggTypes = append(aggTypes, aggSpec.GetFieldType().GetFuncFieldParams(f)...)
	}
	//FIXME(stone1100) need remove it
	if len(aggTypes) == 0 {
		aggTypes = append(aggTypes, field.Sum)
	}
	//storageInterval := queryInterval.Int64() / int64(1)
	//startIdx := calc.CalcSlot(queryTimeRange.Start, segmentTime, storageInterval)
	//endIdx := calc.CalcSlot(queryTimeRange.End, segmentTime, storageInterval) + 1
	agg.aggregator = NewFieldAggregator(aggTypes, queryTimeRange.Start, selector.NewIndexSlotSelector(0, 360, 1))
	return agg
}

// FieldName returns field name
func (a *seriesAggregator) FieldName() field.Name {
	return a.fieldName
}

// GetFieldType returns the field type
func (a *seriesAggregator) GetFieldType() field.Type {
	return a.fieldType
}

func (a *seriesAggregator) GetFiledAggregator() FieldAggregator {
	return a.aggregator
}

// SetFieldType sets field type
func (a *seriesAggregator) SetFieldType(fieldType field.Type) {
	a.fieldType = fieldType
}

// ResultSet returns the result set of series aggregator
func (a *seriesAggregator) ResultSet() series.Iterator {
	return newSeriesIterator(a)
}

// Reset resets the aggregator's context for reusing
//func (a *seriesAggregator) Reset() {
//	a.aggregator.reset()
//	//for _, aggregator := range a.aggregates {
//	//	if aggregator == nil {
//	//		continue
//	//	}
//	//	aggregator.reset()
//	//}
//}
//
//// GetAggregator gets field aggregator by segment start time, if not exist return (nil,false).
//func (a *seriesAggregator) GetAggregateBlock(segmentStartTime int64) (agg series.Block, ok bool) {
//	if segmentStartTime < a.startTime {
//		return
//	}
//	idx := a.calc.CalcTimeWindows(a.startTime, segmentStartTime) - 1
//	return a.aggregator.GetBlock(idx, func() series.Block {
//		storageTimeRange := &timeutil.TimeRange{
//			Start: segmentStartTime,
//			End:   a.calc.CalcFamilyEndTime(segmentStartTime),
//		}
//		timeRange := a.queryTimeRange.Intersect(storageTimeRange)
//		storageInterval := a.queryInterval.Int64() / int64(a.ratio)
//		startIdx := a.calc.CalcSlot(timeRange.Start, segmentStartTime, storageInterval)
//		endIdx := a.calc.CalcSlot(timeRange.End, segmentStartTime, storageInterval) + 1
//		return series.NewBlock(startIdx, endIdx)
//	})
//}
