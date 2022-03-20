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
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./cache.go -destination=./cache_mock.go -package table

// for test
var (
	newMMapStoreReaderFunc   = newMMapStoreReader
	_once4Cache              sync.Once
	_instanceCacheStatistics *cacheStatistics
)

func getCacheStatistics() *cacheStatistics {
	_once4Cache.Do(func() {
		tableCacheScope := linmetric.StorageRegistry.NewScope("lindb.kv.table.cache")
		_instanceCacheStatistics = &cacheStatistics{
			evictCounts: tableCacheScope.NewCounter("evict_counts"),
			cacheHits:   tableCacheScope.NewCounter("cache_hits"),
			cacheMisses: tableCacheScope.NewCounter("cache_misses"),
			CloseCounts: tableCacheScope.NewCounter("close_counts"),
			CloseErrors: tableCacheScope.NewCounter("close_errors"),
		}
	})
	return _instanceCacheStatistics
}

type cacheStatistics struct {
	evictCounts *linmetric.BoundCounter
	cacheHits   *linmetric.BoundCounter
	cacheMisses *linmetric.BoundCounter
	CloseCounts *linmetric.BoundCounter
	CloseErrors *linmetric.BoundCounter
}

// Cache caches table readers
type Cache interface {
	// GetReader returns store reader from cache, create new reader if not exist.
	GetReader(family string, fileName string) (Reader, error)
	// Evict evicts file reader from cache
	Evict(family string, fileName string)
	// Close cleans cache data after closing reader resource firstly
	Close() error
}

// Cache caches table readers based on map
type mapCache struct {
	storePath string
	readers   map[string]Reader
	mutex     sync.Mutex
}

// NewCache creates cache for store readers
func NewCache(storePath string) Cache {
	return &mapCache{
		storePath: storePath,
		readers:   make(map[string]Reader),
	}
}

// Evict evicts file reader from cache
func (c *mapCache) Evict(family, fileName string) {
	filePath := filepath.Join(family, fileName)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	reader, ok := c.readers[filePath]
	if ok {
		if err := reader.Close(); err != nil {
			getCacheStatistics().CloseErrors.Incr()
			tableLogger.Error("close store reader error",
				logger.String("path", c.storePath),
				logger.String("file", filePath), logger.Error(err))
		} else {
			getCacheStatistics().CloseCounts.Incr()
		}
		getCacheStatistics().evictCounts.Incr()
		delete(c.readers, filePath)
	}
}

// GetReader returns store reader from cache, create new reader if not exist
func (c *mapCache) GetReader(family, fileName string) (Reader, error) {
	filePath := filepath.Join(family, fileName)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// find from cache
	reader, ok := c.readers[filePath]
	if ok {
		getCacheStatistics().cacheHits.Incr()
		return reader, nil
	}

	getCacheStatistics().cacheMisses.Incr()
	// create new reader
	path := filepath.Join(c.storePath, filePath)
	newReader, err := newMMapStoreReaderFunc(path)
	if err != nil {
		return nil, err
	}
	c.readers[filePath] = newReader
	return newReader, nil
}

// Close closes reader resource and cleans cache data.
func (c *mapCache) Close() error {
	for k, v := range c.readers {
		if err := v.Close(); err != nil {
			getCacheStatistics().CloseErrors.Incr()
			tableLogger.Error("close store reader error",
				logger.String("path", c.storePath),
				logger.String("file", k), logger.Error(err))
		} else {
			getCacheStatistics().CloseCounts.Incr()
		}
	}
	return nil
}
