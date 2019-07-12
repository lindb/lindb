package models

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/option"
)

// CreateShardTask represents create shard task param
type CreateShardTask struct {
	Database    string             `json:"database"`
	ShardIDs    []int              `json:"shardIDs"`
	ShardOption option.ShardOption `json:"shardOption"`
}

// Bytes returns create shard task binary data using json
func (t CreateShardTask) Bytes() []byte {
	data, err := json.Marshal(t)
	if err != nil {
		logger.GetLogger().Error("marshal create shard task error", zap.Error(err))
		return nil
	}
	return data
}
