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
	"context"
	"encoding/binary"
	"sort"
	"sync"
	"time"

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

// TaskContext represents task execute context.
type TaskContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

// NewTaskContextWithTimeout creates a task context with timeout.
func NewTaskContextWithTimeout(ctx context.Context, timeout time.Duration) *TaskContext {
	c, cancel := context.WithTimeout(ctx, timeout)
	return &TaskContext{
		Ctx:    c,
		Cancel: cancel,
	}
}

// Release releases context's resource after query.
func (ctx *TaskContext) Release() {
	ctx.Cancel()
}

// StorageExecuteContext represents storage level query execute context.
type StorageExecuteContext struct {
	TaskCtx       *TaskContext
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

	Stats *models.StorageStats // storage query stats track for explain query

	mutex sync.Mutex
}

// collectGroupingTagValueIDs collects grouping tag value ids when does grouping operation.
func (ctx *StorageExecuteContext) collectGroupingTagValueIDs(tagValueIDs []uint32) {
	// need add lock, because build group concurrent
	ctx.mutex.Lock()
	for idx, tagValueID := range tagValueIDs {
		tIDs := ctx.GroupingTagValueIDs[idx]
		if tIDs == nil {
			ctx.GroupingTagValueIDs[idx] = roaring.BitmapOf(tagValueID)
		} else {
			ctx.GroupingTagValueIDs[idx].Add(tagValueID)
		}
	}
	ctx.mutex.Unlock()
}

// CalcQuerySlotRange returns slot range for query by family time and query time range.
func (ctx *StorageExecuteContext) CalcQuerySlotRange(familyTime int64) timeutil.SlotRange {
	return ctx.Query.Interval.CalcQuerySlotRange(familyTime, ctx.Query.TimeRange)
}

