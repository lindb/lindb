package table

import (
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/pkg/logger"
)

//FIXME store100 using lru cache?????

//go:generate mockgen -source ./cache.go -destination=./cache_mock.go -package table

// for test
var (
	newMMapStoreReaderFunc = newMMapStoreReader
)

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
func (c *mapCache) Evict(family string, fileName string) {
	filePath := filepath.Join(family, fileName)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	reader, ok := c.readers[filePath]
	if ok {
		if err := reader.Close(); err != nil {
			tableLogger.Error("close store reader error",
				logger.String("path", c.storePath),
				logger.String("file", filePath), logger.Error(err))
		}
		delete(c.readers, filePath)
	}
}

// GetReader returns store reader from cache, create new reader if not exist
func (c *mapCache) GetReader(family string, fileName string) (Reader, error) {
	filePath := filepath.Join(family, fileName)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// find from cache
	reader, ok := c.readers[filePath]
	if ok {
		return reader, nil
	}

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
			tableLogger.Error("close store reader error",
				logger.String("path", c.storePath),
				logger.String("file", k), logger.Error(err))
		}
	}
	return nil
}
