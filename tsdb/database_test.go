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
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestDatabase_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encodeToml = ltoml.EncodeToml
		mkDirIfNotExist = fileutil.MkDirIfNotExist
		newMetadataFunc = metadb.NewMetadata
		kv.InitStoreManager(nil)
		ctrl.Finish()
	}()

	storeMgr := kv.NewMockStoreManager(ctrl)
	store := kv.NewMockStore(ctrl)
	kv.InitStoreManager(storeMgr)
	opt := &option.DatabaseOption{}

	cases := []struct {
		name    string
		cfg     *models.DatabaseConfig
		prepare func()
		wantErr bool
	}{
		{
			name: "create database path err",
			prepare: func() {
				mkDirIfNotExist = func(path string) error {
					return fmt.Errorf("mkdir err")
				}
			},
			wantErr: true,
		},
		{
			name: "dump config err",
			prepare: func() {
				encodeToml = func(fileName string, v interface{}) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create kv store err",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("create store err"))
			},
			wantErr: true,
		},
		{
			name: "create kv family err",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(store, nil)
				store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "create metadata err",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(store, nil)
				store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil)
				newMetadataFunc = func(ctx context.Context, databaseName, parent string,
					tagFamily kv.Family) (metadata metadb.Metadata, err error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name:    "option validation fail",
			cfg:     &models.DatabaseConfig{Option: opt},
			wantErr: true,
		},
		{
			name: "create shard err",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(store, nil)
				store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil)
				newShardFunc = func(db Database, shardID models.ShardID) (s Shard, err error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create database successfully",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(store, nil)
				store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil)
				metadata := metadb.NewMockMetadata(ctrl)
				newMetadataFunc = func(ctx context.Context, databaseName, parent string,
					tagFamily kv.Family) (metadb.Metadata, error) {
					return metadata, nil
				}
				newShardFunc = func(db Database, shardID models.ShardID) (s Shard, err error) {
					return nil, nil
				}
			},
			wantErr: false,
		},
		{
			name: "close metadata err when create database failure",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(store, nil)
				store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil)
				metadata := metadb.NewMockMetadata(ctrl)
				newMetadataFunc = func(ctx context.Context, databaseName, parent string,
					tagFamily kv.Family) (metadb.Metadata, error) {
					return metadata, nil
				}
				newShardFunc = func(db Database, shardID models.ShardID) (s Shard, err error) {
					return nil, fmt.Errorf("err")
				}
				metadata.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				encodeToml = func(fileName string, v interface{}) error {
					return nil
				}
				mkDirIfNotExist = func(path string) error {
					return nil
				}
				newMetadataFunc = func(ctx context.Context, databaseName,
					parent string, tagFamily kv.Family) (metadb.Metadata, error) {
					return nil, nil
				}
				newShardFunc = newShard
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			opt := &option.DatabaseOption{Intervals: option.Intervals{{Interval: 10}}}
			cfg := &models.DatabaseConfig{
				ShardIDs: []models.ShardID{1, 2, 3},
				Option:   opt,
			}
			if tt.cfg != nil {
				cfg = tt.cfg
			}
			db, err := newDatabase("db", cfg, nil)
			if ((err != nil) != tt.wantErr && db == nil) || (!tt.wantErr && db == nil) {
				t.Errorf("newDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}

			if db != nil {
				// assert database information after create successfully
				assert.NotNil(t, db.Metadata())
				assert.NotNil(t, db.ExecutorPool())
				assert.Equal(t, "db", db.Name())
				assert.True(t, db.NumOfShards() >= 0)
				assert.Equal(t, &option.DatabaseOption{Intervals: option.Intervals{{Interval: 10}}}, db.GetOption())
				assert.NotNil(t, db.GetConfig())
			}
		})
	}
}

