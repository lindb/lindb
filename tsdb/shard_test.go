package tsdb

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

var _testShard1Path = filepath.Join(testPath, shardDir, "1")

func TestShard_New(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		mkdirFunc = fileutil.MkDirIfNotExist
		newReplicaSequenceFunc = newReplicaSequence
		newIntervalSegmentFunc = newIntervalSegment
		newKVStoreFunc = kv.NewStore
		newIndexDBFunc = indexdb.NewIndexDatabase
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("db").AnyTimes()
	db.EXPECT().Metadata().Return(nil).AnyTimes()
	// case 1: database option err
	thisShard, err := newShard(db, 1, _testShard1Path, option.DatabaseOption{})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 2: interval err
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "as"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 3: create path err
	mkdirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 4: new replica sequence err
	mkdirFunc = fileutil.MkDirIfNotExist
	newReplicaSequenceFunc = func(dirPath string) (sequence ReplicaSequence, err error) {
		return nil, fmt.Errorf("err")
	}
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 5: new interval segment err
	newReplicaSequenceFunc = newReplicaSequence
	newIntervalSegmentFunc = func(interval timeutil.Interval, path string) (segment IntervalSegment, err error) {
		return nil, fmt.Errorf("err")
	}
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 6: new kv store err
	newIntervalSegmentFunc = newIntervalSegment
	newKVStoreFunc = func(name string, option kv.StoreOption) (store kv.Store, err error) {
		return nil, fmt.Errorf("err")
	}
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 7: create forward family err
	kvStore := kv.NewMockStore(ctrl)
	newKVStoreFunc = func(name string, option kv.StoreOption) (store kv.Store, err error) {
		return kvStore, nil
	}
	kvStore.EXPECT().CreateFamily(forwardIndexDir, gomock.Any()).Return(nil, fmt.Errorf("err"))
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 8: create forward family err
	family := kv.NewMockFamily(ctrl)
	kvStore.EXPECT().CreateFamily(forwardIndexDir, gomock.Any()).Return(family, nil)
	kvStore.EXPECT().CreateFamily(invertedIndexDir, gomock.Any()).Return(nil, fmt.Errorf("err"))
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)
	// case 9: create index db err
	kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(family, nil).AnyTimes()
	newIndexDBFunc = func(ctx context.Context, name, parent string,
		metadata metadb.Metadata, forward kv.Family, inverted kv.Family,
	) (indexDatabase indexdb.IndexDatabase, err error) {
		return nil, fmt.Errorf("err")
	}
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Error(t, err)
	assert.Nil(t, thisShard)

	// case 10: create shard success
	newIndexDBFunc = indexdb.NewIndexDatabase
	thisShard, err = newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.NoError(t, err)
	assert.NotNil(t, thisShard)
	assert.NotNil(t, thisShard.IndexDatabase())
	assert.Equal(t, "db", thisShard.DatabaseName())
	assert.Equal(t, int32(1), thisShard.ShardID())
	s, err := thisShard.GetOrCreateSequence("tes")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	assert.True(t, fileutil.Exist(_testShard1Path))
	assert.False(t, thisShard.IsFlushing())
}

func TestShard_GetDataFamilies(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test-db").AnyTimes()
	db.EXPECT().Metadata().Return(nil).AnyTimes()
	s, _ := newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	assert.Nil(t, s.GetDataFamilies(timeutil.Month, timeutil.TimeRange{}))
	assert.Nil(t, s.GetDataFamilies(timeutil.Day, timeutil.TimeRange{}))
	assert.Equal(t, 0, len(s.GetDataFamilies(timeutil.Day, timeutil.TimeRange{})))
}

