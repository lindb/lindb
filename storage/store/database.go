package store

import (
	"io"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
)

// Database represents abstract database for log/metric/trace etc.
type Database interface {
	io.Closer

	// Name returns time series database's name
	Name() string
	// NumOfShards returns number of families in time series database
	NumOfShards() int
	// GetConfig return the configuration of database.
	GetConfig() *models.DatabaseConfig
	// GetOption returns the database options
	GetOption() *option.DatabaseOption
	// CreateShards creates families for data partition
	CreateShards(shardIDs []models.ShardID) error
	// GetShard returns shard by given shard id
	GetShard(shardID models.ShardID) (Shard, bool)

	// Flush flushes memory data of all families to disk
	Flush() error
	// Drop drops current database include all data.
	Drop() error
	// TTL expires the data of each shard base on time to live.
	TTL()
}
