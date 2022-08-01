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

// shardLookupStage represents shard level metadata lookup.
type shardLookupStage struct {
	baseStage
	executeCtx      *context.LeafMetadataContext
	shardExecuteCtx *flow.ShardExecuteContext
	shard           tsdb.Shard
}

// NewShardLookupStage creates a shardLookupStage instance.
func NewShardLookupStage(executeCtx *context.LeafMetadataContext, shardExecuteCtx *flow.ShardExecuteContext, shard tsdb.Shard) Stage {
	return &shardLookupStage{
		baseStage: baseStage{
			stageType: ShardLookup,
		},
		executeCtx:      executeCtx,
		shardExecuteCtx: shardExecuteCtx,
		shard:           shard,
	}
}

// Plan returns sub execution tree for tag values collect.
func (stage *shardLookupStage) Plan() PlanNode {
	execPlan := NewEmptyPlanNode()
	// add shard level series filtering node
	execPlan.AddChild(NewPlanNodeWithIgnore(operator.NewSeriesFiltering(stage.shardExecuteCtx, stage.shard)))
	// add tag values collect node
	execPlan.AddChild(NewPlanNode(operator.NewTagValueCollect(stage.executeCtx, stage.shardExecuteCtx, stage.shard)))
	return execPlan
}

// Identifier returns identifier value of shard lookup stage.
func (stage *shardLookupStage) Identifier() string {
	return fmt.Sprintf("Shard Lookup[Shard(%d)]", stage.shard.ShardID())
}
