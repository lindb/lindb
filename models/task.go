package models

import (
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/option"
)

// CreateShardTask represents the create shard task param
type CreateShardTask struct {
	Database string              `json:"database"` // database's name
	ShardIDs []int32             `json:"shardIDs"` // shard ids
	Engine   option.EngineOption `json:"engine"`   // time series engine
}

// Bytes returns the create shard task binary data using json
func (t CreateShardTask) Bytes() []byte {
	return encoding.JSONMarshal(t)
}
