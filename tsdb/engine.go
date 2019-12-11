package tsdb

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/monitoring"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
)

var (
	globalMemoryUsageCheckInterval = *atomic.NewDuration(time.Second * 1)
	shardMemoryUsageCheckInterval  = *atomic.NewDuration(time.Minute)
	flushMetaInterval              = *atomic.NewDuration(time.Hour)
)

//go:generate mockgen -source=./engine.go -destination=./engine_mock.go -package=tsdb

var engineLogger = logger.GetLogger("tsdb", "Engine")

// Engine represents a time series engine
type Engine interface {
	// CreateDatabase creates database instance by database's name
	// return success when creating database's path successfully
	CreateDatabase(databaseName string) (Database, error)
	// GetDatabase returns the time series database by given name
	GetDatabase(databaseName string) (Database, bool)
	// Close closes the cached time series databases
	Close()

	// There are 4 flush policies of the Engine as below:
	// 1. FullFlush
	//    highest priority, triggered by external API from the users.
	//    this action will blocks any other flush checkers.
	// 2. GlobalMemoryUsageChecker
	//    This checker will check the global memory usage of the host periodically,
	//    when the metric is above MemoryHighWaterMark, a `watermarkFlusher` will be spawned
	//    whose responsibility is to flush the biggest shard until memory is lower than  MemoryLowWaterMark.
	// 3. ShardMemoryUsageChecker
	//    This checker will check each shard's memory usage periodically,
	//    If this shard is above ShardMemoryUsedThreshold. it will be flushed to disk.
	// 4. DatabaseMetaFlusher
	//    It is a simple checker which flush the meta of database to disk periodically.
	//
	// a). Each shard or database is restricted to flush by one goroutine at the same time via CAS operation;
	// b). The flush workers runs concurrently;
	// c). All unit will be flushed when closing;

	// FLushDatabase produces a signal to workers for flushing memory database by name
	FlushDatabase(ctx context.Context, databaseName string) bool
	// FlushAll produces a signal to workers for flushing all
	FlushAll()
	// flushDatabase is the real method for flushing certain database
	// called by FlushDatabase and flushAllDatabasesAndShards
	flushDatabase(ctx context.Context, db Database) bool
	// globalMemoryUsageChecker checks global memory usage periodically,
	// The biggest shard's will be flushed until memory usage is down MemoryLowWaterMark.
	globalMemoryUsageChecker(ctx context.Context)
	// shardMemoryUsageChecker checks shard memory usage periodically
	shardMemoryUsageChecker(ctx context.Context)
	// watermarkFlusher will be spawned if memory-usage is above high-watermark,
	// it will be pended when memory-usage is lower than low-watermark or full-flushing is enabled
	watermarkFlusher(ctx context.Context)
	// 	databaseMetaFlusher flushes database meta periodically
	databaseMetaFlusher(ctx context.Context)
	// flushShardAboveMemoryUsageThreshold flushes the shard whose memory usage is above threshold
	flushShardAboveMemoryUsageThreshold(ctx context.Context)
	// flushBiggestMemoryUsageShard flushes the biggest shard's memdb
	flushBiggestMemoryUsageShard(ctx context.Context)
	// flushAllDatabases sends all databases to flush
	flushAllDatabases(ctx context.Context)
	// flushAllDatabasesAndShards sends all databases and shards to flush
	flushAllDatabasesAndShards(ctx context.Context)
	// flushWorker is daemon goroutine who flushes data of shard or database
	flushWorker(ctx context.Context)
}

// engine implements Engine
type engine struct {
	cfg                  config.TSDB                 // the common cfg of time series database
	databases            sync.Map                    // databaseName -> Database
	ctx                  context.Context             // context
	cancel               context.CancelFunc          // cancel function of flusher
	shardToFlushCh       chan Shard                  // shard to flush
	memoryStatGetterFunc monitoring.MemoryStatGetter // used for mocking
	databaseToFlushCh    chan Database               // database to flush
	isFullFlushing       atomic.Bool                 // this flag symbols if engine is in full-flushing process
	isWatermarkFlushing  atomic.Bool                 // this flag symbols if engine is in water-mark flushing
}

