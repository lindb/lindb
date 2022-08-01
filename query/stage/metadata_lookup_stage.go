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
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/query/operator"
)

// metadataLookupStage represents metadata lookup stage.
type metadataLookupStage struct {
	baseStage
	leafExecuteCtx *context.LeafExecuteContext
}

// NewMetadataLookupStage creates a metadataLookupStage instance.
func NewMetadataLookupStage(leafExecuteCtx *context.LeafExecuteContext) Stage {
	return &metadataLookupStage{
		baseStage: baseStage{
			stageType: MetadataLookup,
		},
		leafExecuteCtx: leafExecuteCtx,
	}
}

// Plan returns sub execution tree for metadata lookup.
func (stage *metadataLookupStage) Plan() PlanNode {
	execPlan := NewEmptyPlanNode()
	execCtx := stage.leafExecuteCtx.StorageExecuteCtx
	database := stage.leafExecuteCtx.Database
	// add metadata lookup(name/tag/field etc.) node
	execPlan.AddChild(NewPlanNode(operator.NewMetadataLookup(execCtx, database)))
	hasWhereCondition := execCtx.Query.Condition != nil
	if hasWhereCondition {
		// add tag values lookup node if query has where condition
		execPlan.AddChild(NewPlanNode(operator.NewTagValuesLookup(execCtx, database)))
	}
	return execPlan
}

// NextStages returns the next stages after metadata lookup completed.
func (stage *metadataLookupStage) NextStages() (stages []Stage) {
	storageExecuteCtx := stage.leafExecuteCtx.StorageExecuteCtx
	shardIDs := storageExecuteCtx.ShardIDs
	storageExecuteCtx.ShardContexts = make([]*flow.ShardExecuteContext, len(shardIDs))
	for shardIdx := range shardIDs {
		shardExecuteCtx := flow.NewShardExecuteContext(storageExecuteCtx)
		storageExecuteCtx.ShardContexts[shardIdx] = shardExecuteCtx
		if shard, ok := stage.leafExecuteCtx.Database.GetShard(shardIDs[shardIdx]); ok {
			stages = append(stages, NewShardScanStage(stage.leafExecuteCtx, shardExecuteCtx, shard))
		}
	}
	return
}

// Identifier returns identifier value of metadata lookup stage.
func (stage *metadataLookupStage) Identifier() string {
	return "Metadata Lookup"
}
