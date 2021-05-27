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
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/concurrent"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/tagkeymeta"
)

//go:generate mockgen -source=./database.go -destination=./database_mock.go -package=tsdb

// for testing
var (
	newMetadataFunc = metadb.NewMetadata
	newShardFunc    = newShard
	encodeToml      = ltoml.EncodeToml
)

const (
	options       = "OPTIONS"
	shardDir      = "shard"
	metricMetaDir = "metric"
	tagMetaDir    = "tag"
	tagValueDir   = "tag_value"
)

// Database represents an abstract time series database
type Database interface {
	// Name returns time series database's name
	Name() string
	// NumOfShards returns number of shards in time series database
	NumOfShards() int
	// GetOption returns the data base options
	GetOption() option.DatabaseOption
	// CreateShards creates shards for data partition
	CreateShards(option option.DatabaseOption, shardIDs []int32) error
	// GetShard returns shard by given shard id
	GetShard(shardID int32) (Shard, bool)
	// ExecutorPool returns the pool for querying tasks
	ExecutorPool() *ExecutorPool
	// Close closes database's underlying resource
	io.Closer
	// Metadata returns the metadata include metric/tag
	Metadata() metadb.Metadata
	// FlushMeta flushes meta to disk
	FlushMeta() error
	// FLush flushes memory data of all shards to disk
	Flush() error
}

// databaseConfig represents a database configuration about config and shards
type databaseConfig struct {
	ShardIDs []int32               `toml:"shardIDs"`
	Option   option.DatabaseOption `toml:"option"`
}

// database implements Database for storing shards,
// each shard represents a time series storage
type database struct {
	name         string          // database-name
	path         string          // database root path
	config       *databaseConfig // meta configuration
	executorPool *ExecutorPool   // executor pool for querying task
	shards       sync.Map        // shardID(int32)->shard(Shard)
	numOfShards  atomic.Int32    // counter
	mutex        sync.Mutex      // mutex for creating shards
	metadata     metadb.Metadata // underlying metric metadata
	metaStore    kv.Store        // underlying meta kv store
	isFlushing   atomic.Bool     // restrict flusher concurrency

	flushChecker DataFlushChecker
}

// newDatabase creates the database instance
func newDatabase(databaseName string, databasePath string, cfg *databaseConfig,
	flushChecker DataFlushChecker,
) (Database, error) {
	db := &database{
		name:         databaseName,
		path:         databasePath,
		flushChecker: flushChecker,
		config:       cfg,
		numOfShards:  *atomic.NewInt32(0),
		executorPool: &ExecutorPool{
			Filtering: concurrent.NewPool(
				databaseName+"-filtering-pool",
				runtime.NumCPU(), /*nRoutines*/
				time.Second*5),
			Grouping: concurrent.NewPool(
				databaseName+"-grouping-pool",
				runtime.NumCPU(), /*nRoutines*/
				time.Second*5),
			Scanner: concurrent.NewPool(
				databaseName+"-scanner-pool",
				runtime.NumCPU(), /*nRoutines*/
				time.Second*5),
		},
		isFlushing: *atomic.NewBool(false),
	}
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
	// load shards if engine is exist
	var shard Shard
	if len(db.config.ShardIDs) > 0 {
		for _, shardID := range db.config.ShardIDs {
			shard, err = newShardFunc(
				db,
				shardID,
				filepath.Join(databasePath, shardDir, strconv.Itoa(int(shardID))),
				db.config.Option)
			if err != nil {
				return nil, fmt.Errorf("cannot create shard[%d] of database[%s] with error: %s",
					shardID, databaseName, err)
			}
			db.shards.Store(shardID, shard)
			db.numOfShards.Inc()
		}
	}

	return db, nil
}

func (db *database) Metadata() metadb.Metadata {
	return db.metadata
}

func (db *database) Name() string {
	return db.name
}

func (db *database) NumOfShards() int {
	return int(db.numOfShards.Load())
}

