package log

import (
	"fmt"
	"path"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/flush"
	"github.com/lindb/lindb/storage/log/index"
	"github.com/lindb/lindb/storage/store"
)

func init() {
	store.RegisterEngine(option.Log, NewDatabase)
}

type Database struct {
	base.Database

	dir string

	indexDB index.Database

	logger logger.Logger
}

func NewDatabase(name string, opt *models.DatabaseConfig,
	limit *models.Limits, checker flush.Checker,
) (store.Database, error) {
	dbPath, err := store.CreateDatabasePath(name)
	if err != nil {
		return nil, err
	}
	db := &Database{
		dir:    dbPath,
		logger: logger.GetLogger("Log", "Database"),
	}
	db.Database = base.Database{
		Options:       opt,
		DatabaseName:  name,
		ShardSet:      *store.NewShardSet(),
		CreateShardFn: db.createShard,
	}

	if err = db.DumpOption(); err != nil {
		return nil, err
	}
	db.indexDB = index.NewDatabase(path.Join(dbPath, "meta"))

	var shard store.Shard
	if len(db.Options.ShardIDs) > 0 {
		for _, shardID := range db.Options.ShardIDs {
			shard, err = db.createShard(shardID)
			if err != nil {
				return nil, fmt.Errorf("cannot create shard[%d] of database[%s] with error: %s",
					shardID, name, err)
			}
			db.ShardSet.InsertShard(shardID, shard)
		}
	}
	return db, nil
}

func (db *Database) IndexDatabase() index.Database {
	return db.indexDB
}

func (db *Database) FindMatchSmallestInterval(interval timeutil.Interval) timeutil.Interval {
	return minuteInterval
}

func (db *Database) createShard(shardID models.ShardID) (store.Shard, error) {
	return NewShard(shardID, db)
}

// Drop implements store.Database.
func (db *Database) Drop() error {
	panic("unimplemented")
}

// Flush implements store.Database.
func (db *Database) Flush() error {
	panic("unimplemented")
}

// TTL implements store.Database.
func (db *Database) TTL() {
	panic("unimplemented")
}

// Close implements store.Database.
func (db *Database) Close() error {
	db.indexDB.Flush()
	db.indexDB.Close()

	for _, shardEntry := range db.ShardSet.Entries() {
		if err := shardEntry.Shard.Close(); err != nil {
			db.logger.Error(fmt.Sprintf(
				"close shard[%d] of database[%s]", shardEntry.ShardID, db.DatabaseName), logger.Error(err))
		}
	}
	return nil
}
