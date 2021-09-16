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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestDatabase_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	tmpDir := t.TempDir()

	defer func() {
		newMetadataFunc = metadb.NewMetadata
		newKVStoreFunc = kv.NewStore
		newShardFunc = newShard
		encodeToml = ltoml.EncodeToml
		ctrl.Finish()
	}()
	// case 1: dump config err
	encodeToml = func(fileName string, v interface{}) error {
		return fmt.Errorf("err")
	}
	db, err := newDatabase("db", tmpDir, &databaseConfig{
		Option: option.DatabaseOption{},
	}, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
	encodeToml = ltoml.EncodeToml
	// case 2: create kv store err
	newKVStoreFunc = func(name string, option kv.StoreOption) (store kv.Store, err error) {
		return nil, fmt.Errorf("err")
	}
	db, err = newDatabase("db", tmpDir, &databaseConfig{
		Option: option.DatabaseOption{},
	}, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
	// case 3: create family err
	kvStore := kv.NewMockStore(ctrl)
	newKVStoreFunc = func(name string, option kv.StoreOption) (store kv.Store, err error) {
		return kvStore, nil
	}
	kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	db, err = newDatabase("db", tmpDir, &databaseConfig{
		Option: option.DatabaseOption{},
	}, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
	// case 4: new metadata err
	kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	newMetadataFunc = func(ctx context.Context, databaseName, parent string,
		tagFamily kv.Family) (metadata metadb.Metadata, err error) {
		return nil, fmt.Errorf("err")
	}
	db, err = newDatabase("db", tmpDir, &databaseConfig{
		Option: option.DatabaseOption{},
	}, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
	// case 5: create shard err
	newMetadataFunc = metadb.NewMetadata
	newShardFunc = func(db Database, shardID models.ShardID, shardPath string, option option.DatabaseOption) (s Shard, err error) {
		return nil, fmt.Errorf("err")
	}
	db, err = newDatabase("db", tmpDir, &databaseConfig{
		ShardIDs: []models.ShardID{1, 2, 3},
		Option:   option.DatabaseOption{},
	}, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
	// case 6: create db success
	newShardFunc = newShard
	db, err = newDatabase("db", tmpDir, &databaseConfig{
		ShardIDs: []models.ShardID{1, 2, 3},
		Option:   option.DatabaseOption{Interval: "10s"},
	}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, db.ExecutorPool())
	assert.Equal(t, option.DatabaseOption{Interval: "10s"}, db.GetOption())
	assert.Equal(t, 3, db.NumOfShards())
	kvStore.EXPECT().Close().Return(nil).AnyTimes() // include shard close
	err = db.Close()
	assert.NoError(t, err)
	// case 7: close metadata err when create db
	metadata := metadb.NewMockMetadata(ctrl)
	newMetadataFunc = func(ctx context.Context, databaseName, parent string, tagFamily kv.Family) (metadb.Metadata, error) {
		return metadata, nil
	}
	newShardFunc = func(db Database, shardID models.ShardID, shardPath string, option option.DatabaseOption) (s Shard, err error) {
		return nil, fmt.Errorf("err")
	}
	metadata.EXPECT().Close().Return(fmt.Errorf("err"))
	db, err = newDatabase("db", tmpDir, &databaseConfig{
		ShardIDs: []models.ShardID{1, 2, 3},
		Option:   option.DatabaseOption{},
	}, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestDatabase_CreateShards(t *testing.T) {
	tmpDir := t.TempDir()

	ctrl := gomock.NewController(t)
	defer func() {
		newShardFunc = newShard
		encodeToml = ltoml.EncodeToml
		ctrl.Finish()
	}()
	db, err := newDatabase("db", tmpDir, &databaseConfig{
		ShardIDs: []models.ShardID{1, 2, 3},
		Option:   option.DatabaseOption{Interval: "10s"},
	}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// case 1: shard ids cannot be empty
	err = db.CreateShards(option.DatabaseOption{}, nil)
	assert.Error(t, err)
	// case 2: create shard err
	newShardFunc = func(db Database, shardID models.ShardID,
		shardPath string, option option.DatabaseOption) (s Shard, err error) {
		return nil, fmt.Errorf("err")
	}
	err = db.CreateShards(option.DatabaseOption{}, []models.ShardID{4, 5, 6})
	assert.Error(t, err)
	// case 3: create exist shard
	err = db.CreateShards(option.DatabaseOption{}, []models.ShardID{1, 2, 3})
	assert.NoError(t, err)
	// case 4: create shard success
	newShardFunc = func(db Database, shardID models.ShardID,
		shardPath string, option option.DatabaseOption) (s Shard, err error) {
		return nil, nil
	}
	err = db.CreateShards(option.DatabaseOption{}, []models.ShardID{4, 5, 6})
	assert.NoError(t, err)
	// case 5: dump option err
	newShardFunc = func(db Database, shardID models.ShardID,
		shardPath string, option option.DatabaseOption) (s Shard, err error) {
		return nil, nil
	}
	encodeToml = func(fileName string, v interface{}) error {
		return fmt.Errorf("err")
	}
	err = db.CreateShards(option.DatabaseOption{}, []models.ShardID{9})
	assert.Error(t, err)
	// case 6: create exist shard
	db1 := db.(*database)
	err = db1.createShard(1, option.DatabaseOption{})
	assert.NoError(t, err)
}

func TestDatabase_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := kv.NewMockStore(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadata.EXPECT().Flush().Return(nil).AnyTimes()
	db := &database{
		metadata:  metadata,
		shardSet:  *newShardSet(),
		metaStore: mockStore}
	// case 1: close metadata err
	metadata.EXPECT().Close().Return(fmt.Errorf("err"))
	err := db.Close()
	assert.Error(t, err)
	// case 2: close meta store err
	metadata.EXPECT().Close().Return(nil).AnyTimes()
	mockStore.EXPECT().Close().Return(fmt.Errorf("err"))
	err = db.Close()
	assert.Error(t, err)

	mockStore.EXPECT().Close().Return(nil)

	// mock shard close error
	mockShard := NewMockShard(ctrl)
	mockShard.EXPECT().Close().Return(fmt.Errorf("error"))
	db.shardSet.InsertShard(models.ShardID(1), mockShard)
	assert.Nil(t, db.Close())
}

func TestDatabase_FlushMeta(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metadata := metadb.NewMockMetadata(ctrl)
	db := &database{
		metadata:   metadata,
		isFlushing: *atomic.NewBool(false)}
	// case 1: flushing
	db.isFlushing.Store(true)
	err := db.FlushMeta()
	assert.NoError(t, err)
	// case 2: need flush meta
	metadata.EXPECT().Flush().Return(nil)
	db.isFlushing.Store(false)
	err = db.FlushMeta()
	assert.NoError(t, err)
}

func TestDatabase_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	checker := NewMockDataFlushChecker(ctrl)

	db, err := newDatabase("db", t.TempDir(), &databaseConfig{
		Option: option.DatabaseOption{},
	}, checker)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db1 := db.(*database)
	shard1 := NewMockShard(ctrl)
	shard2 := NewMockShard(ctrl)
	shard1.EXPECT().Indicator().Return("shard1").AnyTimes()
	shard2.EXPECT().Indicator().Return("shard2").AnyTimes()
	db1.shardSet.InsertShard(1, shard1)
	db1.shardSet.InsertShard(2, shard2)
	checker.EXPECT().requestFlushJob(gomock.Any())
	checker.EXPECT().requestFlushJob(gomock.Any())
	err = db.Flush()
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

func Benchmark_SharSet_binarySearch(b *testing.B) {
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
