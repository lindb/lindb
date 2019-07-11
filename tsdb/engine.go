package tsdb

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/util"
)

const options = "OPTIONS"
const shardPath = "shard"

//go:generate mockgen -source ./engine.go -destination=./engine_mock.go -package tsdb

// Engine represents a time series storage engine
type Engine interface {
	// Name returns tsdb engine's name, engine's name is database's name for user
	Name() string
	// NumOfShards returns number of shards in tsdb engine
	NumOfShards() int
	// CreateShards creates shards for data partition
	CreateShards(option option.ShardOption, shardIDs ...int) error
	// GetShard returns shard by given shard id, if not exist returns nil
	GetShard(shardID int) Shard
	// Close closed engine then release resource
	Close() error
}

// info represents a engine information about config and shards
type info struct {
	ShardIDs    []int              `toml:"shardIds"`
	ShardOption option.ShardOption `toml:"shardOption"`
}

// engine implements Engine for storing shards, each shard represents a time series storage
type engine struct {
	name   string
	path   string
	shards sync.Map
	info   *info

	numOfShards int

	mutex sync.Mutex
}

// NewEngine creates engine instance if create engine's path successfully
func NewEngine(name string, path string) (Engine, error) {
	enginePath := filepath.Join(path, name)
	// create engine path
	if err := util.MkDirIfNotExist(enginePath); err != nil {
		return nil, fmt.Errorf("create path of tsdb engine[%s] erorr: %s", name, err)
	}
	infoPath := infoPath(enginePath)
	info := &info{}
	if util.Exist(infoPath) {
		if err := util.DecodeToml(infoPath, info); err != nil {
			return nil, fmt.Errorf("load engine option from file[%s] error:%s", infoPath, err)
		}
	}
	e := &engine{
		name: name,
		path: enginePath,
		info: info,
	}
	// load shards if engine is exist
	if len(e.info.ShardIDs) > 0 {
		for _, shardID := range e.info.ShardIDs {
			shard, err := newShard(shardID, filepath.Join(enginePath, shardPath, string(shardID)), info.ShardOption)
			if err != nil {
				return nil, fmt.Errorf("cannot create shard[%d] for engine[%s] error:%s", shardID, name, err)
			}
			e.shards.Store(shardID, shard)
			e.numOfShards++
		}
	}
	return e, nil
}

// Name returns tsdb engine's name, engine's name is database's name for user
func (e *engine) Name() string {
	return e.name
}

// NumOfShards returns number of shards in tsdb engine
func (e *engine) NumOfShards() int {
	return e.numOfShards
}

// CreateShards creates shards for data partition
func (e *engine) CreateShards(option option.ShardOption, shardIDs ...int) error {
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
				shard, err := newShard(shardID, filepath.Join(e.path, shardPath, fmt.Sprintf("%d", shardID)), option)
				if err != nil {
					e.mutex.Unlock()
					return fmt.Errorf("cannot create shard[%d] for engine[%s] error:%s", shardID, e.name, err)
				}
				// using new shard option
				newInfo := &info{ShardOption: option, ShardIDs: e.info.ShardIDs}
				// add new shard id
				newInfo.ShardIDs = append(newInfo.ShardIDs, shardID)
				if err := e.dumpEningeInfo(newInfo); err != nil {
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

// GetShard returns shard by given shard id, if not exist returns nil
func (e *engine) GetShard(shardID int) Shard {
	shard, _ := e.shards.Load(shardID)
	s, ok := shard.(Shard)
	if ok {
		return s
	}
	return nil
}

// Close closed engine then release resource
func (e *engine) Close() error {
	//TODO impl close logic
	return nil
}

// dumpEningeInfo persists option info to OPTIONS file
func (e *engine) dumpEningeInfo(newInfo *info) error {
	infoPath := infoPath(e.path)
	// write store info using toml format
	if err := util.EncodeToml(infoPath, newInfo); err != nil {
		return fmt.Errorf("write engine info to file[%s] error:%s", infoPath, err)
	}
	e.info = newInfo
	return nil
}

// infoPath returns options file path
func infoPath(path string) string {
	return filepath.Join(path, options)
}
