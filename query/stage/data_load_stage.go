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
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/operator"
)

// dataLoadStage represents data load stage.
type dataLoadStage struct {
	baseStage
	leafExecuteCtx *context.LeafExecuteContext
	executeCtx     *flow.DataLoadContext
	segmentRS      *flow.TimeSegmentResultSet
}

// NewDataLoadStage creates a dataLoadStage instance.
func NewDataLoadStage(leafExecuteCtx *context.LeafExecuteContext,
	executeCtx *flow.DataLoadContext, segmentRS *flow.TimeSegmentResultSet,
) Stage {
	return &dataLoadStage{
		baseStage: baseStage{
			ctx:       leafExecuteCtx.TaskCtx.Ctx,
			execPool:  leafExecuteCtx.Database.ExecutorPool().Scanner,
			stageType: DataLoad,
		},
		leafExecuteCtx: leafExecuteCtx,
		executeCtx:     executeCtx,
		segmentRS:      segmentRS,
	}
}

// Plan returns sub execution plan tree for data load.
func (stage *dataLoadStage) Plan() PlanNode {
	execPlan := NewEmptyPlanNode()
	shardExecuteCtx := stage.executeCtx.ShardExecuteCtx
	stage.segmentRS.IntervalRatio = uint16(shardExecuteCtx.StorageExecuteCtx.Query.IntervalRatio)
	// calc base slot based on query interval and family time of storage
	queryInterval := shardExecuteCtx.StorageExecuteCtx.Query.Interval
	calc := queryInterval.Calculator()
	familyTimeForQuery := calc.CalcFamilyTime(stage.segmentRS.FamilyTime)
	stage.segmentRS.BaseTime = uint16(calc.CalcSlot(stage.segmentRS.FamilyTime, familyTimeForQuery, queryInterval.Int64()))
	stage.segmentRS.TargetRange = shardExecuteCtx.StorageExecuteCtx.CalcTargetSlotRange(familyTimeForQuery)

	for idx := range stage.segmentRS.FilterRS {
		execPlan.AddChild(NewPlanNode(
			operator.NewDataLoad(stage.executeCtx, stage.segmentRS, stage.segmentRS.FilterRS[idx])))
	}
	execPlan.AddChild(NewPlanNode(operator.NewLeafReduce(stage.leafExecuteCtx, stage.executeCtx)))
	return execPlan
}

// Identifier returns identifier value of data load stage.
func (stage *dataLoadStage) Identifier() string {
	return fmt.Sprintf("Data Load[%s]",
		timeutil.FormatTimestamp(stage.segmentRS.FamilyTime, timeutil.DataTimeFormat2))
}
