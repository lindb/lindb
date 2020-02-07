package tsdb

import (
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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
	thisShard, err := newShard("db", 1, _testShard1Path, mockIDSequencer, option.DatabaseOption{})
	assert.Error(t, err)
	assert.Nil(t, thisShard)

	thisShard, err = newShard("db", 1, _testShard1Path, mockIDSequencer, option.DatabaseOption{Interval: "as"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)

	thisShard, err = newShard("db", 1, _testShard1Path, mockIDSequencer, option.DatabaseOption{Interval: "10s"})
	assert.NoError(t, err)
	assert.NotNil(t, thisShard)
	assert.Nil(t, thisShard.IndexDatabase())
	assert.Equal(t, "db", thisShard.DatabaseName())
	assert.Equal(t, int32(1), thisShard.ShardID())
	s, err := thisShard.GetOrCreateSequence("tes")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.True(t, fileutil.Exist(_testShard1Path))
}

func TestGetSegments(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDSequencer := metadb.NewMockIDSequencer(ctrl)
	s, _ := newShard("db", 1, _testShard1Path, mockIDSequencer, option.DatabaseOption{Interval: "10s"})
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

	shardINTF, _ := newShard("db", 1, _testShard1Path, mockIDSequencer, option.DatabaseOption{Interval: "10s"})
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
	shardINTF.(*shard).cancel()
}

func TestShard_Write_Accept(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDSequencer := metadb.NewMockIDSequencer(ctrl)

	shardINTF, _ := newShard(
		"db",
		1,
		_testShard1Path,
		mockIDSequencer,
		option.DatabaseOption{Interval: "10s", Ahead: "1h", Behind: "1h"})
	assert.Nil(t, shardINTF.IndexFilter())
	assert.Nil(t, shardINTF.MemoryFilter())

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
	shardINTF.(*shard).cancel()
}

//
//func Test_Shard_Close_Flush_error(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	defer func() {
//		_ = fileutil.RemoveDir(testPath)
//	}()
//	mockStore := kv.NewMockStore(ctrl)
//
//	// prepare mocked segment
//	mockIntervalSegment := NewMockIntervalSegment(ctrl)
//	replicaSequence, err := newReplicaSequence(filepath.Join(testPath, replicaDir))
//	assert.NoError(t, err)
//	s := &shard{
//		segment:  mockIntervalSegment,
//		interval: timeutil.Interval(timeutil.OneSecond * 10),
//		sequence: replicaSequence,
//	}
//	_, cancel := context.WithCancel(context.Background())
//	s.cancel = cancel
//
//	s.indexStore = mockStore
//	mockFlusher := kv.NewMockFlusher(ctrl)
//	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
//	mockFlusher.EXPECT().Commit().Return(nil).AnyTimes()
//
//	mockFamily := kv.NewMockFamily(ctrl)
//	mockFamily.EXPECT().NewFlusher().Return(mockFlusher).AnyTimes()
//	s.invertedFamily = mockFamily
//
//	mockMemDB := memdb.NewMockMemoryDatabase(ctrl)
//	s.memDB = mockMemDB
//	// mock flush ok
//	mockMemDB.EXPECT().Families().Return(nil).AnyTimes()
//	//mockMemDB.EXPECT().FlushInvertedIndexTo(gomock.Any()).Return(nil)
//	mockStore.EXPECT().Close().Return(fmt.Errorf("error")).AnyTimes()
//	assert.NotNil(t, s.Close())
//	// mock flush inverted index error
//	//mockMemDB.EXPECT().FlushInvertedIndexTo(gomock.Any()).Return(fmt.Errorf("error"))
//	assert.NotNil(t, s.Close())
//
//	// mock flush families error
//	mockMemDB.EXPECT().Families().Return([]int64{1}).AnyTimes()
//	//mockMemDB.EXPECT().FlushInvertedIndexTo(gomock.Any()).Return(nil).AnyTimes()
//	// mock GetOrCreateSegment error
//	mockIntervalSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(nil, fmt.Errorf("error"))
//	assert.NotNil(t, s.Close())
//	// mock GetDataFamily error
//	mockSegment := NewMockSegment(ctrl)
//	mockIntervalSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(mockSegment, nil).AnyTimes()
//	mockSegment.EXPECT().GetDataFamily(gomock.Any()).Return(nil, fmt.Errorf("error"))
//	assert.NotNil(t, s.Close())
//	// mock FlushFamilyTo ok
//	mockDataFamily := NewMockDataFamily(ctrl)
//	mockDataFamily.EXPECT().Family().Return(mockFamily).AnyTimes()
//	mockMemDB.EXPECT().FlushFamilyTo(gomock.Any(), gomock.Any()).Return(nil)
//	mockSegment.EXPECT().GetDataFamily(gomock.Any()).Return(mockDataFamily, nil).AnyTimes()
//	assert.NotNil(t, s.Close())
//	// mock FlushFamilyTo error
//	mockMemDB.EXPECT().FlushFamilyTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
//	assert.NotNil(t, s.Close())
//
//	// mock isFlushing CAS false
//	assert.False(t, s.IsFlushing())
//	s.isFlushing.Store(true)
//	assert.Nil(t, s.Flush())
//}
