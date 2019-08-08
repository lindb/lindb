package tsdb

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/option"
)

var testPath = "test_data"
var validOption = option.ShardOption{Interval: time.Second * 10, IntervalType: interval.Day}

func TestNew(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	factory := NewEngineFactory()
	engine, _ := factory.CreateEngine("test_db", testPath)
	assert.NotNil(t, engine)
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db")))

	assert.Equal(t, 0, engine.NumOfShards())

	err := engine.CreateShards(option.ShardOption{})
	assert.NotNil(t, err)

	err = engine.CreateShards(option.ShardOption{}, 1, 2, 3)
	assert.NotNil(t, err)

	err = engine.CreateShards(validOption, 1, 2, 3)
	assert.Nil(t, err)
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db", "OPTIONS")))

	assert.NotNil(t, engine.GetShard(1))
	assert.NotNil(t, engine.GetShard(2))
	assert.NotNil(t, engine.GetShard(3))
	assert.Nil(t, engine.GetShard(10))
	assert.Equal(t, 3, engine.NumOfShards())
	_ = engine.Close()

	// re-open engine test load exist data
	engine, _ = factory.CreateEngine("test_db", testPath)
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db")))
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db", "OPTIONS")))

	assert.NotNil(t, engine.GetShard(1))
	assert.NotNil(t, engine.GetShard(2))
	assert.NotNil(t, engine.GetShard(3))
	assert.Nil(t, engine.GetShard(10))
	assert.Equal(t, 3, engine.NumOfShards())
	_ = engine.Close()
}
