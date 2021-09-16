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

import (
	"sort"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
)

type shardEntry struct {
	shardID models.ShardID
	shard   Shard
}

type shardEntries []shardEntry

func (se shardEntries) Len() int           { return len(se) }
func (se shardEntries) Less(i, j int) bool { return se[i].shardID < se[j].shardID }
func (se shardEntries) Swap(i, j int)      { se[i], se[j] = se[j], se[i] }

// shardSet is a immutable data structure in database to provide lock-free lookup
type shardSet struct {
	value atomic.Value // shardEntries
	num   atomic.Int32 // number of families
}

// newShardSet returns a default empty shardSet
func newShardSet() *shardSet {
	set := &shardSet{
		num: *atomic.NewInt32(0),
	}
	// initialize it with a new empty entry slice
	set.value.Store(shardEntries{})
	return set
}

// InsertShard appends a new shard into the slice,
// then changes atomic.Value to the new sorted set
func (ss *shardSet) InsertShard(shardID models.ShardID, shard Shard) {
	oldEntries := ss.value.Load().(shardEntries)
	var (
		newEntries shardEntries
		newEntry   = shardEntry{shardID: shardID, shard: shard}
	)

	newEntries = make([]shardEntry, oldEntries.Len()+1)
	copy(newEntries, oldEntries)
	newEntries[len(newEntries)-1] = newEntry
	sort.Sort(newEntries)

	ss.value.Store(newEntries)
	ss.num.Inc()
}

// GetShard searches the shard by shardID from the shardSet
// BinarySearch is not always faster than iterating
func (ss *shardSet) GetShard(shardID models.ShardID) (Shard, bool) {
	entries := ss.value.Load().(shardEntries)
	// fast path when length < 20
	if entries.Len() <= 20 {
		for idx := range entries {
			if entries[idx].shardID == shardID {
				return entries[idx].shard, true
			}
		}
		return nil, false
	}
	index := sort.Search(entries.Len(), func(i int) bool {
		return entries[i].shardID >= shardID
	})
	if index < 0 || index >= entries.Len() {
		return nil, false
	}
	return entries[index].shard, entries[index].shardID == shardID
}

// GetShardNum returns the shard number
func (ss *shardSet) GetShardNum() int {
	return int(ss.num.Load())
}

func (ss *shardSet) Entries() shardEntries {
	return ss.value.Load().(shardEntries)
}
