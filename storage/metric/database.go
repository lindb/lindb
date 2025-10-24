package metric

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/ltoml"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/flush"
	"github.com/lindb/lindb/storage/metric/memdb"
	"github.com/lindb/lindb/storage/store"
)

func init() {
	store.RegisterEngine(option.Metric, newDatabase)
}

// Database implements Database for storing families,
// each shard represents a time series storage
type Database struct {
	base.Database

	metaDB         index.MetricMetaDatabase
	executorPool   *ExecutorPool // executor pool for querying task
	flushCondition *sync.Cond    // flush condition

	memMetaDB memdb.MetadataDatabase

	statistics *metrics.DatabaseStatistics

	dir string

	mutex      sync.Mutex  // mutex for creating families
	isFlushing atomic.Bool // restrict flusher concurrency

	logger logger.Logger
}

// newDatabase creates the database instance
func newDatabase(
	databaseName string,
	cfg *models.DatabaseConfig,
	limits *models.Limits,
	check flush.Checker,
) (store.Database, error) {
	if err := cfg.Option.Validate(); err != nil {
		return nil, fmt.Errorf("database option is invalid, err: %s", err)
	}
	db := &Database{
		Database: base.Database{
			DatabaseName: databaseName,
			Options:      cfg,
			ShardSet:     *store.NewShardSet(),
		},
		// flushChecker: flushChecker,
		executorPool: &ExecutorPool{
			MetaFetcher: concurrent.NewPool(
				databaseName+"-meta-fetcher-pool",
				runtime.GOMAXPROCS(-1), /*nRoutines*/
				time.Second*5,
				metrics.NewConcurrentStatistics(databaseName+"-meta-fecher", linmetric.StorageRegistry),
			),
			DataFetcher: concurrent.NewPool(
				databaseName+"-data-fetcher-pool",
				runtime.GOMAXPROCS(-1), /*nRoutines*/
				time.Second*5,
				metrics.NewConcurrentStatistics(databaseName+"-data-fetcher", linmetric.StorageRegistry),
			),
			Reducer: concurrent.NewPool(
				databaseName+"-reducer-pool",
				runtime.GOMAXPROCS(-1), /*nRoutines*/
				time.Second*5,
				metrics.NewConcurrentStatistics(databaseName+"-reducer", linmetric.StorageRegistry),
			),
		},
		isFlushing:     *atomic.NewBool(false),
		flushCondition: sync.NewCond(&sync.Mutex{}),
		statistics:     metrics.NewDatabaseStatistics(databaseName),
		logger:         logger.GetLogger("Metric", "Database"),
	}
	dbPath, err0 := store.CreateDatabasePath(databaseName)
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
				db.logger.Error("close metric metadata database err will create database",
					logger.Error(e), logger.String("db", databaseName))
			}
		}
	}()
	models.SetDatabaseLimits(databaseName, limits)

	db.CreateShardFn = db.createShard

	db.memMetaDB = memdb.NewMetadataDatabase(db.Options, db.metaDB)
	// load shards if engine is existed
	if len(db.Options.ShardIDs) > 0 {
		if err = db.CreateShards(db.Options.ShardIDs); err != nil {
			return nil, fmt.Errorf("cannot create shards of database[%s] with error: %s",
				databaseName, err)
		}
	}
	return db, nil
}

// MetaDB returns the metric metadata database include metric/tag/schema etc.
func (db *Database) MetaDB() index.MetricMetaDatabase {
	return db.metaDB
}

// MemMetaDB returns memory metadata database.
func (db *Database) MemMetaDB() memdb.MetadataDatabase {
	return db.memMetaDB
}

func (db *Database) FindMatchSmallestInterval(interval timeutil.Interval) timeutil.Interval {
	return db.Options.Option.FindMatchSmallestInterval(interval)
}

// createShard creates a new shard based on option
func (db *Database) createShard(shardID models.ShardID) (store.Shard, error) {
	return newShard(
		db,
		shardID)
}

// ExecutorPool returns the query task execute pool
func (db *Database) ExecutorPool() *ExecutorPool {
	return db.executorPool
}