func TestDatabase_CreateShards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encodeToml = ltoml.EncodeToml
		ctrl.Finish()
	}()
	db := &database{
		config:   &models.DatabaseConfig{},
		shardSet: *newShardSet(),
	}
	type args struct {
		option   option.DatabaseOption
		shardIDs []models.ShardID
	}
	cases := []struct {
		name    string
		args    args
		prepare func()
		wantErr bool
	}{
		{
			name:    "shard ids cannot be empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "create shard err",
			args: args{option.DatabaseOption{}, []models.ShardID{4, 5, 6}},
			prepare: func() {
				newShardFunc = func(db Database, shardID models.ShardID) (s Shard, err error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create exist shard",
			args: args{option.DatabaseOption{}, []models.ShardID{4}},
			prepare: func() {
				db.shardSet.InsertShard(models.ShardID(4), nil)
			},
			wantErr: false,
		},
		{
			name: "create shard successfully",
			args: args{option.DatabaseOption{}, []models.ShardID{5}},
			prepare: func() {
				newShardFunc = func(db Database, shardID models.ShardID) (s Shard, err error) {
					return nil, nil
				}
			},
			wantErr: false,
		},
		{
			name: "dump option err",
			args: args{option.DatabaseOption{}, []models.ShardID{6}},
			prepare: func() {
				newShardFunc = func(db Database, shardID models.ShardID) (s Shard, err error) {
					return nil, nil
				}
				encodeToml = func(fileName string, v interface{}) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newShardFunc = newShard
				encodeToml = func(fileName string, v interface{}) error {
					return nil
				}
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			if err := db.CreateShards(tt.args.shardIDs); (err != nil) != tt.wantErr {
				t.Errorf("CreateShards() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("create exist shard", func(t *testing.T) {
		db.shardSet.InsertShard(1, nil)
		err := db.createShard(1)
		assert.NoError(t, err)
	})
}

func TestDatabase_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		kv.InitStoreManager(nil)
		ctrl.Finish()
	}()

	storeMgr := kv.NewMockStoreManager(ctrl)
	kv.InitStoreManager(storeMgr)
	metadata := metadb.NewMockMetadata(ctrl)
	metadata.EXPECT().Flush().Return(nil).AnyTimes()
	store := kv.NewMockStore(ctrl)
	store.EXPECT().Name().Return("metaStore").AnyTimes()
	db := &database{
		metadata:       metadata,
		shardSet:       *newShardSet(),
		metaStore:      store,
		flushCondition: sync.NewCond(&sync.Mutex{}),
	}
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "close metadata err",
			prepare: func() {
				metadata.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "close meta store err",
			prepare: func() {
				gomock.InOrder(
					metadata.EXPECT().Close().Return(nil),
					storeMgr.EXPECT().CloseStore("metaStore").Return(fmt.Errorf("err")),
				)
			},
			wantErr: true,
		},
		{
			name: "close meta store err",
			prepare: func() {
				mockShard := NewMockShard(ctrl)
				db.shardSet.InsertShard(models.ShardID(1), mockShard)
				gomock.InOrder(
					metadata.EXPECT().Close().Return(nil),
					storeMgr.EXPECT().CloseStore("metaStore").Return(nil),
					mockShard.EXPECT().Close().Return(fmt.Errorf("err")),
				)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			if err := db.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_FlushMeta(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metadata := metadb.NewMockMetadata(ctrl)
	db := &database{
		metadata:       metadata,
		flushCondition: sync.NewCond(&sync.Mutex{}),
		isFlushing:     *atomic.NewBool(false),
		statistics:     metrics.NewDatabaseStatistics("test"),
	}
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "meta flushing",
			prepare: func() {
				db.isFlushing.Store(true)
			},
			wantErr: false,
		},
		{
			name: "flush meta failure",
			prepare: func() {
				db.isFlushing.Store(false)
				metadata.EXPECT().Flush().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "flush meta successfully",
			prepare: func() {
				db.isFlushing.Store(false)
				metadata.EXPECT().Flush().Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				db.isFlushing.Store(false)
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			if err := db.FlushMeta(); (err != nil) != tt.wantErr {
				t.Errorf("FlushMeta() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	checker := NewMockDataFlushChecker(ctrl)

	db := &database{
		shardSet:     *newShardSet(),
		isFlushing:   *atomic.NewBool(false),
		flushChecker: checker,
	}
	shard1 := NewMockShard(ctrl)
	shard2 := NewMockShard(ctrl)
	shard1.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	shard2.EXPECT().ShardID().Return(models.ShardID(2)).AnyTimes()
	db.shardSet.InsertShard(1, shard1)
	db.shardSet.InsertShard(2, shard2)
	checker.EXPECT().requestFlushJob(gomock.Any())
	checker.EXPECT().requestFlushJob(gomock.Any())
	err := db.Flush()
	assert.NoError(t, err)
}

func Test_ShardSet_multi(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	set := newShardSet()
	shard1 := NewMockShard(ctrl)
	for i := 0; i < 100; i += 2 {
		set.InsertShard(models.ShardID(i), shard1)
	}
	assert.Equal(t, set.GetShardNum(), 50)
	_, ok := set.GetShard(0)
	assert.True(t, ok)
	_, ok = set.GetShard(11)
	assert.False(t, ok)
	_, ok = set.GetShard(101)
	assert.False(t, ok)
}

func TestDatabase_WaitFlushMetaCompleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := timeutil.Now()
	metadata := metadb.NewMockMetadata(ctrl)
	db := &database{
		metadata:       metadata,
		isFlushing:     *atomic.NewBool(false),
		flushCondition: sync.NewCond(&sync.Mutex{}),
		statistics:     metrics.NewDatabaseStatistics("test"),
	}

	metadata.EXPECT().Flush().DoAndReturn(func() error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})
	var wait sync.WaitGroup
	wait.Add(2)
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
		err := db.FlushMeta()
		assert.NoError(t, err)
	}()
	<-ch
	time.Sleep(10 * time.Millisecond)
	go func() {
		db.WaitFlushMetaCompleted()
		wait.Done()
	}()
	go func() {
		db.WaitFlushMetaCompleted()
		wait.Done()
	}()
	wait.Wait()
	assert.True(t, timeutil.Now()-now >= 100*time.Millisecond.Milliseconds())
}

func TestDatabase_Drop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		removeDir = fileutil.RemoveDir
		kv.InitStoreManager(nil)
		ctrl.Finish()
	}()
	storeMgr := kv.NewMockStoreManager(ctrl)
	kv.InitStoreManager(storeMgr)
	storeMgr.EXPECT().CloseStore(gomock.Any()).Return(nil).AnyTimes()
	metadata := metadb.NewMockMetadata(ctrl)
	store := kv.NewMockStore(ctrl)
	db := &database{
		metadata:       metadata,
		metaStore:      store,
		shardSet:       *newShardSet(),
		isFlushing:     *atomic.NewBool(false),
		flushCondition: sync.NewCond(&sync.Mutex{}),
	}
	store.EXPECT().Name().Return("test").AnyTimes()
	metadata.EXPECT().Close().Return(fmt.Errorf("err"))
	assert.Error(t, db.Drop())
	removeDir = func(path string) error {
		return fmt.Errorf("err")
	}
	metadata.EXPECT().Close().Return(nil)
	assert.Error(t, db.Drop())
	removeDir = func(path string) error {
		return nil
	}
	metadata.EXPECT().Close().Return(nil)
	assert.NoError(t, db.Drop())
}

func TestDatabase_TTL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	set := newShardSet()
	shard1 := NewMockShard(ctrl)
	set.InsertShard(models.ShardID(0), shard1)
	db := &database{
		shardSet: *set,
	}
	shard1.EXPECT().TTL()
	db.TTL()
}

func TestDatabase_EvictSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	set := newShardSet()
	shard1 := NewMockShard(ctrl)
	set.InsertShard(models.ShardID(0), shard1)
	db := &database{
		shardSet: *set,
	}
	shard1.EXPECT().EvictSegment()
	db.EvictSegment()
}

func Benchmark_LoadSyncMap(b *testing.B) {
	var m sync.Map
	for i := 0; i < boundaryShardSetLen; i++ {
		m.Store(i, &shard{})
	}
	// 8.435 ns
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			item, ok := m.Load(boundaryShardSetLen - 1)
			if ok {
				_, _ = item.(*shard)
			}
		}
	})
}

