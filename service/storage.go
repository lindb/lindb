package service

import (
	"fmt"
	"sync"

	"github.com/lindb/lindb/config"

	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source ./storage.go -destination=./storage_mock.go -package service

// StorageService represents a storage manage interface for time series engine
type StorageService interface {
	// CreateShards creates shards for data partition
	CreateShards(db string, option option.ShardOption, shardIDs ...int32) error
	// GetEngine returns engine by given db name, if not exist return nil
	GetEngine(db string) tsdb.Engine
	// GetShard returns shard by given db and shard id, if not exist return nil
	GetShard(db string, shardID int32) tsdb.Shard
}

// storageService implements StorageService interface
type storageService struct {
	engines sync.Map

	factory tsdb.EngineFactory
	config  config.Engine
	mutex   sync.Mutex
}

// NewStorageService creates storage service instance for managing time series engine
func NewStorageService(config config.Engine, factory tsdb.EngineFactory) StorageService {
	return &storageService{
		factory: factory,
		config:  config,
	}
}

// CreateShards creates shards for data partition by given options
// 1) dump engine option into local disk
// 2) create shard storage struct
func (s *storageService) CreateShards(db string, option option.ShardOption, shardIDs ...int32) error {
	if len(shardIDs) == 0 {
		return fmt.Errorf("cannot create empty shard for db[%s]", db)
	}
	engine := s.GetEngine(db)
	if engine == nil {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		// double check
		engine = s.GetEngine(db)
		if engine == nil {
			// create time series engine
			var err error
			engine, err = s.factory.CreateEngine(db, s.config.Path)
			if err != nil {
				return err
			}
			s.engines.Store(db, engine)
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
	engine := s.GetEngine(db)
	if engine == nil {
		return nil
	}
	return engine.GetShard(shardID)
}

// GetEngine returns engine by given db name, if not exist return nil
func (s *storageService) GetEngine(db string) tsdb.Engine {
	engine, _ := s.engines.Load(db)
	e, ok := engine.(tsdb.Engine)
	if ok {
		return e
	}
	return nil
}
