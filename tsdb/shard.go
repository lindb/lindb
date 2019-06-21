package tsdb

import (
	"context"
	"time"

	"github.com/eleme/lindb/tsdb/memdb"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/option"
)

// Shard is a horizontal partition of metrics for LinDB.
type Shard interface {
}

type shard struct {
	id     int32
	option option.ShardOption
	memDb  memdb.MemoryDatabase
	cancel func()
}

func newShard(shardID int32, option option.ShardOption) Shard {
	ctx, cancel := context.WithCancel(context.Background())
	shard := &shard{
		id:     shardID,
		option: option,
		memDb:  memdb.NewMemoryDatabase(ctx),
		cancel: cancel}
	return shard
}

// Close closes the memDatabase and spawned goroutines.
func (s *shard) Close() {
	s.cancel()
}

// Write writes the metric-point into memory-database.
func (s *shard) Write(point models.Point) error {
	timestamp := point.Timestamp()
	now := time.Now().Unix()

	if timestamp < now-s.option.Behind {
		return nil
	}

	if timestamp > now+s.option.Ahead {
		return nil
	}
	//use family base time for memory store
	segmentTime := s.option.IntervalType.CalFamilyBaseTime(point.Timestamp())
	//slot time for ts data store
	slotTime := s.option.IntervalType.CalSlot(point.Timestamp())

	return s.memDb.Write(point, segmentTime, slotTime)
}