func (db *database) GetOption() option.DatabaseOption {
	return db.config.Option
}

// CreateShards creates shards for data partition
func (db *database) CreateShards(
	option option.DatabaseOption,
	shardIDs []int32,
) error {
	if len(shardIDs) == 0 {
		return fmt.Errorf("shardIDs list is empty")
	}
	for _, shardID := range shardIDs {
		_, ok := db.GetShard(shardID)
		if ok {
			continue
		}
		if err := db.createShard(shardID, option); err != nil {
			return err
		}
	}
	return nil
}

// createShard creates a new shard based on option
func (db *database) createShard(shardID int32, option option.DatabaseOption) error {
	// be careful need do mutex unlock
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// double check
	_, ok := db.GetShard(shardID)
	if ok {
		return nil
	}
	// new shard
	createdShard, err := newShardFunc(
		db,
		shardID,
		filepath.Join(db.path, shardDir, strconv.Itoa(int(shardID))),
		option)
	if err != nil {
		return fmt.Errorf("create shard[%d] for engine[%s] with error: %s", shardID, db.name, err)
	}
	// using new engine option
	newCfg := &databaseConfig{Option: option, ShardIDs: db.config.ShardIDs}
	// add new shard id
	newCfg.ShardIDs = append(newCfg.ShardIDs, shardID)
	if err := db.dumpDatabaseConfig(newCfg); err != nil {
		return err
	}
	db.shards.Store(shardID, createdShard)
	db.numOfShards.Inc()
	return nil
}

// GetShard returns shard by given shard id,
func (db *database) GetShard(shardID int32) (Shard, bool) {
	item, ok := db.shards.Load(shardID)
	if !ok {
		return nil, false
	}
	return item.(Shard), true
}

// ExecutorPool returns the query task execute pool
func (db *database) ExecutorPool() *ExecutorPool {
	return db.executorPool
}

// Close closes database's underlying resource
func (db *database) Close() error {
	if err := db.metadata.Close(); err != nil {
		return err
	}
	if err := db.metaStore.Close(); err != nil {
		return err
	}
	db.shards.Range(func(key, value interface{}) bool {
		thisShard := value.(Shard)
		if err := thisShard.Close(); err != nil {
			engineLogger.Error(fmt.Sprintf(
				"close shard[%d] of database[%s]", key.(int32), db.name), logger.Error(err))
		}
		return true
	})
	return nil
}

// dumpDatabaseConfig persists option info to OPTIONS file
func (db *database) dumpDatabaseConfig(newConfig *databaseConfig) error {
	cfgPath := optionsPath(db.path)
	// write store info using toml format
	if err := encodeToml(cfgPath, newConfig); err != nil {
		return fmt.Errorf("write engine info to file[%s] error:%s", cfgPath, err)
	}
	db.config = newConfig
	return nil
}

// initMetadata initializes metadata backend storage
func (db *database) initMetadata() error {
	metaStoreOption := kv.DefaultStoreOption(filepath.Join(db.path, metaDir, tagMetaDir))
	//FIXME close kv store if err??
	metaStore, err := newKVStoreFunc(metaStoreOption.Path, metaStoreOption)
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
	metadata, err := newMetadataFunc(context.TODO(), db.name, filepath.Join(db.path, metaDir, metricMetaDir), tagMetaFamily)
	if err != nil {
		return err
	}
	db.metadata = metadata
	return nil
}

func (db *database) FlushMeta() (err error) {
	// another flush process is running
	if !db.isFlushing.CAS(false, true) {
		return nil
	}
	defer db.isFlushing.Store(false)

	return db.metadata.Flush()
}

// FLush flushes memory data of all shards to disk
func (db *database) Flush() error {
	db.shards.Range(func(key, value interface{}) bool {
		shard := value.(Shard)
		db.flushChecker.requestFlushJob(shard, false)
		return true
	})
	return nil
}

// optionsPath returns options file path
func optionsPath(path string) string {
	return filepath.Join(path, options)
}
