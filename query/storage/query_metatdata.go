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
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

var (
	newTagSearchFunc    = newTagSearch
	newSeriesSearchFunc = newSeriesSearch
)

// metadataStorageExecutor represents the executor which executes metric metadata suggest in storage side
type metadataStorageExecutor struct {
	database tsdb.Database
	request  *stmt.MetricMetadata
	shardIDs []models.ShardID

	logger *logger.Logger
}

// newMetadataStorageExecutor creates a metadata suggest executor in storage side
func newStorageMetadataQuery(
	database tsdb.Database,
	shardIDs []models.ShardID,
	request *stmt.MetricMetadata,
) storageMetadataQuery {
	return &metadataStorageExecutor{
		database: database,
		request:  request,
		shardIDs: shardIDs,
		logger:   logger.GetLogger("Query", "Metadata"),
	}
}

// Execute executes the metadata suggest query based on query type
func (e *metadataStorageExecutor) Execute() (result []string, err error) {
	req := e.request
	limit := req.Limit
	if limit == 0 || limit > constants.MaxSuggestions {
		// if limit = 0 or > max suggestion items, need reset limit
		limit = constants.MaxSuggestions
	}

	switch req.Type {
	case stmt.Namespace:
		return e.database.Metadata().MetadataDatabase().SuggestNamespace(req.Prefix, limit)
	case stmt.Metric:
		return e.database.Metadata().MetadataDatabase().SuggestMetrics(req.Namespace, req.Prefix, limit)
	case stmt.TagKey:
		tagKeys, err := e.database.Metadata().MetadataDatabase().GetAllTagKeys(req.Namespace, req.MetricName)
		if err != nil {
			return nil, err
		}
		for _, tagKey := range tagKeys {
			result = append(result, tagKey.Key)
		}
		return result, nil
	case stmt.Field:
		fields, err := e.database.Metadata().MetadataDatabase().GetAllFields(req.Namespace, req.MetricName)
		if err != nil {
			return nil, err
		}
		result = append(result, string(encoding.JSONMarshal(fields)))
		return result, nil
	case stmt.TagValue:
		tagKeyID, err := e.database.Metadata().MetadataDatabase().GetTagKeyID(req.Namespace, req.MetricName, req.TagKey)
		if err != nil {
			return nil, err
		}
		if req.Condition == nil {
			// if not tag filter condition, just get tag value by tag key
			result = e.database.Metadata().TagMetadata().SuggestTagValues(tagKeyID, req.Prefix, limit)
		} else {
			// 1. do tag filter
			ctx := &executeContext{
				database: e.database,
				storageExecuteCtx: &flow.StorageExecuteContext{
					Query: &stmt.Query{
						Namespace:  req.Namespace,
						MetricName: req.MetricName,
						Condition:  req.Condition,
					},
				},
			}
			tagSearch := newTagSearchFunc(ctx)
			err := tagSearch.Filter()
			if err != nil {
				return nil, err
			}
			if len(ctx.storageExecuteCtx.TagFilterResult) == 0 {
				// filter not match, return not found
				return nil, fmt.Errorf("%w , namespace: %s, metricName: %s",
					constants.ErrTagFilterResultNotFound, req.Namespace, req.MetricName)
			}
			// get shard by given query shard id list
			for _, shardID := range e.shardIDs {
				tagValues, err := e.findTagValuesFromShard(ctx, req, shardID, tagKeyID, limit)
				if err != nil {
					// ignore shard level err
					e.logger.Warn("find tag values failure",
						logger.Any("db", e.database), logger.Any("shard", shardID),
						logger.String("metric", req.MetricName), logger.String("tagKey", req.TagKey),
						logger.Error(err))
					continue
				}
				for _, tagValue := range tagValues {
					result = append(result, tagValue)
					if len(result) >= limit {
						return result, nil
					}
				}
			}
		}
	}
	return result, nil
}

// findTagValuesFromShard returns tag values from shard by given condition.
func (e *metadataStorageExecutor) findTagValuesFromShard(ctx *executeContext,
	req *stmt.MetricMetadata, shardID models.ShardID, tagKeyID tag.KeyID, limit int,
) (result []string, err error) {
	shard, ok := e.database.GetShard(shardID)
	if !ok {
		return
	}
	// if shard exist, do series search
	// if it gets tag filter result do series ids searching
	seriesSearch := newSeriesSearchFunc(shard.IndexDatabase(), ctx.storageExecuteCtx.TagFilterResult, req.Condition)
	seriesIDs, err := seriesSearch.Search()
	if err != nil {
		return nil, err
	}
	ctx.storageExecuteCtx.GroupByTagKeyIDs = []tag.KeyID{tagKeyID}
	shardExecuteCtx := &flow.ShardExecuteContext{
		StorageExecuteCtx:       ctx.storageExecuteCtx,
		SeriesIDsAfterFiltering: seriesIDs,
	}
	// get grouping based on tag keys and series ids
	err = shard.IndexDatabase().GetGroupingContext(shardExecuteCtx)
	if err != nil {
		return nil, err
	}
	highKeys := seriesIDs.GetHighKeys()
	for i, highKey := range highKeys {
		// get tag value ids
		tagValueIDs := shardExecuteCtx.GroupingContext.ScanTagValueIDs(highKey, seriesIDs.GetContainerAtIndex(i))
		tagValues := make(map[uint32]string)
		// get tag value
		err = e.database.Metadata().TagMetadata().CollectTagValues(tagKeyID, tagValueIDs[0], tagValues)
		if err != nil {
			return nil, err
		}
		for _, tagValue := range tagValues {
			result = append(result, tagValue)
			if len(result) >= limit {
				return result, nil
			}
		}
	}
	return result, nil
}
