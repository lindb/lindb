package service

import (
	"fmt"
	"sync"

	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source ./storage.go -destination=./storage_mock.go -package service

// StorageService represents a storage manage interface for time series engine
type StorageService interface {
	// CreateShards creates shards for data partition by given options
	// 1) dump engine option into local disk
	// 2) create shard storage struct
	CreateShards(
		databaseName string,
		databaseOption option.DatabaseOption,
		shardIDs ...int32,
	) error

	// GetDatabase returns database by given db-name
	GetDatabase(databaseName string) (tsdb.Database, bool)

	// GetShard returns shard by given db and shard id
	GetShard(databaseName string, shardID int32) (tsdb.Shard, bool)
}

// storageService implements StorageService interface
type storageService struct {
	engine tsdb.Engine
	mutex  sync.Mutex
}

// NewStorageService creates storage service instance for managing time series engine
func NewStorageService(engine tsdb.Engine) StorageService {
	return &storageService{
		engine: engine,
	}
}

func (s *storageService) CreateShards(
	databaseName string,
	databaseOption option.DatabaseOption,
	shardIDs ...int32,
) error {
	if len(shardIDs) == 0 {
		return fmt.Errorf("cannot create empty shard for database[%s]", databaseName)
	}
	db, ok := s.GetDatabase(databaseName)
	if !ok {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		// double check
		db, ok = s.GetDatabase(databaseName)
		if !ok {
			// create a time series engine
			var err error
			db, err = s.engine.CreateDatabase(databaseName)
			if err != nil {
				return err
			}
		}
	}

	// create shards for database
	if err := db.CreateShards(databaseOption, shardIDs...); err != nil {
		return err
	}
	return nil
}

func (s *storageService) GetShard(databaseName string, shardID int32) (tsdb.Shard, bool) {
	db, ok := s.GetDatabase(databaseName)
	if !ok {
		return nil, false
	}
	return db.GetShard(shardID)
}

func (s *storageService) GetDatabase(databaseName string) (tsdb.Database, bool) {
	return s.engine.GetDatabase(databaseName)
}
