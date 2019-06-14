package tsdb

import "github.com/eleme/lindb/pkg/option"

// time series storage engine
type Engine interface {
	// create shard for data partition
	CreateShard(shardId int32, option option.ShardOption) error
}
