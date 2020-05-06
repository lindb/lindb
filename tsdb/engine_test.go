package tsdb

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/ltoml"
)

var testPath = "test_data"
var engineCfg = config.TSDB{Dir: testPath}

func TestNew(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		mkDirIfNotExist = fileutil.MkDirIfNotExist
		listDir = fileutil.ListDir
	}()

	// test new error
	mkDirIfNotExist = func(path string) error {
		return fmt.Errorf("err")
	}
	e, err := NewEngine(engineCfg)
	assert.Error(t, err)
	assert.Nil(t, e)
	mkDirIfNotExist = fileutil.MkDirIfNotExist

	// test new err when load engine err
	listDir = func(path string) (strings []string, e error) {
		return nil, fmt.Errorf("err")
	}
	e, err = NewEngine(engineCfg)
	assert.Error(t, err)
	assert.Nil(t, e)
	listDir = fileutil.ListDir

	e, err = NewEngine(engineCfg)
	assert.NoError(t, err)

	db, err := e.CreateDatabase("test_db")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db")))
	assert.Equal(t, 0, db.NumOfShards())
	e.Close()

	// test load db error
	mkDirIfNotExist = func(path string) error {
		if path == filepath.Join(testPath, "test_db") {
			return fmt.Errorf("err")
		}
		return fileutil.MkDirIfNotExist(path)
	}
	e, err = NewEngine(engineCfg)
	assert.Error(t, err)
	assert.Nil(t, e)
}

func TestEngine_CreateDatabase(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		mkDirIfNotExist = fileutil.MkDirIfNotExist
		decodeToml = ltoml.DecodeToml
		newDatabaseFunc = newDatabase
	}()

	e, err := NewEngine(engineCfg)
	assert.NoError(t, err)

	db, err := e.CreateDatabase("test_db")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db")))

	_, ok := e.GetDatabase("inexist")
	assert.False(t, ok)
	assert.NotNil(t, db.ExecutorPool())

	e.Close()

	// re-open engine, err
	decodeToml = func(fileName string, v interface{}) error {
		return fmt.Errorf("err")
	}
	e, err = NewEngine(engineCfg)
	assert.Error(t, err)
	assert.Nil(t, e)
	decodeToml = ltoml.DecodeToml

	// re-open engine
	e, err = NewEngine(engineCfg)
	assert.NoError(t, err)
	db, ok = e.GetDatabase("test_db")
	assert.True(t, ok)
	assert.NotNil(t, db)
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db")))
	assert.True(t, fileutil.Exist(filepath.Join(testPath, "test_db", "OPTIONS")))

	// mkdir database path err
	mkDirIfNotExist = func(path string) error {
		return fmt.Errorf("err")
	}
	db, err = e.CreateDatabase("test_db_err")
	assert.Error(t, err)
	assert.Nil(t, db)
	mkDirIfNotExist = fileutil.MkDirIfNotExist
	// create db err
	newDatabaseFunc = func(databaseName string, databasePath string, cfg *databaseConfig,
		checker DataFlushChecker) (d Database, err error) {
		return nil, fmt.Errorf("err")
	}
	db, err = e.CreateDatabase("test_db_err")
	assert.Error(t, err)
	assert.Nil(t, db)
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
	// case 1: flush success
	mockDatabase.EXPECT().Flush().Return(nil)
	engineImpl.databases.Store("test_db_1", mockDatabase)
	ok = e.FlushDatabase(context.TODO(), "test_db_1")
	assert.True(t, ok)
	// case 2: flush err
	mockDatabase.EXPECT().Flush().Return(fmt.Errorf("err"))
	ok = e.FlushDatabase(context.TODO(), "test_db_1")
	assert.False(t, ok)
}
