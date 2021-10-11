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

	"github.com/golang/mock/gomock"
	"github.com/shirou/gopsutil/mem"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/tsdb/memdb"
)

func TestDataFlushChecker_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		memoryUsageCheckInterval.Store(time.Second)
		ctrl.Finish()
	}()
	shard := NewMockShard(ctrl)
	bufferMgr := memdb.NewMockBufferManager(ctrl)
	bufferMgr.EXPECT().GarbageCollect().AnyTimes()
	shard.EXPECT().BufferManager().Return(bufferMgr).AnyTimes()

	shard.EXPECT().Indicator().Return("shard").AnyTimes()
	family1 := NewMockDataFamily(ctrl)
	family1.EXPECT().Indicator().Return("family1").AnyTimes()
	family1.EXPECT().IsFlushing().Return(true).AnyTimes()
	family2 := NewMockDataFamily(ctrl)
	family2.EXPECT().Indicator().Return("family2").AnyTimes()
	family2.EXPECT().IsFlushing().Return(true).AnyTimes()
	GetFamilyManager().AddFamily(family1)
	defer GetFamilyManager().RemoveFamily(family1)
	GetFamilyManager().AddFamily(family2)
	defer GetFamilyManager().RemoveFamily(family2)

	memoryUsageCheckInterval.Store(10 * time.Millisecond)
	checker := newDataFlushChecker(context.TODO())
	family1.EXPECT().NeedFlush().Return(true)
	family2.EXPECT().NeedFlush().Return(true)
	shard.EXPECT().Flush().Return(nil)
	family1.EXPECT().Shard().Return(shard)
	family2.EXPECT().Shard().Return(shard)
	family1.EXPECT().Flush().Return(fmt.Errorf("err"))
	family2.EXPECT().Flush().Return(fmt.Errorf("err"))
	family1.EXPECT().NeedFlush().Return(false).AnyTimes()
	family2.EXPECT().NeedFlush().Return(false).AnyTimes()
	checker.Start()

	time.Sleep(100 * time.Millisecond)
	checker.Stop()
}

func TestDataFlushChecker_check_high_memory_waterMark(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		memoryUsageCheckInterval.Store(time.Second)
		ctrl.Finish()
	}()
	shard := NewMockShard(ctrl)
	bufferMgr := memdb.NewMockBufferManager(ctrl)
	bufferMgr.EXPECT().GarbageCollect().AnyTimes()
	shard.EXPECT().BufferManager().Return(bufferMgr).AnyTimes()

	shard.EXPECT().Indicator().Return("shard").AnyTimes()
	shard.EXPECT().Flush().Return(nil).AnyTimes()
	// case 1: family is flushing
	family := NewMockDataFamily(ctrl)
	family.EXPECT().NeedFlush().Return(false).AnyTimes()
	family.EXPECT().Indicator().Return("family1").AnyTimes()
	family.EXPECT().IsFlushing().Return(true).AnyTimes()
	GetFamilyManager().AddFamily(family)
	defer GetFamilyManager().RemoveFamily(family)
	memoryUsageCheckInterval.Store(10 * time.Millisecond)
	checker := newDataFlushChecker(context.TODO())
	check := checker.(*dataFlushChecker)
	check.memoryStatGetterFunc = func() (stat *mem.VirtualMemoryStat, err error) {
		return &mem.VirtualMemoryStat{UsedPercent: config.GlobalStorageConfig().TSDB.MaxMemUsageBeforeFlush + 0.1}, nil
	}
	checker.Start()

	time.Sleep(100 * time.Millisecond)
	checker.Stop()

	// case 2: pick biggest family data
	family1 := NewMockDataFamily(ctrl)
	family1.EXPECT().NeedFlush().Return(false).AnyTimes()
	family1.EXPECT().Indicator().Return("family").AnyTimes()
	family1.EXPECT().IsFlushing().Return(false).AnyTimes()
	family1.EXPECT().MemDBSize().Return(int64(100)).AnyTimes()
	GetFamilyManager().AddFamily(family1)
	defer GetFamilyManager().RemoveFamily(family1)

	family2 := NewMockDataFamily(ctrl)
	family2.EXPECT().NeedFlush().Return(false).AnyTimes()
	family2.EXPECT().Indicator().Return("family2").AnyTimes()
	family2.EXPECT().IsFlushing().Return(false).AnyTimes()
	family2.EXPECT().Flush().Return(nil).AnyTimes()
	family2.EXPECT().MemDBSize().Return(int64(ignoreMemorySize) + 100).AnyTimes()
	family2.EXPECT().Shard().Return(shard).AnyTimes()
	GetFamilyManager().AddFamily(family2)
	defer GetFamilyManager().RemoveFamily(family2)

	memoryUsageCheckInterval.Store(10 * time.Millisecond)
	checker = newDataFlushChecker(context.TODO())
	check = checker.(*dataFlushChecker)
	check.memoryStatGetterFunc = func() (stat *mem.VirtualMemoryStat, err error) {
		return &mem.VirtualMemoryStat{
			UsedPercent: config.GlobalStorageConfig().TSDB.MaxMemUsageBeforeFlush + 0.1}, nil
	}
	checker.Start()

	time.Sleep(100 * time.Millisecond)
	checker.Stop()
}

