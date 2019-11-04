package tsdb

import (
	"path/filepath"
	"testing"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/option"

	"github.com/stretchr/testify/assert"
)

var testPath = "test_data"
var validOption = option.DatabaseOption{Interval: "10s"}
var engineCfg = config.Engine{Dir: testPath}

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
