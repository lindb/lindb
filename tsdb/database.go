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
	"io"
	"runtime"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/tagkeymeta"
)

//go:generate mockgen -source=./database.go -destination=./database_mock.go -package=tsdb

// Database represents an abstract time series database
type Database interface {
	// Name returns time series database's name
	Name() string
	// NumOfShards returns number of families in time series database
	NumOfShards() int
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
	// Metadata returns the metadata include metric/tag
	Metadata() metadb.Metadata
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
}

// databaseConfig represents a database configuration about config and families
type databaseConfig struct {
	ShardIDs []models.ShardID       `toml:"shardIDs"`
	Option   *option.DatabaseOption `toml:"option"`
}

// database implements Database for storing families,
// each shard represents a time series storage
type database struct {
	name           string // database-name
	dir            string
	config         *databaseConfig // meta configuration
	executorPool   *ExecutorPool   // executor pool for querying task
	mutex          sync.Mutex      // mutex for creating families
	shardSet       shardSet        // atomic value
	metadata       metadb.Metadata // underlying metric metadata
	metaStore      kv.Store        // underlying meta kv store
	isFlushing     atomic.Bool     // restrict flusher concurrency
	flushCondition *sync.Cond      // flush condition

	statistics *metrics.DatabaseStatistics

	flushChecker DataFlushChecker
}

// newDatabase creates the database instance
func newDatabase(
	databaseName string,
	cfg *databaseConfig,
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
		if err != nil && db.metadata != nil {
			if e := db.metadata.Close(); e != nil {
				engineLogger.Error("close metadata err will create database",
					logger.Error(e), logger.String("db", databaseName))
			}
		}
	}()
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

// Metadata returns the metadata include metric/tag
func (db *database) Metadata() metadb.Metadata {
	return db.metadata
}

// Name returns time series database's name
func (db *database) Name() string {
	return db.name
}

// NumOfShards returns number of families in time series database
func (db *database) NumOfShards() int {
	return db.shardSet.GetShardNum()
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
	newCfg := &databaseConfig{Option: db.config.Option, ShardIDs: db.config.ShardIDs}
	// add new shard id
	newCfg.ShardIDs = append(newCfg.ShardIDs, shardID)
	if err := db.dumpDatabaseConfig(newCfg); err != nil {
		// TODO if dump config err, need close shard??
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

	if err := db.metadata.Close(); err != nil {
		return err
	}
	if err := kv.GetStoreManager().CloseStore(db.metaStore.Name()); err != nil {
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
func (db *database) dumpDatabaseConfig(newConfig *databaseConfig) error {
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
	// FIXME close kv store if err??
	metaStore, err := kv.GetStoreManager().CreateStore(tagMetaIndicator(db.name), kv.DefaultStoreOption())
	if err != nil {
		return err
	}
	tagMetaFamily, err := metaStore.CreateFamily(
		tagValueDir,
		kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           string(tagkeymeta.MergerName)})
	if err != nil {
		return err
	}
	db.metaStore = metaStore
	metadata, err := newMetadataFunc(context.TODO(), db.name, metricsMetaPath(db.name), tagMetaFamily)
	if err != nil {
		return err
	}
	db.metadata = metadata
	return nil
}

// FlushMeta flushes meta to disk.
func (db *database) FlushMeta() (err error) {
	// another flush process is running
	if !db.isFlushing.CAS(false, true) {
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
	if err := db.metadata.Flush(); err != nil {
		db.statistics.MetaDBFlushFailures.Incr()
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
