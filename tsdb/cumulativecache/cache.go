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

package cumulativecache

import (
	"time"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/tsdb/memdb"
)

// Cache is a simplified cache used to transform cumulative metric value into delta
// It is assumed that cumulative is from other data-source like prometheus.
// Thus cumulative metric is a non-multi field metric.
// If metric comes from different host but contains the same tags,
// cumulative metric is not monotonic, this may results in a point losing problem
type Cache struct {
	shards    []*cacheShard
	ttl       uint32
	shardMask uint64
	closeCh   chan struct{}
	metrics   cacheMetrics
}

type cacheMetrics struct {
	// Nums is a number of keys stored
	Nums *linmetric.BoundGauge
	// Hits is a number of successfully found keys
	Hits *linmetric.BoundDeltaCounter
	// Misses is a number of not found keys
	Misses *linmetric.BoundDeltaCounter
	// Evicts is a number of successfully deleted keys
	Evicts *linmetric.BoundDeltaCounter
}

func newCacheMetrics(scope linmetric.Scope) *cacheMetrics {
	return &cacheMetrics{
		Nums:   scope.NewGauge("key_nums"),
		Hits:   scope.NewDeltaCounter("key_hits"),
		Misses: scope.NewDeltaCounter("key_misses"),
		Evicts: scope.NewDeltaCounter("key_evicts"),
	}
}

func NewCache(
	shardsCount int,
	ttl time.Duration,
	checkInterval time.Duration,
	scope linmetric.Scope,
) *Cache {
	if shardsCount&(shardsCount-1) != 0 {
		panic("shards count must be a power of two")
	}
	cache := &Cache{
		shards:    make([]*cacheShard, shardsCount),
		shardMask: uint64(shardsCount - 1),
		ttl:       uint32(ttl.Seconds()),
		metrics:   *newCacheMetrics(scope),
		closeCh:   make(chan struct{}),
	}
	for i := 0; i < shardsCount; i++ {
		cache.shards[i] = newCacheShard(cache.ttl, &cache.metrics)
	}

	if checkInterval.Seconds() > 0 {
		go func() {
			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					cache.clean()
				case <-cache.closeCh:
					return
				}
			}
		}()
	}
	return cache
}

func (c *Cache) getShard(hashedKey uint64) (shard *cacheShard) {
	return c.shards[hashedKey&c.shardMask]
}

// loadAndStore picks the old data and stored with new
func (c *Cache) loadAndStore(key uint64, newValue []byte) (oldValue []byte, found bool) {
	shard := c.getShard(key)
	return shard.getAndSet(key, newValue)
}

func (c *Cache) clean() {
	for _, shard := range c.shards {
		shard.cleanUp()
	}
}

// Capacity returns amount of bytes store in the cache.
func (c *Cache) Capacity() int {
	return int(c.metrics.Nums.Get())
}

// Close closes the cache
func (c *Cache) Close() {
	close(c.closeCh)
}

func (c *Cache) CumulativePointToDelta(mp *memdb.MetricPoint) (updated bool) {
	key := uint64(mp.MetricID)<<32 + uint64(mp.SeriesID)
	newData := encodeCumulativeFields(mp)
	oldData, ok := c.loadAndStore(key, newData)
	if !ok {
		return false
	}
	return decodeCumulativeFieldsInto(mp, oldData)
}
