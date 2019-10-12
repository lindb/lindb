package tsdb

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/concurrent"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb/diskdb"
)

//go:generate mockgen -source=./engine.go -destination=./engine_mock.go -package=tsdb

const options = "OPTIONS"
const shardPath = "shard"

var log = logger.GetLogger("tsdb", "Engine")

// EngineFactory represents a time series engine create factory
type EngineFactory interface {
	// CreateEngine creates engine instance if create engine's path successfully
	CreateEngine(name string) (Engine, error)
	// GetEngine returns the time series engine by given name
	GetEngine(db string) Engine
	// Close closed the cached time series engines
	Close()
}

// Engine represents a time series storage engine
type Engine interface {
	// Name returns time series engine's name, engine's name is database's name for user
	Name() string
	// NumOfShards returns number of shards in time series engine
	NumOfShards() int
	// CreateShards creates shards for data partition
	CreateShards(option option.EngineOption, shardIDs ...int32) error
	// GetShard returns shard by given shard id, if not exist returns nil
	GetShard(shardID int32) Shard
	// GetIDGetter returns id getter for metric level metadata
	GetIDGetter() diskdb.IDGetter
	// GetExecutePool returns the query task execute pool
	GetExecutePool() *ExecutePool
	// Close closed engine then release resource
	Close() error

	// getIndex returns the index
	getIndex() Index
}

// info represents a engine information about config and shards
type info struct {
	ShardIDs []int32             `toml:"shardIds"`
	Engine   option.EngineOption `toml:"engine"`
}

// engine implements Engine for storing shards, each shard represents a time series storage
type engine struct {
	name   string
	path   string
	shards sync.Map
	info   *info
	index  Index

	executePool *ExecutePool // execute query task

	numOfShards int

	mutex sync.Mutex
}

// engineFactory implements engine factory interface
type engineFactory struct {
	cfg config.Engine // the common cfg of time series engine

	engines sync.Map
}

// NewEngineFactory creates an engine factory for creating time series engine
func NewEngineFactory(cfg config.Engine) (EngineFactory, error) {
	// create time series storage path
	if err := fileutil.MkDirIfNotExist(cfg.Dir); err != nil {
		return nil, fmt.Errorf("create time sereis storage path[%s] erorr: %s", cfg.Dir, err)
	}

	f := &engineFactory{cfg: cfg}
	if err := f.load(); err != nil {
		// close opened engine
		f.Close()
		return nil, err
	}
	return f, nil
}

// CreateEngine creates an engine instance if create engine's path successfully
func (f *engineFactory) CreateEngine(name string) (Engine, error) {
	enginePath := filepath.Join(f.cfg.Dir, name)
	// create engine path
	if err := fileutil.MkDirIfNotExist(enginePath); err != nil {
		return nil, fmt.Errorf("create path of tsdb engine[%s] erorr: %s", name, err)
	}
	infoPath := infoPath(enginePath)
	info := &info{}
	if fileutil.Exist(infoPath) {
		if err := fileutil.DecodeToml(infoPath, info); err != nil {
			return nil, fmt.Errorf("load engine option from file[%s] error:%s", infoPath, err)
		}
	}
	//FIXME store1100 add cfg
	index, err := newIndex(name, f.cfg)
	if err != nil {
		return nil, err
	}
	e := &engine{
		name:  name,
		path:  enginePath,
		info:  info,
		index: index,
		//TODO add pool config
		executePool: &ExecutePool{
			Scan:  concurrent.NewPool(name+"-"+"executor-pool", 100 /*nRoutines*/, 10 /*queueSize*/),
			Merge: concurrent.NewPool(name+"-"+"executor-pool", 100 /*nRoutines*/, 10 /*queueSize*/),
		},
	}
	// load shards if engine is exist
	if len(e.info.ShardIDs) > 0 {
		for _, shardID := range e.info.ShardIDs {
			shard, err := newShard(shardID, filepath.Join(enginePath, shardPath, fmt.Sprintf("%d", shardID)),
				e.getIndex(), info.Engine)
			if err != nil {
				return nil, fmt.Errorf("cannot create shard[%d] for engine[%s] error:%s", shardID, name, err)
			}
			e.shards.Store(shardID, shard)
			e.numOfShards++
		}
	}
	f.engines.Store(name, e)
	return e, nil
}

