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

package write

import (
	"context"
	"slices"
	"sync/atomic"

	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/samber/lo"

	"github.com/lindb/lindb/app/broker/write/writer"
	"github.com/lindb/lindb/models"
)

type database[V larrow.EntryType] struct {
	ctx context.Context
	cfg models.Database

	writer writer.Writer

	numOfShards atomic.Value // number of shards(int)TODO: remove it?
	shards      atomic.Value // []Shard
}

func NewDatabase[V larrow.EntryType](ctx context.Context, cfg models.Database) writer.DatabaseAccessor[V] {
	db := &database[V]{
		ctx: ctx,
		cfg: cfg,
	}
	return db
}

func (db *database[V]) LeaderChanged(shardStates map[models.ShardID]models.ShardState,
	liveNodes map[models.NodeID]models.StatefulNode,
) {
	value := db.shards.Load()
	if value == nil {
		return
	}
	shards := value.([]writer.Shard[V])
	for _, shard := range shards {
		if shardState, ok := shardStates[shard.ID()]; ok {
			shard.LeaderChanged(shardState, liveNodes)
		}
	}
}

func (db *database[V]) Name() string {
	return db.cfg.Name
}

func (db *database[V]) GetShards() []writer.Shard[V] {
	value := db.shards.Load()
	if value == nil {
		return nil
	}
	return value.([]writer.Shard[V])
}

func (db *database[V]) CreateShards(shardStates map[models.ShardID]models.ShardState, liveNodes map[models.NodeID]models.StatefulNode) {
	shards := db.GetShards()
	var newShards []writer.Shard[V]
	newShards = append(newShards, shards...)

	shardMap := lo.Associate(shards, func(shard writer.Shard[V]) (models.ShardID, writer.Shard[V]) {
		return shard.ID(), shard
	})

	for _, shardState := range shardStates {
		if _, ok := shardMap[shardState.ID]; !ok {
			// create new shard if not exist
			newShards = append(newShards, NewShard[V](db.ctx, db, shardState, liveNodes))
		}
	}

	if len(newShards) == len(shards) {
		// no new shard created
		return
	}

	// sort by shard id
	slices.SortFunc(newShards, func(i, j writer.Shard[V]) int {
		return int(i.ID() - j.ID())
	})
	db.shards.Store(newShards)
	db.numOfShards.Store(len(newShards))
}
