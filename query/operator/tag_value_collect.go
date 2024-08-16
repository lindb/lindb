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

package operator

import (
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/query/context"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb"
)

// tagValueCollect represents tag value collect operator.
type tagValueCollect struct {
	executeCtx      *context.LeafMetadataContext
	shardExecuteCtx *flow.ShardExecuteContext
	shard           tsdb.Shard

	logger logger.Logger
}

// NewTagValueCollect create a tagValueCollect instance.
func NewTagValueCollect(executeCtx *context.LeafMetadataContext, shardExecuteCtx *flow.ShardExecuteContext, shard tsdb.Shard) Operator {
	return &tagValueCollect{
		executeCtx:      executeCtx,
		shardExecuteCtx: shardExecuteCtx,
		shard:           shard,
		logger:          logger.GetLogger("Operator", "TagValueCollect"),
	}
}

// Execute collects tag values with condition, if it has error ignore it.
func (op *tagValueCollect) Execute() error {
	if err := op.execute(); err != nil {
		req := op.executeCtx.Request
		// ignore shard level err
		op.logger.Warn("collect tag values failure",
			logger.Any("db", op.executeCtx.Database.Name()), logger.Any("shard", op.shard.ShardID()),
			logger.String("metric", req.MetricName), logger.String("tagKey", req.TagKey),
			logger.Error(err))
	}
	return nil
}

func (op *tagValueCollect) execute() error {
	tagKeyID := op.executeCtx.TagKeyID
	op.executeCtx.StorageExecuteCtx.GroupByTagKeyIDs = []tag.KeyID{tagKeyID}
	// get grouping based on tag keys and series ids
	// if err := op.shard.IndexDB().GetGroupingContext(op.shardExecuteCtx); err != nil {
	// 	return err
	// }
	seriesIDs := op.shardExecuteCtx.SeriesIDsAfterFiltering
	highKeys := seriesIDs.GetHighKeys()
	for i, highKey := range highKeys {
		// get tag value ids
		tagValueIDs := op.shardExecuteCtx.GroupingContext.ScanTagValueIDs(highKey, seriesIDs.GetContainerAtIndex(i))
		tagValues := make(map[uint32]string)
		// get tag value
		err := op.executeCtx.Database.MetaDB().CollectTagValues(tagKeyID, tagValueIDs[0], tagValues)
		if err != nil {
			return err
		}
		for _, tagValue := range tagValues {
			op.executeCtx.AddValue(tagValue)
		}
	}
	return nil
}

// Identifier returns identifier value of tag value collect operator.
func (op *tagValueCollect) Identifier() string {
	return "Tag Value Collect"
}
