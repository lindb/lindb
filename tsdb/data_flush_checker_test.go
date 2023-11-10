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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/tsdb/memdb"
)

func TestDataFlushChecker_Lifecycle(t *testing.T) {
	checker := newDataFlushChecker(context.TODO())
	checker1 := checker.(*dataFlushChecker)
	assert.False(t, checker1.running.Load())
	checker.Stop() // ignore stop
	assert.False(t, checker1.running.Load())
	checker.Start()
	assert.True(t, checker1.running.Load())
	checker.Start() // dup
	time.Sleep(50 * time.Millisecond)
	assert.True(t, checker1.running.Load())
	checker.Stop()
	assert.False(t, checker1.running.Load())
}

func TestDataFlushCheck_startCheckDataFlush(t *testing.T) {
	defer memoryUsageCheckInterval.Store(time.Minute)

	memoryUsageCheckInterval.Store(10 * time.Millisecond)
	checker := newDataFlushChecker(context.TODO())
	checker1 := checker.(*dataFlushChecker)
	checker1.running.Store(true)
	go func() {
		time.Sleep(50 * time.Millisecond)
		checker1.Stop()
	}()
	checker1.startCheckDataFlush()
}

func TestDataFamilyCheck_check(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		memoryUsageCheckInterval.Store(time.Minute)
		ctrl.Finish()
	}()
	shard := NewMockShard(ctrl)
	bufferMgr := memdb.NewMockBufferManager(ctrl)
	bufferMgr.EXPECT().GarbageCollect().AnyTimes()
	shard.EXPECT().BufferManager().Return(bufferMgr).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	db := NewMockDatabase(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	db.EXPECT().Name().Return("db").AnyTimes()

	family1 := NewMockDataFamily(ctrl)
	family2 := NewMockDataFamily(ctrl)
	family1.EXPECT().Indicator().Return("family1").AnyTimes()
	family1.EXPECT().Shard().Return(shard).AnyTimes()
	family2.EXPECT().Indicator().Return("family2").AnyTimes()
	family2.EXPECT().Shard().Return(shard).AnyTimes()

	cases := []struct {
		name    string
		prepare func(c *dataFlushChecker)
		assert  func(c *dataFlushChecker)
	}{
		{
			name: "no family",
			assert: func(c *dataFlushChecker) {
				v, ok := c.dbInFlushing.Load("db")
				assert.False(t, ok)
				assert.Nil(t, v)
			},
		},
		{
			name: "no family need flush",
			prepare: func(_ *dataFlushChecker) {
				GetFamilyManager().AddFamily(family1)
				family1.EXPECT().NeedFlush().Return(false)
			},
			assert: func(c *dataFlushChecker) {
				v, ok := c.dbInFlushing.Load("db")
				assert.False(t, ok)
				assert.Nil(t, v)
			},
		},
		{
			name: "family need flush",
			prepare: func(_ *dataFlushChecker) {
				GetFamilyManager().AddFamily(family2)
				family2.EXPECT().NeedFlush().Return(true)
				GetFamilyManager().AddFamily(family1)
				family1.EXPECT().NeedFlush().Return(true)
			},
			assert: func(c *dataFlushChecker) {
				v, ok := c.dbInFlushing.Load("db")
				assert.True(t, ok)
				assert.Equal(t, &flushRequest{
					db: db,
					shards: map[models.ShardID]*flushShard{
						models.ShardID(1): {
							families: []DataFamily{family1, family2},
							shard:    shard,
						},
					},
					global: false,
				}, v)
			},
		},
		{
			name: "pick family for Global memory limit",
			prepare: func(_ *dataFlushChecker) {
				GetFamilyManager().AddFamily(family2)
				family2.EXPECT().NeedFlush().Return(false)
				family2.EXPECT().MemDBSize().Return(int64(2 * ignoreMemorySize))
				family2.EXPECT().IsFlushing().Return(false)
				GetFamilyManager().AddFamily(family1)
				family1.EXPECT().NeedFlush().Return(false)
				family1.EXPECT().MemDBSize().Return(int64(199 * ignoreMemorySize))
				family1.EXPECT().IsFlushing().Return(false)
				GetFamilyManager().AddFamily(family1)
				cfg := config.GlobalStorageConfig()
				cfg.TSDB.MaxMemUsageBeforeFlush = 0.0001
				config.SetGlobalStorageConfig(cfg)
			},
			assert: func(c *dataFlushChecker) {
				v, ok := c.dbInFlushing.Load("db")
				assert.True(t, ok)
				assert.Equal(t, &flushRequest{
					db: db,
					shards: map[models.ShardID]*flushShard{
						models.ShardID(1): {
							families: []DataFamily{family2},
							shard:    shard,
						},
					},
					global: true,
				}, v)
			},
		},
		{
			name: "pick family for Global memory limit, but no match family",
			prepare: func(_ *dataFlushChecker) {
				GetFamilyManager().AddFamily(family2)
				family2.EXPECT().NeedFlush().Return(false)
				family2.EXPECT().MemDBSize().Return(int64(199))
				family2.EXPECT().IsFlushing().Return(false)
				GetFamilyManager().AddFamily(family1)
				family1.EXPECT().NeedFlush().Return(false)
				family1.EXPECT().IsFlushing().Return(true)
				GetFamilyManager().AddFamily(family1)
				cfg := config.GlobalStorageConfig()
				cfg.TSDB.MaxMemUsageBeforeFlush = 0.0001
				config.SetGlobalStorageConfig(cfg)
			},
			assert: func(c *dataFlushChecker) {
				v, ok := c.dbInFlushing.Load("db")
				assert.False(t, ok)
				assert.Nil(t, v)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			memoryUsageCheckInterval.Store(10 * time.Millisecond)
			checker := newDataFlushChecker(context.TODO())
			checker1 := checker.(*dataFlushChecker)
			checker1.running.Store(true)
			defer func() {
				GetFamilyManager().RemoveFamily(family1)
				GetFamilyManager().RemoveFamily(family2)
				config.SetGlobalStorageConfig(config.NewDefaultStorageBase())
				checker1.dbInFlushing.Delete("db")
			}()
			cfg := config.GlobalStorageConfig()
			cfg.TSDB.MaxMemUsageBeforeFlush = 1.0

			if tt.prepare != nil {
				tt.prepare(checker1)
			}
			go func() {
				for range checker1.flushRequestCh {
				}
			}()
			checker1.check()
			if tt.assert != nil {
				tt.assert(checker1)
			}
		})
	}
}

