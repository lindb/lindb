package tsdb

import (
	"github.com/eleme/lindb/pkg/option"
)

// Engine that time series storage engine
type Engine interface {
	// Create shard for data partition
	CreateShard(shardID int32, option option.ShardOption) error
}
