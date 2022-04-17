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

package storagequery

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

// baseQueryTask represents base query task stats track for task execute cost
type baseQueryTask struct {
	start time.Time
	cost  time.Duration
}

// BeforeRun invokes before task run function
func (t *baseQueryTask) BeforeRun() {
	t.start = time.Now()
}

// Run executes task logic
func (t *baseQueryTask) Run() error {
	return nil
}

// AfterRun invokes after task run function
func (t *baseQueryTask) AfterRun() {
	t.cost = time.Since(t.start)
}

// queryStatTask represents the query stat task
type queryStatTask struct {
	task flow.QueryTask
}

// BeforeRun invokes before task run function
func (t *queryStatTask) BeforeRun() {
}

// Run executes query cost stat
func (t *queryStatTask) Run() error {
	t.task.BeforeRun()
	defer func() {
		t.task.AfterRun()
	}()
	return t.task.Run()
}

// AfterRun invokes after task run function
func (t *queryStatTask) AfterRun() {
}

// storagePlanTask represents storage execute plan task
type storagePlanTask struct {
	baseQueryTask
	ctx  *executeContext
	plan *storageExecutePlan
}

// newStoragePlanTask creates storage execute plan task
func newStoragePlanTask(ctx *executeContext, plan *storageExecutePlan) flow.QueryTask {
	task := &storagePlanTask{
		ctx:  ctx,
		plan: plan,
	}
	if ctx.storageExecuteCtx.Query.Explain {
		// if it needs to explain query, use queryStatTask
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes storage execute plan
func (t *storagePlanTask) Run() error {
	return t.plan.Plan()
}

// AfterRun invokes after execute plan, collects plan stats
func (t *storagePlanTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.ctx.storageExecuteCtx.Stats.SetPlanCost(t.cost)
}

// tagFilterTask represents tag filtering task based on where condition
type tagFilterTask struct {
	baseQueryTask
	ctx       *executeContext
	tagSearch TagSearch
}

// newTagFilterTask creates tag filtering task
func newTagFilterTask(ctx *executeContext, tagSearch TagSearch) flow.QueryTask {
	task := &tagFilterTask{
		ctx:       ctx,
		tagSearch: tagSearch,
	}
	if ctx.storageExecuteCtx.Query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes tag filtering based on where condition
func (t *tagFilterTask) Run() error {
	err := t.tagSearch.Filter()
	if err != nil {
		return err
	}
	if len(t.ctx.storageExecuteCtx.TagFilterResult) == 0 {
		// filter not match, return not found
		return constants.ErrNotFound
	}
	return nil
}

// AfterRun invokes after tag filtering, collects tag filtering stats
func (t *tagFilterTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.ctx.storageExecuteCtx.Stats.SetTagFilterCost(t.cost)
}

// seriesIDsSearchTask represents series ids search task based on tag filtering result set
type seriesIDsSearchTask struct {
	baseQueryTask

	shardExecuteContext *flow.ShardExecuteContext
	shard               tsdb.Shard
}

// newSeriesIDsSearchTask creates series ids search task
func newSeriesIDsSearchTask(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
	task := &seriesIDsSearchTask{
		shardExecuteContext: shardExecuteContext,
		shard:               shard,
	}
	if shardExecuteContext.StorageExecuteCtx.Query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes series ids search based on tag filtering result
func (t *seriesIDsSearchTask) Run() (err error) {
	queryStmt := t.shardExecuteContext.StorageExecuteCtx.Query
	condition := queryStmt.Condition
	var seriesIDs *roaring.Bitmap
	if condition != nil {
		// if it gets tag filter result do series ids searching
		seriesSearch := newSeriesSearchFunc(t.shard.IndexDatabase(), t.shardExecuteContext.StorageExecuteCtx.TagFilterResult, condition)
		seriesIDs, err = seriesSearch.Search()
	} else {
		// get series ids for metric level
		seriesIDs, err = t.shard.IndexDatabase().GetSeriesIDsForMetric(queryStmt.Namespace, queryStmt.MetricName)
		if err == nil && !queryStmt.HasGroupBy() {
			// add series id without tags, maybe metric has too many series, but one series without tags
			seriesIDs.Add(series.IDWithoutTags)
		}
	}
	if err == nil && seriesIDs != nil {
		t.shardExecuteContext.SeriesIDsAfterFiltering.Or(seriesIDs)
	}
	return
}

// AfterRun invokes after series ids search, collects the series ids search stats
func (t *seriesIDsSearchTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.shardExecuteContext.StorageExecuteCtx.Stats.SetShardSeriesIDsSearchStats(t.shard.ShardID(),
		t.shardExecuteContext.SeriesIDsAfterFiltering.GetCardinality(),
		t.cost)
}

// familyFilterTask represents family data filtering task
type familyFilterTask struct {
	baseQueryTask

	shardExecuteContext *flow.ShardExecuteContext
	shard               tsdb.Shard
}

// newFamilyFilterTask creates family data filtering task
func newFamilyFilterTask(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
	task := &familyFilterTask{
		shardExecuteContext: shardExecuteContext,
		shard:               shard,
	}
	if shardExecuteContext.StorageExecuteCtx.Query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes file data filtering based on series ids and time range for each data family
func (t *familyFilterTask) Run() error {
	families := t.shard.GetDataFamilies(t.shardExecuteContext.StorageExecuteCtx.Query.StorageInterval.Type(),
		t.shardExecuteContext.StorageExecuteCtx.Query.TimeRange)
	if len(families) == 0 {
		return nil
	}
	for idx := range families {
		family := families[idx]
		// execute data family search in background goroutine
		resultSet, err := family.Filter(t.shardExecuteContext)
		if errors.Is(err, constants.ErrNotFound) {
			continue
		}
		if err != nil {
			return err
		}
		for _, rs := range resultSet {
			t.shardExecuteContext.TimeSegmentContext.AddFilterResultSet(family.Interval(), rs)
		}
	}
	return nil
}

// AfterRun invokes after file data filtering, collects the file data filtering stats
func (t *familyFilterTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.shardExecuteContext.StorageExecuteCtx.Stats.SetShardKVDataFilterCost(t.shard.ShardID(), t.cost)
}

// groupingContextFindTask represents group by context find task
type groupingContextFindTask struct {
	baseQueryTask

	executeCtx *flow.ShardExecuteContext
	shard      tsdb.Shard
}

// newGroupingContextFindTask creates the group by context find task
func newGroupingContextFindTask(executeCtx *flow.ShardExecuteContext, shard tsdb.Shard) flow.QueryTask {
	task := &groupingContextFindTask{
		executeCtx: executeCtx,
		shard:      shard,
	}
	if executeCtx.StorageExecuteCtx.Query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes group by context finding based on group by tag key ids
func (t *groupingContextFindTask) Run() error {
	err := t.shard.IndexDatabase().GetGroupingContext(t.executeCtx)
	if err != nil {
		return err
	}

	return nil
}

// AfterRun invokes after group by context, collects the find group by context stats
func (t *groupingContextFindTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.executeCtx.StorageExecuteCtx.Stats.SetShardGroupingCost(t.shard.ShardID(), t.cost)
}

// buildGroupTask represents build grouped tag value ids => series ids mapping
type buildGroupTask struct {
	baseQueryTask
	shard   tsdb.Shard
	loadCtx *flow.DataLoadContext
}

// newBuildGroupTask creates build group task
func newBuildGroupTask(shard tsdb.Shard, loadCtx *flow.DataLoadContext) flow.QueryTask {
	task := &buildGroupTask{
		shard:   shard,
		loadCtx: loadCtx,
	}
	if loadCtx.ShardExecuteCtx.StorageExecuteCtx.Query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes grouped series ids(tag value ids=>series ids mapping)
func (t *buildGroupTask) Run() error {
	t.loadCtx.Grouping()
	if t.loadCtx.ShardExecuteCtx.GroupingContext != nil {
		// build group by data, grouped series: tags => series IDs
		t.loadCtx.ShardExecuteCtx.GroupingContext.BuildGroup(t.loadCtx)
	} else {
		t.loadCtx.PrepareAggregatorWithoutGrouping()
	}
	return nil
}

// AfterRun invokes after build grouped series, collects build stats
func (t *buildGroupTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.loadCtx.ShardExecuteCtx.StorageExecuteCtx.Stats.SetShardGroupBuildStats(t.shard.ShardID(), t.cost)
}

// dataLoadTask represents data load task based on filtering result set
type dataLoadTask struct {
	baseQueryTask
	dataLoadCtx *flow.DataLoadContext
	shard       tsdb.Shard
	queryFlow   flow.StorageQueryFlow
	segmentIdx  int
	segmentCtx  *flow.TimeSegmentResultSet

	costs []time.Duration
}

// newDataLoadTask creates the data load task
func newDataLoadTask(
	shard tsdb.Shard,
	queryFlow flow.StorageQueryFlow,
	dataLoadCtx *flow.DataLoadContext,
	segmentIdx int,
	segmentCtx *flow.TimeSegmentResultSet,
) flow.QueryTask {
	task := &dataLoadTask{
		shard:       shard,
		queryFlow:   queryFlow,
		segmentIdx:  segmentIdx,
		dataLoadCtx: dataLoadCtx,
		segmentCtx:  segmentCtx,
	}
	if dataLoadCtx.ShardExecuteCtx.StorageExecuteCtx.Query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes data load based on filtering result set
func (t *dataLoadTask) Run() error {
	explain := t.dataLoadCtx.ShardExecuteCtx.StorageExecuteCtx.Query.Explain
	if explain {
		t.costs = make([]time.Duration, len(t.segmentCtx.FilterRS))
	}
	queryIntervalRatio := t.dataLoadCtx.ShardExecuteCtx.StorageExecuteCtx.Query.IntervalRatio
	seriesIDs := t.dataLoadCtx.ShardExecuteCtx.SeriesIDsAfterFiltering // after group result
	targetSlotRange := t.dataLoadCtx.ShardExecuteCtx.StorageExecuteCtx.CalcTargetSlotRange(t.segmentCtx.FamilyTime)

	for idx, rs := range t.segmentCtx.FilterRS {
		// double filtering, maybe some series ids be filtered out when do grouping.
		// filter logic: forward_reader.go -> GetGroupingScanner
		if roaring.FastAnd(seriesIDs, rs.SeriesIDs()).IsEmpty() {
			continue
		}
		// maybe return nil loader
		var start time.Time
		if explain {
			start = time.Now()
		}

		loader := rs.Load(t.dataLoadCtx)
		if loader == nil {
			continue
		}

		// load field series data by series ids
		tsdDecoder := encoding.GetTSDDecoder()
		t.dataLoadCtx.DownSampling = func(slotRange timeutil.SlotRange, lowSeriesIdx uint16, fieldIdx int, fieldData []byte) {
			var agg aggregation.FieldAggregator
			seriesAggregator := t.dataLoadCtx.GetSeriesAggregator(lowSeriesIdx, fieldIdx)

			var ok bool
			agg, ok = seriesAggregator.GetAggregator(t.segmentCtx.FamilyTime)
			if !ok {
				return
			}
			tsdDecoder.ResetWithTimeRange(fieldData, slotRange.Start, slotRange.End)
			aggregation.DownSamplingSeries(
				targetSlotRange, uint16(queryIntervalRatio), 0, // same family, base slot = 0
				tsdDecoder,
				agg.AggregateBySlot,
			)
		}

		// loads the metric data by given series id from load result.
		// if found data need to do down sampling aggregate.
		loader.Load(t.dataLoadCtx)
		// release tsd decoder back to pool for re-use.
		encoding.ReleaseTSDDecoder(tsdDecoder)
		// after load, need to reduce the aggregator's result to query flow.
		t.dataLoadCtx.Reduce(t.queryFlow.Reduce)

		if explain {
			t.costs[idx] = time.Since(start)
			start = time.Now()
		}
	}
	return nil
}

// AfterRun invokes after data load, collects the data load stats
func (t *dataLoadTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	for idx, rs := range t.segmentCtx.FilterRS {
		identifiers := strings.Split(rs.Identifier(), fmt.Sprintf("shard/%d/segment", t.shard.ShardID()))
		var identifier string
		if len(identifiers) > 1 {
			identifier = identifiers[1]
		} else {
			identifier = identifiers[0]
		}
		foundSeries := 0
		lowContainer := rs.SeriesIDs().GetContainer(t.dataLoadCtx.SeriesIDHighKey)
		if lowContainer != nil {
			foundSeries = lowContainer.GetCardinality()
		}
		t.dataLoadCtx.ShardExecuteCtx.StorageExecuteCtx.Stats.
			SetShardScanStats(t.shard.ShardID(), identifier, t.costs[idx], foundSeries)
	}
}

// collectTagValuesTask represents collect tag values by tag value ids
type collectTagValuesTask struct {
	baseQueryTask
	ctx         *executeContext
	metadata    metadb.Metadata
	tagKey      tag.Meta
	tagValueIDs *roaring.Bitmap
	tagValues   map[uint32]string
}

// newCollectTagValuesTask creates the collect tag values task
func newCollectTagValuesTask(ctx *executeContext, metadata metadb.Metadata,
	tagKey tag.Meta, tagValueIDs *roaring.Bitmap, tagValues map[uint32]string,
) flow.QueryTask {
	task := &collectTagValuesTask{
		ctx:         ctx,
		metadata:    metadata,
		tagKey:      tagKey,
		tagValueIDs: tagValueIDs,
		tagValues:   tagValues,
	}
	if ctx.storageExecuteCtx.Query.Explain {
		return &queryStatTask{
			task: task,
		}
	}
	return task
}

// Run executes collect tag values by ids
func (t *collectTagValuesTask) Run() error {
	return t.metadata.TagMetadata().CollectTagValues(t.tagKey.ID, t.tagValueIDs, t.tagValues)
}

// AfterRun invokes after tag value collect, collects execution stats
func (t *collectTagValuesTask) AfterRun() {
	t.baseQueryTask.AfterRun()
	t.ctx.storageExecuteCtx.Stats.SetCollectTagValuesStats(t.tagKey.Key, t.cost)
}
