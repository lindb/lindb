package tsdb

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb/memdb"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

var testPath = "test_data"
var validOption = option.DatabaseOption{Interval: "10s"}
var engineCfg = config.TSDB{Dir: testPath}

func TestNew(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	e, err := NewEngine(engineCfg)
	assert.NoError(t, err)

	db, _ := e.CreateDatabase("test_db")
	assert.NotNil(t, db)
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db")))

	assert.Equal(t, 0, db.NumOfShards())

	err = db.CreateShards(option.DatabaseOption{})
	assert.NotNil(t, err)

	err = db.CreateShards(option.DatabaseOption{}, 1, 2, 3)
	assert.NotNil(t, err)

	err = db.CreateShards(validOption, 1, 2, 3)
	assert.Nil(t, err)
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db", "OPTIONS")))
	assert.Equal(t, "test_db", db.Name())

	_, ok := db.GetShard(1)
	assert.True(t, ok)
	_, ok = db.GetShard(2)
	assert.True(t, ok)
	_, ok = db.GetShard(3)
	assert.True(t, ok)
	_, ok = db.GetShard(10)
	assert.False(t, ok)
	assert.Equal(t, 3, db.NumOfShards())

	_, ok = e.GetDatabase("inexist")
	assert.False(t, ok)
	assert.NotNil(t, db.ExecutorPool())

	e.Close()

	// re-open factory
	e, err = NewEngine(engineCfg)
	assert.NoError(t, err)

	db, ok = e.GetDatabase("test_db")
	assert.True(t, ok)
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db")))
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db", "OPTIONS")))

	_, ok = db.GetShard(1)
	assert.True(t, ok)
	_, ok = db.GetShard(2)
	assert.True(t, ok)
	_, ok = db.GetShard(3)
	assert.True(t, ok)
	_, ok = db.GetShard(10)
	assert.False(t, ok)
	assert.Equal(t, 3, db.NumOfShards())
}

func Test_Engine_Close(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e, _ := NewEngine(engineCfg)
	engineImpl := e.(*engine)
	defer engineImpl.cancel()

	mockDatabase := NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Close().Return(fmt.Errorf("error")).AnyTimes()
	engineImpl.databases.Store("1", mockDatabase)
	engineImpl.databases.Store("2", mockDatabase)

	e.Close()
}

func Test_Engine_Flush_Database(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	e, _ := NewEngine(engineCfg)
	engineImpl := e.(*engine)
	defer engineImpl.cancel()
	ok := e.FlushDatabase(context.TODO(), "test_db_3")
	assert.False(t, ok)

	mockDatabase := NewMockDatabase(ctrl)
	mockDatabase.EXPECT().FlushMeta()
	mockDatabase.EXPECT().Range(gomock.Any())
	engineImpl.databases.Store("test_db_1", mockDatabase)
	ok = e.FlushDatabase(context.TODO(), "test_db_1")
	assert.True(t, ok)
}

func Test_Engine_FlushAll(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	e, _ := NewEngine(engineCfg)
	engineImpl := e.(*engine)
	defer engineImpl.cancel()

	engineImpl.isFullFlushing.Store(true)
	e.FlushAll()

	engineImpl.isFullFlushing.Store(false)
	e.FlushAll()
}

func Test_Engine_globalMemoryUsageChecker_LowerThanHighWaterMark(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	engineImpl, _ := newEngine(engineCfg)
	engineImpl.memoryStatGetterFunc = func() (*models.MemoryStat, error) {
		return &models.MemoryStat{UsedPercent: constants.MemoryHighWaterMark - 0.1}, nil
	}
	engineImpl.run()
	defer engineImpl.cancel()

	globalMemoryUsageCheckInterval.Store(time.Millisecond)

	go engineImpl.globalMemoryUsageChecker(engineImpl.ctx)
	time.Sleep(time.Second)
}

func Test_Engine_globalMemoryUsageChecker_HigherThan_MemoryHighWaterMark(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	engineImpl, _ := newEngine(engineCfg)
	engineImpl.memoryStatGetterFunc = func() (*models.MemoryStat, error) {
		return &models.MemoryStat{UsedPercent: constants.MemoryHighWaterMark + 0.1}, nil
	}
	engineImpl.run()
	defer engineImpl.cancel()

	globalMemoryUsageCheckInterval.Store(time.Millisecond)

	go engineImpl.globalMemoryUsageChecker(engineImpl.ctx)
	time.Sleep(time.Second)
}

func Test_Engine_watermarkFlusher_LowerThan_MemoryLowWaterMark(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	engineImpl, _ := newEngine(engineCfg)
	engineImpl.memoryStatGetterFunc = func() (*models.MemoryStat, error) {
		return &models.MemoryStat{UsedPercent: constants.MemoryLowWaterMark - 0.1}, nil
	}
	engineImpl.run()
	defer engineImpl.cancel()

	go engineImpl.watermarkFlusher(engineImpl.ctx)

	time.Sleep(time.Second)
}

