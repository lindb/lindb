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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/common/pkg/logger"
	commontimeutil "github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestSegment_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		kv.InitStoreManager(nil)
		ctrl.Finish()
	}()

	storeMgr := kv.NewMockStoreManager(ctrl)
	kv.InitStoreManager(storeMgr)
	store := kv.NewMockStore(ctrl)

	database := NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test").AnyTimes()
	shard := NewMockShard(ctrl)
	interval := timeutil.Interval(commontimeutil.OneSecond * 10)
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	shard.EXPECT().CurrentInterval().Return(interval).AnyTimes()
	segmentName := "20190904"
	cases := []struct {
		name        string
		segmentName string
		prepare     func()
		wantErr     bool
	}{
		{
			name:        "create segment successfully",
			segmentName: segmentName,
			prepare: func() {
				database.EXPECT().GetOption().Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: interval}}})
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(store, nil)
			},
		},
		{
			name:        "create segment successfully and set rollup",
			segmentName: segmentName,
			prepare: func() {
				database.EXPECT().GetOption().Return(&option.DatabaseOption{
					Intervals: option.Intervals{
						{Interval: interval},
						{Interval: timeutil.Interval(5 * commontimeutil.OneMinute)},
					},
				})
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(store, nil)
			},
		},
		{
			name:        "parse segment name err",
			segmentName: "xx",
			wantErr:     true,
		},
		{
			name:        "create store err",
			segmentName: segmentName,
			prepare: func() {
				database.EXPECT().GetOption().Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: interval}}})
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newDataFamilyFunc = newDataFamily
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			s, err := newSegment(shard, tt.segmentName, interval)
			if ((err != nil) != tt.wantErr && s == nil) || (!tt.wantErr && s == nil) {
				t.Errorf("newSegment() error = %v, wantErr %v", err, tt.wantErr)
			}

			if s != nil {
				// check base time
				now, _ := commontimeutil.ParseTimestamp("20190904 00:00:00", "20060102 15:04:05")
				assert.Equal(t, now, s.BaseTime())
			}
		})
	}
}