// GetEngine returns the time series engine by given db name, if not exist return nil
func (f *engineFactory) GetEngine(db string) Engine {
	engine, _ := f.engines.Load(db)
	e, ok := engine.(Engine)
	if ok {
		return e
	}
	return nil
}

// Close closed the cached time series engines when system shutdown
func (f *engineFactory) Close() {
	f.engines.Range(func(key, value interface{}) bool {
		engine, ok := value.(Engine)
		if ok {
			if err := engine.Close(); err != nil {
				log.Error("close engine", logger.Error(err))
			}
		}

		return true
	})
}

// load loads the time series engines if exist
func (f *engineFactory) load() error {
	names, err := fileutil.ListDir(f.cfg.Dir)
	if err != nil {
		return err
	}
	for _, name := range names {
		_, err := f.CreateEngine(name)
		if err != nil {
			return err
		}
	}
	return nil
}

// Name returns time series engine's name, engine's name is database's name for user
func (e *engine) Name() string {
	return e.name
}

// NumOfShards returns number of shards in time series engine
func (e *engine) NumOfShards() int {
	return e.numOfShards
}

// CreateShards creates shards for data partition
func (e *engine) CreateShards(option option.EngineOption, shardIDs ...int32) error {
	if len(shardIDs) == 0 {
		return fmt.Errorf("shard is list is empty")
	}
	for _, shardID := range shardIDs {
		shard := e.GetShard(shardID)

		if shard == nil {
			// be careful need do mutex unlock
			e.mutex.Lock()
			// double check
			shard = e.GetShard(shardID)
			if shard == nil {
				// new shard
				shard, err := newShard(shardID, filepath.Join(e.path, shardPath, fmt.Sprintf("%d", shardID)),
					e.getIndex(), option)
				if err != nil {
					e.mutex.Unlock()
					return fmt.Errorf("cannot create shard[%d] for engine[%s] error:%s", shardID, e.name, err)
				}
				// using new engine option
				newInfo := &info{Engine: option, ShardIDs: e.info.ShardIDs}
				// add new shard id
				newInfo.ShardIDs = append(newInfo.ShardIDs, shardID)
				if err := e.dumpEngineInfo(newInfo); err != nil {
					e.mutex.Unlock()
					return err
				}
				e.shards.Store(shardID, shard)
				e.numOfShards++
				e.mutex.Unlock()
			}
		}
	}
	return nil
}

// getIndex returns the index
func (e *engine) getIndex() Index {
	return e.index
}

// GetShard returns shard by given shard id, if not exist returns nil
func (e *engine) GetShard(shardID int32) Shard {
	shard, _ := e.shards.Load(shardID)
	s, ok := shard.(Shard)
	if ok {
		return s
	}
	return nil
}

// GetIDGetter returns id getter for metric level metadata
func (e *engine) GetIDGetter() diskdb.IDGetter {
	return e.index.GetIDSequencer()
}

// GetExecutePool returns the query task execute pool
func (e *engine) GetExecutePool() *ExecutePool {
	return e.executePool
}

// Close closed engine then release resource
func (e *engine) Close() error {
	e.index.Close()
	//TODO impl close logic
	return nil
}

// dumpEngineInfo persists option info to OPTIONS file
func (e *engine) dumpEngineInfo(newInfo *info) error {
	infoPath := infoPath(e.path)
	// write store info using toml format
	if err := fileutil.EncodeToml(infoPath, newInfo); err != nil {
		return fmt.Errorf("write engine info to file[%s] error:%s", infoPath, err)
	}
	e.info = newInfo
	return nil
}

// infoPath returns options file path
func infoPath(path string) string {
	return filepath.Join(path, options)
}
