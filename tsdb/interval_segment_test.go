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

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
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
			name: "list segment dir err",
			prepare: func() {
				listDir = func(path string) ([]string, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create segment err",
			prepare: func() {
				listDir = func(path string) ([]string, error) {
					return []string{"20190707"}, nil
				}
				newSegmentFunc = func(shard Shard, segmentName string, interval timeutil.Interval) (Segment, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create segment successfully",
			prepare: func() {
				listDir = func(path string) ([]string, error) {
					return []string{"20190707"}, nil
				}
				newSegmentFunc = func(shard Shard, segmentName string, interval timeutil.Interval) (Segment, error) {
					return nil, nil
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
				newSegmentFunc = newSegment
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			s, err := newIntervalSegment(shard, timeutil.Interval(timeutil.OneSecond*10))
			if ((err != nil) != tt.wantErr && s == nil) || (!tt.wantErr && s == nil) {
				t.Errorf("newIntervalSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIntervalSegment_GetOrCreateSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newSegmentFunc = newSegment
		ctrl.Finish()
	}()

	shard := NewMockShard(ctrl)
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test-db").AnyTimes()
	shard.EXPECT().Database().Return(db).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	db.EXPECT().GetOption().Return(option.DatabaseOption{}).AnyTimes()
	segment := NewMockSegment(ctrl)
	segment.EXPECT().Close()
	newSegmentFunc = func(shard Shard, segmentName string, interval timeutil.Interval) (Segment, error) {
		if segmentName == "201907-a" {
			return nil, fmt.Errorf("er")
		}
		return segment, nil
	}

	s, _ := newIntervalSegment(shard, timeutil.Interval(timeutil.OneSecond*10))
	seg, err := s.GetOrCreateSegment("20190702")
	assert.Nil(t, err)
	assert.NotNil(t, seg)

	seg1, err1 := s.GetOrCreateSegment("20190702")
	assert.NoError(t, err1)
	assert.Equal(t, seg, seg1)

	// test create fail
	seg, err = s.GetOrCreateSegment("201907-a")
	assert.Nil(t, seg)
	assert.NotNil(t, err)

	s.Close()

	// test re-open
	listDir = func(path string) ([]string, error) {
		return []string{"20190702"}, nil
	}
	s, _ = newIntervalSegment(shard, timeutil.Interval(timeutil.OneSecond*10))

	s1, ok := s.(*intervalSegment)
	if ok {
		seg, ok = s1.getSegment("20190702")
		assert.NotNil(t, seg)
		assert.True(t, ok)
	} else {
		t.Fail()
	}
}

func TestIntervalSegment_getDataFamilies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newSegmentFunc = newSegment
		ctrl.Finish()
	}()
	database := NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test").AnyTimes()
	shard := NewMockShard(ctrl)
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	interval := timeutil.Interval(timeutil.OneSecond * 10)
	newSegmentFunc = func(shard Shard, segmentName string, interval timeutil.Interval) (Segment, error) {
		segment := NewMockSegment(ctrl)
		segment.EXPECT().GetOrCreateDataFamily(gomock.Any()).Return(nil, nil).AnyTimes()
		baseTime, _ := interval.Calculator().ParseSegmentTime(segmentName)
		segment.EXPECT().BaseTime().Return(baseTime).AnyTimes()
		start, _ := timeutil.ParseTimestamp("20190902 19:05:48", "20060102 15:04:05")
		end, _ := timeutil.ParseTimestamp("20190905 22:10:48", "20060102 15:04:05")
		segment.EXPECT().getDataFamilies(timeutil.TimeRange{Start: start, End: end}).Return([]DataFamily{nil}).AnyTimes()
		return segment, nil
	}

	s, _ := newIntervalSegment(shard, interval)
	segment1, _ := s.GetOrCreateSegment("20190902")
	now, _ := timeutil.ParseTimestamp("20190902 19:10:48", "20060102 15:04:05")
	_, _ = segment1.GetOrCreateDataFamily(now)
	now, _ = timeutil.ParseTimestamp("20190902 20:10:48", "20060102 15:04:05")
	_, _ = segment1.GetOrCreateDataFamily(now)
	now, _ = timeutil.ParseTimestamp("20190902 22:10:48", "20060102 15:04:05")
	_, _ = segment1.GetOrCreateDataFamily(now)
	segment2, _ := s.GetOrCreateSegment("20190904")
	now, _ = timeutil.ParseTimestamp("20190904 22:10:48", "20060102 15:04:05")
	_, _ = segment2.GetOrCreateDataFamily(now)
	now, _ = timeutil.ParseTimestamp("20190904 20:10:48", "20060102 15:04:05")
	_, _ = segment2.GetOrCreateDataFamily(now)

	start, _ := timeutil.ParseTimestamp("20190901 20:10:48", "20060102 15:04:05")
	end, _ := timeutil.ParseTimestamp("20190901 22:10:48", "20060102 15:04:05")
	families := s.getDataFamilies(timeutil.TimeRange{Start: start, End: end})
	assert.Equal(t, 0, len(families))

	start, _ = timeutil.ParseTimestamp("20190905 20:10:48", "20060102 15:04:05")
	end, _ = timeutil.ParseTimestamp("20190905 22:10:48", "20060102 15:04:05")
	families = s.getDataFamilies(timeutil.TimeRange{Start: start, End: end})
	assert.Equal(t, 0, len(families))

	start, _ = timeutil.ParseTimestamp("20190902 19:05:48", "20060102 15:04:05")
	end, _ = timeutil.ParseTimestamp("20190905 22:10:48", "20060102 15:04:05")
	families = s.getDataFamilies(timeutil.TimeRange{Start: start, End: end})
	assert.Equal(t, 2, len(families))
}
