package tsdb

import (
	"path/filepath"
	"testing"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/series"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var path = filepath.Join(testPath, shardPath, "1")

func TestNewShard(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	shard, err := newShard(1, path, option.EngineOption{})
	assert.NotNil(t, err)
	assert.Nil(t, shard)

	shard, err = newShard(1, path, option.EngineOption{Interval: "as"})
	assert.NotNil(t, err)
	assert.Nil(t, shard)

	shard, err = newShard(1, path, option.EngineOption{Interval: "10s"})
	assert.Nil(t, err)
	assert.NotNil(t, shard)

	assert.True(t, fileutil.Exist(path))
}

func TestGetSegments(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	shard, _ := newShard(1, path, option.EngineOption{Interval: "10s"})
	assert.Nil(t, shard.GetSegments(interval.Month, timeutil.TimeRange{}))
	assert.Nil(t, shard.GetSegments(interval.Day, timeutil.TimeRange{}))
	assert.Equal(t, 0, len(shard.GetSegments(interval.Day, timeutil.TimeRange{})))
}

func TestWrite(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemDB := memdb.NewMockMemoryDatabase(ctrl)
	gomock.InOrder(
		mockMemDB.EXPECT().Write(gomock.Any()).Return(nil),
		mockMemDB.EXPECT().Write(gomock.Any()).Return(series.ErrTooManyTags),
	)

	shardINTF, _ := newShard(1, path, option.EngineOption{Interval: "10s"})
	shardIns := shardINTF.(*shard)
	shardIns.memDB = mockMemDB

	assert.Nil(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{
			{Name: "f1", Field: &pb.Field_Sum{Sum: 1.0}},
		},
	}))
	assert.NotNil(t, shardINTF.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{
			{Name: "f1", Field: &pb.Field_Sum{Sum: 1.0}},
		},
	}))
}
