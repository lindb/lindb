package log

import (
	"fmt"
	"path"

	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/log/index"
	"github.com/lindb/lindb/storage/store"
	"github.com/lindb/lindb/storage/utils"
)

type Database struct {
	base.Database

	dir string

	indexDB index.Database
}

func NewDatabase(name string, opt *models.DatabaseConfig) (store.Database, error) {
	dbPath, err := utils.CreateDatabasePath(name)
	if err != nil {
		return nil, err
	}
	db := &Database{
		dir: dbPath,
	}
	db.Database = base.Database{
		Options:       opt,
		DatabaseName:  name,
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
			db.Shards.Store(shardID, shard)
		}
	}
	return db, nil
}

func (db *Database) IndexDatabase() index.Database {
	return db.indexDB
}

func (db *Database) createShard(shardID models.ShardID) (store.Shard, error) {
	shard, err := NewShard(shardID, db)
	if err == nil {
		shard.GetOrCreatePartition(timeutil.Now())
	}
	return shard, err
}

// Drop implements store.Database.
func (db *Database) Drop() error {
	panic("unimplemented")
}

// Flush implements store.Database.
func (db *Database) Flush() error {
	panic("unimplemented")
}

// GetConfig implements store.Database.
func (db *Database) GetConfig() *models.DatabaseConfig {
	panic("unimplemented")
}

// GetOption implements store.Database.
func (db *Database) GetOption() *option.DatabaseOption {
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

	db.Shards.Range(func(key, value any) bool {
		value.(store.Shard).Close()
		return true
	})
	return nil
}
