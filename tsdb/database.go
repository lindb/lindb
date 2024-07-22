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
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb/memdb"
)

//go:generate mockgen -source=./database.go -destination=./database_mock.go -package=tsdb

// Database represents an abstract time series database
type Database interface {
	// Name returns time series database's name
	Name() string
	// NumOfShards returns number of families in time series database
	NumOfShards() int
	// GetConfig return the configuration of database.
	GetConfig() *models.DatabaseConfig
	// GetOption returns the database options
	GetOption() *option.DatabaseOption
	// CreateShards creates families for data partition
	CreateShards(shardIDs []models.ShardID) error
	// GetShard returns shard by given shard id
	GetShard(shardID models.ShardID) (Shard, bool)
	// ExecutorPool returns the pool for querying tasks
	ExecutorPool() *ExecutorPool
	// Closer closes database's underlying resource
	io.Closer
	// MetaDB returns the metric metadata database include metric/tag/schema etc.
	MetaDB() index.MetricMetaDatabase
	// MemMetaDB returns memory metadata database.
	MemMetaDB() memdb.MetadataDatabase
	// FlushMeta flushes meta to disk
	FlushMeta() error
	// WaitFlushMetaCompleted waits flush metadata job completed.
	WaitFlushMetaCompleted()
	// Flush flushes memory data of all families to disk
	Flush() error
	// Drop drops current database include all data.
	Drop() error
	// TTL expires the data of each shard base on time to live.
	TTL()
	// EvictSegment evicts segment which long term no read operation.
	EvictSegment()
	// SetLimits sets database's limits.
	SetLimits(limits *models.Limits)
	// GetLimits returns database's limits.
	GetLimits() *models.Limits
}

// database implements Database for storing families,
// each shard represents a time series storage
type database struct {
	metaDB         index.MetricMetaDatabase
	config         *models.DatabaseConfig // meta configuration
	executorPool   *ExecutorPool          // executor pool for querying task
	shardSet       shardSet               // atomic value
	flushCondition *sync.Cond             // flush condition

	memMetaDB memdb.MetadataDatabase

	statistics   *metrics.DatabaseStatistics
	flushChecker DataFlushChecker

	name string // database-name
	dir  string

	mutex      sync.Mutex  // mutex for creating families
	isFlushing atomic.Bool // restrict flusher concurrency
}

// newDatabase creates the database instance
func newDatabase(
	databaseName string,
	cfg *models.DatabaseConfig,
	limits *models.Limits,
	flushChecker DataFlushChecker,
) (Database, error) {
	if err := cfg.Option.Validate(); err != nil {
		return nil, fmt.Errorf("database option is invalid, err: %s", err)
	}
	db := &database{
		name:         databaseName,
		flushChecker: flushChecker,
		config:       cfg,
		shardSet:     *newShardSet(),
		executorPool: &ExecutorPool{
			Filtering: concurrent.NewPool(
				databaseName+"-filtering-pool",
				runtime.GOMAXPROCS(-1), /*nRoutines*/
				time.Second*5,
				metrics.NewConcurrentStatistics(databaseName+"-filtering", linmetric.StorageRegistry),
			),
			Grouping: concurrent.NewPool(
				databaseName+"-grouping-pool",
				runtime.GOMAXPROCS(-1), /*nRoutines*/
				time.Second*5,
				metrics.NewConcurrentStatistics(databaseName+"-grouping", linmetric.StorageRegistry),
			),
			Scanner: concurrent.NewPool(
				databaseName+"-scanner-pool",
				runtime.GOMAXPROCS(-1), /*nRoutines*/
				time.Second*5,
				metrics.NewConcurrentStatistics(databaseName+"-scanner", linmetric.StorageRegistry),
			),
		},
		isFlushing:     *atomic.NewBool(false),
		flushCondition: sync.NewCond(&sync.Mutex{}),
		statistics:     metrics.NewDatabaseStatistics(databaseName),
	}
	dbPath, err0 := createDatabasePath(databaseName)
	if err0 != nil {
		return nil, err0
	}
	db.dir = dbPath
	if err := db.dumpDatabaseConfig(cfg); err != nil {
		return nil, err
	}
	if err := db.initMetadata(); err != nil {
		return nil, err
	}
	var err error
	defer func() {
		if err != nil && db.metaDB != nil {
			if e := db.metaDB.Close(); e != nil {
				engineLogger.Error("close metric metadata database err will create database",
					logger.Error(e), logger.String("db", databaseName))
			}
		}
	}()
	models.SetDatabaseLimits(databaseName, limits)

	db.memMetaDB = memdb.NewMetadataDatabase(db.config, db.metaDB)
	// load families if engine is existed
	var shard Shard
	if len(db.config.ShardIDs) > 0 {
		for _, shardID := range db.config.ShardIDs {
			shard, err = newShardFunc(db, shardID)
			if err != nil {
				return nil, fmt.Errorf("cannot create shard[%d] of database[%s] with error: %s",
					shardID, databaseName, err)
			}
			db.shardSet.InsertShard(shardID, shard)
		}
	}
	return db, nil
}

