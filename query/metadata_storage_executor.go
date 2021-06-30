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

package query

import (
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

// metadataStorageExecutor represents the executor which executes metric metadata suggest in storage side
type metadataStorageExecutor struct {
	database tsdb.Database
	request  *stmt.Metadata
	shardIDs []int32
}

// newMetadataStorageExecutor creates a metadata suggest executor in storage side
func newMetadataStorageExecutor(database tsdb.Database, shardIDs []int32,
	request *stmt.Metadata,
) parallel.MetadataExecutor {
	return &metadataStorageExecutor{
		database: database,
		request:  request,
		shardIDs: shardIDs,
	}
}

// Execute executes the metadata suggest query based on query type
func (e *metadataStorageExecutor) Execute() (result []string, err error) {
	req := e.request
	limit := req.Limit

	switch req.Type {
	case stmt.Namespace:
		return e.database.Metadata().MetadataDatabase().SuggestNamespace(req.Prefix, limit)
	case stmt.Metric:
		return e.database.Metadata().MetadataDatabase().SuggestMetrics(req.Namespace, req.Prefix, limit)
	case stmt.TagKey:
		return e.database.Metadata().MetadataDatabase().SuggestTagKeys(req.Namespace, req.MetricName, req.Prefix, limit)
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
			tagSearch := newTagSearchFunc(req.Namespace, req.MetricName,
				req.Condition, e.database.Metadata())
			tagFilterResult, err := tagSearch.Filter()
			if err != nil {
				return nil, err
			}
			if len(tagFilterResult) == 0 {
				// filter not match, return not found
				return nil, fmt.Errorf("%w , namespace: %s, metricName: %s",
					constants.ErrTagFilterResultNotFound, req.Namespace, req.MetricName)
			}
			groupByTagKeyIDs := []uint32{tagKeyID}
			// get shard by given query shard id list
			for _, shardID := range e.shardIDs {
				shard, ok := e.database.GetShard(shardID)
				// if shard exist, do series search
				if ok {
					// if get tag filter result do series ids searching
					seriesSearch := newSeriesSearchFunc(shard.IndexDatabase(), tagFilterResult, req.Condition)
					seriesIDs, err := seriesSearch.Search()
					if err != nil {
						return nil, err
					}
					// get grouping based on tag keys and series ids
					gCtx, err := shard.IndexDatabase().GetGroupingContext(groupByTagKeyIDs, seriesIDs)
					if err != nil {
						return nil, err
					}
					highKeys := seriesIDs.GetHighKeys()
					for i, highKey := range highKeys {
						// get tag value ids
						tagValueIDs := gCtx.ScanTagValueIDs(highKey, seriesIDs.GetContainerAtIndex(i))
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
				}
			}
		}
	}
	return result, nil
}
