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

package tsdb

import "sync"

var (
	sManager          ShardManager
	once4ShardManager sync.Once
)

// GetShardManager returns the shard manager singleton instance
func GetShardManager() ShardManager {
	once4ShardManager.Do(func() {
		sManager = newShardManager()
	})
	return sManager
}

// ShardManager represents the shard manager
type ShardManager interface {
	// AddShard adds the shard
	AddShard(shard Shard)
	// RemoveShard removes the shard
	RemoveShard(shard Shard)
	// WalkEntry walks each shard entry via fn.
	WalkEntry(fn func(shard Shard))
}

// shardManager implements ShardManager interface
type shardManager struct {
	shards sync.Map
}

// newStorageManager creates the storage manager
func newShardManager() ShardManager {
	return &shardManager{}
}

// AddShard adds the shard
func (sm *shardManager) AddShard(shard Shard) {
	sm.shards.Store(shard.ShardInfo(), shard)
}

// RemoveShard adds the shard
func (sm *shardManager) RemoveShard(shard Shard) {
	sm.shards.Delete(shard.ShardInfo())
}

// WalkEntry walks each shard entry via fn.
func (sm *shardManager) WalkEntry(fn func(shard Shard)) {
	sm.shards.Range(func(key, value interface{}) bool {
		shard := value.(Shard)
		fn(shard)
		return true
	})
}