// Close closes database's underlying resource
func (db *Database) Close() error {
	// wait previous flush job completed
	db.WaitFlushMetaCompleted()

	if err := db.flushMeta(); err != nil {
		return err
	}

	db.memMetaDB.Close()
	for _, shardEntry := range db.ShardSet.Entries() {
		thisShard := shardEntry.Shard.(*Shard)
		if err := thisShard.FlushIndex(); err != nil {
			db.logger.Error(fmt.Sprintf(
				"flush shard[%d] of database[%s]", shardEntry.ShardID, db.DatabaseName), logger.Error(err))
		}
	}
	if err := db.metaDB.Close(); err != nil {
		return err
	}
	for _, shardEntry := range db.ShardSet.Entries() {
		thisShard := shardEntry.Shard
		if err := thisShard.Close(); err != nil {
			db.logger.Error(fmt.Sprintf(
				"close shard[%d] of database[%s]", shardEntry.ShardID, db.DatabaseName), logger.Error(err))
		}
	}
	return nil
}

// TTL expires the data of each shard base on time to live.
func (db *Database) TTL() {
	// for _, shardEntry := range db.shardSet.Entries() {
	// 	thisShard := shardEntry.shard
	// 	thisShard.TTL()
	// }
}

// EvictSegment evicts segment which long term no read operation.
func (db *Database) EvictSegment() {
	// for _, shardEntry := range db.ShardSet.Entries() {
	// 	thisShard := shardEntry.shard
	// 	thisShard.EvictSegment()
	// }
}

// dumpDatabaseConfig persists option info to OPTIONS file
func (db *Database) dumpDatabaseConfig(newConfig *models.DatabaseConfig) error {
	cfgPath := store.OptionsPath(db.DatabaseName)
	// write store info using toml format
	if err := ltoml.EncodeToml(cfgPath, newConfig); err != nil {
		return fmt.Errorf("write engine options to file[%s] error:%s", cfgPath, err)
	}
	db.Options = newConfig
	return nil
}

// initMetadata initializes metadata backend storage
func (db *Database) initMetadata() error {
	metaDB, err := index.NewMetricMetaDatabase(db.DatabaseName, metricsMetaPath(db.DatabaseName))
	if err != nil {
		return err
	}
	db.metaDB = metaDB
	return nil
}

// FlushMeta flushes meta to disk.
func (db *Database) FlushMeta() (err error) {
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
func (db *Database) WaitFlushMetaCompleted() {
	db.flushCondition.L.Lock()
	if db.isFlushing.Load() {
		db.flushCondition.Wait()
	}
	db.flushCondition.L.Unlock()
}

// Flush flushes memory data of all families to disk.
func (db *Database) Flush() error {
	// for _, shardEntry := range db.ShardSet.Entries() {
	// 	shard := shardEntry.shard
	// 	db.flushChecker.requestFlushJob(&flushRequest{
	// 		db: db,
	// 		shards: map[models.ShardID]*flushShard{
	// 			shard.ShardID(): {
	// 				shard:    shard,
	// 				families: GetFamilyManager().GetFamiliesByShard(shard),
	// 			},
	// 		},
	// 		global: false,
	// 	})
	// }
	return nil
}

func (db *Database) flushMeta() error {
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
func (db *Database) Drop() error {
	// if err := db.Close(); err != nil {
	// 	return err
	// }
	// if err := removeDir(db.dir); err != nil {
	// 	return err
	// }
	return nil
}

func (db *Database) GetStream() flow.Stream {
	return newStream(6 * 60)
}

type stream struct {
	values []float64
}

func newStream(size int) flow.Stream {
	return &stream{
		values: make([]float64, size),
	}
}

func (s *stream) SetAtStep(step int, value float64, fn func(a, b float64) float64) {
	s.values[step] = fn(s.values[step], value)
}

func (s *stream) GetAtStep(step int) float64 {
	return s.values[step]
}

func (s *stream) Merge(stream flow.Stream, fn func(a, b float64) float64) {
	panic("not implemented")
}

func (s *stream) Values() []float64 {
	panic("not implemented")
}
