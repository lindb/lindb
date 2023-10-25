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

package table

import (
	"container/list"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/metrics"
)

//go:generate mockgen -source ./cache.go -destination=./cache_mock.go -package table

// for test
var (
	newMMapStoreReaderFunc = newMMapStoreReader
)

// Cache caches table readers.
type Cache interface {
	// GetReader returns store reader from cache, create new reader if not exist.
	GetReader(family string, fileName string) (Reader, error)
	// ReleaseReaders releases reader after read completed.
	ReleaseReaders(readers []Reader)
	// Evict evicts file reader from cache.
	Evict(fileName string)
	// Cleanup cleans the expired reader from cache.
	Cleanup()
	// Close cleans cache data after closing reader resource firstly.
	Close() error
}

// Cache caches table readers based on lru cache.
type storeCache struct {
	ttl       time.Duration
	storePath string
	families  map[string]map[string]struct{} // family name => files
	cache     *LRUCache
	mutex     sync.Mutex
}

// NewCache creates cache for store readers.
func NewCache(storePath string, ttl time.Duration) Cache {
	return &storeCache{
		ttl:       ttl,
		storePath: storePath,
		families:  make(map[string]map[string]struct{}),
		cache:     NewLRUCache(),
	}
}

// Evict evicts file reader from cache.
func (c *storeCache) Evict(fileName string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if entry, ok := c.cache.Get(fileName); ok {
		c.evict(entry)
		c.cache.Remove(fileName)
	}
}

// ReleaseReaders releases reader after read completed.
func (c *storeCache) ReleaseReaders(readers []Reader) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, r := range readers {
		if entry, ok := c.cache.Get(r.FileName()); ok {
			entry.release()
		}
	}
}

// GetReader returns store reader from cache, create new reader if not exist.
func (c *storeCache) GetReader(family, fileName string) (Reader, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// find from cache
	if entry, ok := c.cache.Get(fileName); ok {
		entry.retain()
		metrics.TableCacheStatistics.Hit.Incr()
		return entry.reader, nil
	}

	metrics.TableCacheStatistics.Miss.Incr()
	metrics.TableCacheStatistics.ActiveReaders.Incr()
	// create new reader
	path := filepath.Join(c.storePath, family, fileName)
	newReader, err := newMMapStoreReaderFunc(path, fileName)
	if err != nil {
		return nil, err
	}
	entry := &cacheEntry{
		key:      fileName,
		reader:   newReader,
		family:   family,
		fileName: fileName,
	}
	entry.retain()
	c.cache.Add(fileName, entry)

	if files, ok := c.families[family]; ok {
		files[fileName] = struct{}{}
	} else {
		c.families[family] = map[string]struct{}{fileName: {}}
	}

	return newReader, nil
}

// Cleanup cleans the expired reader from cache.
func (c *storeCache) Cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ttl := c.ttl.Milliseconds()
	c.cache.Walk(func(entry *cacheEntry) bool {
		if entry.ref.Load() == 0 && timeutil.Now()-entry.last > ttl {
			c.evict(entry)
			metrics.TableCacheStatistics.Evict.Incr()
			return true
		}
		return false
	})
}

// Close closes reader resource and cleans cache data.
func (c *storeCache) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache.Purge(func(entry *cacheEntry) {
		c.closeReader(entry)
		metrics.TableCacheStatistics.Evict.Incr()
	})
	return nil
}

func (c *storeCache) closeReader(entry *cacheEntry) {
	metrics.TableCacheStatistics.ActiveReaders.Decr()
	if err := entry.reader.Close(); err != nil {
		metrics.TableCacheStatistics.CloseFailures.Incr()
		tableLogger.Error("close store reader error",
			logger.String("path", c.storePath),
			logger.String("family", entry.family),
			logger.String("file", entry.fileName), logger.Error(err))
	} else {
		metrics.TableCacheStatistics.Close.Incr()
	}
}

func (c *storeCache) evict(entry *cacheEntry) {
	c.closeReader(entry)

	files := c.families[entry.family]
	delete(files, entry.fileName)
	if len(files) == 0 {
		delete(c.families, entry.family)
	}
	metrics.TableCacheStatistics.Evict.Incr()
}

// cacheEntry represents entry in lru cache.
type cacheEntry struct {
	key              string       // file name
	reader           Reader       // file reader
	ref              atomic.Int32 // how many read
	family, fileName string
	last             int64 // last read timestamp
}

func (e *cacheEntry) retain() {
	e.ref.Inc()
	e.last = timeutil.Now()
}

func (e *cacheEntry) release() {
	e.ref.Dec()
}

type LRUCache struct {
	items     map[string]*list.Element
	evictList *list.List
}

// NewLRUCache constructs an LRU cache.
func NewLRUCache() *LRUCache {
	return &LRUCache{
		evictList: list.New(),
		items:     make(map[string]*list.Element),
	}
}

// Add adds a value to the cache.
func (c *LRUCache) Add(key string, value *cacheEntry) {
	entry := c.evictList.PushFront(value)
	c.items[key] = entry
}

// Get looks up a key's value from the cache.
func (c *LRUCache) Get(key string) (value *cacheEntry, ok bool) {
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		value := ent.Value.(*cacheEntry)
		return value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *LRUCache) Remove(key string) {
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
	}
}

// Walk walks the old entries, if it's expired, need to remove.
func (c *LRUCache) Walk(fn func(entry *cacheEntry) bool) {
	size := len(c.items)
	for i := 0; i < size; i++ {
		// get oldest entry
		ent := c.evictList.Back()
		if ent != nil {
			entry := ent.Value.(*cacheEntry)
			if fn(entry) {
				c.removeElement(ent)
			} else {
				break
			}
		}
	}
}

// Purge is used to completely clear the cache.
func (c *LRUCache) Purge(fn func(entry *cacheEntry)) {
	for k, v := range c.items {
		entry := v.Value.(*cacheEntry)
		fn(entry)
		delete(c.items, k)
	}
	c.evictList.Init()
}

// removeElement is used to remove a given list element from the cache
func (c *LRUCache) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*cacheEntry)
	delete(c.items, kv.key)
}
