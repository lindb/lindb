package tsdb

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

var _testShard1Path = filepath.Join(testPath, shardDir, "1")

func TestNewShard(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDSequencer := metadb.NewMockIDSequencer(ctrl)
	thisShard, err := newShard(1, _testShard1Path, mockIDSequencer, option.DatabaseOption{})
	assert.NotNil(t, err)
	assert.Nil(t, thisShard)

	thisShard, err = newShard(1, _testShard1Path, mockIDSequencer, option.DatabaseOption{Interval: "as"})
	assert.NotNil(t, err)
	assert.Nil(t, thisShard)

	thisShard, err = newShard(1, _testShard1Path, mockIDSequencer, option.DatabaseOption{Interval: "10s"})
	assert.Nil(t, err)
	assert.NotNil(t, thisShard)
	assert.NotNil(t, thisShard.IndexDatabase())

	assert.True(t, fileutil.Exist(_testShard1Path))
}

func TestGetSegments(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDSequencer := metadb.NewMockIDSequencer(ctrl)
	s, _ := newShard(1, _testShard1Path, mockIDSequencer, option.DatabaseOption{Interval: "10s"})
	assert.Nil(t, s.GetDataFamilies(timeutil.Month, timeutil.TimeRange{}))
	assert.Nil(t, s.GetDataFamilies(timeutil.Day, timeutil.TimeRange{}))
	assert.Equal(t, 0, len(s.GetDataFamilies(timeutil.Day, timeutil.TimeRange{})))
}

func TestWrite(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemDB := memdb.NewMockMemoryDatabase(ctrl)
	mockIDSequencer := metadb.NewMockIDSequencer(ctrl)
	gomock.InOrder(
		mockMemDB.EXPECT().Write(gomock.Any()).Return(nil),
		mockMemDB.EXPECT().Write(gomock.Any()).Return(series.ErrTooManyTags),
	)

	shardINTF, _ := newShard(1, _testShard1Path, mockIDSequencer, option.DatabaseOption{Interval: "10s"})
	shardIns := shardINTF.(*shard)
	shardIns.memDB = mockMemDB

	assert.NotNil(t, shardINTF.Write(nil))
	assert.NotNil(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
	}))

	assert.Nil(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{
			{Name: "f1", Field: &pb.Field_Sum{Sum: &pb.Sum{Value: 1.0}}},
		},
	}))
	assert.NotNil(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{
			{Name: "f1", Field: &pb.Field_Sum{Sum: &pb.Sum{Value: 1.0}}},
		},
	}))

	assert.NotNil(t, shardINTF.MemoryDatabase())
	shardINTF.Close()
}

func TestShard_Write_Accept(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDSequencer := metadb.NewMockIDSequencer(ctrl)

	shardINTF, _ := newShard(
		1,
		_testShard1Path,
		mockIDSequencer,
		option.DatabaseOption{Interval: "10s", Ahead: "1h", Behind: "1h"})
	assert.NotNil(t, shardINTF.IndexFilter())
	assert.NotNil(t, shardINTF.IndexMetaGetter())
	assert.NotNil(t, shardINTF.MemoryFilter())
	assert.NotNil(t, shardINTF.MemoryMetaGetter())

	assert.Nil(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now() + timeutil.OneHour + 10000,
		Fields: []*pb.Field{
			{Name: "f1", Field: &pb.Field_Sum{Sum: &pb.Sum{Value: 1.0}}},
		},
	}))
	assert.Nil(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now() - timeutil.OneHour - 10000,
		Fields: []*pb.Field{
			{Name: "f1", Field: &pb.Field_Sum{Sum: &pb.Sum{Value: 1.0}}},
		},
	}))
	shardINTF.Close()
}

func Test_Shard_Close_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStore := kv.NewMockStore(ctrl)

	s := &shard{}
	_, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	s.indexStore = mockStore
	mockStore.EXPECT().Close().Return(fmt.Errorf("error")).AnyTimes()
	s.Close()
}
