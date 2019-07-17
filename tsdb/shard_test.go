package tsdb

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	pb "github.com/eleme/lindb/rpc/proto/field"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/interval"
	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/timeutil"
	"github.com/eleme/lindb/pkg/util"
)

var path = filepath.Join(testPath, shardPath, "1")

func TestNewShard(t *testing.T) {
	defer func() {
		_ = util.RemoveDir(testPath)
	}()
	shard, err := newShard(1, path, option.ShardOption{})
	assert.NotNil(t, err)
	assert.Nil(t, shard)

	shard, err = newShard(1, path, option.ShardOption{Interval: time.Second * 10})
	assert.NotNil(t, err)
	assert.Nil(t, shard)

	shard, err = newShard(1, path, option.ShardOption{Interval: time.Second * 10, IntervalType: interval.Day})
	assert.Nil(t, err)
	assert.NotNil(t, shard)

	assert.True(t, util.Exist(path))
}

func TestGetSegments(t *testing.T) {
	defer func() {
		_ = util.RemoveDir(testPath)
	}()
	shard, _ := newShard(1, path, option.ShardOption{Interval: time.Second * 10, IntervalType: interval.Day})
	assert.Nil(t, shard.GetSegments(interval.Month, models.TimeRange{}))
	assert.Nil(t, shard.GetSegments(interval.Day, models.TimeRange{}))
	assert.Equal(t, 0, len(shard.GetSegments(interval.Day, models.TimeRange{})))
}

func TestWrite(t *testing.T) {
	defer func() {
		_ = util.RemoveDir(testPath)
	}()
	shard, _ := newShard(1, path, option.ShardOption{Interval: time.Second * 10, IntervalType: interval.Day})

	err := shard.Write(&pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{
			{Name: "f1", Field: &pb.Field_Sum{Sum: 1.0}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