func TestDataFlushChecker_requestFlush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		memoryUsageCheckInterval.Store(time.Second)
		ctrl.Finish()
	}()
	shard := NewMockShard(ctrl)
	bufferMgr := memdb.NewMockBufferManager(ctrl)
	bufferMgr.EXPECT().GarbageCollect().AnyTimes()
	shard.EXPECT().BufferManager().Return(bufferMgr).AnyTimes()

	shard.EXPECT().Indicator().Return("shard").AnyTimes()
	shard.EXPECT().Flush().Return(nil).AnyTimes()
	var families []DataFamily
	for i := 0; i < 2; i++ {
		family := NewMockDataFamily(ctrl)
		family.EXPECT().Shard().Return(shard).AnyTimes()
		family.EXPECT().NeedFlush().Return(true).AnyTimes()
		family.EXPECT().Indicator().Return(fmt.Sprintf("family-%d", i)).AnyTimes()
		family.EXPECT().Flush().DoAndReturn(func() error {
			time.Sleep(200 * time.Millisecond)
			return fmt.Errorf("err")
		}).AnyTimes()
		GetFamilyManager().AddFamily(family)
		families = append(families, family)
	}
	defer func() {
		for _, family := range families {
			GetFamilyManager().RemoveFamily(family)
		}
	}()
	memoryUsageCheckInterval.Store(10 * time.Millisecond)
	checker := newDataFlushChecker(context.TODO())
	check := checker.(*dataFlushChecker)
	check.memoryStatGetterFunc = func() (stat *mem.VirtualMemoryStat, err error) {
		return &mem.VirtualMemoryStat{
			UsedPercent: config.GlobalStorageConfig().TSDB.MaxMemUsageBeforeFlush + 0.1}, nil
	}
	checker.Start()

	time.Sleep(100 * time.Millisecond)

	checker1 := checker.(*dataFlushChecker)
	family := NewMockDataFamily(ctrl)
	family.EXPECT().Indicator().Return("family").MaxTimes(3)
	checker1.requestFlushJob(&flushRequest{
		shard:    shard,
		families: []DataFamily{family},
		global:   false,
	}) // request success
	checker1.requestFlushJob(&flushRequest{
		shard:    shard,
		families: []DataFamily{family},
		global:   true,
	}) // reject, because has pending flush job

	checker.Stop()
	checker1.requestFlushJob(&flushRequest{
		shard:    shard,
		families: []DataFamily{family},
		global:   false,
	}) // reject, because not running
	time.Sleep(100 * time.Millisecond)
}

func TestDataFlushChecker_cancel_ctx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	shard := NewMockShard(ctrl)
	shard.EXPECT().Indicator().Return("shard").AnyTimes()
	checker := newDataFlushChecker(context.TODO())
	check := checker.(*dataFlushChecker)
	go func() {
		check.requestFlushJob(&flushRequest{shard: shard})
	}()
	checker.Stop()
}