func Test_Engine_watermarkFlusher_HigherThan_MemoryLowWaterMark(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	engineImpl, _ := newEngine(engineCfg)
	engineImpl.memoryStatGetterFunc = func() (*models.MemoryStat, error) {
		return &models.MemoryStat{UsedPercent: constants.MemoryLowWaterMark + 0.1}, nil
	}
	engineImpl.run()
	defer engineImpl.cancel()

	go engineImpl.watermarkFlusher(engineImpl.ctx)

	time.Sleep(time.Second)
}

func Test_Engine_databaseMetaFlusher_shardMemoryUsageChecker(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	flushMetaInterval.Store(time.Millisecond)
	shardMemoryUsageCheckInterval.Store(time.Millisecond)

	e, _ := NewEngine(engineCfg)
	engineImpl := e.(*engine)
	defer engineImpl.cancel()

	time.Sleep(time.Second)
}

func Test_Engine_flushBiggestMemoryUsageShard(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e, _ := NewEngine(engineCfg)
	engineImpl := e.(*engine)
	defer engineImpl.cancel()

	mockMemoryDatabase := memdb.NewMockMemoryDatabase(ctrl)
	mockMemoryDatabase.EXPECT().MemSize().Return(1024 * 1024 * 1024).AnyTimes()
	mockShard := NewMockShard(ctrl)
	mockShard.EXPECT().Close().Return(nil).AnyTimes()
	mockShard.EXPECT().MemoryDatabase().Return(mockMemoryDatabase).AnyTimes()
	mockShard.EXPECT().Flush().Return(nil).AnyTimes()
	mockDatabase := &database{isFlushing: *atomic.NewBool(true)}
	mockDatabase.shards.Store(int32(1), mockShard)
	engineImpl.databases.Store("1", mockDatabase)

	// mock all shards is flushing
	mockShard.EXPECT().IsFlushing().Return(true).Times(1)
	e.flushBiggestMemoryUsageShard(engineImpl.ctx)

	// mock engine is full flushing
	mockShard.EXPECT().IsFlushing().Return(false).Times(1)
	engineImpl.isFullFlushing.Store(true)
	e.flushBiggestMemoryUsageShard(engineImpl.ctx)

	// mock biggest-shard available
	mockShard.EXPECT().IsFlushing().Return(false).AnyTimes()
	engineImpl.isFullFlushing.Store(false)
	e.flushBiggestMemoryUsageShard(engineImpl.ctx)

	time.Sleep(time.Second)
}

func Test_Engine_flushShardAboveMemoryUsageThreshold_flushAllDatabasesAndShards(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e, _ := NewEngine(engineCfg)
	engineImpl := e.(*engine)
	defer engineImpl.cancel()

	mockMemoryDatabase := memdb.NewMockMemoryDatabase(ctrl)
	mockShard := NewMockShard(ctrl)
	mockShard.EXPECT().Close().Return(nil).AnyTimes()
	mockShard.EXPECT().MemoryDatabase().Return(mockMemoryDatabase).AnyTimes()
	mockShard.EXPECT().Flush().Return(nil).AnyTimes()
	mockDatabase := &database{isFlushing: *atomic.NewBool(true)}
	mockDatabase.shards.Store(int32(1), mockShard)
	engineImpl.databases.Store("1", mockDatabase)

	// mock isFullFlushing
	engineImpl.isFullFlushing.Store(true)
	engineImpl.flushShardAboveMemoryUsageThreshold(engineImpl.ctx)

	engineImpl.isFullFlushing.Store(false)

	// mock no available shard to flush
	mockMemoryDatabase.EXPECT().MemSize().Return(0).Times(1)
	engineImpl.flushShardAboveMemoryUsageThreshold(engineImpl.ctx)

	// mock shard available to flush
	mockMemoryDatabase.EXPECT().MemSize().Return(1024 * 1024 * 1024).Times(1)
	engineImpl.flushShardAboveMemoryUsageThreshold(engineImpl.ctx)

	// flushAllDatabasesAndShards
	engineImpl.flushAllDatabasesAndShards(engineImpl.ctx)
}

func Test_Engine_flushWorker_error(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e, _ := NewEngine(engineCfg)
	engineImpl := e.(*engine)
	defer engineImpl.cancel()

	mockShard := NewMockShard(ctrl)
	mockShard.EXPECT().Close().Return(nil).AnyTimes()
	mockShard.EXPECT().Flush().Return(fmt.Errorf("error")).AnyTimes()

	mockDatabase := NewMockDatabase(ctrl)
	mockDatabase.EXPECT().FlushMeta().Return(fmt.Errorf("error")).AnyTimes()

	engineImpl.shardToFlushCh <- mockShard
	engineImpl.databaseToFlushCh <- mockDatabase

}
