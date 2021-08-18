// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package tsdb

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/ltoml"
)

var testPath = "test_data"

func withTestPath() {
	cfg := config.GlobalStorageConfig()
	cfg.TSDB.Dir = testPath
}

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
	withTestPath()

	e, err := NewEngine()
	assert.Error(t, err)
	assert.Nil(t, e)
	mkDirIfNotExist = fileutil.MkDirIfNotExist

	// test new err when load engine err
	listDir = func(path string) (strings []string, e error) {
		return nil, fmt.Errorf("err")
	}
	e, err = NewEngine()
	assert.Error(t, err)
	assert.Nil(t, e)
	listDir = fileutil.ListDir

	e, err = NewEngine()
	assert.NoError(t, err)

	db, err := e.createDatabase("test_db")
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
	e, err = NewEngine()
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
	withTestPath()

	e, err := NewEngine()
	assert.NoError(t, err)

	db, err := e.createDatabase("test_db")
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
	e, err = NewEngine()
	assert.Error(t, err)
	assert.Nil(t, e)
	decodeToml = ltoml.DecodeToml

	// re-open engine
	e, err = NewEngine()
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
	db, err = e.createDatabase("test_db_err")
	assert.Error(t, err)
	assert.Nil(t, db)
	mkDirIfNotExist = fileutil.MkDirIfNotExist
	// create db err
	newDatabaseFunc = func(databaseName string, databasePath string, cfg *databaseConfig,
		checker DataFlushChecker) (d Database, err error) {
		return nil, fmt.Errorf("err")
	}
	db, err = e.createDatabase("test_db_err")
	assert.Error(t, err)
	assert.Nil(t, db)
}

func Test_Engine_Close(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	withTestPath()

	e, _ := NewEngine()
	engineImpl := e.(*engine)
	defer engineImpl.cancel()

	mockDatabase := NewMockDatabase(ctrl)
	mockDatabase.EXPECT().Close().Return(fmt.Errorf("error")).AnyTimes()
	engineImpl.dbSet.PutDatabase("1", mockDatabase)
	engineImpl.dbSet.PutDatabase("2", mockDatabase)

	e.Close()
}

func Test_Engine_Flush_Database(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	withTestPath()

	e, _ := NewEngine()
	engineImpl := e.(*engine)
	defer engineImpl.cancel()
	ok := e.FlushDatabase(context.TODO(), "test_db_3")
	assert.False(t, ok)

	mockDatabase := NewMockDatabase(ctrl)
	// case 1: flush success
	mockDatabase.EXPECT().Flush().Return(nil)
	engineImpl.dbSet.PutDatabase("test_db_1", mockDatabase)
	ok = e.FlushDatabase(context.TODO(), "test_db_1")
	assert.True(t, ok)
	// case 2: flush err
	mockDatabase.EXPECT().Flush().Return(fmt.Errorf("err"))
	ok = e.FlushDatabase(context.TODO(), "test_db_1")
	assert.False(t, ok)
}

var testDatabaseNames = []string{
	"_internal", "system", "docker", "network", "java",
	"runtime", "go", "php", "k8s", "infra", "prometheus",
	"application", "nginx", "frontend", "kernel", "other",
	"trace", "test", "test2", "test3", "test4",
}

func BenchmarkEngine_DatabaseWithSyncMap(b *testing.B) {
	var sm sync.Map
	for _, dn := range testDatabaseNames {
		sm.Store(dn, &database{})
	}
	// 9.365 ns
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			item, _ := sm.Load("application")
			_ = item.(*database)
		}
	})
}

func BenchmarkEngine_DatabaseWithLockFreeMap(b *testing.B) {
	var v atomic.Value
	var lm = make(map[string]*database)
	for _, dn := range testDatabaseNames {
		lm[dn] = &database{}
	}
	v.Store(lm)
	// 3.895 ns
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			thisMap := v.Load().(map[string]*database)
			_ = thisMap["application"]
		}
	})
}

func BenchmarkEngine_DatabaseWithLockFreeSlice(b *testing.B) {
	type entry struct {
		name string
		db   *database
	}
	var v atomic.Value
	var entries []entry
	for _, dn := range testDatabaseNames {
		entries = append(entries, entry{name: dn, db: &database{}})
	}
	v.Store(entries)
	// 2.534 ns
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sl := v.Load().([]entry)
			for idx := range sl {
				if sl[idx].name == "nginx" {
					break
				}
			}
		}
	})
}
