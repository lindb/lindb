package tsdb

import (
	"github.com/eleme/lindb/pkg/option"
)

// Engine represents a time series storage engine
type Engine interface {
	// CreateShard creates shard for data partition
	CreateShard(shardID int32, option option.ShardOption) error
}
