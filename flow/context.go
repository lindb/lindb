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

package flow

import (
	"sort"
	"sync"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
)

// StorageExecuteContext represents storage query execute context.
type StorageExecuteContext struct {
	Query         *stmt.Query
	ShardIDs      []models.ShardID
	ShardContexts []*ShardExecuteContext

	// set value in plan stage when lookup table.
	MetricID metric.ID

	// set value in plan stage when lookup select fields.
	Fields            field.Metas
	DownSamplingSpecs aggregation.AggregatorSpecs
	AggregatorSpecs   aggregation.AggregatorSpecs

	// result which after tag condition metadata filter
	// set value in tag search, the where clause condition that user input
	// first find all tag values in where clause, then do tag match
	TagFilterResult map[string]*TagFilterResult

	// set value in plan stage when lookup group by tags.
	GroupByTags      tag.Metas
	GroupByTagKeyIDs []tag.KeyID
	// for group by query store tag value ids for each group tag key
	GroupingTagValueIDs []*roaring.Bitmap
	mutex               sync.Mutex

	// Time range and interval
	QueryTimeRange     timeutil.TimeRange
	QueryInterval      timeutil.Interval
	QueryIntervalRatio int

	Stats *models.StorageStats // storage query stats track for explain query
}

// CollectGroupingTagValueIDs collects grouping tag value ids when does grouping operation.
func (ctx *StorageExecuteContext) CollectGroupingTagValueIDs(groupingTagValueIDs []*roaring.Bitmap) {
	// need add lock, because build group concurrent
	ctx.mutex.Lock()
	for idx, tagValueIDs := range groupingTagValueIDs {
		tIDs := ctx.GroupingTagValueIDs[idx]
		if tIDs == nil {
			ctx.GroupingTagValueIDs[idx] = tagValueIDs
		} else {
			ctx.GroupingTagValueIDs[idx].Or(tagValueIDs)
		}
	}
	ctx.mutex.Unlock()
}

// HasWhereCondition returns if query has where clause condition.
func (ctx *StorageExecuteContext) HasWhereCondition() bool {
	return ctx.Query.Condition != nil
}

// SortFields sorts fields by field ids for reading data in order.
func (ctx *StorageExecuteContext) SortFields() {
	sort.Slice(ctx.Fields, func(i, j int) bool {
		return ctx.Fields[i].ID < ctx.Fields[j].ID
	})
}

// QueryStats returns the storage query stats.
func (ctx *StorageExecuteContext) QueryStats() *models.StorageStats {
	if ctx.Stats != nil {
		ctx.Stats.Complete()
	}
	return ctx.Stats
}

func (ctx *StorageExecuteContext) Release() {
	for idx := range ctx.ShardContexts {
		ctx.ShardContexts[idx].Release()
	}
}

// TagFilterResult represents the tag filter result, include tag key id and tag value ids.
type TagFilterResult struct {
	TagKeyID    tag.KeyID
	TagValueIDs *roaring.Bitmap
}

type TimeSegmentResultSet struct {
	TimeSegments map[int64]*TimeSegmentContext // familyTime -> timeSpan todo need fix
	SeriesIDs    *roaring.Bitmap
}

func NewTimeSegmentResultSet() *TimeSegmentResultSet {
	return &TimeSegmentResultSet{
		TimeSegments: make(map[int64]*TimeSegmentContext),
		SeriesIDs:    roaring.New(),
	}
}

func (ts *TimeSegmentResultSet) AddFilterResultSet(interval timeutil.Interval, rs FilterResultSet) {
	familyTime := rs.FamilyTime()
	segment, ok := ts.TimeSegments[familyTime]
	if !ok {
		segment = &TimeSegmentContext{
			FamilyTime: familyTime,
			Source:     rs.SlotRange(),
			Target:     rs.SlotRange(),
			Interval:   interval,
		}
		ts.TimeSegments[familyTime] = segment
	} else {
		// calc source slot range
		segment.Source = segment.Source.Union(rs.SlotRange())
	}

	segment.FilterRS = append(segment.FilterRS, rs)

	// merge all series ids after filtering => final series ids
	ts.SeriesIDs.Or(rs.SeriesIDs())
}

func (ts *TimeSegmentResultSet) GetTimeSegments() (rs TimeSegmentContexts) {
	for _, segment := range ts.TimeSegments {
		rs = append(rs, segment)
	}
	sort.Sort(rs)
	return rs
}

func (ts *TimeSegmentResultSet) Release() {
	for idx := range ts.TimeSegments {
		ts.TimeSegments[idx].Release()
	}
}

type ShardExecuteContext struct {
	StorageExecuteCtx *StorageExecuteContext

	TimeSegmentRS *TimeSegmentResultSet // result set for each time segment

	SeriesIDsAfterFiltering *roaring.Bitmap

	GroupingContext GroupingContext
}

func NewShardExecuteContext(storageExecuteCtx *StorageExecuteContext) *ShardExecuteContext {
	return &ShardExecuteContext{
		StorageExecuteCtx:       storageExecuteCtx,
		SeriesIDsAfterFiltering: roaring.New(),
		TimeSegmentRS:           NewTimeSegmentResultSet(),
	}
}

func (ctx *ShardExecuteContext) IsSeriesIDsEmpty() bool {
	// maybe some series ids not write data in query time range,
	// so need reset series ids using ids which after data filtering.
	ctx.SeriesIDsAfterFiltering = ctx.TimeSegmentRS.SeriesIDs
	return ctx.SeriesIDsAfterFiltering.IsEmpty()
}

