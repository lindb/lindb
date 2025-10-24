package store

import (
	"fmt"
	"io"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/flush"
)

var creatingFns = make(map[option.EngineType]CreateDatabaseFn)

// CreateDatabaseFn represents create database function.
type CreateDatabaseFn func(
	databaseName string,
	cfg *models.DatabaseConfig,
	limits *models.Limits,
	flushChecker flush.Checker,
) (Database, error)

// RegisterEngine registers a database creating function.
func RegisterEngine(engine option.EngineType, fn CreateDatabaseFn) {
	creatingFns[engine] = fn
}

// CreateDatabase represents create database.
func CreateDatabase(
	databaseName string,
	cfg *models.DatabaseConfig,
	limits *models.Limits,
	checker flush.Checker,
) (Database, error) {
	fn, ok := creatingFns[cfg.Option.Engine]
	if !ok {
		return nil, fmt.Errorf("not support engine: %s", cfg.Option.Engine)
	}
	return fn(databaseName, cfg, limits, checker)
}

// Database represents abstract database for log/metric/trace etc.
type Database interface {
	io.Closer

	// Name returns time series database's name
	Name() string
	// NumOfShards returns number of families in time series database
	NumOfShards() int
	// GetOption returns the database options
	GetOption() *models.DatabaseConfig
	FindMatchSmallestInterval(interval timeutil.Interval) timeutil.Interval
	// CreateShards creates families for data partition
	CreateShards(shardIDs []models.ShardID) error
	// GetShard returns shard by given shard id
	GetShard(shardID models.ShardID) (Shard, bool)
	EvictSegment()
	SetLimits(limits *models.Limits)
	GetLimits() *models.Limits
	// Flush flushes memory data of all families to disk
	Flush() error
	// Drop drops current database include all data.
	Drop() error
	// TTL expires the data of each shard base on time to live.
	TTL()
}