func Benchmark_LoadAtomicValue(b *testing.B) {
	var v atomic.Value
	l := make([]*shard, boundaryShardSetLen)
	v.Store(l)

	// 2.631ns
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			list := v.Load().([]*shard)
			for i := 0; i < boundaryShardSetLen; i++ {
				if i == boundaryShardSetLen-1 {
					_ = list[boundaryShardSetLen-1]
				}
			}
		}
	})
}

func Benchmark_SyncRWMutex(b *testing.B) {
	var lock sync.RWMutex
	m := make(map[int]*shard)
	for i := 0; i < boundaryShardSetLen; i++ {
		m[i] = &shard{}
	}

	// 34.75 ns
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lock.RLock()
			_ = m[boundaryShardSetLen-1]
			lock.RUnlock()
		}
	})
}

func Benchmark_MapWithoutLock(b *testing.B) {
	m := make(map[int]*shard)
	for i := 0; i < boundaryShardSetLen; i++ {
		m[i] = &shard{}
	}
	var v atomic.Value
	v.Store(m)
	// 3.066 ns
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			item := v.Load().(map[int]*shard)
			_ = item[boundaryShardSetLen-1]
		}
	})
}

var (
	boundaryShardSetLen = 20
)

func Benchmark_ShardSet_iterating(b *testing.B) {
	set := newShardSet()
	for i := 0; i < boundaryShardSetLen; i++ {
		set.InsertShard(models.ShardID(i), nil)
	}
	// 2.8ns
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			set.GetShard(models.ShardID(boundaryShardSetLen - 1))
		}
	})
}

func Benchmark_ShardSet_binarySearch(b *testing.B) {
	set := newShardSet()
	for i := 0; i < boundaryShardSetLen+1; i++ {
		set.InsertShard(models.ShardID(i), nil)
	}
	// 4.68ns
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			set.GetShard(models.ShardID(boundaryShardSetLen))
		}
	})
}
