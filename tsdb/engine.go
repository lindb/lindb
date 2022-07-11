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

package tsdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
)

//go:generate mockgen -source=./engine.go -destination=./engine_mock.go -package=tsdb

var engineLogger = logger.GetLogger("TSDB", "Engine")

// Engine represents a time series engine
type Engine interface {
	// createDatabase creates database instance by database's name
	// return success when creating database's path successfully
	// called when CreateShards without database created
	createDatabase(databaseName string, dbOption *option.DatabaseOption) (Database, error)
	// CreateShards creates families for data partition by given options
	// 1) dump engine option into local disk
	// 2) create shard storage struct
	CreateShards(
		databaseName string,
		databaseOption *option.DatabaseOption,
		shardIDs ...models.ShardID,
	) error
	// GetShard returns shard by given db and shard id
	GetShard(databaseName string, shardID models.ShardID) (Shard, bool)
	// GetDatabase returns the time series database by given name
	GetDatabase(databaseName string) (Database, bool)
	// FlushDatabase produces a signal to workers for flushing memory database by name
	FlushDatabase(ctx context.Context, databaseName string) bool
	// DropDatabases drops databases, keep active database.
	DropDatabases(activeDatabases map[string]struct{})
	// TTL expires the data of each database base on time to live.
	TTL()
	// EvictSegment evicts segment which long term no read operation.
	EvictSegment()
	// Close closes the cached time series databases
	Close()
}

// engine implements Engine
type engine struct {
	mutex            sync.Mutex         // mutex for creating database
	dbSet            databaseSet        // atomic value, holding databaseName -> Database
	ctx              context.Context    // context
	cancel           context.CancelFunc // cancel function of flusher
	dataFlushChecker DataFlushChecker
}

// NewEngine creates an engine for manipulating the databases
func NewEngine() (Engine, error) {
	// create time series storage path
	if err := mkDirIfNotExist(config.GlobalStorageConfig().TSDB.Dir); err != nil {
		return nil, fmt.Errorf("create time sereis storage path[%s] erorr: %s",
			config.GlobalStorageConfig().TSDB.Dir, err)
	}
	e := &engine{
		dbSet: *newDatabaseSet(),
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())
	e.dataFlushChecker = newDataFlushChecker(e.ctx)
	e.dataFlushChecker.Start()

	if err := e.load(); err != nil {
		engineLogger.Error("load engine data error when create a new engine", logger.Error(err))
		// close opened engine
		e.Close()
		return nil, err
	}
	return e, nil
}

// createDatabase creates database instance by database's name
// return success when creating database's path successfully
func (e *engine) createDatabase(databaseName string, dbOption *option.DatabaseOption) (Database, error) {
	cfgPath := optionsPath(databaseName)
	cfg := &databaseConfig{Option: dbOption}
	engineLogger.Info("load database option from local storage", logger.String("path", cfgPath))
	if fileExist(cfgPath) {
		if err := decodeToml(cfgPath, cfg); err != nil {
			return nil, fmt.Errorf("load database[%s] config from file[%s] with error: %s",
				databaseName, cfgPath, err)
		}
	}
	db, err := newDatabaseFunc(databaseName, cfg, e.dataFlushChecker)
	if err != nil {
		return nil, err
	}
	e.dbSet.PutDatabase(databaseName, db)
	return db, nil
}

func (e *engine) CreateShards(
	databaseName string,
	databaseOption *option.DatabaseOption,
	shardIDs ...models.ShardID,
) error {
	if len(shardIDs) == 0 {
		return fmt.Errorf("cannot create empty shard for database[%s]", databaseName)
	}
	db, ok := e.GetDatabase(databaseName)
	if !ok {
		e.mutex.Lock()
		defer e.mutex.Unlock()
		if db, ok = e.GetDatabase(databaseName); !ok {
			// double check
			var err error
			db, err = e.createDatabase(databaseName, databaseOption)
			if err != nil {
				engineLogger.Error("failed to create database",
					logger.Error(err))
				return err
			}
			engineLogger.Info("create database successfully",
				logger.String("database", databaseName))
		}
	}

	// create families for database
	shardIDData := encoding.JSONMarshal(shardIDs)
	if err := db.CreateShards(shardIDs); err != nil {
		engineLogger.Error("failed to create shard", logger.String("shardIDs", string(shardIDData)))
		return err
	}
	engineLogger.Info("create shard successfully", logger.String("shardIDs", string(shardIDData)))
	return nil
}

// GetDatabase returns the time series database by given name
func (e *engine) GetDatabase(databaseName string) (Database, bool) {
	return e.dbSet.GetDatabase(databaseName)
}

// GetShard returns shard by given db and shard id
func (e *engine) GetShard(databaseName string, shardID models.ShardID) (Shard, bool) {
	if db, ok := e.GetDatabase(databaseName); ok {
		return db.GetShard(shardID)
	}
	return nil, false
}

// Close closes the cached time series databases
func (e *engine) Close() {
	if e.dataFlushChecker != nil {
		e.dataFlushChecker.Stop()
	}
	for dbName, db := range e.dbSet.Entries() {
		if err := db.Close(); err != nil {
			engineLogger.Error("close database",
				logger.String("name", dbName),
				logger.Error(err))
		}
	}
}

// FlushDatabase produces a signal to workers for flushing memory database by name
func (e *engine) FlushDatabase(_ context.Context, name string) bool {
	if db, ok := e.dbSet.GetDatabase(name); ok {
		if err := db.Flush(); err != nil {
			return false
		}
		return true
	}
	return false
}

// DropDatabases drops databases, keep active database.
func (e *engine) DropDatabases(activeDatabases map[string]struct{}) {
	for dbName, db := range e.dbSet.Entries() {
		_, ok := activeDatabases[dbName]
		if ok {
			continue
		}
		if err := db.Drop(); err != nil {
			engineLogger.Warn("drop database failure", logger.String("database", dbName), logger.Error(err))
			continue
		}
		e.dbSet.DropDatabase(dbName)
		engineLogger.Info("drop database successfully", logger.String("database", dbName))
	}
}

// TTL expires the data of each database base on time to live.
func (e *engine) TTL() {
	for _, db := range e.dbSet.Entries() {
		db.TTL()
	}
}

// EvictSegment evicts segment which long term no read operation.
func (e *engine) EvictSegment() {
	for _, db := range e.dbSet.Entries() {
		db.EvictSegment()
	}
}

// load the time series engines if exist
func (e *engine) load() error {
	databaseNames, err := listDir(config.GlobalStorageConfig().TSDB.Dir)
	if err != nil {
		return err
	}
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for _, databaseName := range databaseNames {
		_, err := e.createDatabase(databaseName, &option.DatabaseOption{}) // need load config from local file
		if err != nil {
			return err
		}
	}
	return nil
}