func TestDataFlushChecker_requestFlush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("db").AnyTimes()
	cases := []struct {
		name    string
		prepare func(c *dataFlushChecker)
		assert  func(c *dataFlushChecker)
	}{
		{
			name: "not running",
			prepare: func(c *dataFlushChecker) {
				c.running.Store(false)
			},
			assert: func(c *dataFlushChecker) {
				_, ok := c.dbInFlushing.Load("db")
				assert.False(t, ok)
			},
		},
		{
			name: "db is flushing",
			prepare: func(c *dataFlushChecker) {
				c.dbInFlushing.Store("db", &flushRequest{})
			},
			assert: func(c *dataFlushChecker) {
				_, ok := c.dbInFlushing.Load("db")
				assert.True(t, ok)
			},
		},
		{
			name: "checker is stopped",
			prepare: func(c *dataFlushChecker) {
				c.flushRequestCh = make(chan *flushRequest)
				c.Stop()
				c.running.Store(true)
			},
			assert: func(c *dataFlushChecker) {
				_, ok := c.dbInFlushing.Load("db")
				assert.False(t, ok)
			},
		},
		{
			name: "request flush successfully",
			prepare: func(c *dataFlushChecker) {
				go func() {
					select {
					case <-c.ctx.Done():
						return
					case <-c.flushRequestCh:
						return
					}
				}()
			},
			assert: func(c *dataFlushChecker) {
				_, ok := c.dbInFlushing.Load("db")
				assert.True(t, ok)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			checker := newDataFlushChecker(context.TODO())
			checker1 := checker.(*dataFlushChecker)
			checker1.running.Store(true)
			if tt.prepare != nil {
				tt.prepare(checker1)
			}
			checker1.requestFlushJob(&flushRequest{
				db: db,
			})
			tt.assert(checker1)
		})
	}
}

func TestDataFlushChecker_flushWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	db.EXPECT().FlushMeta().Return(fmt.Errorf("err"))
	checker := newDataFlushChecker(context.TODO())
	checker1 := checker.(*dataFlushChecker)
	checker1.running.Store(true)
	ch := make(chan struct{})
	go func() {
		time.Sleep(100 * time.Millisecond)
		checker.Stop()
		ch <- struct{}{}
	}()
	go func() {
		checker1.flushWorker()
	}()
	checker1.flushRequestCh <- &flushRequest{
		db: db,
	}
	<-ch
}

func TestDataFlushChecker_doFlush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	shard := NewMockShard(ctrl)
	family := NewMockDataFamily(ctrl)
	family.EXPECT().Indicator().Return("family").AnyTimes()
	bufMgr := memdb.NewMockBufferManager(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	shard.EXPECT().BufferManager().Return(bufMgr).AnyTimes()
	shard.EXPECT().Indicator().Return("shard").AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	bufMgr.EXPECT().GarbageCollect().AnyTimes()

	cases := []struct {
		name    string
		prepare func(c *dataFlushChecker)
	}{
		{
			name: "flush meta db failure",
			prepare: func(_ *dataFlushChecker) {
				db.EXPECT().FlushMeta().Return(fmt.Errorf("err"))
			},
		},
		{
			name: "flush index db failure",
			prepare: func(_ *dataFlushChecker) {
				db.EXPECT().FlushMeta().Return(nil)
				db.EXPECT().WaitFlushMetaCompleted()
				shard.EXPECT().FlushIndex().Return(fmt.Errorf("err"))
			},
		},
		{
			name: "flush family failure",
			prepare: func(_ *dataFlushChecker) {
				db.EXPECT().FlushMeta().Return(nil)
				db.EXPECT().WaitFlushMetaCompleted()
				shard.EXPECT().FlushIndex().Return(nil)
				shard.EXPECT().WaitFlushIndexCompleted()
				family.EXPECT().Flush().Return(fmt.Errorf("err"))
			},
		},
		{
			name: "flush family successfully",
			prepare: func(_ *dataFlushChecker) {
				db.EXPECT().FlushMeta().Return(nil)
				db.EXPECT().WaitFlushMetaCompleted()
				shard.EXPECT().FlushIndex().Return(nil)
				shard.EXPECT().WaitFlushIndexCompleted()
				family.EXPECT().Flush().Return(nil)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			checker := newDataFlushChecker(context.TODO())
			checker1 := checker.(*dataFlushChecker)
			checker1.running.Store(true)
			if tt.prepare != nil {
				tt.prepare(checker1)
			}
			checker1.doFlush(&flushRequest{
				db: db,
				shards: map[models.ShardID]*flushShard{
					1: {
						shard:    shard,
						families: []DataFamily{family},
					},
				},
				global: true,
			})
		})
	}
}
