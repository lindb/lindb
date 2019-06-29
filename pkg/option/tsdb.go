package option

import (
	"time"

	"github.com/eleme/lindb/pkg/interval"
)

// EngineOption represents a engine option include shard ids and shard's option
type EngineOption struct {
	ShardOption ShardOption `toml:"shardOption"`
	ShardIDs    []int32     `toml:"shardIDs"`
}

// ShardOption represents a shard storage configuration
type ShardOption struct {
	Behind       int64         `toml:"behind"`       // allowed timestamp write behind
	Ahead        int64         `toml:"ahead"`        // allowed timestamp write ahead
	Interval     time.Duration `toml:"interval"`     // interval duration
	IntervalType interval.Type `toml:"intervalType"` // interval type
}
