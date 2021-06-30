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

package service

import (
	"context"
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

	// FlushDatabase produces a signal to workers for flushing memory database by name
	FlushDatabase(ctx context.Context, databaseName string) bool

	// Close closes the time series engine
	Close()
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
	if err := db.CreateShards(databaseOption, shardIDs); err != nil {
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

func (s *storageService) FlushDatabase(ctx context.Context, databaseName string) bool {
	return s.engine.FlushDatabase(ctx, databaseName)
}

func (s *storageService) Close() {
	s.engine.Close()
}
