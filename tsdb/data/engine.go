package data

import (
	"fmt"

	"github.com/eleme/lindb/pkg/option"
)

type Engine struct {
	name   string
	shards map[int32]*Shard
}

func NewEngine(name string) *Engine {
	return &Engine{
		name:   name,
		shards: make(map[int32]*Shard),
	}
}

func (e *Engine) CreateShard(shardID int32, option option.ShardOption) error {
	var shard, ok = e.shards[shardID]
	if ok {
		return fmt.Errorf("engine[%s] exist shard[%d]", e.name, shardID)
	}
	shard = NewShard(shardID, option)
	e.shards[shardID] = shard
	return nil
}
