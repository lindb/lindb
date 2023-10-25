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

	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	commontimeutil "github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestIntervalSegment_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := NewMockShard(ctrl)
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test-db").AnyTimes()
	shard.EXPECT().Database().Return(db).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "make segment dir err",
			prepare: func() {
				mkDirIfNotExist = func(path string) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create segment successfully",
			prepare: func() {
				mkDirIfNotExist = func(path string) error {
					return nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				mkDirIfNotExist = fileutil.MkDirIfNotExist
				listDir = fileutil.ListDir
				newSegmentFunc = newSegment
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			s, err := newIntervalSegment(shard, option.Interval{Interval: timeutil.Interval(commontimeutil.OneSecond * 10)})
			if ((err != nil) != tt.wantErr && s == nil) || (!tt.wantErr && s == nil) {
				t.Errorf("newIntervalSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntervalSegment_GetOrCreateSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	segment := NewMockSegment(ctrl)

	cases := []struct {
		name        string
		segmentName string
		prepare     func()
		wantErr     bool
	}{
		{
			name:        "get from memory",
			segmentName: "test",
			wantErr:     false,
		},
		{
			name:        "create segment failure",
			segmentName: "test-2",
			prepare: func() {
				newSegmentFunc = func(shard Shard, segmentName string, interval timeutil.Interval) (Segment, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name:        "create segment failure",
			segmentName: "test-2",
			prepare: func() {
				newSegmentFunc = func(shard Shard, segmentName string, interval timeutil.Interval) (Segment, error) {
					return segment, nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newSegmentFunc = newSegment
			}()
			s := &intervalSegment{
				interval: option.Interval{Interval: timeutil.Interval(commontimeutil.OneSecond * 10)},
				segments: map[string]Segment{
					"test": segment,
				},
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			seg, err := s.GetOrCreateSegment(tt.segmentName)
			if ((err != nil) != tt.wantErr && seg == nil) || (!tt.wantErr && seg == nil) {
				t.Errorf("GetOrCreateSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntervalSegment_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segment := NewMockSegment(ctrl)
	s := &intervalSegment{
		interval: option.Interval{Interval: timeutil.Interval(commontimeutil.OneSecond * 10)},
		segments: map[string]Segment{
			"test": segment,
		},
	}
	segment.EXPECT().Close()
	s.Close()
}

func TestIntervalSegment_GetDataFamilies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := commontimeutil.Now() - 4*commontimeutil.OneDay
	segmentDir := commontimeutil.FormatTimestamp(now, "20060102")
	segment := NewMockSegment(ctrl)

	cases := []struct {
		name      string
		timeRange timeutil.TimeRange
		prepare   func(s *intervalSegment)
		len       int
	}{
		{
			name: "list segment path failure",
			prepare: func(s *intervalSegment) {
				listDir = func(path string) ([]string, error) {
					return nil, fmt.Errorf("err")
				}
			},
			len: 0,
		},
		{
			name: "list empty segment path",
			prepare: func(s *intervalSegment) {
				listDir = func(path string) ([]string, error) {
					return nil, nil
				}
			},
			len: 0,
		},
		{
			name: "segment is expired",
			prepare: func(s *intervalSegment) {
				s.interval = option.Interval{
					Interval:  timeutil.Interval(commontimeutil.OneSecond * 10),
					Retention: timeutil.Interval(commontimeutil.OneDay * 2),
				}
				listDir = func(path string) ([]string, error) {
					return []string{segmentDir}, nil
				}
			},
			len: 0,
		},
		{
			name: "time range not match",
			timeRange: timeutil.TimeRange{
				Start: commontimeutil.Now() - 4*commontimeutil.OneHour,
				End:   commontimeutil.Now(),
			},
			prepare: func(s *intervalSegment) {
				listDir = func(path string) ([]string, error) {
					return []string{segmentDir}, nil
				}
			},
			len: 0,
		},
		{
			name: "get segment from memory",
			timeRange: timeutil.TimeRange{
				Start: now - 4*commontimeutil.OneHour,
				End:   now,
			},
			prepare: func(s *intervalSegment) {
				listDir = func(path string) ([]string, error) {
					return []string{segmentDir}, nil
				}
				segment.EXPECT().GetDataFamilies(gomock.Any()).Return([]DataFamily{nil})
			},
			len: 1,
		},
		{
			name: "load segment failure",
			timeRange: timeutil.TimeRange{
				Start: now - 2*commontimeutil.OneHour,
				End:   now,
			},
			prepare: func(s *intervalSegment) {
				listDir = func(path string) ([]string, error) {
					return []string{segmentDir}, nil
				}
				delete(s.segments, segmentDir)
				newSegmentFunc = func(shard Shard, segmentName string, interval timeutil.Interval) (Segment, error) {
					return nil, fmt.Errorf("err")
				}
			},
			len: 0,
		},
		{
			name: "load segment successfully",
			timeRange: timeutil.TimeRange{
				Start: now - 2*commontimeutil.OneHour,
				End:   now,
			},
			prepare: func(s *intervalSegment) {
				listDir = func(path string) ([]string, error) {
					return []string{segmentDir}, nil
				}
				delete(s.segments, segmentDir)
				newSegmentFunc = func(shard Shard, segmentName string, interval timeutil.Interval) (Segment, error) {
					return segment, nil
				}
				segment.EXPECT().GetDataFamilies(gomock.Any()).Return([]DataFamily{nil})
			},
			len: 1,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				listDir = fileutil.ListDir
				newSegmentFunc = newSegment
			}()
			s := &intervalSegment{
				segments: map[string]Segment{
					segmentDir: segment,
				},
				interval: option.Interval{
					Interval:  timeutil.Interval(commontimeutil.OneSecond * 10),
					Retention: timeutil.Interval(commontimeutil.OneDay * 20),
				},
				logger: logger.GetLogger("test", "Segment"),
			}
			if tt.prepare != nil {
				tt.prepare(s)
			}

			families := s.GetDataFamilies(tt.timeRange)
			assert.Len(t, families, tt.len)
		})
	}
}

func TestIntervalSegment_TTL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := commontimeutil.Now() - 40*commontimeutil.OneDay
	segmentDir := commontimeutil.FormatTimestamp(now, "20060102")

	segment := NewMockSegment(ctrl)
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "list segment dir failure",
			prepare: func() {
				listDir = func(path string) ([]string, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "parse segment base time failure",
			prepare: func() {
				listDir = func(path string) ([]string, error) {
					return []string{"abc"}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "remove data dir failure",
			prepare: func() {
				listDir = func(path string) ([]string, error) {
					return []string{segmentDir}, nil
				}
				segment.EXPECT().Close()
				removeDir = func(path string) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: false,
		},
		{
			name: "remove data dir successfully",
			prepare: func() {
				listDir = func(path string) ([]string, error) {
					return []string{segmentDir}, nil
				}
				segment.EXPECT().Close()
				removeDir = func(path string) error {
					return nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				listDir = fileutil.ListDir
			}()
			s := &intervalSegment{
				interval: option.Interval{
					Interval:  timeutil.Interval(10 * commontimeutil.OneSecond),
					Retention: timeutil.Interval(30 * commontimeutil.OneDay),
				},
				segments: map[string]Segment{
					segmentDir: segment,
				},
				logger: logger.GetLogger("TSDB", "segment"),
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			err := s.TTL()
			if (err != nil) != tt.wantErr {
				t.Errorf("TTL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntervalSegment_EvictSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	segment := NewMockSegment(ctrl)
	s := &intervalSegment{
		interval: option.Interval{
			Interval:  timeutil.Interval(10 * commontimeutil.OneSecond),
			Retention: timeutil.Interval(30 * commontimeutil.OneDay),
		},
		segments: map[string]Segment{
			segmentDir: segment,
		},
		logger: logger.GetLogger("TSDB", "Segment"),
	}
	segment.EXPECT().NeedEvict().Return(false)
	s.EvictSegment()
	assert.Len(t, s.segments, 1)

	segment.EXPECT().NeedEvict().Return(true)
	segment.EXPECT().Close()
	s.EvictSegment()
	assert.Len(t, s.segments, 0)
}
