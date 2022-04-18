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

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/tsdb"
)

// for testing
var (
	newTagSearchFunc               = newTagSearch
	newStorageExecutePlanFunc      = newStorageExecutePlan
	newSeriesSearchFunc            = newSeriesSearch
	newBuildGroupTaskFunc          = newBuildGroupTask
	newDataLoadTaskFunc            = newDataLoadTask
	newStoragePlanTaskFunc         = newStoragePlanTask
	newSeriesIDsSearchTaskFunc     = newSeriesIDsSearchTask
	newTagFilterTaskFunc           = newTagFilterTask
	newFamilyFilterTaskFunc        = newFamilyFilterTask
	newGroupingContextFindTaskFunc = newGroupingContextFindTask
	newCollectTagValuesTaskFunc    = newCollectTagValuesTask
)

var (
	errNoShardID         = errors.New("there is no shard id in search condition")
	errNoShardInDatabase = errors.New("there is no shard in database storage engine")
	errShardNotFound     = errors.New("shard not found in database storage engine")
	errShardNumNotMatch  = errors.New("got shard size not equals input shard size")
)

// storageExecutor represents execution search logic in storage level,
// does query task async, then merge result, such as map-reduce job.
// 1) Filtering
// 2) Grouping if it needs
// 3) Scanning and Loading
// 4) Down sampling
// 5) Simple aggregation
type storageExecutor struct {
	ctx *executeContext

	queryFlow flow.StorageQueryFlow
	track     *groupingExecuteTrack
}

// newStorageMetricQuery creates the execution which queries the data of storage engine
func newStorageMetricQuery(
	queryFlow flow.StorageQueryFlow,
	executeCtx *executeContext,
) storageMetricQuery {
	return &storageExecutor{
		ctx:       executeCtx,
		queryFlow: queryFlow,
		track:     newGroupingExecuteTrack(executeCtx, queryFlow),
	}
}

// Execute executes search logic in storage level,
// 1) validation input params
// 2) build execute plan
// 3) build execute pipeline
// 4) run pipeline
func (e *storageExecutor) Execute() {
	if err := e.ctx.prepare(); err != nil {
		e.queryFlow.Complete(err)
		return
	}

	plan := newStorageExecutePlanFunc(e.ctx)
	t := newStoragePlanTaskFunc(e.ctx, plan)

	if err := t.Run(); err != nil {
		e.queryFlow.Complete(err)
		return
	}

	if e.ctx.storageExecuteCtx.HasWhereCondition() {
		tagSearch := newTagSearchFunc(e.ctx)
		t = newTagFilterTaskFunc(e.ctx, tagSearch)
		if err := t.Run(); err != nil {
			if errors.Is(constants.ErrNotFound, err) {
				e.queryFlow.Complete(nil)
			} else {
				e.queryFlow.Complete(err)
			}
			return
		}
	}

	// prepare storage query flow
	e.queryFlow.Prepare()

	// execute query flow
	e.executeQuery()
}

// executeQuery executes query flow for each shard
func (e *storageExecutor) executeQuery() {
	shards := e.ctx.shards
	e.ctx.storageExecuteCtx.ShardContexts = make([]*flow.ShardExecuteContext, len(shards))
	for idx := range shards {
		shardIdx := idx
		shard := shards[shardIdx]
		e.queryFlow.Submit(flow.FilteringStage, func() {
			// 1. get series ids by query condition
			shardExecuteCtx := flow.NewShardExecuteContext(e.ctx.storageExecuteCtx)
			e.ctx.storageExecuteCtx.ShardContexts[shardIdx] = shardExecuteCtx
			t := newSeriesIDsSearchTaskFunc(shardExecuteCtx, shard)
			err := t.Run()
			if err != nil && !errors.Is(err, constants.ErrNotFound) {
				// maybe series ids not found in shard, so ignore not found err
				e.queryFlow.Complete(err)
				return
			}
			// if series ids not found
			if shardExecuteCtx.SeriesIDsAfterFiltering.IsEmpty() {
				return
			}

			// 2. filter data each data family in shard
			t = newFamilyFilterTaskFunc(shardExecuteCtx, shard)
			err = t.Run()
			if err != nil && !errors.Is(err, constants.ErrNotFound) {
				// maybe data not exist in shard, so ignore not found err
				e.queryFlow.Complete(err)
				return
			}
			if shardExecuteCtx.IsSeriesIDsEmpty() {
				// data not found
				return
			}

			// 3. execute group by
			e.queryFlow.Submit(flow.GroupingStage, func() {
				e.executeGroupBy(shardExecuteCtx, shard)
			})
		})
	}
}

// executeGroupBy executes the query flow, step as below:
// 1. grouping
// 2. loading
func (e *storageExecutor) executeGroupBy(shardExecuteContext *flow.ShardExecuteContext, shard tsdb.Shard) {
	// time segments sorted by family
	timeSegments := shardExecuteContext.TimeSegmentContext.GetTimeSegments()

	if e.ctx.storageExecuteCtx.Query.HasGroupBy() {
		// 1. grouping, if it has grouping, do group by tag keys, else just split series ids as batch first,
		// get grouping context if it needs
		// group context find task maybe change shardExecuteContext.SeriesIDsAfterFiltering value.
		t := newGroupingContextFindTaskFunc(shardExecuteContext, shard)
		err := t.Run()
		if err != nil && !errors.Is(err, constants.ErrNotFound) {
			// maybe group by not found, so ignore not found err
			e.queryFlow.Complete(err)
			return
		}
	}
	// if not grouping found, series id is empty.
	seriesIDs := shardExecuteContext.SeriesIDsAfterFiltering
	seriesIDsHighKeys := seriesIDs.GetHighKeys()

	for seriesIDHighKeyIdx := range seriesIDsHighKeys {
		// be carefully, need use new variable for variable scope problem(closures)
		// ref: https://go.dev/doc/faq#closures_and_goroutines
		highSeriesIDIdx := seriesIDHighKeyIdx

		// grouping based on group by tag keys for each low series container
		e.track.submitTask(flow.GroupingStage, func() {
			lowSeriesIDs := seriesIDs.GetContainerAtIndex(highSeriesIDIdx)
			dataLoadCtx := &flow.DataLoadContext{
				ShardExecuteCtx:       shardExecuteContext,
				LowSeriesIDsContainer: lowSeriesIDs,
				SeriesIDHighKey:       seriesIDsHighKeys[highSeriesIDIdx],
				IsMultiField:          len(shardExecuteContext.StorageExecuteCtx.Fields) > 1,
				IsGrouping:            shardExecuteContext.StorageExecuteCtx.Query.HasGroupBy(),
			}

			t := newBuildGroupTaskFunc(shard, dataLoadCtx)
			if err := t.Run(); err != nil {
				e.queryFlow.Complete(err)
				return
			}
			if !dataLoadCtx.HasGroupingData() {
				// if it hasn't grouping result, after grouping build.
				return
			}

			e.queryFlow.Submit(flow.ScannerStage, func() {
				for segmentIdx := range timeSegments {
					// 3.load data by grouped lowSeriesIDs
					t := newDataLoadTaskFunc(shard, e.queryFlow, dataLoadCtx, segmentIdx, timeSegments[segmentIdx])
					if err := t.Run(); err != nil {
						e.queryFlow.Complete(err)
						return
					}
				}
			})
		})
	}
}
