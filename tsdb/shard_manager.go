package tsdb

import "sync"

var (
	sManager          ShardManager
	once4ShardManager sync.Once
)

// GetShardManager returns the shard manager singleton instance
func GetShardManager() ShardManager {
	once4ShardManager.Do(func() {
		sManager = newShardManager()
	})
	return sManager
}

// ShardManager represents the shard manager
type ShardManager interface {
	// AddShard adds the shard
	AddShard(shard Shard)
	// RemoveShard removes the shard
	RemoveShard(shard Shard)
	// WalkEntry walks each shard entry via fn.
	WalkEntry(fn func(shard Shard))
}

// shardManager implements ShardManager interface
type shardManager struct {
	shards sync.Map
}

// newStorageManager creates the storage manager
func newShardManager() ShardManager {
	return &shardManager{}
}

// AddShard adds the shard
func (sm *shardManager) AddShard(shard Shard) {
	sm.shards.Store(shard.ShardInfo(), shard)
}

// RemoveShard adds the shard
func (sm *shardManager) RemoveShard(shard Shard) {
	sm.shards.Delete(shard.ShardInfo())
}

// WalkEntry walks each shard entry via fn.
func (sm *shardManager) WalkEntry(fn func(shard Shard)) {
	sm.shards.Range(func(key, value interface{}) bool {
		shard := value.(Shard)
		fn(shard)
		return true
	})
}
