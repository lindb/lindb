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

//go:generate mockgen -source=./group_agg.go -destination=./group_agg_mock.go -package=aggregation

// GroupingAggregator represents an aggregator which merges time series and does grouping if need.
type GroupingAggregator interface {
	// Aggregate aggregates the time series data
	Aggregate(it series.GroupedIterator)
	// ResultSet returns the result set of aggregator
	ResultSet() series.GroupedIterators
	// TimeRange returns the time range of aggregator.
	TimeRange() timeutil.TimeRange
	// Interval returns the time interval of aggregator.
	Interval() timeutil.Interval
	// Fields returns all fields.
	Fields() []field.Name
}

// groupingAggregator implements GroupingAggregator interface.
type groupingAggregator struct {
	aggSpecs      AggregatorSpecs
	interval      timeutil.Interval
	intervalRatio int
	timeRange     timeutil.TimeRange
	aggregates    map[string]FieldAggregates // tag values => field aggregates
	fields        map[field.Name]field.Name
}

// NewGroupingAggregator creates a grouping aggregator
func NewGroupingAggregator(
	interval timeutil.Interval,
	intervalRatio int,
	timeRange timeutil.TimeRange,
	aggSpecs AggregatorSpecs,
) GroupingAggregator {
	return &groupingAggregator{
		aggSpecs:      aggSpecs,
		interval:      interval,
		intervalRatio: intervalRatio,
		timeRange:     timeRange,
		aggregates:    make(map[string]FieldAggregates),
		fields:        make(map[field.Name]field.Name),
	}
}

// Aggregate aggregates the time series data.
func (ga *groupingAggregator) Aggregate(it series.GroupedIterator) {
	seriesAgg := ga.getAggregator(it.Tags())
	var sAgg SeriesAggregator
	for it.HasNext() {
		seriesIt := it.Next()
		fieldName := seriesIt.FieldName()
		ga.fields[fieldName] = fieldName
		// 1. find field aggregator
		sAgg = nil
		for _, aggregator := range seriesAgg {
			if aggregator.FieldName() == fieldName {
				sAgg = aggregator
				break
			}
		}
		if sAgg == nil {
			continue
		}
		// 2. merge the field series data
		for seriesIt.HasNext() {
			startTime, fieldIt := seriesIt.Next()
			if fieldIt == nil {
				continue
			}
			aggregator := sAgg.getAggregator(startTime)
			aggregator.Aggregate(fieldIt)
		}
	}
}

// ResultSet returns the result set of aggregator.
func (ga *groupingAggregator) ResultSet() series.GroupedIterators {
	length := len(ga.aggregates)
	if length == 0 {
		return nil
	}
	seriesList := make([]series.GroupedIterator, length)
	idx := 0
	for tags, aggregator := range ga.aggregates {
		seriesList[idx] = aggregator.ResultSet(tags)
		idx++
	}
	return seriesList
}

// TimeRange return the time range of aggregator.
func (ga *groupingAggregator) TimeRange() timeutil.TimeRange {
	return ga.timeRange
}

// Interval return the time interval of aggregator.
func (ga *groupingAggregator) Interval() timeutil.Interval {
	return ga.interval
}

// Fields returns all fields.
func (ga *groupingAggregator) Fields() []field.Name {
	length := len(ga.fields)
	if length == 0 {
		return nil
	}
	rs := make([]field.Name, length)
	idx := 0
	for fieldName := range ga.fields {
		rs[idx] = fieldName
		idx++
	}
	return rs
}

// getAggregator returns the time series aggregator by the tag of time series.
func (ga *groupingAggregator) getAggregator(tags string) (agg FieldAggregates) {
	// get series aggregator
	if agg0, ok := ga.aggregates[tags]; ok {
		return agg0
	}
	agg = NewFieldAggregates(ga.interval, ga.intervalRatio, ga.timeRange, ga.aggSpecs)
	ga.aggregates[tags] = agg
	return
}
