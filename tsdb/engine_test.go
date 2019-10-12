package tsdb

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/option"
)

var testPath = "test_data"
var validOption = option.EngineOption{Interval: "10s"}
var engineCfg = config.Engine{Dir: testPath}

func TestNew(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	factory, err := NewEngineFactory(engineCfg)
	assert.NoError(t, err)

	engine, _ := factory.CreateEngine("test_db")
	assert.NotNil(t, engine)
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db")))

	assert.Equal(t, 0, engine.NumOfShards())

	err = engine.CreateShards(option.EngineOption{})
	assert.NotNil(t, err)

	err = engine.CreateShards(option.EngineOption{}, 1, 2, 3)
	assert.NotNil(t, err)

	err = engine.CreateShards(validOption, 1, 2, 3)
	assert.Nil(t, err)
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db", "OPTIONS")))
	assert.Equal(t, "test_db", engine.Name())

	assert.NotNil(t, engine.GetShard(1))
	assert.NotNil(t, engine.GetShard(2))
	assert.NotNil(t, engine.GetShard(3))
	assert.Nil(t, engine.GetShard(10))
	assert.Equal(t, 3, engine.NumOfShards())

	assert.Nil(t, factory.GetEngine("no_exist"))
	assert.NotNil(t, engine.GetIDGetter())
	assert.NotNil(t, engine.GetExecutePool())

	factory.Close()

	// re-open factory
	factory, err = NewEngineFactory(engineCfg)
	assert.NoError(t, err)

	engine = factory.GetEngine("test_db")
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db")))
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db", "OPTIONS")))

	assert.NotNil(t, engine.GetShard(1))
	assert.NotNil(t, engine.GetShard(2))
	assert.NotNil(t, engine.GetShard(3))
	assert.Nil(t, engine.GetShard(10))
	assert.Equal(t, 3, engine.NumOfShards())
}
