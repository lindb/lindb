package service

import (
	"fmt"
	"sync"

	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/storage/config"
	"github.com/eleme/lindb/tsdb"
)

// StorageService represents a storage manage interface for tsdb engine
type StorageService interface {
	// CreateShard creates shard for data partition
	CreateShard(db string, option option.ShardOption, shardIDs ...int32) error
}

var (
	storageSvc StorageService
	once       sync.Once
	// EngineConfig must set when system init
	EngineConfig config.Engine
)

// GetStorageService returns singleton storage service instance
func GetStorageService() StorageService {
	once.Do(func() {
		storageSvc = newStorageService()
	})
	return storageSvc
}

// newStorageService creates storage service instance for managing tsdb engine
func newStorageService() StorageService {
	return &storageService{
		engines: make(map[string]tsdb.Engine),
	}
}

// storageService implements StorageService interface
type storageService struct {
	engines map[string]tsdb.Engine

	mu sync.RWMutex
}

// CreateShard creates shard for data partition by given options
// 1) dump engine option into local disk
// 2) create shard storage struct
func (s *storageService) CreateShard(db string, option option.ShardOption, shardIDs ...int32) error {
	if len(shardIDs) <= 0 {
		return fmt.Errorf("cannot create empty shard for db[%s]", db)
	}

	return nil
}
