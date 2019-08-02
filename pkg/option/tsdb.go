package option

import (
	"time"

	"github.com/lindb/lindb/pkg/interval"
)

// EngineOption represents a engine option include shard ids and shard's option
type EngineOption struct {
	ShardOption ShardOption `toml:"shardOption"`
	ShardIDs    []int32     `toml:"shardIDs"`
}

// ShardOption represents a shard storage configuration
type ShardOption struct {
	TimeWindow   int           `toml:"timeWindow" json:"timeWindow"`     // time window of memory database block
	Behind       int64         `toml:"behind" json:"behind"`             // allowed timestamp write behind
	Ahead        int64         `toml:"ahead" json:"ahead"`               // allowed timestamp write ahead
	Interval     time.Duration `toml:"interval" json:"interval"`         // interval duration
	IntervalType interval.Type `toml:"intervalType" json:"intervalType"` // interval type
}
