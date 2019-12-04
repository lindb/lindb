package models

import (
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/option"
)

// CreateShardTask represents the create shard task's param
type CreateShardTask struct {
	DatabaseName   string                `json:"databaseName"`   // database's name
	ShardIDs       []int32               `json:"shardIDs"`       // shard ids
	DatabaseOption option.DatabaseOption `json:"databaseOption"` // time series database
}

// Bytes returns the create shard task's  binary data using json
func (t CreateShardTask) Bytes() []byte {
	return encoding.JSONMarshal(t)
}

// DatabaseFlushTask represents the database flush task's param
type DatabaseFlushTask struct {
	DatabaseName string `json:"databaseName"` // database's name
}

// Bytes returns the database flush task's binary data using json
func (t DatabaseFlushTask) Bytes() []byte {
	return encoding.JSONMarshal(t)
}