// SetLimits sets database's limits.
func (db *database) SetLimits(limits *models.Limits) {
	models.SetDatabaseLimits(db.name, limits)
}

// GetLimits returns database's limits.
func (db *database) GetLimits() *models.Limits {
	return models.GetDatabaseLimits(db.name)
}

// MetaDB returns the metric metadata database include metric/tag/schema etc.
func (db *database) MetaDB() index.MetricMetaDatabase {
	return db.metaDB
}

// MemMetaDB returns memory metadata database.
func (db *database) MemMetaDB() memdb.MetadataDatabase {
	return db.memMetaDB
}

// Name returns time series database's name
func (db *database) Name() string {
	return db.name
}

// NumOfShards returns number of families in time series database
func (db *database) NumOfShards() int {
	return db.shardSet.GetShardNum()
}

// GetConfig return the configuration of database.
func (db *database) GetConfig() *models.DatabaseConfig {
	return db.config
}

// GetOption returns the database options
func (db *database) GetOption() *option.DatabaseOption {
	return db.config.Option
}

// CreateShards creates families for data partition
func (db *database) CreateShards(
	shardIDs []models.ShardID,
) error {
	if len(shardIDs) == 0 {
		return fmt.Errorf("shardIDs list is empty")
	}
	for _, shardID := range shardIDs {
		if _, ok := db.GetShard(shardID); ok {
			continue
		}
		if err := db.createShard(shardID); err != nil {
			return err
		}
	}
	return nil
}

// createShard creates a new shard based on option
func (db *database) createShard(shardID models.ShardID) error {
	// be careful need do mutex unlock
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// double check
	if _, ok := db.GetShard(shardID); ok {
		return nil
	}
	// new shard
	createdShard, err := newShardFunc(
		db,
		shardID)
	if err != nil {
		return fmt.Errorf("create shard[%d] for engine[%s] with error: %s", shardID, db.name, err)
	}
	// using new engine option
	newCfg := &models.DatabaseConfig{Name: db.name, Option: db.config.Option, ShardIDs: db.config.ShardIDs}
	// add new shard id
	newCfg.ShardIDs = append(newCfg.ShardIDs, shardID)
	if err := db.dumpDatabaseConfig(newCfg); err != nil {
		// TODO: if dump config err, need close shard??
		return err
	}
	db.shardSet.InsertShard(shardID, createdShard)
	return nil
}

// GetShard returns shard by given shard id,
func (db *database) GetShard(shardID models.ShardID) (Shard, bool) {
	return db.shardSet.GetShard(shardID)
}

// ExecutorPool returns the query task execute pool
func (db *database) ExecutorPool() *ExecutorPool {
	return db.executorPool
}

