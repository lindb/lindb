package tsdb

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/util"
)

const options = "OPTIONS"

// Engine represents a time series storage engine
type Engine interface {
	// CreateShards creates shards for data partition
	CreateShards(option option.ShardOption, shardIDs ...int32) error
	// GetShard returns shard by given shard id, if not exist returns nil
	GetShard(shardID int32) Shard
	// Close closed engine then release resource
	Close() error
}

// info represents a engine information about config and shards
type info struct {
	ShardIDs    []int32            `toml:"shardIds"`
	ShardOption option.ShardOption `toml:"shardOption"`
}

// engine implements Engine for storing shards, each shard represents a time series storage
type engine struct {
	name   string
	path   string
	shards map[int32]Shard
	info   *info

	mutex sync.RWMutex
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
		name:   name,
		path:   enginePath,
		shards: make(map[int32]Shard),
		info:   info,
	}
	// load shards if engine is exist
	if len(e.info.ShardIDs) > 0 {
		for _, shardID := range e.info.ShardIDs {
			shard := newShard(shardID, info.ShardOption)
			e.shards[shardID] = shard
		}
	}
	return e, nil
}

// CreateShards creates shards for data partition
func (e *engine) CreateShards(option option.ShardOption, shardIDs ...int32) error {
	if len(shardIDs) == 0 {
		return fmt.Errorf("shard is list is empty")
	}
	for _, shardID := range shardIDs {
		e.mutex.RLock()
		_, ok := e.shards[shardID]
		e.mutex.RUnlock()

		if !ok {
			// be careful need do mutex unlock
			e.mutex.Lock()
			_, ok = e.shards[shardID]
			if !ok {
				// using new shard option
				newInfo := &info{ShardOption: option, ShardIDs: e.info.ShardIDs}
				// add new shard id
				newInfo.ShardIDs = append(newInfo.ShardIDs, shardID)
				if err := e.dumpEningeInfo(newInfo); err != nil {
					e.mutex.Unlock()
					return err
				}

				shard := newShard(shardID, option)
				e.shards[shardID] = shard
				e.mutex.Unlock()
			}
		}
	}
	return nil
}

// GetShard returns shard by given shard id, if not exist returns nil
func (e *engine) GetShard(shardID int32) Shard {
	e.mutex.RLock()
	shard := e.shards[shardID]
	e.mutex.RUnlock()
	return shard
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
