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

package context

import (
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

type LeafMetadataContext struct {
	Request  *stmt.MetricMetadata
	Database tsdb.Database
	ShardIDs []models.ShardID

	StorageExecuteCtx *flow.StorageExecuteContext

	ResultSet []string
	TagKeyID  tag.KeyID // for tag values suggest

	Limit int
}

func NewLeafMetadataContext(request *stmt.MetricMetadata, database tsdb.Database, shardIDs []models.ShardID) *LeafMetadataContext {
	ctx := &LeafMetadataContext{
		Request:  request,
		Database: database,
		ShardIDs: shardIDs,
	}
	ctx.Limit = ctx.getLimit()
	return ctx
}

func (ctx *LeafMetadataContext) getLimit() int {
	req := ctx.Request
	limit := req.Limit
	if limit == 0 || limit > constants.MaxSuggestions {
		// if limit = 0 or > max suggestion items, need reset limit
		limit = constants.MaxSuggestions
	}
	return limit
}

func (ctx *LeafMetadataContext) AddValue(val string) {
	if len(ctx.ResultSet) >= ctx.Limit {
		return
	}
	ctx.ResultSet = append(ctx.ResultSet, val)
}