func TestSegment_GetOrCreateDataFamily(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := kv.NewMockStore(ctrl)
	interval := timeutil.Interval(10 * 1000)
	baseTime, _ := commontimeutil.ParseTimestamp("20190904 00:00:00", "20060102 15:04:05")
	cases := []struct {
		name      string
		timestamp string
		prepare   func(seg *segment)
		wantErr   bool
	}{
		{
			name:      "segment time != base time",
			timestamp: "20190905 19:10:48",
			wantErr:   true,
		},
		{
			name:      "create new family",
			timestamp: "20190904 19:10:48",
			prepare: func(_ *segment) {
				newDataFamilyFunc = func(shard Shard, _ Segment,
					interval timeutil.Interval, timeRange timeutil.TimeRange,
					familyTime int64, family kv.Family) DataFamily {
					return NewMockDataFamily(ctrl)
				}
				store.EXPECT().GetFamily(gomock.Any()).Return(nil)
				store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:      "family exist in kv store",
			timestamp: "20190904 19:10:48",
			prepare: func(_ *segment) {
				newDataFamilyFunc = func(shard Shard, _ Segment,
					interval timeutil.Interval, timeRange timeutil.TimeRange,
					familyTime int64, family kv.Family) DataFamily {
					return NewMockDataFamily(ctrl)
				}
				store.EXPECT().GetFamily(gomock.Any()).Return(kv.NewMockFamily(ctrl))
			},
		},
		{
			name:      "create new family err",
			timestamp: "20190904 20:10:48",
			prepare: func(_ *segment) {
				store.EXPECT().GetFamily(gomock.Any()).Return(nil)
				store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "get exist family",
			timestamp: "20190904 22:10:48",
			prepare: func(seg *segment) {
				now, _ := commontimeutil.ParseTimestamp("20190904 22:10:48", "20060102 15:04:05")
				familyTime := interval.Calculator().CalcFamily(now, seg.baseTime)
				seg.mutex.Lock()
				seg.families[familyTime] = NewMockDataFamily(ctrl)
				seg.mutex.Unlock()
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newDataFamilyFunc = newDataFamily
			}()
			seg := &segment{
				baseTime: baseTime,
				kvStore:  store,
				interval: interval,
				families: make(map[int]DataFamily),
			}
			if tt.prepare != nil {
				tt.prepare(seg)
			}
			now, _ := commontimeutil.ParseTimestamp(tt.timestamp, "20060102 15:04:05")
			dataFamily, err := seg.GetOrCreateDataFamily(now)
			if ((err != nil) != tt.wantErr && dataFamily == nil) || (!tt.wantErr && dataFamily == nil) {
				t.Errorf("GetOrCreateDataFamily() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSegment_GetDataFamilies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := kv.NewMockStore(ctrl)
	now, _ := commontimeutil.ParseTimestamp("20220326 10:30:00", "20060102 15:04:05")
	timeRange := timeutil.TimeRange{
		Start: now - commontimeutil.OneHour,
		End:   now,
	}
	cases := []struct {
		name      string
		timeRange timeutil.TimeRange
		prepare   func(seg *segment)
		len       int
	}{
		{
			name: "family empty",
			prepare: func(_ *segment) {
				store.EXPECT().ListFamilyNames().Return(nil)
			},
			len: 0,
		},
		{
			name: "parse family name failure",
			prepare: func(_ *segment) {
				store.EXPECT().ListFamilyNames().Return([]string{"a"})
			},
			len: 0,
		},
		{
			name:      "get family from memory",
			timeRange: timeRange,
			prepare: func(seq *segment) {
				family := NewMockDataFamily(ctrl)
				seq.families[10] = family
				store.EXPECT().ListFamilyNames().Return([]string{"10"})
				family.EXPECT().TimeRange().Return(timeRange)
			},
			len: 1,
		},
		{
			name:      "get family from storage",
			timeRange: timeRange,
			prepare: func(_ *segment) {
				dataFamily := NewMockDataFamily(ctrl)
				family := kv.NewMockFamily(ctrl)
				store.EXPECT().GetFamily(gomock.Any()).Return(family)
				newDataFamilyFunc = func(shard Shard, _ Segment, interval timeutil.Interval,
					timeRange timeutil.TimeRange, familyTime int64, family kv.Family) DataFamily {
					return dataFamily
				}
				store.EXPECT().ListFamilyNames().Return([]string{"10"})
				dataFamily.EXPECT().TimeRange().Return(timeRange)
			},
			len: 1,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := &segment{
				kvStore:  store,
				interval: timeutil.Interval(10 * 1000),
				families: make(map[int]DataFamily),
			}
			if tt.prepare != nil {
				tt.prepare(s)
			}
			families := s.GetDataFamilies(tt.timeRange)
			assert.Len(t, families, tt.len)
		})
	}
}

func TestSegment_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		kv.InitStoreManager(nil)
		ctrl.Finish()
	}()
	storeMgr := kv.NewMockStoreManager(ctrl)
	kv.InitStoreManager(storeMgr)
	store := kv.NewMockStore(ctrl)
	seg := &segment{
		kvStore:  store,
		families: make(map[int]DataFamily),
		logger:   logger.GetLogger("TSDB", "Test"),
	}
	assertFamilyEmpty := func() {
		seg.mutex.Lock()
		assert.Empty(t, seg.families)
		seg.mutex.Unlock()
	}
	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "no family",
			prepare: func() {
				gomock.InOrder(
					store.EXPECT().Name().Return("test"),
					storeMgr.EXPECT().CloseStore(gomock.Any()).Return(nil),
				)
			},
		},
		{
			name: "close kv store err",
			prepare: func() {
				gomock.InOrder(
					store.EXPECT().Name().Return("test"),
					storeMgr.EXPECT().CloseStore(gomock.Any()).Return(fmt.Errorf("err")),
				)
			},
		},
		{
			name: "close family err",
			prepare: func() {
				family := NewMockDataFamily(ctrl)
				seg.mutex.Lock()
				seg.families[1] = family
				assert.Len(t, seg.families, 1)
				seg.mutex.Unlock()
				gomock.InOrder(
					family.EXPECT().Close().Return(fmt.Errorf("err")),
					family.EXPECT().Indicator().Return("family"),
					store.EXPECT().Name().Return("test"),
					storeMgr.EXPECT().CloseStore(gomock.Any()).Return(fmt.Errorf("err")),
				)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			seg.Close()
			assertFamilyEmpty()
		})
	}
}

func TestSegment_NeedEvict(t *testing.T) {
	interval := timeutil.Interval(10 * 1000)
	s := &segment{interval: interval}
	assert.True(t, s.NeedEvict())
	s.EvictFamily(commontimeutil.Now())
}
