package models

import (
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/option"
)

// CreateShardTask represents the create shard task param
type CreateShardTask struct {
	DatabaseName   string                `json:"databaseName"`   // database's name
	ShardIDs       []int32               `json:"shardIDs"`       // shard ids
	DatabaseOption option.DatabaseOption `json:"databaseOption"` // time series database
}

// Bytes returns the create shard task binary data using json
func (t CreateShardTask) Bytes() []byte {
	return encoding.JSONMarshal(t)
}
