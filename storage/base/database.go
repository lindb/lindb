package base

import (
	"fmt"
	"sync"

	"github.com/lindb/common/pkg/ltoml"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/storage/store"
)

type CreateShardFn func(shardID models.ShardID) (store.Shard, error)

type Database struct {
	DatabaseName string
	Options      *models.DatabaseConfig

	ShardSet      store.ShardSet // atomic value
	CreateShardFn CreateShardFn

	mutex sync.Mutex
}

func (db *Database) DumpOption() error {
	optionsPath := store.OptionsPath(db.DatabaseName)
	fmt.Println(optionsPath)
	// write options using toml format
	if err := ltoml.EncodeToml(optionsPath, db.Options); err != nil {
		return fmt.Errorf("dump database options to file[%s] error:%s", optionsPath, err)
	}
	return nil
}

// GetShard returns shard by given shard id,
func (db *Database) GetShard(shardID models.ShardID) (store.Shard, bool) {
	return db.ShardSet.GetShard(shardID)
}

func (db *Database) CreateShards(shardIDs []models.ShardID) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	for _, shardID := range shardIDs {
		// double check
		if _, ok := db.GetShard(shardID); ok {
			return nil
		}
		// new shard
		createdShard, err := db.CreateShardFn(shardID)
		if err != nil {
			return fmt.Errorf("create shard[%d] for database[%s] with error: %s", shardID, db.DatabaseName, err)
		}
		db.ShardSet.InsertShard(shardID, createdShard)
	}

	// using new engine option
	db.Options.ShardIDs = append(db.Options.ShardIDs, shardIDs...)
	if err := db.DumpOption(); err != nil {
		// TODO: if dump config err, need close shard??
		return err
	}

	return nil
}

func (db *Database) Name() string {
	return db.DatabaseName
}

func (db *Database) GetOption() *models.DatabaseConfig {
	return db.Options
}

// SetLimits sets database's limits.
func (db *Database) SetLimits(limits *models.Limits) {
	models.SetDatabaseLimits(db.DatabaseName, limits)
}

// GetLimits returns database's limits.
func (db *Database) GetLimits() *models.Limits {
	return models.GetDatabaseLimits(db.DatabaseName)
}

// NumOfShards implements store.Database.
func (db *Database) NumOfShards() int {
	return db.ShardSet.GetShardNum()
}

func (db *Database) EvictSegment() {
	panic("need implements")
}
