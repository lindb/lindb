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
	"path"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"
)

var writeConfigTestLock sync.Mutex

func withTestPath(dir string) {
	cfg := config.GlobalStorageConfig()
	cfg.TSDB.Dir = dir
}

func TestEngine_New(t *testing.T) {
	writeConfigTestLock.Lock()
	defer writeConfigTestLock.Unlock()

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "make engine path err",
			prepare: func() {
				mkDirIfNotExist = func(path string) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "load engine err",
			prepare: func() {
				listDir = func(path string) (strings []string, e error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create engine successfully",
			prepare: func() {
				listDir = func(path string) (strings []string, e error) {
					return nil, nil
				}
			},
		},
		{
			name: "create engine err because load database err",
			prepare: func() {
				listDir = func(path string) (strings []string, e error) {
					return []string{"db"}, nil
				}
				newDatabaseFunc = func(databaseName string, cfg *models.DatabaseConfig,
					flushChecker DataFlushChecker) (Database, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				mkDirIfNotExist = fileutil.MkDirIfNotExist
				listDir = fileutil.ListDir
				newDatabaseFunc = newDatabase
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			e, err := NewEngine()
			if ((err != nil) != tt.wantErr && e == nil) || (!tt.wantErr && e == nil) {
				t.Errorf("NewEngine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEngine_createDatabase(t *testing.T) {
	writeConfigTestLock.Lock()
	defer writeConfigTestLock.Unlock()
	ctrl := gomock.NewController(t)
	tmpDir := t.TempDir()
	defer func() {
		mkDirIfNotExist = fileutil.MkDirIfNotExist
		fileExist = fileutil.Exist
		decodeToml = ltoml.DecodeToml
		ctrl.Finish()
	}()

	t.Run("create database successfully", func(t *testing.T) {
		defer func() {
			newDatabaseFunc = newDatabase
		}()
		mockDB := NewMockDatabase(ctrl)
		newDatabaseFunc = func(databaseName string, cfg *models.DatabaseConfig, flushChecker DataFlushChecker) (Database, error) {
			return mockDB, nil
		}
		withTestPath(path.Join(tmpDir, "new"))
		e, err := NewEngine()
		assert.NoError(t, err)
		assert.NotNil(t, e)
		db, err := e.createDatabase("test_db", &option.DatabaseOption{})
		assert.NotNil(t, db)
		assert.NoError(t, err)

		db, ok := e.GetDatabase("test_db")
		assert.NotNil(t, db)
		assert.True(t, ok)

		db, ok = e.GetDatabase("db_not_exist")
		assert.Nil(t, db)
		assert.False(t, ok)

		shard, ok := e.GetShard("db_not_exist", models.ShardID(1))
		assert.False(t, ok)
		assert.Nil(t, shard)

		mockDB.EXPECT().GetShard(gomock.Any()).Return(NewMockShard(ctrl), true)
		shard, ok = e.GetShard("test_db", models.ShardID(1))
		assert.True(t, ok)
		assert.NotNil(t, shard)

		assert.Equal(t, map[string]Database{"test_db": mockDB}, e.GetAllDatabases())

		mockDB.EXPECT().Close()
		e.Close()
	})

	t.Run("re-open", func(t *testing.T) {
		defer func() {
			newDatabaseFunc = newDatabase
			listDir = fileutil.ListDir
		}()
		mockDB := NewMockDatabase(ctrl)
		newDatabaseFunc = func(databaseName string, cfg *models.DatabaseConfig, flushChecker DataFlushChecker) (Database, error) {
			return mockDB, nil
		}
		withTestPath(path.Join(tmpDir, "re-open"))
		e, err := NewEngine()
		assert.NoError(t, err)
		assert.NotNil(t, e)
		db, err := e.createDatabase("test_reopen_db", &option.DatabaseOption{})
		assert.NotNil(t, db)
		assert.NoError(t, err)
		mockDB.EXPECT().Close()
		e.Close()

		listDir = func(path string) ([]string, error) {
			return []string{"test_reopen_db"}, nil
		}
		e, err = NewEngine()
		assert.NoError(t, err)
		assert.NotNil(t, e)
		db, ok := e.GetDatabase("test_reopen_db")
		assert.NotNil(t, db)
		assert.True(t, ok)

		db, ok = e.GetDatabase("db_not_exist")
		assert.Nil(t, db)
		assert.False(t, ok)
	})
	t.Run("option file decode failure", func(t *testing.T) {
		defer func() {
			fileExist = fileutil.Exist
			decodeToml = ltoml.DecodeToml
		}()
		e := &engine{}
		fileExist = func(file string) bool {
			return true
		}
		decodeToml = func(fileName string, v interface{}) error {
			return fmt.Errorf("err")
		}
		db, err := e.createDatabase("test_reopen_db", &option.DatabaseOption{})
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}

func Test_Engine_Close(t *testing.T) {
	writeConfigTestLock.Lock()
	defer writeConfigTestLock.Unlock()

	tmpDir := t.TempDir()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	withTestPath(tmpDir)

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
	writeConfigTestLock.Lock()
	defer writeConfigTestLock.Unlock()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	withTestPath(t.TempDir())

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

func TestEngine_DropDatabases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e, _ := NewEngine()
	engineImpl := e.(*engine)
	mockDatabase1 := NewMockDatabase(ctrl)
	engineImpl.dbSet.PutDatabase("test_db_1", mockDatabase1)
	mockDatabase2 := NewMockDatabase(ctrl)
	engineImpl.dbSet.PutDatabase("test_db_2", mockDatabase2)

	// drop fail
	mockDatabase1.EXPECT().Drop().Return(fmt.Errorf("err"))
	e.DropDatabases(map[string]struct{}{"test_db_2": {}})
	assert.Len(t, engineImpl.dbSet.Entries(), 2)
	// drop ok
	mockDatabase1.EXPECT().Drop().Return(nil)
	e.DropDatabases(map[string]struct{}{"test_db_2": {}})
	assert.Len(t, engineImpl.dbSet.Entries(), 1)
}

func TestEngine_TTL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e, _ := NewEngine()
	engineImpl := e.(*engine)
	mockDatabase1 := NewMockDatabase(ctrl)
	engineImpl.dbSet.PutDatabase("test_db_1", mockDatabase1)
	mockDatabase1.EXPECT().TTL()
	e.TTL()
}

func TestEngine_EvictSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e, _ := NewEngine()
	engineImpl := e.(*engine)
	mockDatabase1 := NewMockDatabase(ctrl)
	engineImpl.dbSet.PutDatabase("test_db_1", mockDatabase1)
	mockDatabase1.EXPECT().EvictSegment()
	e.EvictSegment()
}

func TestEngine_CreateShards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		fileExist = fileutil.Exist
		ctrl.Finish()
	}()
	fileExist = func(file string) bool {
		return false
	}
	mockDatabase := NewMockDatabase(ctrl)

	cases := []struct {
		name     string
		db       string
		shardIDs []models.ShardID
		prepare  func(e *engine)
		wantErr  bool
	}{
		{
			name:    "shard ids is empty",
			wantErr: true,
		},
		{
			name:     "create shard failure",
			db:       "test",
			shardIDs: []models.ShardID{1},
			prepare: func(e *engine) {
				mockDatabase.EXPECT().CreateShards(gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:     "create shard successfully",
			db:       "test",
			shardIDs: []models.ShardID{1},
			prepare: func(e *engine) {
				mockDatabase.EXPECT().CreateShards(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "create db failure",
			db:       "test-2",
			shardIDs: []models.ShardID{1},
			prepare: func(e *engine) {
				newDatabaseFunc = func(databaseName string, cfg *models.DatabaseConfig,
					flushChecker DataFlushChecker) (Database, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name:     "create db/shard successfully",
			db:       "test-2",
			shardIDs: []models.ShardID{1},
			prepare: func(e *engine) {
				newDatabaseFunc = func(databaseName string, cfg *models.DatabaseConfig,
					flushChecker DataFlushChecker) (Database, error) {
					return mockDatabase, nil
				}
				mockDatabase.EXPECT().CreateShards(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &engine{
				dbSet: *newDatabaseSet(),
			}
			e.dbSet.PutDatabase("test", mockDatabase)
			if tt.prepare != nil {
				tt.prepare(e)
			}
			err := e.CreateShards(tt.db, &option.DatabaseOption{}, tt.shardIDs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("newShard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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
