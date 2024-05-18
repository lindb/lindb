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
	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	commontimeutil "github.com/lindb/common/pkg/timeutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"

	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSegmentPartition_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := NewMockShard(ctrl)
	intervals := option.Intervals{option.Interval{
		Interval: timeutil.Interval(commontimeutil.OneMinute),
	}}

	sp := newSegmentPartition(shard, intervals)
	p := sp.(*segmentPartition)
	assert.NotNil(t, p)
	assert.Equal(t, shard, p.shard)
	assert.Equal(t, intervals, p.intervals)
}

func TestSegmentPartition_GetOrCreateSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		getMonthTimestampFunc = timeutil.GetMonthTimestamp
		newSegmentFunc = newSegment
		ctrl.Finish()
	}()

	segment := NewMockSegment(ctrl)
	segment.EXPECT().GetTimestamp().Return(int64(1)).AnyTimes()
	segment2 := NewMockSegment(ctrl)
	segment2.EXPECT().GetTimestamp().Return(int64(2)).AnyTimes()

	sp := segmentPartition{}

	cases := []struct {
		name       string
		familyTime int64
		prepare    func()
		segmentLen int
		wantErr    bool
	}{
		{
			name:       "GetOrCreateSegment err",
			familyTime: 1,
			prepare: func() {
				newSegmentFunc = func(shard Shard, timestamp int64, intervals option.Intervals) (Segment, error) {
					return nil, fmt.Errorf("err")
				}
			},
			segmentLen: 0,
			wantErr:    true,
		},
		{
			name:       "GetOrCreateSegment timestamp already exist",
			familyTime: 1,
			prepare: func() {
				getMonthTimestampFunc = func(familyTime int64) int64 {
					return 1
				}
				sp.segments = []Segment{segment}
			},
			segmentLen: 1,
		},
		{
			name:       "successfully",
			familyTime: 2,
			prepare: func() {
				getMonthTimestampFunc = func(familyTime int64) int64 {
					return 2
				}
				newSegmentFunc = func(shard Shard, timestamp int64, intervals option.Intervals) (Segment, error) {
					return segment2, nil
				}
			},
			segmentLen: 2,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			_, err := sp.GetOrCreateSegment(tt.familyTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrCreateSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, len(sp.segments), tt.segmentLen)
		})
	}
}

func TestSegmentPartition_GetSegments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer func() {

	}()
	segment := NewMockSegment(ctrl)
	sp := segmentPartition{
		segments: []Segment{segment},
	}
	assert.Equal(t, 1, len(sp.segments))
}

func TestSegmentPartition_Recovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	segment0 := NewMockSegment(ctrl)
	segment0.EXPECT().GetTimestamp().Return(int64(2)).AnyTimes()
	segment := NewMockSegment(ctrl)
	segment.EXPECT().GetTimestamp().Return(int64(1)).AnyTimes()
	database := NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("database name").AnyTimes()
	shard := NewMockShard(ctrl)
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()

	sp := segmentPartition{
		shard: shard,
	}

	cases := []struct {
		name       string
		prepare    func()
		segmentLen int
		wantErr    bool
	}{
		{
			name: "file not exists",
			prepare: func() {
				fileExist = func(file string) bool {
					return false
				}
			},
			segmentLen: 0,
			wantErr:    false,
		},
		{
			name: "listDir err",
			prepare: func() {
				fileExist = func(file string) bool {
					return true
				}
				listDirName = func(path string) ([]string, error) {
					return nil, fmt.Errorf("err")
				}
			},
			segmentLen: 0,
			wantErr:    true,
		},
		{
			name: "parseTimestamp err",
			prepare: func() {
				fileExist = func(file string) bool {
					return true
				}
				listDirName = func(path string) ([]string, error) {
					return []string{""}, nil
				}
				parseTimestamp = func(timestampStr string, layout ...string) (int64, error) {
					return 0, fmt.Errorf("err")
				}
			},
			segmentLen: 0,
			wantErr:    true,
		},
		{
			name: "new segment err",
			prepare: func() {
				fileExist = func(file string) bool {
					return true
				}
				listDirName = func(path string) ([]string, error) {
					return []string{""}, nil
				}
				parseTimestamp = func(timestampStr string, layout ...string) (int64, error) {
					return 1, nil
				}
				newSegmentFunc = func(shard Shard, timestamp int64, intervals option.Intervals) (Segment, error) {
					return nil, fmt.Errorf("err")
				}
			},
			segmentLen: 0,
			wantErr:    true,
		},
		{
			name: "new segment no err",
			prepare: func() {
				sp.segments = []Segment{segment0}
				fileExist = func(file string) bool {
					return true
				}
				listDirName = func(path string) ([]string, error) {
					return []string{""}, nil
				}
				parseTimestamp = func(timestampStr string, layout ...string) (int64, error) {
					return 1, nil
				}
				newSegmentFunc = func(shard Shard, timestamp int64, intervals option.Intervals) (Segment, error) {
					return segment, nil
				}
			},
			segmentLen: 2,
			wantErr:    false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				fileExist = fileutil.Exist
				listDirName = fileutil.ListDir
				parseTimestamp = commontimeutil.ParseTimestamp
				newSegmentFunc = newSegment
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			err := sp.Recover()
			if (err != nil) != tt.wantErr {
				t.Errorf("Recovery() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, len(sp.GetSegments()), tt.segmentLen)
		})
	}
}

