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

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/operator"
	"github.com/lindb/lindb/tsdb"
)

// groupingStage represents grouping stage.
type groupingStage struct {
	baseStage
	leafExecuteCtx *context.LeafExecuteContext
	executeCtx     *flow.DataLoadContext
	shard          tsdb.Shard
}

// NewGroupingStage creates a groupingStage instance.
func NewGroupingStage(leafExecuteCtx *context.LeafExecuteContext, executeCtx *flow.DataLoadContext, shard tsdb.Shard) Stage {
	leafExecuteCtx.GroupingCtx.ForkGroupingTask()
	return &groupingStage{
		baseStage: baseStage{
			ctx:       leafExecuteCtx.TaskCtx.Ctx,
			execPool:  leafExecuteCtx.Database.ExecutorPool().Grouping,
			stageType: Grouping,
		},
		leafExecuteCtx: leafExecuteCtx,
		executeCtx:     executeCtx,
		shard:          shard,
	}
}

// Plan returns sub execution plan tree for grouping.
func (stage *groupingStage) Plan() PlanNode {
	// add find grouping node
	return NewPlanNode(operator.NewGroupingTagsLookup(stage.executeCtx))
}

// NextStages returns the stages after grouping.
func (stage *groupingStage) NextStages() (stages []Stage) {
	if stage.executeCtx.IsGrouping && len(stage.executeCtx.GroupingSeriesAgg) == 0 {
		// if not found any grouping tags, terminal.
		return
	}
	// time segments sorted by family time
	// timeSegments := stage.executeCtx.ShardExecuteCtx.TimeSegmentContext.GetTimeSegments()
	// dlCtx := stage.executeCtx
	// for segmentIdx := range timeSegments {
	// 	dataLoadCtx := *dlCtx // copy data load context, because data load context not thread safe
	// 	// add data load stage based on time segment, one by one
	// 	stages = append(stages, NewDataLoadStage(stage.leafExecuteCtx, &dataLoadCtx, timeSegments[segmentIdx]))
	// 	// track if all data load tasks completed
	// 	stage.executeCtx.PendingDataLoadTasks.Add(int32(len(timeSegments[segmentIdx].FilterRS)))
	// }
	return
}

// Complete completes grouping task.
func (stage *groupingStage) Complete() {
	stage.leafExecuteCtx.GroupingCtx.CompleteGroupingTask()
}

// Identifier returns identifier value of grouping stage.
func (stage *groupingStage) Identifier() string {
	return fmt.Sprintf("Grouping[Shard(%d)]", stage.shard.ShardID())
}