func TestShard_Write(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	db.EXPECT().Name().Return("test-db").AnyTimes()
	db.EXPECT().Metadata().Return(metadata).AnyTimes()

	mockMemDB := memdb.NewMockMemoryDatabase(ctrl)
	mockMemDB.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	shardINTF, _ := newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s", Behind: "1m", Ahead: "1m"})
	shardIns := shardINTF.(*shard)
	shardIns.mutable = mockMemDB
	shardIns.indexDB = indexDB

	// case 1: metric nil
	assert.Error(t, shardINTF.Write(nil))
	// case 2: metric name is empty
	assert.Error(t, shardINTF.Write(&pb.Metric{
		Timestamp: timeutil.Now(),
	}))
	// case 3: field is empty
	assert.Error(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
	}))
	// case 4: reject before
	assert.NoError(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now() - 2*timeutil.OneMinute,
		Fields: []*pb.Field{{
			Name:   "f1",
			Type:   pb.FieldType_Sum,
			Fields: []*pb.PrimitiveField{{Value: 1.0, PrimitiveID: int32(field.SimpleFieldPFieldID)}},
		}},
	}))
	// case 5: reject ahead
	assert.NoError(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now() + 2*timeutil.OneMinute,
		Fields: []*pb.Field{{
			Name:   "f1",
			Fields: []*pb.PrimitiveField{{Value: 1.0, PrimitiveID: int32(field.SimpleFieldPFieldID)}},
		}},
	}))
	// case 6: gen metric id err
	metadataDB.EXPECT().GenMetricID(constants.DefaultNamespace, "test").Return(uint32(0), fmt.Errorf("err"))
	assert.Error(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:   "f1",
			Fields: []*pb.PrimitiveField{{Value: 1.0, PrimitiveID: int32(field.SimpleFieldPFieldID)}},
		}},
	}))
	// case 7: gen series id err
	metadataDB.EXPECT().GenMetricID(constants.DefaultNamespace, "test").Return(uint32(10), nil).AnyTimes()
	indexDB.EXPECT().GetOrCreateSeriesID(uint32(10), uint64(10)).Return(uint32(0), false, fmt.Errorf("err"))
	assert.Error(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		TagsHash:  10,
		Tags:      map[string]string{"ip": "1.1.1.1"},
		Fields: []*pb.Field{{
			Name:   "f1",
			Fields: []*pb.PrimitiveField{{Value: 1.0, PrimitiveID: int32(field.SimpleFieldPFieldID)}},
		}},
	}))
	// case 7: get old series id
	indexDB.EXPECT().GetOrCreateSeriesID(uint32(10), uint64(10)).Return(uint32(10), false, nil)
	assert.NoError(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		TagsHash:  10,
		Tags:      map[string]string{"ip": "1.1.1.1"},
		Fields: []*pb.Field{{
			Name:   "f1",
			Fields: []*pb.PrimitiveField{{Value: 1.0, PrimitiveID: int32(field.SimpleFieldPFieldID)}},
		}},
	}))
	// case 8: create new series id
	indexDB.EXPECT().GetOrCreateSeriesID(uint32(10), uint64(10)).Return(uint32(10), true, nil)
	indexDB.EXPECT().BuildInvertIndex(constants.DefaultNamespace, "test", map[string]string{"ip": "1.1.1.1"}, uint32(10))
	assert.NoError(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		TagsHash:  10,
		Tags:      map[string]string{"ip": "1.1.1.1"},
		Fields: []*pb.Field{{
			Name:   "f1",
			Fields: []*pb.PrimitiveField{{Value: 1.0, PrimitiveID: int32(field.SimpleFieldPFieldID)}},
		}},
	}))

	assert.NotNil(t, shardINTF.MemoryDatabase())
}

func TestShard_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()

	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test-db").AnyTimes()
	db.EXPECT().Metadata().Return(nil).AnyTimes()
	s, _ := newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	index := indexdb.NewMockIndexDatabase(ctrl)
	s1 := s.(*shard)
	s1.indexDB = index
	// case 1: close index db err
	index.EXPECT().Close().Return(fmt.Errorf("err"))
	err := s.Close()
	assert.Error(t, err)
	// case 2: flush family err
	index.EXPECT().Close().Return(nil).AnyTimes()
	mutable := memdb.NewMockMemoryDatabase(ctrl)
	s1.mutable = mutable
	mutable.EXPECT().Families().Return([]int64{1, 2})
	mutable.EXPECT().FlushFamilyTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = s.Close()
	assert.Error(t, err)
	// case 3: close success
	mutable.EXPECT().Families().Return(nil)
	err = s.Close()
	assert.NoError(t, err)
}

func TestShard_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test-db").AnyTimes()
	db.EXPECT().Metadata().Return(nil).AnyTimes()
	s, _ := newShard(db, 1, _testShard1Path, option.DatabaseOption{Interval: "10s"})
	s1 := s.(*shard)
	mutable := memdb.NewMockMemoryDatabase(ctrl)
	s1.mutable = mutable
	// case 1: flush is doing
	s1.isFlushing.Store(true)
	err := s.Flush()
	assert.NoError(t, err)
	// case 2: flush err
	s1.isFlushing.Store(false)
	mutable.EXPECT().Families().Return([]int64{1, 2})
	mutable.EXPECT().FlushFamilyTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = s.Flush()
	assert.Error(t, err)
	// case 3: get segment err
	intervalSegment := NewMockIntervalSegment(ctrl)
	s1.segment = intervalSegment
	mutable.EXPECT().Families().Return([]int64{1, 2}).AnyTimes()
	intervalSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = s.Flush()
	assert.Error(t, err)
	// case 4: ack replica sequence err
	seq := NewMockReplicaSequence(ctrl)
	s1.sequence = seq
	seq.EXPECT().getAllHeads().Return(nil).AnyTimes()
	seq.EXPECT().ack(gomock.Any()).Return(fmt.Errorf("err")).AnyTimes()
	intervalSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err = s.Flush()
	assert.Error(t, err)
	// case 5: get family err
	segment := NewMockSegment(ctrl)
	intervalSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(segment, nil).AnyTimes()
	segment.EXPECT().GetDataFamily(gomock.Any()).Return(nil, fmt.Errorf("err")).Times(2)
	err = s.Flush()
	assert.NoError(t, err)
}