// HasGroupingTagValueIDs returns if it needs collect grouping tag value.
func (ctx *StorageExecuteContext) HasGroupingTagValueIDs() bool {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	for idx := range ctx.GroupingTagValueIDs {
		tIDs := ctx.GroupingTagValueIDs[idx]
		if tIDs != nil && !tIDs.IsEmpty() {
			return true
		}
	}
	return false
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

// Release releases context's resource after query.
func (ctx *StorageExecuteContext) Release() {
	for idx := range ctx.ShardContexts {
		shardCtx := ctx.ShardContexts[idx]
		if shardCtx != nil {
			shardCtx.Release()
		}
	}
	ctx.TaskCtx.Release()
}

// TagFilterResult represents the tag filter result, include tag key id and tag value ids.
type TagFilterResult struct {
	TagKeyID    tag.KeyID
	TagValueIDs *roaring.Bitmap
}

// TimeSegmentContext represents time segment context
type TimeSegmentContext struct {
	TimeSegments map[int64]*TimeSegmentResultSet // familyTime -> time segment result set list
	SeriesIDs    *roaring.Bitmap                 // matched series ids after data filter
}

// NewTimeSegmentContext creates a time segment context.
func NewTimeSegmentContext() *TimeSegmentContext {
	return &TimeSegmentContext{
		TimeSegments: make(map[int64]*TimeSegmentResultSet),
		SeriesIDs:    roaring.New(),
	}
}

// AddFilterResultSet adds a result set after data filtering.
func (ts *TimeSegmentContext) AddFilterResultSet(interval timeutil.Interval, rs FilterResultSet) {
	familyTime := rs.FamilyTime()
	segment, ok := ts.TimeSegments[familyTime]
	if !ok {
		segment = &TimeSegmentResultSet{
			FamilyTime: familyTime,
			Source:     rs.SlotRange(),
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

// GetTimeSegments returns
func (ts *TimeSegmentContext) GetTimeSegments() (rs TimeSegmentContexts) {
	for _, segment := range ts.TimeSegments {
		rs = append(rs, segment)
	}
	sort.Sort(rs)
	return rs
}

// Release releases time segment's data resource after query.
func (ts *TimeSegmentContext) Release() {
	for idx := range ts.TimeSegments {
		ts.TimeSegments[idx].Release()
	}
}

// ShardExecuteContext represents shard level query execute context.
type ShardExecuteContext struct {
	StorageExecuteCtx  *StorageExecuteContext
	TimeSegmentContext *TimeSegmentContext // result set for each time segment

	GroupingContext         GroupingContext // after get grouping context if it has grouping query
	SeriesIDsAfterFiltering *roaring.Bitmap // after data filter
}

// NewShardExecuteContext creates a shard execute context.
func NewShardExecuteContext(storageExecuteCtx *StorageExecuteContext) *ShardExecuteContext {
	return &ShardExecuteContext{
		StorageExecuteCtx:       storageExecuteCtx,
		SeriesIDsAfterFiltering: roaring.New(),
		TimeSegmentContext:      NewTimeSegmentContext(),
	}
}

func (ctx *ShardExecuteContext) IsSeriesIDsEmpty() bool {
	// maybe some series ids not write data in query time range,
	// so need reset series ids using ids which after data filtering.
	ctx.SeriesIDsAfterFiltering = ctx.TimeSegmentContext.SeriesIDs
	return ctx.SeriesIDsAfterFiltering.IsEmpty()
}

// Release releases shard context's resource after query.
func (ctx *ShardExecuteContext) Release() {
	if ctx.TimeSegmentContext != nil {
		ctx.TimeSegmentContext.Release()
	}
}

// GroupingSeriesAgg represents grouping series aggregator.
type GroupingSeriesAgg struct {
	Key         string
	Aggregator  aggregation.SeriesAggregator // for single field query
	Aggregators aggregation.FieldAggregates  // for multi fields query
}

// reduce aggregator's result set.
func (agg *GroupingSeriesAgg) reduce(reduceFn func(it series.GroupedIterator)) {
	if agg.Aggregator != nil {
		reduceFn(aggregation.FieldAggregates{agg.Aggregator}.ResultSet(agg.Key))
		// reset aggregate context
		agg.Aggregator.Reset()
	} else {
		reduceFn(agg.Aggregators.ResultSet(agg.Key))
		// reset aggregate context
		agg.Aggregators.Reset()
	}
}

// DataLoadContext represents data load level query execute context.
type DataLoadContext struct {
	ShardExecuteCtx *ShardExecuteContext

	MinSeriesID, MaxSeriesID uint16
	// range of min/max low series id
	// if no grouping value is low series ids
	// if grouping value is index of GroupingSeriesAgg
	LowSeriesIDs          []uint16
	SeriesIDHighKey       uint16
	LowSeriesIDsContainer roaring.Container

	GroupingSeriesAggRefs    []uint16 // series id => GroupingSeriesAgg index
	WithoutGroupingSeriesAgg *GroupingSeriesAgg
	GroupingSeriesAgg        []*GroupingSeriesAgg
	groupingSeriesAggRefIdx  uint16

	DownSampling func(slotRange timeutil.SlotRange, seriesIdx uint16, fieldIdx int, fieldData []byte)
}

// PrepareAggregatorWithoutGrouping prepares context for without grouping query.
func (ctx *DataLoadContext) PrepareAggregatorWithoutGrouping() {
	ctx.WithoutGroupingSeriesAgg = &GroupingSeriesAgg{
		Key: "",
	}
	if ctx.IsMultiField() {
		ctx.WithoutGroupingSeriesAgg.Aggregators = ctx.newSeriesAggregators()
	} else {
		ctx.WithoutGroupingSeriesAgg.Aggregator = ctx.newSeriesAggregator(0)
	}
}

// NewSeriesAggregator creates the series aggregator with grouping key for grouping query,
// returns index of grouping aggregator.
func (ctx *DataLoadContext) NewSeriesAggregator(groupingKey string) uint16 {
	rs := ctx.groupingSeriesAggRefIdx
	groupingSeriesAgg := &GroupingSeriesAgg{
		Key: groupingKey,
	}
	tagsData := []byte(groupingKey)
	var tagValueIDs []uint32
	for idx := range ctx.ShardExecuteCtx.StorageExecuteCtx.GroupByTagKeyIDs {
		offset := idx * 4
		tagValueID := binary.LittleEndian.Uint32(tagsData[offset:])
		tagValueIDs = append(tagValueIDs, tagValueID)
	}
	ctx.ShardExecuteCtx.StorageExecuteCtx.collectGroupingTagValueIDs(tagValueIDs)

	if ctx.IsMultiField() {
		groupingSeriesAgg.Aggregators = ctx.newSeriesAggregators()
	} else {
		groupingSeriesAgg.Aggregator = ctx.newSeriesAggregator(0)
	}
	ctx.GroupingSeriesAgg = append(ctx.GroupingSeriesAgg, groupingSeriesAgg)
	ctx.groupingSeriesAggRefIdx++
	return rs
}

// newSeriesAggregators creates the series aggregators for multi field.
func (ctx *DataLoadContext) newSeriesAggregators() []aggregation.SeriesAggregator {
	rs := make([]aggregation.SeriesAggregator, len(ctx.ShardExecuteCtx.StorageExecuteCtx.Fields))
	for fieldIdx := range ctx.ShardExecuteCtx.StorageExecuteCtx.Fields {
		rs[fieldIdx] = aggregation.NewSeriesAggregator(
			ctx.ShardExecuteCtx.StorageExecuteCtx.Query.Interval,
			ctx.ShardExecuteCtx.StorageExecuteCtx.Query.IntervalRatio,
			ctx.ShardExecuteCtx.StorageExecuteCtx.Query.TimeRange,
			ctx.ShardExecuteCtx.StorageExecuteCtx.DownSamplingSpecs[fieldIdx])
	}
	return rs
}

// newSeriesAggregator creates a series aggregator with field index.
func (ctx *DataLoadContext) newSeriesAggregator(fieldIdx int) aggregation.SeriesAggregator {
	return aggregation.NewSeriesAggregator(
		ctx.ShardExecuteCtx.StorageExecuteCtx.Query.Interval,
		ctx.ShardExecuteCtx.StorageExecuteCtx.Query.IntervalRatio,
		ctx.ShardExecuteCtx.StorageExecuteCtx.Query.TimeRange,
		ctx.ShardExecuteCtx.StorageExecuteCtx.DownSamplingSpecs[fieldIdx])
}

// IsMultiField returns if it is multi field select for query.
func (ctx *DataLoadContext) IsMultiField() bool {
	return len(ctx.ShardExecuteCtx.StorageExecuteCtx.Fields) > 1
}

// IsGrouping returns if it is grouping query.
func (ctx *DataLoadContext) IsGrouping() bool {
	return ctx.ShardExecuteCtx.StorageExecuteCtx.Query.HasGroupBy()
}

// HasGroupingData returns if it is grouping data.
func (ctx *DataLoadContext) HasGroupingData() bool {
	if ctx.IsGrouping() {
		return len(ctx.GroupingSeriesAgg) > 0
	}
	return true
}

// GetSeriesAggregator gets series aggregator with low series id and field index.
func (ctx *DataLoadContext) GetSeriesAggregator(lowSeriesIdx uint16, fieldIdx int) (result aggregation.SeriesAggregator) {
	var groupingSeriesAgg *GroupingSeriesAgg
	if ctx.IsGrouping() {
		groupingSeriesIdx := ctx.GroupingSeriesAggRefs[lowSeriesIdx]
		groupingSeriesAgg = ctx.GroupingSeriesAgg[groupingSeriesIdx]
	} else {
		groupingSeriesAgg = ctx.WithoutGroupingSeriesAgg
	}
	if ctx.IsMultiField() {
		return groupingSeriesAgg.Aggregators[fieldIdx]
	}
	return groupingSeriesAgg.Aggregator
}

// Grouping prepares context for grouping query.
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
			// match low series invoke callback
			fn(seriesIdxFromQuery, seriesIdxFromStorage)
		}
		seriesIdxFromStorage++
	}
}

// Reduce reduces down sampling result.
func (ctx *DataLoadContext) Reduce(reduceFn func(it series.GroupedIterator)) {
	if ctx.IsGrouping() {
		for _, groupAgg := range ctx.GroupingSeriesAgg {
			groupAgg.reduce(reduceFn)
		}
	} else {
		ctx.WithoutGroupingSeriesAgg.reduce(reduceFn)
	}
}

// TimeSegmentContexts represents the time segment slice in query time range.
type TimeSegmentContexts []*TimeSegmentResultSet

func (f TimeSegmentContexts) Len() int           { return len(f) }
func (f TimeSegmentContexts) Less(i, j int) bool { return f[i].FamilyTime < f[j].FamilyTime }
func (f TimeSegmentContexts) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

// TimeSegmentResultSet represents the time segment in query time range.
type TimeSegmentResultSet struct {
	FamilyTime int64
	Source     timeutil.SlotRange
	Interval   timeutil.Interval

	FilterRS []FilterResultSet
}

// Release releases filter result set's resource after query.
func (ctx *TimeSegmentResultSet) Release() {
	for idx := range ctx.FilterRS {
		ctx.FilterRS[idx].Close()
	}
}
