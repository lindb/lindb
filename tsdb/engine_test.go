package tsdb

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/util"
)

var testPath = "test_data"

func TestNew(t *testing.T) {
	defer util.RemoveDir(testPath)
	engine, _ := NewEngine("test_db", testPath)
	assert.NotNil(t, engine)
	assert.True(t, util.Exist(filepath.Join(testPath, "test_db")))

	err := engine.CreateShards(option.ShardOption{})
	assert.NotNil(t, err)

	err = engine.CreateShards(option.ShardOption{}, 1, 2, 3)
	assert.Nil(t, err)
	assert.True(t, util.Exist(filepath.Join(testPath, "test_db", "OPTIONS")))

	assert.NotNil(t, engine.GetShard(1))
	assert.NotNil(t, engine.GetShard(2))
	assert.NotNil(t, engine.GetShard(3))
	engine.Close()

	// re-open engine test load exist data
	engine, _ = NewEngine("test_db", testPath)
	assert.True(t, util.Exist(filepath.Join(testPath, "test_db")))
	assert.True(t, util.Exist(filepath.Join(testPath, "test_db", "OPTIONS")))

	assert.NotNil(t, engine.GetShard(1))
	assert.NotNil(t, engine.GetShard(2))
	assert.NotNil(t, engine.GetShard(3))
	engine.Close()
}
