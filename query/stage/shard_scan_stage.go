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

package stage

import (
	"fmt"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/operator"
	"github.com/lindb/lindb/tsdb"
)

// shardScanStage represents shard scan stage.
type shardScanStage struct {
	baseStage
	leafExecuteCtx *context.LeafExecuteContext

	shardExecuteCtx *flow.ShardExecuteContext

	shard tsdb.Shard
}

// NewShardScanStage creates a shardScanStage instance.
func NewShardScanStage(leafExecuteCtx *context.LeafExecuteContext,
	shardExecuteCtx *flow.ShardExecuteContext, shard tsdb.Shard,
) Stage {
	leafExecuteCtx.GroupingCtx.ForkGroupingTask()
	return &shardScanStage{
		baseStage: baseStage{
			ctx:       leafExecuteCtx.TaskCtx.Ctx,
			execPool:  leafExecuteCtx.Database.ExecutorPool().Filtering,
			stageType: ShardScan,
		},
		leafExecuteCtx:  leafExecuteCtx,
		shardExecuteCtx: shardExecuteCtx,
		shard:           shard,
	}
}

// Plan returns sub execution tree for shard scan.
func (stage *shardScanStage) Plan() PlanNode {
	shardExecuteCtx := stage.shardExecuteCtx
	queryStmt := shardExecuteCtx.StorageExecuteCtx.Query
	shard := stage.shard
	// if shard exist, add shard to query list
	families := shard.GetDataFamilies(queryStmt.StorageInterval.Type(), queryStmt.TimeRange)
	if len(families) == 0 {
		// no data family found
		return nil
	}
	execPlan := NewEmptyPlanNode()
	if queryStmt.Condition != nil {
		// add shard level series filtering node
		execPlan.AddChild(NewPlanNodeWithIgnore(operator.NewSeriesFiltering(shardExecuteCtx, shard)))
	} else {
		// add shard level all series lookup node
		execPlan.AddChild(NewPlanNodeWithIgnore(operator.NewMetricAllSeries(shardExecuteCtx, shard)))
	}

	for idx := range families {
		family := families[idx]
		// add data family reader node, found series ids which match condition.
		execPlan.AddChild(NewPlanNodeWithIgnore(operator.NewDataFamilyRead(shardExecuteCtx, family)))
	}

	if shardExecuteCtx.StorageExecuteCtx.Query.HasGroupBy() {
		// if it has grouping, do group by tag keys, else just split series ids as batch first,
		// get grouping context if it needs
		// group context find task maybe change shardExecuteContext.SeriesIDsAfterFiltering value.
		execPlan.AddChild(NewPlanNodeWithIgnore(operator.NewGroupingContextBuild(shardExecuteCtx, shard)))
	}
	execPlan.AddChild(NewPlanNodeWithIgnore(operator.NewSeriesLimit(shardExecuteCtx, shard)))
	return execPlan
}

// NextStages returns the next stages after shard scan completed.
func (stage *shardScanStage) NextStages() (stages []Stage) {
	// if not grouping found, series id is empty.
	shardExecuteContext := stage.shardExecuteCtx
	seriesIDs := shardExecuteContext.SeriesIDsAfterFiltering
	seriesIDsHighKeys := seriesIDs.GetHighKeys()

	for seriesIDHighKeyIdx := range seriesIDsHighKeys {
		// be carefully, need use new variable for variable scope problem(closures)
		// ref: https://go.dev/doc/faq#closures_and_goroutines
		highSeriesIDIdx := seriesIDHighKeyIdx
		// grouping based on group by tag keys for each low series container
		lowSeriesIDs := seriesIDs.GetContainerAtIndex(highSeriesIDIdx)
		dataLoadCtx := &flow.DataLoadContext{
			// ShardExecuteCtx:       shardExecuteContext,
			LowSeriesIDsContainer: lowSeriesIDs,
			SeriesIDHighKey:       seriesIDsHighKeys[highSeriesIDIdx],
			IsMultiField:          len(shardExecuteContext.StorageExecuteCtx.Fields) > 1,
			IsGrouping:            shardExecuteContext.StorageExecuteCtx.Query.HasGroupBy(),
			PendingDataLoadTasks:  atomic.NewInt32(0),
		}

		stages = append(stages, NewGroupingStage(stage.leafExecuteCtx, dataLoadCtx, stage.shard))
	}
	return stages
}

// Complete completes shard scan stage, dec grouping task counter.
func (stage *shardScanStage) Complete() {
	stage.leafExecuteCtx.GroupingCtx.CompleteGroupingTask()
}

// Identifier returns identifier value of shard scan stage.
func (stage *shardScanStage) Identifier() string {
	return fmt.Sprintf("Shard Scan[Shard(%d)]", stage.shard.ShardID())
}