// NewEngine creates an engine for manipulating the databases
func NewEngine(cfg config.TSDB) (Engine, error) {
	engine, err := newEngine(cfg)
	if err != nil {
		return nil, err
	}
	engine.run()
	return engine, nil
}

func newEngine(cfg config.TSDB) (*engine, error) {
	// create time series storage path
	if err := fileutil.MkDirIfNotExist(cfg.Dir); err != nil {
		return nil, fmt.Errorf("create time sereis storage path[%s] erorr: %s", cfg.Dir, err)
	}
	e := &engine{
		cfg:                  cfg,
		shardToFlushCh:       make(chan Shard),
		databaseToFlushCh:    make(chan Database),
		isFullFlushing:       *atomic.NewBool(false),
		isWatermarkFlushing:  *atomic.NewBool(false),
		memoryStatGetterFunc: monitoring.GetMemoryStat,
	}
	if err := e.load(); err != nil {
		// close opened engine
		e.Close()
	}
	return e, nil
}

// run spawns the flusher of engine.
func (e *engine) run() {
	e.ctx, e.cancel = context.WithCancel(context.Background())
	for i := 0; i < constants.FlushConcurrency; i++ {
		go e.flushWorker(e.ctx)
	}
	go e.globalMemoryUsageChecker(e.ctx)
	go e.shardMemoryUsageChecker(e.ctx)
	go e.databaseMetaFlusher(e.ctx)
}

func (e *engine) CreateDatabase(databaseName string) (Database, error) {
	dbPath := filepath.Join(e.cfg.Dir, databaseName)
	if err := fileutil.MkDirIfNotExist(dbPath); err != nil {
		return nil, fmt.Errorf("create database[%s]'s path with error: %s", databaseName, err)
	}
	cfgPath := optionsPath(dbPath)
	cfg := &databaseConfig{}
	if fileutil.Exist(cfgPath) {
		if err := ltoml.DecodeToml(cfgPath, cfg); err != nil {
			return nil, fmt.Errorf("load database[%s] config from file[%s] with error: %s",
				databaseName, cfgPath, err)
		}
	}
	db, err := newDatabase(databaseName, dbPath, cfg)
	if err != nil {
		return nil, err
	}
	e.databases.Store(databaseName, db)
	return db, nil
}

func (e *engine) GetDatabase(databaseName string) (Database, bool) {
	item, _ := e.databases.Load(databaseName)
	db, ok := item.(Database)
	return db, ok
}

func (e *engine) Close() {
	e.isFullFlushing.Store(true)
	e.cancel()
	e.databases.Range(func(key, value interface{}) bool {
		db := value.(Database)
		if err := db.Close(); err != nil {
			engineLogger.Error("close database", logger.Error(err))
		}
		return true
	})
}

func (e *engine) FlushDatabase(ctx context.Context, name string) bool {
	item, ok := e.databases.Load(name)
	if !ok {
		return false
	}
	return e.flushDatabase(ctx, item.(Database))
}

func (e *engine) flushDatabase(ctx context.Context, db Database) bool {
	select {
	case <-ctx.Done():
		return false
	case e.databaseToFlushCh <- db:
	}
	// iterate shards
	db.Range(func(key, value interface{}) bool {
		theShard := value.(Shard)
		select {
		case <-ctx.Done():
			return false
		case e.shardToFlushCh <- theShard:
		}
		return true
	})
	return true
}

// load loads the time series engines if exist
func (e *engine) load() error {
	databaseNames, err := fileutil.ListDir(e.cfg.Dir)
	if err != nil {
		return err
	}
	for _, databaseName := range databaseNames {
		_, err := e.CreateDatabase(databaseName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *engine) FlushAll() {
	if e.isFullFlushing.CAS(false, true) {
		e.flushAllDatabasesAndShards(e.ctx)
	} else {
		return
	}
	_ = e.isFullFlushing.CAS(true, false)
}

func (e *engine) globalMemoryUsageChecker(ctx context.Context) {
	ticker := time.NewTicker(globalMemoryUsageCheckInterval.Load())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// memory is lower than the high-watermark
			stat, _ := e.memoryStatGetterFunc()
			if stat.UsedPercent < constants.MemoryHighWaterMark {
				continue
			}
			// restrict watermarkFlusher concurrency thread-safe
			if e.isWatermarkFlushing.CAS(false, true) {
				go e.watermarkFlusher(ctx)
			}
		}
	}
}

