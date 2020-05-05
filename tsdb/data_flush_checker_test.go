package tsdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shirou/gopsutil/mem"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/tsdb/memdb"
)

func TestDataFlushChecker_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		memoryUsageCheckInterval.Store(time.Second)
		ctrl.Finish()
	}()
	shard := NewMockShard(ctrl)
	shard.EXPECT().NeedFlush().Return(true).AnyTimes()
	shard.EXPECT().ShardInfo().Return("shardInfo").AnyTimes()
	shard.EXPECT().Flush().Return(fmt.Errorf("err")).AnyTimes()
	GetShardManager().AddShard(shard)
	memoryUsageCheckInterval.Store(10 * time.Millisecond)
	checker := newDataFlushChecker(context.TODO())
	checker.Start()

	time.Sleep(100 * time.Millisecond)
	checker.Stop()
	GetShardManager().RemoveShard(shard)
}

func TestDataFlushChecker_check_high_memory_waterMark(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		memoryUsageCheckInterval.Store(time.Second)
		ctrl.Finish()
	}()
	// case 1: shard is flushing
	shard := NewMockShard(ctrl)
	shard.EXPECT().NeedFlush().Return(false).AnyTimes()
	shard.EXPECT().ShardInfo().Return("shardInfo").AnyTimes()
	shard.EXPECT().IsFlushing().Return(true).AnyTimes()
	GetShardManager().AddShard(shard)
	memoryUsageCheckInterval.Store(10 * time.Millisecond)
	checker := newDataFlushChecker(context.TODO())
	check := checker.(*dataFlushChecker)
	check.memoryStatGetterFunc = func() (stat *mem.VirtualMemoryStat, err error) {
		return &mem.VirtualMemoryStat{UsedPercent: constants.MemoryHighWaterMark + 0.1}, nil
	}
	checker.Start()

	time.Sleep(100 * time.Millisecond)
	checker.Stop()
	GetShardManager().RemoveShard(shard)

	// case 2: pick biggest shard data
	shard1 := NewMockShard(ctrl)
	shard1.EXPECT().NeedFlush().Return(false).AnyTimes()
	shard1.EXPECT().ShardInfo().Return("shardInfo").AnyTimes()
	shard1.EXPECT().IsFlushing().Return(false).AnyTimes()
	mDB1 := memdb.NewMockMemoryDatabase(ctrl)
	mDB1.EXPECT().MemSize().Return(int32(100)).AnyTimes()
	shard1.EXPECT().MemoryDatabase().Return(mDB1).AnyTimes()
	GetShardManager().AddShard(shard1)

	shard2 := NewMockShard(ctrl)
	shard2.EXPECT().NeedFlush().Return(false).AnyTimes()
	shard2.EXPECT().ShardInfo().Return("shardInfo").AnyTimes()
	shard2.EXPECT().IsFlushing().Return(false).AnyTimes()
	shard2.EXPECT().Flush().Return(nil).AnyTimes()
	mDB2 := memdb.NewMockMemoryDatabase(ctrl)
	mDB2.EXPECT().MemSize().Return(int32(1000)).AnyTimes()
	shard2.EXPECT().MemoryDatabase().Return(mDB2).AnyTimes()
	GetShardManager().AddShard(shard2)

	memoryUsageCheckInterval.Store(10 * time.Millisecond)
	checker = newDataFlushChecker(context.TODO())
	check = checker.(*dataFlushChecker)
	check.memoryStatGetterFunc = func() (stat *mem.VirtualMemoryStat, err error) {
		return &mem.VirtualMemoryStat{UsedPercent: constants.MemoryHighWaterMark + 0.1}, nil
	}
	checker.Start()

	time.Sleep(100 * time.Millisecond)
	checker.Stop()
	GetShardManager().RemoveShard(shard1)
	GetShardManager().RemoveShard(shard2)
}

func TestDataFlushChecker_requestFlush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		memoryUsageCheckInterval.Store(time.Second)
		ctrl.Finish()
	}()
	var shards []Shard
	for i := 0; i < 2; i++ {
		shard := NewMockShard(ctrl)
		shard.EXPECT().NeedFlush().Return(true).AnyTimes()
		shard.EXPECT().ShardInfo().Return(fmt.Sprintf("shard-%d", i)).AnyTimes()
		shard.EXPECT().Flush().DoAndReturn(func() error {
			time.Sleep(200 * time.Millisecond)
			return fmt.Errorf("err")
		}).AnyTimes()
		GetShardManager().AddShard(shard)
		shards = append(shards, shard)
	}
	memoryUsageCheckInterval.Store(10 * time.Millisecond)
	checker := newDataFlushChecker(context.TODO())
	check := checker.(*dataFlushChecker)
	check.memoryStatGetterFunc = func() (stat *mem.VirtualMemoryStat, err error) {
		return &mem.VirtualMemoryStat{UsedPercent: constants.MemoryHighWaterMark + 0.1}, nil
	}
	checker.Start()

	time.Sleep(100 * time.Millisecond)
	checker.Stop()
	time.Sleep(100 * time.Millisecond)

	for _, shard := range shards {
		GetShardManager().AddShard(shard)
	}
}
