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

func (e *Engine) CreateShard(shardId int32, option option.ShardOption) error {
	var shard, ok = e.shards[shardId]
	if ok {
		return fmt.Errorf("engine[%s] exist shard[%d]", e.name, shardId)
	}
	shard = NewShard(shardId, option)
	e.shards[shardId] = shard
	return nil
}