func (e *engine) watermarkFlusher(ctx context.Context) {
	// if watermarkFlusher cancels, marks the flag to false
	defer e.isWatermarkFlushing.Store(false)
	// sleep interval between flushing last shard
	const sleepInterval = time.Millisecond * 50
	timer := time.NewTimer(sleepInterval)
	defer timer.Stop()

	for {
		select {
		// cancel-case1
		case <-ctx.Done():
			return
		default:
			// cancel-case2: memory is lower than MemoryLowWaterMark
			stat, _ := e.memoryStatGetterFunc()
			if stat.UsedPercent < constants.MemoryLowWaterMark {
				return
			}
			// prevent entering dead loop
			select {
			case <-timer.C:
				e.flushBiggestMemoryUsageShard(ctx)
				timer.Reset(sleepInterval)
			case <-ctx.Done():
				return
			}
		}
	}
}

func (e *engine) databaseMetaFlusher(ctx context.Context) {
	ticker := time.NewTicker(flushMetaInterval.Load())
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return
	case <-ticker.C:
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.flushAllDatabases(ctx)
		}
	}
}

func (e *engine) shardMemoryUsageChecker(ctx context.Context) {
	ticker := time.NewTicker(shardMemoryUsageCheckInterval.Load())
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return
	case <-ticker.C:
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.flushShardAboveMemoryUsageThreshold(ctx)
		}
	}
}

func (e *engine) flushBiggestMemoryUsageShard(ctx context.Context) {
	var (
		biggestShard   Shard
		biggestMemSize int
	)
	// iterate databases;
	e.databases.Range(func(key, value interface{}) bool {
		db := value.(Database)
		// iterate shards
		db.Range(func(key, value interface{}) bool {
			theShard := value.(Shard)
			// skip shard in flushing
			if theShard.IsFlushing() {
				return true
			}
			theShardSize := theShard.MemoryDatabase().MemSize()
			if theShardSize > biggestMemSize {
				biggestMemSize = theShardSize
				biggestShard = theShard
			}
			return true
		})
		return true
	})
	if biggestMemSize == 0 {
		return
	}
	// engine is already in full flushing process
	if e.isFullFlushing.Load() {
		return
	}
	select {
	case <-ctx.Done():
		return
	case e.shardToFlushCh <- biggestShard:
	}
}

func (e *engine) flushShardAboveMemoryUsageThreshold(ctx context.Context) {
	// iterate databases;
	e.databases.Range(func(key, value interface{}) bool {
		db := value.(Database)
		// iterate shards
		db.Range(func(key, value interface{}) bool {
			theShard := value.(Shard)
			if e.isFullFlushing.Load() {
				return false
			}
			if theShard.MemoryDatabase().MemSize() > constants.ShardMemoryUsedThreshold {
				select {
				case <-ctx.Done():
					return false
				case e.shardToFlushCh <- theShard:
				}
			}
			return true
		})
		return !e.isFullFlushing.Load()
	})
}

func (e *engine) flushAllDatabasesAndShards(ctx context.Context) {
	// iterate databases
	e.databases.Range(func(key, value interface{}) bool {
		return e.flushDatabase(ctx, value.(Database))
	})
}

func (e *engine) flushAllDatabases(ctx context.Context) {
	// iterate databases;
	e.databases.Range(func(key, value interface{}) bool {
		db := value.(Database)
		// in full flushing, break
		if e.isFullFlushing.Load() {
			return false
		}
		select {
		case <-ctx.Done():
			return false
		case e.databaseToFlushCh <- db:
		}
		return true
	})
}

func (e *engine) flushWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case theShard := <-e.shardToFlushCh:
			if err := theShard.Flush(); err != nil {
				engineLogger.Error("flush shard with error", logger.Error(err))
			}
		case theDatabase := <-e.databaseToFlushCh:
			if err := theDatabase.FlushMeta(); err != nil {
				engineLogger.Error("flush database metadata with error", logger.Error(err))
			}
		}
	}
}
