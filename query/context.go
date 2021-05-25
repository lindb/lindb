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
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/stmt"
)

// storageExecuteContext represents storage query execute context
type storageExecuteContext struct {
	query    *stmt.Query
	shardIDs []int32

	tagFilterResult map[string]*tagFilterResult

	stats *models.StorageStats // storage query stats track for explain query
}

// newStorageExecuteContext creates storage execute context
func newStorageExecuteContext(shardIDs []int32, query *stmt.Query) *storageExecuteContext {
	ctx := &storageExecuteContext{
		query:    query,
		shardIDs: shardIDs,
	}
	if query.Explain {
		// if explain query, create storage query stats
		ctx.stats = models.NewStorageStats()
	}
	return ctx
}

// QueryStats returns the storage query stats
func (ctx *storageExecuteContext) QueryStats() *models.StorageStats {
	if ctx.stats != nil {
		ctx.stats.Complete()
	}
	return ctx.stats
}

// setTagFilterResult sets tag filter result
func (ctx *storageExecuteContext) setTagFilterResult(tagFilterResult map[string]*tagFilterResult) {
	ctx.tagFilterResult = tagFilterResult
}
