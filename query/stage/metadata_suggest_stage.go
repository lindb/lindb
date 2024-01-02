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
	"github.com/lindb/lindb/sql/tree"
)

// metadataSuggestStage represents metadata suggest stage.
type metadataSuggestStage struct {
	baseStage
	ctx *context.LeafMetadataContext
}

// NewMetadataSuggestStage creates a metadataSuggestStage instance.
func NewMetadataSuggestStage(ctx *context.LeafMetadataContext) Stage {
	return &metadataSuggestStage{
		baseStage: baseStage{
			stageType: MetadataSuggest,
		},
		ctx: ctx,
	}
}

// Plan returns sub execution tree for metadata suggest.
func (stage *metadataSuggestStage) Plan() PlanNode {
	req := stage.ctx.Request
	switch req.Type {
	case tree.Namespace:
		return NewPlanNode(operator.NewNamespaceSuggest(stage.ctx))
	case tree.Metric:
		return NewPlanNode(operator.NewMetricSuggest(stage.ctx))
	case tree.TagKey:
		return NewPlanNode(operator.NewTagKeySuggest(stage.ctx))
	case tree.Field:
		return NewPlanNode(operator.NewFieldSuggest(stage.ctx))
	case tree.TagValue:
		execPlan := NewEmptyPlanNode()
		execPlan.AddChild(NewPlanNode(operator.NewTagKeyIDLookup(stage.ctx)))
		stage.ctx.StorageExecuteCtx = &flow.StorageExecuteContext{
			Query: &tree.Query1{
				Namespace:  req.Namespace,
				MetricName: req.MetricName,
				Condition:  req.Condition,
			},
		}
		if req.Condition == nil {
			// if not tag filter condition, just get tag value by tag key
			execPlan.AddChild(NewPlanNode(operator.NewTagValueSuggest(stage.ctx)))
		} else {
			// 1. do tag values lookup
			execPlan.AddChild(NewPlanNode(operator.NewTagValuesLookup(stage.ctx.StorageExecuteCtx, stage.ctx.Database)))
		}
		return execPlan
	}
	return nil
}

// NextStages returns the next stages.
func (stage *metadataSuggestStage) NextStages() (stages []Stage) {
	req := stage.ctx.Request
	if req.Type != tree.TagValue {
		return
	}
	if len(stage.ctx.StorageExecuteCtx.TagFilterResult) == 0 {
		// filter not match, return not found
		return
	}
	// get shard by given query shard id list
	for _, shardID := range stage.ctx.ShardIDs {
		shard, ok := stage.ctx.Database.GetShard(shardID)
		if !ok {
			continue
		}
		// if shard exist, add shard lookup stage
		shardExecuteCtx := flow.NewShardExecuteContext(stage.ctx.StorageExecuteCtx)
		stages = append(stages, NewShardLookupStage(stage.ctx, shardExecuteCtx, shard))
	}
	return
}

// Identifier returns identifier value of metadata suggest stage.
func (stage *metadataSuggestStage) Identifier() string {
	return "Metadata Suggest"
}