func (ctx *ShardExecuteContext) Release() {
	ctx.TimeSegmentRS.Release()
}

type GroupingSeriesAgg struct {
	Key        string
	Aggregator aggregation.SeriesAggregator
}

type DataLoadContext struct {
	ShardExecuteCtx *ShardExecuteContext

	MinSeriesID, MaxSeriesID uint16
	// range of min/max low series id
	// if no grouping value is low series ids
	// if grouping value is index of GroupingSeriesAgg
	LowSeriesIDs          []uint16
	SeriesIDHighKey       uint16
	LowSeriesIDsContainer roaring.Container

	// time segment => a list of DataLoader(each family)
	Loaders [][]DataLoader // item maybe DataLoader is nil

	GroupingSeriesAggRefs []uint16 // series id => GroupingSeriesAgg index
	SingleFieldAgg        aggregation.SeriesAggregator
	GroupingSeriesAgg     []*GroupingSeriesAgg

	DownSampling func(slotRange timeutil.SlotRange, seriesIdx uint16, fieldIdx int, fieldData []byte)
}

func (ctx *DataLoadContext) PrepareAggregator() {
	if ctx.IsMultiField() {
		panic("need impl")
	} else {
		ctx.SingleFieldAgg = ctx.NewSeriesAggregator()
	}
}

func (ctx *DataLoadContext) NewSeriesAggregator() aggregation.SeriesAggregator {
	return aggregation.NewSeriesAggregator(
		ctx.ShardExecuteCtx.StorageExecuteCtx.QueryInterval,
		ctx.ShardExecuteCtx.StorageExecuteCtx.QueryIntervalRatio,
		ctx.ShardExecuteCtx.StorageExecuteCtx.QueryTimeRange,
		ctx.ShardExecuteCtx.StorageExecuteCtx.DownSamplingSpecs[0])
}

func (ctx *DataLoadContext) IsMultiField() bool {
	return len(ctx.ShardExecuteCtx.StorageExecuteCtx.Fields) > 1
}

func (ctx *DataLoadContext) IsGrouping() bool {
	return ctx.ShardExecuteCtx.StorageExecuteCtx.Query.HasGroupBy()
}

func (ctx *DataLoadContext) Grouping() {
	min := ctx.LowSeriesIDsContainer.Minimum()
	ctx.MinSeriesID = min
	ctx.MaxSeriesID = ctx.LowSeriesIDsContainer.Maximum()
	lengthOfSeriesIDs := int(ctx.MaxSeriesID-ctx.MinSeriesID) + 1
	ctx.LowSeriesIDs = make([]uint16, lengthOfSeriesIDs)
	if ctx.IsGrouping() {
		ctx.GroupingSeriesAggRefs = make([]uint16, lengthOfSeriesIDs)
	}
	it := ctx.LowSeriesIDsContainer.PeekableIterator()
	for it.HasNext() {
		lowSeriesID := it.Next()
		seriesIdx := lowSeriesID - min
		ctx.LowSeriesIDs[seriesIdx] = lowSeriesID
	}
}

// IterateLowSeriesIDs iterates low series ids from storage, then found low series id which query need.
func (ctx *DataLoadContext) IterateLowSeriesIDs(lowSeriesIDsFromStorage roaring.Container,
	fn func(seriesIdxFromQuery uint16, seriesIdxFromStorage int),
) {
	min := ctx.MinSeriesID
	max := ctx.MaxSeriesID
	lowSeriesIDs := ctx.LowSeriesIDs
	it := lowSeriesIDsFromStorage.PeekableIterator()
	seriesIdxFromStorage := 0
	for it.HasNext() {
		seriesID := it.Next()
		if seriesID > max {
			break
		}
		if seriesID < min {
			seriesIdxFromStorage++
			continue
		}
		seriesIdxFromQuery := seriesID - min
		if lowSeriesIDs[seriesIdxFromQuery] == seriesID {
			// load data by series id index
			fn(seriesIdxFromQuery, seriesIdxFromStorage)
		}
		seriesIdxFromStorage++
	}
}

// Reduce reduces down sampling result.
func (ctx *DataLoadContext) Reduce(reduceFn func(it series.GroupedIterator)) {
	if ctx.IsGrouping() {
		for _, groupAgg := range ctx.GroupingSeriesAgg {
			reduceFn(aggregation.FieldAggregates{groupAgg.Aggregator}.ResultSet(groupAgg.Key))
			// reset aggregate context
			groupAgg.Aggregator.Reset()
		}
	} else {
		reduceFn(aggregation.FieldAggregates{ctx.SingleFieldAgg}.ResultSet(""))
		// reset aggregate context
		ctx.SingleFieldAgg.Reset()
	}
}

// TimeSegmentContexts represents the time segment slice in query time range.
type TimeSegmentContexts []*TimeSegmentContext

func (f TimeSegmentContexts) Len() int           { return len(f) }
func (f TimeSegmentContexts) Less(i, j int) bool { return f[i].FamilyTime < f[j].FamilyTime }
func (f TimeSegmentContexts) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

// TimeSegmentContext represents the time segment in query time range.
type TimeSegmentContext struct {
	FamilyTime     int64
	Source, Target timeutil.SlotRange
	Interval       timeutil.Interval

	FilterRS []FilterResultSet
}

func (ctx *TimeSegmentContext) Release() {
	for idx := range ctx.FilterRS {
		ctx.FilterRS[idx].Close()
	}
}
