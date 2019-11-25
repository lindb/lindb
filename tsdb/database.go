package tsdb

import (
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
)

//go:generate mockgen -source=./database.go -destination=./database_mock.go -package=tsdb

const (
	options          = "OPTIONS"
	shardDir         = "shard"
	metaDir          = "meta"
	metricNameIDsDir = "metric_nameid"
	metricMetaDir    = "metric_meta"
)

// Database represents an abstract time series database
type Database interface {
	// Name returns time series database's name
	Name() string
	// NumOfShards returns number of shards in time series database
	NumOfShards() int
	// CreateShards creates shards for data partition
	CreateShards(option option.DatabaseOption, shardIDs ...int32) error
	// GetShard returns shard by given shard id
	GetShard(shardID int32) (Shard, bool)
	// ExecutorPool returns the pool for querying tasks
	ExecutorPool() *ExecutorPool
	// Close closes database's underlying resource
	io.Closer
	// IDGetter returns the id getter
	IDGetter() metadb.IDGetter
	// Flush flushes meta to disk
	FlushMeta() error
	// Range is the proxy method for iterating shards
	Range(f func(key, value interface{}) bool)

	// initIDSequencer loads the meta store to initialize id sequencer
	initIDSequencer() error
}

// databaseConfig represents a database configuration about config and shards
type databaseConfig struct {
	ShardIDs []int32               `toml:"shardIDs"`
	Option   option.DatabaseOption `toml:"databaseOption"`
}

// database implements Database for storing shards,
// each shard represents a time series storage
type database struct {
	name         string             // database-name
	path         string             // database root path
	config       *databaseConfig    // meta configuration
	executorPool *ExecutorPool      // executor pool for querying task
	shards       sync.Map           // shardID(int32)->shard(Shard)
	numOfShards  atomic.Int32       // counter
	mutex        sync.Mutex         // mutex for creating shards
	idSequencer  metadb.IDSequencer // database-level reused object
	metaStore    kv.Store           // underlying meta kv store
	isFlushing   atomic.Bool        // restrict flusher concurrency
}

func newDatabase(
	databaseName string,
	databasePath string,
	cfg *databaseConfig,
) (
	db *database,
	err error,
) {
	db = &database{
		name:        databaseName,
		path:        databasePath,
		config:      cfg,
		numOfShards: *atomic.NewInt32(0),
		executorPool: &ExecutorPool{
			Scanners: concurrent.NewPool(
				runtime.NumCPU(), /*nRoutines*/
				time.Second*5),
			Mergers: concurrent.NewPool(
				runtime.NumCPU(),
				time.Second*5),
		},
		isFlushing: *atomic.NewBool(false),
	}
	if err = db.initIDSequencer(); err != nil {
		return nil, err
	}
	// load shards if engine is exist
	if len(db.config.ShardIDs) > 0 {
		for _, shardID := range db.config.ShardIDs {
			shard, err := newShard(
				shardID,
				filepath.Join(databasePath, shardDir, strconv.Itoa(int(shardID))),
				db.idSequencer,
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

func (db *database) Name() string {
	return db.name
}

func (db *database) NumOfShards() int {
	return int(db.numOfShards.Load())
}

func (db *database) IDGetter() metadb.IDGetter {
	return db.idSequencer
}

func (db *database) CreateShards(
	option option.DatabaseOption,
	shardIDs ...int32,
) error {
	if len(shardIDs) == 0 {
		return fmt.Errorf("shardIDs list is empty")
	}
	for _, shardID := range shardIDs {
		_, ok := db.GetShard(shardID)
		if ok {
			continue
		}
		// be careful need do mutex unlock
		db.mutex.Lock()
		// double check
		_, ok = db.GetShard(shardID)
		if ok {
			continue
		}
		// new shard
		createdShard, err := newShard(
			shardID,
			filepath.Join(db.path, shardDir, strconv.Itoa(int(shardID))),
			db.idSequencer,
			option)
		if err != nil {
			db.mutex.Unlock()
			return fmt.Errorf("create shard[%d] for engine[%s] with error: %s", shardID, db.name, err)
		}
		// using new engine option
		newCfg := &databaseConfig{Option: option, ShardIDs: db.config.ShardIDs}
		// add new shard id
		newCfg.ShardIDs = append(newCfg.ShardIDs, shardID)
		if err := db.dumpDatabaseConfig(newCfg); err != nil {
			db.mutex.Unlock()
			return err
		}
		db.shards.Store(shardID, createdShard)
		db.numOfShards.Inc()
		db.mutex.Unlock()
	}
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
	db.shards.Range(func(key, value interface{}) bool {
		thisShard := value.(Shard)
		if err := thisShard.Close(); err != nil {
			engineLogger.Error(fmt.Sprintf(
				"close shard[%d] of database[%s]", key.(int32), db.name), logger.Error(err))
		}
		return true
	})
	if err := db.FlushMeta(); err != nil {
		engineLogger.Error(fmt.Sprintf(
			"flush meta database[%s]", db.name), logger.Error(err))
	}
	return db.metaStore.Close()
}

// dumpDatabaseConfig persists option info to OPTIONS file
func (db *database) dumpDatabaseConfig(newConfig *databaseConfig) error {
	cfgPath := optionsPath(db.path)
	// write store info using toml format
	if err := ltoml.EncodeToml(cfgPath, newConfig); err != nil {
		return fmt.Errorf("write engine info to file[%s] error:%s", cfgPath, err)
	}
	db.config = newConfig
	return nil
}

func (db *database) initIDSequencer() error {
	metaStoreOption := kv.DefaultStoreOption(filepath.Join(db.path, metaDir))
	metaStore, err := kv.NewStore(metaStoreOption.Path, metaStoreOption)
	if err != nil {
		return err
	}
	metricMetaFamily, err := metaStore.CreateFamily(
		metricMetaDir,
		kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           metricMetaMerger})
	if err != nil {
		return err
	}
	metricNameIDsFamily, err := metaStore.CreateFamily(
		metricNameIDsDir,
		kv.FamilyOption{
			CompactThreshold: 0,
			Merger:           metricNameIDsMerger})
	if err != nil {
		return err
	}
	db.metaStore = metaStore
	db.idSequencer = metadb.NewIDSequencer(metricNameIDsFamily, metricMetaFamily)
	return db.idSequencer.Recover()
}

func (db *database) FlushMeta() (err error) {
	// another flush process is running
	if !db.isFlushing.CAS(false, true) {
		return nil
	}
	defer db.isFlushing.Store(false)

	if err = db.idSequencer.FlushMetricsMeta(); err != nil {
		return err
	}
	if err = db.idSequencer.FlushNameIDs(); err != nil {
		return err
	}
	return nil
}

func (db *database) Range(f func(key, value interface{}) bool) {
	db.shards.Range(f)
}

// optionsPath returns options file path
func optionsPath(path string) string {
	return filepath.Join(path, options)
}
