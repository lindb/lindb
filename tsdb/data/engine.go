package data

import (
	"fmt"

	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/tsdb"
)
// engine 
type engine struct {
	name   string
	shards map[int32]*Shard
}

// NewEngine creates
func NewEngine(name string) tsdb.Engine {
	return &engine{
		name:   name,
		shards: make(map[int32]*Shard),
	}
}

func (e *engine) CreateShard(shardID int32, option option.ShardOption) error {
	var shard, ok = e.shards[shardID]
	if ok {
		return fmt.Errorf("engine[%s] exist shard[%d]", e.name, shardID)
	}
	shard = NewShard(shardID, option)
	e.shards[shardID] = shard
	return nil
}
