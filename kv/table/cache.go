package table

import (
	"path/filepath"
	"sync"

	"go.uber.org/zap"

	"github.com/eleme/lindb/kv/version"
	"github.com/eleme/lindb/pkg/logger"
)

//TODO using lur cache?????

// Cache caches table readers
type Cache interface {
	// GetReader returns store reader from cache, create new reader if not exist.
	GetReader(family string, fileNumber int64) (Reader, error)
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

// GetReader returns store reader from cache, create new reader if not exist
func (c *mapCache) GetReader(family string, fileNumber int64) (Reader, error) {
	filePath := filepath.Join(family, version.Table(fileNumber))
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// find from cache
	reader, ok := c.readers[filePath]
	if ok {
		return reader, nil
	}

	// create new reader
	path := filepath.Join(c.storePath, filePath)
	newReader, err := newMMapStoreReader(path)
	if err != nil {
		return nil, err
	}
	c.readers[filePath] = newReader
	return newReader, nil
}

// Close closes reader resource and cleans cache data.
func (c *mapCache) Close() error {
	log := logger.GetLogger()
	for k, v := range c.readers {
		if err := v.Close(); err != nil {
			log.Error("close store reader error",
				zap.String("path", c.storePath), zap.String("file", k), zap.Error(err))
		}
	}
	return nil
}
