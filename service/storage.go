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
	// CreateShards creates shards for data partition
	CreateShards(db string, option option.ShardOption, shardIDs ...int32) error
	// GetShard returns shard by given db and shard id, if not exist return nil
	GetShard(db string, shardID int32) tsdb.Shard
}

var (
	storageSvc StorageService
	once       sync.Once
	// EngineConfig must set when system init
	EngineConfig *config.Engine
)

// GetStorageService returns singleton storage service instance
func GetStorageService() (StorageService, error) {
	if EngineConfig == nil {
		return nil, fmt.Errorf("cannot get storage service, because storage config is nil")
	}
	once.Do(func() {
		storageSvc = newStorageService()
	})
	return storageSvc, nil
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

	mutex sync.RWMutex
}

// CreateShards creates shards for data partition by given options
// 1) dump engine option into local disk
// 2) create shard storage struct
func (s *storageService) CreateShards(db string, option option.ShardOption, shardIDs ...int32) error {
	if len(shardIDs) == 0 {
		return fmt.Errorf("cannot create empty shard for db[%s]", db)
	}
	s.mutex.RLock()
	engine, ok := s.engines[db]
	s.mutex.RUnlock()

	if !ok {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		// double check
		engine, ok = s.engines[db]
		if !ok {
			// check engine config if nil
			// 1) not set when system init
			// 2) clean up when runtime
			if EngineConfig == nil {
				return fmt.Errorf("cannot create engine, because storage config is nil")
			}

			// create tsdb engine
			var err error
			engine, err = tsdb.NewEngine(db, EngineConfig.Path)
			if err != nil {
				return err
			}
			s.engines[db] = engine
		}
	}

	// create shards for database
	if err := engine.CreateShards(option, shardIDs...); err != nil {
		return err
	}

	return nil
}

// GetShard returns shard by given db and shard id, if not exist return nil
func (s *storageService) GetShard(db string, shardID int32) tsdb.Shard {
	s.mutex.RLock()
	engine, ok := s.engines[db]
	s.mutex.RUnlock()
	if !ok {
		return nil
	}
	return engine.GetShard(shardID)
}
