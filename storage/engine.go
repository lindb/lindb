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

package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/ltoml"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/storage/store"
)

//go:generate mockgen -source=./engine.go -destination=./engine_mock.go -package=storage

var engineLogger = logger.GetLogger("Storage", "Engine")

// Engine represents a time series engine
type Engine interface {
	// createDatabase creates database instance by database's name
	// return success when creating database's path successfully
	// called when CreateShards without database created
	createDatabase(databaseName string, dbOption *option.DatabaseOption) (store.Database, error)
	// CreateShards creates families for data partition by given options
	// 1) dump engine option into local disk
	// 2) create shard storage struct
	CreateShards(
		databaseName string,
		databaseOption *option.DatabaseOption,
		shardIDs ...models.ShardID,
	) error
	// SetDatabaseLimits sets database's limits.
	SetDatabaseLimits(database string, limits *models.Limits)
	// GetShard returns shard by given db and shard id
	GetShard(databaseName string, shardID models.ShardID) (store.Shard, bool)
	// GetDatabase returns the time series database by given name
	GetDatabase(databaseName string) (store.Database, bool)
	// GetAllDatabases returns all databases.
	GetAllDatabases() map[string]store.Database
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
	dbSet store.DatabaseSet // atomic value, holding databaseName -> Database

	ctx    context.Context    // context
	cancel context.CancelFunc // cancel function of flusher
	// dataFlushChecker DataFlushChecker
	mutex sync.Mutex // mutex for creating database

	databases map[string]store.Database
}

// NewEngine creates an engine for manipulating the databases
func NewEngine() (Engine, error) {
	// create time series storage path
	if err := fileutil.MkDirIfNotExist(config.GlobalStorageConfig().TSDB.Dir); err != nil {
		return nil, fmt.Errorf("create time sereis storage path[%s] erorr: %s",
			config.GlobalStorageConfig().TSDB.Dir, err)
	}
	e := &engine{
		dbSet: *store.NewDatabaseSet(),
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())
	e.databases = make(map[string]store.Database)
	// e.dataFlushChecker = newDataFlushChecker(e.ctx)
	// e.dataFlushChecker.Start()

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
func (e *engine) createDatabase(databaseName string, dbOption *option.DatabaseOption) (store.Database, error) {
	cfgPath := store.OptionsPath(databaseName)
	cfg := &models.DatabaseConfig{Name: databaseName, Option: dbOption}
	engineLogger.Info("load database option from local storage", logger.String("path", cfgPath))
	if fileutil.Exist(cfgPath) {
		if err := ltoml.DecodeToml(cfgPath, cfg); err != nil {
			return nil, fmt.Errorf("load database[%s] config from file[%s] with error: %s",
				databaseName, cfgPath, err)
		}
	}
	limits := store.LimitsPath(databaseName)
	limitCfg := models.NewDefaultLimits()
	if fileutil.Exist(limits) {
		if err := ltoml.DecodeToml(limits, limitCfg); err != nil {
			return nil, fmt.Errorf("load database[%s] limits config from file[%s] with error: %s",
				databaseName, cfgPath, err)
		}
	}
	db, err := store.CreateDatabase(databaseName, cfg, limitCfg, nil)
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

// SetDatabaseLimits sets database's limits.
func (e *engine) SetDatabaseLimits(database string, limits *models.Limits) {
	db, ok := e.dbSet.GetDatabase(database)
	if ok {
		if err := ltoml.WriteConfig(store.LimitsPath(database), limits.TOML()); err != nil {
			engineLogger.Warn("write limits config failure", logger.Error(err))
		}
		db.SetLimits(limits)
	}
}

// GetDatabase returns the time series database by given name
func (e *engine) GetDatabase(databaseName string) (store.Database, bool) {
	return e.dbSet.GetDatabase(databaseName)
}

func (e *engine) GetDatabase2(databaseName string) (store.Database, bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	db, ok := e.databases[databaseName]
	return db, ok
}

// GetAllDatabases returns all databases.
func (e *engine) GetAllDatabases() map[string]store.Database {
	return e.dbSet.Entries()
}

// GetShard returns shard by given db and shard id
func (e *engine) GetShard(databaseName string, shardID models.ShardID) (store.Shard, bool) {
	if db, ok := e.GetDatabase(databaseName); ok {
		return db.GetShard(shardID)
	}
	return nil, false
}

// Close closes the cached time series databases
func (e *engine) Close() {
	// if e.dataFlushChecker != nil {
	// 	e.dataFlushChecker.Stop()
	// }
	for dbName, db := range e.dbSet.Entries() {
		if err := db.Close(); err != nil {
			engineLogger.Error("close database",
				logger.String("name", dbName),
				logger.Error(err))
		}
	}

	for _, db := range e.databases {
		db.Close()
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
	databaseNames, err := fileutil.GetDirectoryList(config.GlobalStorageConfig().TSDB.Dir)
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
