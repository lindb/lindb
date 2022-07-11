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
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

// executeContext represents storage query execute context.
type executeContext struct {
	database tsdb.Database
	shards   []tsdb.Shard

	storageExecuteCtx *flow.StorageExecuteContext
}

// newStorageExecuteContext creates storage execute context
func newStorageExecuteContext(database tsdb.Database, shardIDs []models.ShardID, query *stmt.Query) *executeContext {
	ctx := &executeContext{
		database: database,
		storageExecuteCtx: &flow.StorageExecuteContext{
			Query:    query,
			ShardIDs: shardIDs,
		},
	}
	if query.Explain {
		// if explain query, create storage query stats
		ctx.storageExecuteCtx.Stats = models.NewStorageStats()
	}
	return ctx
}

// getMetadata returns the database's metadata.
func (ctx *executeContext) getMetadata() metadb.Metadata {
	return ctx.database.Metadata()
}

// prepare the execution context, and validates params.
func (ctx *executeContext) prepare() error {
	// do query validation
	if err := ctx.validation(); err != nil {
		return err
	}

	// get shard by given query shard id list
	for _, shardID := range ctx.storageExecuteCtx.ShardIDs {
		// if shard exist, add shard to query list
		if shard, ok := ctx.database.GetShard(shardID); ok {
			ctx.shards = append(ctx.shards, shard)
		}
	}
	// check got shards if valid
	if err := ctx.checkShards(); err != nil {
		return err
	}
	return nil
}

// validation validates query input params are valid.
func (ctx *executeContext) validation() error {
	// check input shardIDs if empty
	if len(ctx.storageExecuteCtx.ShardIDs) == 0 {
		return errNoShardID
	}
	numOfShards := ctx.database.NumOfShards()
	// check engine has shard
	if numOfShards == 0 {
		return errNoShardInDatabase
	}

	return nil
}

// checkShards checks got shards if valid.
func (ctx *executeContext) checkShards() error {
	numOfShards := len(ctx.shards)
	if numOfShards == 0 {
		return errShardNotFound
	}
	numOfShardIDs := len(ctx.storageExecuteCtx.ShardIDs)
	if numOfShards != numOfShardIDs {
		return errShardNumNotMatch
	}
	return nil
}