func TestSegmentPartition_FlushIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	database := NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test").AnyTimes()
	shard := NewMockShard(ctrl)
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	segment := NewMockSegment(ctrl)
	segment.EXPECT().GetName().Return("test").AnyTimes()

	s := &segmentPartition{
		shard:          shard,
		segments:       []Segment{segment},
		isFlushing:     *atomic.NewBool(false),
		flushCondition: sync.NewCond(&sync.Mutex{}),
		logger:         logger.GetLogger("TSDB", "Test"),
	}

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "flush is doing",
			prepare: func() {
				s.isFlushing.Store(true)
			},
		},
		{
			name: "flush index db err",
			prepare: func() {
				segment.EXPECT().FlushIndex().Return(fmt.Errorf("err"))
			},
			wantErr: false,
		},
		{
			name: "flush successfully",
			prepare: func() {
				segment.EXPECT().FlushIndex().Return(nil)
			},
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				s.isFlushing.Store(false)
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			if err := s.FlushIndex(); (err != nil) != tt.wantErr {
				t.Errorf("FlushIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSegmentPartition_WaitFlushIndexCompleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := commontimeutil.Now()

	database := NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test").AnyTimes()
	shard := NewMockShard(ctrl)
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	segment := NewMockSegment(ctrl)
	segment.EXPECT().GetName().Return("test").AnyTimes()

	s := &segmentPartition{
		shard:          shard,
		segments:       []Segment{segment},
		isFlushing:     *atomic.NewBool(false),
		flushCondition: sync.NewCond(&sync.Mutex{}),
		logger:         logger.GetLogger("TSDB", "Test"),
	}
	s.isFlushing.Store(false)
	segment.EXPECT().FlushIndex().DoAndReturn(func() error {
		// simulate flush time
		time.Sleep(90 * time.Millisecond)
		return nil
	})
	var wait sync.WaitGroup
	wait.Add(2)
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
		err := s.FlushIndex()
		assert.NoError(t, err)
	}()
	<-ch
	time.Sleep(10 * time.Millisecond)
	go func() {
		s.WaitFlushIndexCompleted()
		wait.Done()
	}()
	go func() {
		s.WaitFlushIndexCompleted()
		wait.Done()
	}()
	wait.Wait()
	assert.True(t, commontimeutil.Now()-now >= 90*time.Millisecond.Milliseconds())
}

func TestSegmentPartition_TTL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	database := NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test").AnyTimes()
	shard := NewMockShard(ctrl)
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	segment := NewMockSegment(ctrl)
	segment.EXPECT().GetName().Return("test").AnyTimes()

	s := &segmentPartition{
		shard:  shard,
		logger: logger.GetLogger("TSDB", "Test"),
	}
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "segments nil",
			prepare: func() {
			},
		},
		{
			name: "ttl err",
			prepare: func() {
				s.segments = []Segment{segment}
				segment.EXPECT().TTL().Return(fmt.Errorf("err"))
			},
		},
		{
			name: "ttl no err",
			prepare: func() {
				segment.EXPECT().TTL().Return(nil)
			},
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				s.isFlushing.Store(false)
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			if err := s.TTL(); (err != nil) != tt.wantErr {
				t.Errorf("TTL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSegmentPartition_EvictSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	segment := NewMockSegment(ctrl)
	s := &segmentPartition{
		segments: []Segment{segment},
		logger:   logger.GetLogger("TSDB", "Test"),
	}
	segment.EXPECT().EvictSegment()
	s.EvictSegment()

	s.segments = nil
	s.EvictSegment()
}

func TestSegmentPartition_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	database := NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test").AnyTimes()
	shard := NewMockShard(ctrl)
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	segment := NewMockSegment(ctrl)
	segment.EXPECT().GetName().Return("test").AnyTimes()

	s := &segmentPartition{
		shard:    shard,
		segments: []Segment{segment},
		logger:   logger.GetLogger("TSDB", "Test"),
	}
	segment.EXPECT().Close().Return(fmt.Errorf("err"))
	err := s.Close()
	assert.NoError(t, err)

	segment.EXPECT().Close().Return(nil)
	err = s.Close()
	assert.NoError(t, err)

	s.segments = nil
	err = s.Close()
	assert.NoError(t, err)
}

func TestSegmentPartition_currentDo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &segmentPartition{}
	s.concurrentDo(nil)
}