// Close closes database's underlying resource
func (db *database) Close() error {
	// wait previous flush job completed
	db.WaitFlushMetaCompleted()

	if err := db.flushMeta(); err != nil {
		return err
	}

	db.memMetaDB.Close()
	for _, shardEntry := range db.shardSet.Entries() {
		thisShard := shardEntry.shard
		if err := thisShard.FlushIndex(); err != nil {
			engineLogger.Error(fmt.Sprintf(
				"close shard[%d] of database[%s]", shardEntry.shardID, db.name), logger.Error(err))
		}
	}
	if err := db.metaDB.Close(); err != nil {
		return err
	}
	for _, shardEntry := range db.shardSet.Entries() {
		thisShard := shardEntry.shard
		if err := thisShard.Close(); err != nil {
			engineLogger.Error(fmt.Sprintf(
				"close shard[%d] of database[%s]", shardEntry.shardID, db.name), logger.Error(err))
		}
	}
	return nil
}

// TTL expires the data of each shard base on time to live.
func (db *database) TTL() {
	for _, shardEntry := range db.shardSet.Entries() {
		thisShard := shardEntry.shard
		thisShard.TTL()
	}
}

// EvictSegment evicts segment which long term no read operation.
func (db *database) EvictSegment() {
	for _, shardEntry := range db.shardSet.Entries() {
		thisShard := shardEntry.shard
		thisShard.EvictSegment()
	}
}

// dumpDatabaseConfig persists option info to OPTIONS file
func (db *database) dumpDatabaseConfig(newConfig *models.DatabaseConfig) error {
	cfgPath := optionsPath(db.name)
	// write store info using toml format
	if err := encodeToml(cfgPath, newConfig); err != nil {
		return fmt.Errorf("write engine options to file[%s] error:%s", cfgPath, err)
	}
	db.config = newConfig
	return nil
}

// initMetadata initializes metadata backend storage
func (db *database) initMetadata() error {
	metaDB, err := newMetaDBFunc(db.name, metricsMetaPath(db.name))
	if err != nil {
		return err
	}
	db.metaDB = metaDB
	return nil
}

// FlushMeta flushes meta to disk.
func (db *database) FlushMeta() (err error) {
	// another flush process is running
	if !db.isFlushing.CompareAndSwap(false, true) {
		return nil
	}
	start := time.Now()
	defer func() {
		db.flushCondition.L.Lock()
		db.isFlushing.Store(false)
		db.flushCondition.L.Unlock()
		db.flushCondition.Broadcast()
		db.statistics.MetaDBFlushDuration.UpdateSince(start)
	}()
	if err := db.flushMeta(); err != nil {
		return err
	}
	return nil
}

// WaitFlushMetaCompleted waits flush metadata job completed.
func (db *database) WaitFlushMetaCompleted() {
	db.flushCondition.L.Lock()
	if db.isFlushing.Load() {
		db.flushCondition.Wait()
	}
	db.flushCondition.L.Unlock()
}

// Flush flushes memory data of all families to disk.
func (db *database) Flush() error {
	for _, shardEntry := range db.shardSet.Entries() {
		shard := shardEntry.shard
		db.flushChecker.requestFlushJob(&flushRequest{
			db: db,
			shards: map[models.ShardID]*flushShard{
				shard.ShardID(): {
					shard:    shard,
					families: GetFamilyManager().GetFamiliesByShard(shard),
				},
			},
			global: false,
		})
	}
	return nil
}

func (db *database) flushMeta() error {
	ch := make(chan error, 1)
	db.memMetaDB.Notify(&memdb.FlushEvent{
		Callback: func(err error) {
			ch <- err
		},
	})
	if err := <-ch; err != nil {
		db.statistics.MetaDBFlushFailures.Incr()
		return err
	}
	return nil
}

// Drop drops current database include all data.
func (db *database) Drop() error {
	if err := db.Close(); err != nil {
		return err
	}
	if err := removeDir(db.dir); err != nil {
		return err
	}
	return nil
}
