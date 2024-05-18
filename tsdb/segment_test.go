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
	"time"

	"github.com/lindb/common/pkg/logger"

	commonTimeutil "github.com/lindb/common/pkg/timeutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestSegment_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		formatTimestamp = commonTimeutil.FormatTimestamp
		newIndexDBFunc = index.NewMetricIndexDatabase
		newIntervalDataSegmentFunc = newIntervalDataSegment
		ctrl.Finish()
	}()

	formatTimestamp = func(timestamp int64, layout string) string {
		return "test"
	}

	cases := []struct {
		name       string
		prepare    func()
		assertFunc func(*segment)
		wantErr    bool
	}{
		{
			name: "new index db error",
			prepare: func() {
				newIndexDBFunc = func(dir string, metaDB index.MetricMetaDatabase) (index.MetricIndexDatabase, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "new interval data segment error",
			prepare: func() {
				newIndexDBFunc = func(dir string, metaDB index.MetricMetaDatabase) (index.MetricIndexDatabase, error) {
					return nil, nil
				}
				newIntervalDataSegmentFunc = func(shard Shard, interval option.Interval) (segment IntervalDataSegment, err error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "successfully",
			prepare: func() {
				newIntervalDataSegmentFunc = func(shard Shard, interval option.Interval) (segment IntervalDataSegment, err error) {
					return NewMockIntervalDataSegment(ctrl), nil
				}
			},
			assertFunc: func(segment *segment) {
				assert.NotNil(t, segment)
				assert.NotNil(t, segment.writableDataSegment)
				assert.Equal(t, 3, len(segment.rollupTargets))
			},
		},
	}

	database := NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test").AnyTimes()
	database.EXPECT().MetaDB().Return(nil).AnyTimes()

	shard := NewMockShard(ctrl)
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()

	a := timeutil.Interval(commonTimeutil.OneSecond)
	b := timeutil.Interval(5 * commonTimeutil.OneMinute)
	c := timeutil.Interval(commonTimeutil.OneHour)

	intervals := option.Intervals{
		option.Interval{
			Interval: a,
		},
		option.Interval{
			Interval: b,
		},
		option.Interval{
			Interval: c,
		},
	}

	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			seg, err := newSegment(shard, 1, intervals)
			if (err != nil) != tt.wantErr {
				t.Fatal(tt.name)
			}
			if tt.assertFunc != nil {
				tt.assertFunc(seg.(*segment))
			}
		})
	}
}

func TestSegment_GetName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &segment{
		name: "test",
	}
	assert.Equal(t, s.name, s.GetName())
}

func TestSegment_IndexDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	s := &segment{
		indexDB: indexDB,
	}
	assert.Equal(t, indexDB, s.IndexDB())
}

func TestSegment_GetOrCreateDataFamily(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	a := timeutil.Interval(commonTimeutil.OneSecond)
	b := timeutil.Interval(5 * commonTimeutil.OneMinute)
	c := timeutil.Interval(commonTimeutil.OneHour)

	intervals := option.Intervals{
		option.Interval{
			Interval: a,
		},
		option.Interval{
			Interval: b,
		},
		option.Interval{
			Interval: c,
		},
	}

	dataSegment := NewMockDataSegment(ctrl)
	intervalDataSegment := NewMockIntervalDataSegment(ctrl)
	writableDataSegment := NewMockIntervalDataSegment(ctrl)

	s := &segment{
		interval:            intervals[0].Interval,
		writableDataSegment: writableDataSegment,
		rollupTargets: map[timeutil.Interval]IntervalDataSegment{
			b: intervalDataSegment,
		},
	}

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "get or create dataSegment error",
			prepare: func() {
				writableDataSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "rollup segment get or create data family error",
			prepare: func() {
				writableDataSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(nil, nil)
				intervalDataSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "get or create data family error",
			prepare: func() {
				writableDataSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(dataSegment, nil)
				intervalDataSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(dataSegment, nil)
				dataSegment.EXPECT().GetOrCreateDataFamily(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "successfully",
			prepare: func() {
				writableDataSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(dataSegment, nil)
				intervalDataSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(dataSegment, nil)
				dataSegment.EXPECT().GetOrCreateDataFamily(gomock.Any()).Return(nil, nil)
			},
		},
	}
	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			_, err := s.GetOrCreateDataFamily(time.Now().UnixMilli())
			if (err != nil) != tt.wantErr {
				t.Fatal(tt.name)
			}
		})
	}
}

func TestSegment_GetDataFamilies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writableIntervalDataSegment := NewMockIntervalDataSegment(ctrl)
	rollupSeg := NewMockIntervalDataSegment(ctrl)
	s := &segment{
		writableDataSegment: writableIntervalDataSegment,
		interval:            timeutil.Interval(10 * 1000), // 10s
		rollupTargets: map[timeutil.Interval]IntervalDataSegment{
			timeutil.Interval(10 * 1000):      rollupSeg, // 10s
			timeutil.Interval(10 * 60 * 1000): rollupSeg, // 10min
		},
	}
	cases := []struct {
		name         string
		intervalType timeutil.IntervalType
		prepare      func()
		assert       func(families []DataFamily)
	}{
		{
			name:         "match writable dataSegment",
			intervalType: timeutil.Day,
			prepare: func() {
				writableIntervalDataSegment.EXPECT().GetDataFamilies(gomock.Any()).Return([]DataFamily{nil})
			},
			assert: func(families []DataFamily) {
				assert.Len(t, families, 1)
			},
		},
		{
			name:         "match rollup dataSegment",
			intervalType: timeutil.Month,
			prepare: func() {
				rollupSeg.EXPECT().GetDataFamilies(gomock.Any()).Return([]DataFamily{nil})
			},
			assert: func(families []DataFamily) {
				assert.Len(t, families, 1)
			},
		},
		{
			name:         "not match dataSegment",
			intervalType: timeutil.Year,
			assert: func(families []DataFamily) {
				assert.Len(t, families, 0)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			families := s.GetDataFamilies(tt.intervalType, timeutil.TimeRange{})
			if tt.assert != nil {
				tt.assert(families)
			}
		})
	}
	// test no rollup
	s = &segment{
		writableDataSegment: writableIntervalDataSegment,
		rollupTargets: map[timeutil.Interval]IntervalDataSegment{
			timeutil.Interval(10 * 1000): rollupSeg, // 10s
		},
	}
	writableIntervalDataSegment.EXPECT().GetDataFamilies(gomock.Any()).Return([]DataFamily{nil})
	families := s.GetDataFamilies(timeutil.Year, timeutil.TimeRange{})
	assert.Len(t, families, 1)
}

func TestSegment_GetTimestamp(t *testing.T) {
	s := &segment{timestamp: 1}
	assert.Equal(t, int64(1), s.GetTimestamp())
}

func TestSegment_FlushIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	indexDB := index.NewMockMetricIndexDatabase(ctrl)

	s := &segment{
		indexDB:    indexDB,
		statistics: metrics.NewSegmentStatistics("database name", "1", "segment name"),
	}

	indexDB.EXPECT().Notify(gomock.Any()).Do(func(notifier index.Notifier) {
		mn := notifier.(*index.FlushNotifier)
		mn.Callback(fmt.Errorf("err"))
	})

	indexDB.EXPECT().Notify(gomock.Any()).Do(func(notifier index.Notifier) {
		mn := notifier.(*index.FlushNotifier)
		mn.Callback(nil)
	})

	err := s.FlushIndex()
	assert.Error(t, err)
	err = s.FlushIndex()
	assert.NoError(t, err)
}

func TestSegment_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	a := NewMockIntervalDataSegment(ctrl)
	b := NewMockIntervalDataSegment(ctrl)
	c := NewMockIntervalDataSegment(ctrl)

	a.EXPECT().Close()
	b.EXPECT().Close()
	c.EXPECT().Close()

	s := &segment{
		rollupTargets: map[timeutil.Interval]IntervalDataSegment{
			timeutil.Interval(commonTimeutil.OneSecond):     a,
			timeutil.Interval(5 * commonTimeutil.OneMinute): b,
			timeutil.Interval(commonTimeutil.OneHour):       c,
		},
	}

	err := s.Close()
	assert.NoError(t, err)
}

func TestSegment_TTL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	database := NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test")

	shard := NewMockShard(ctrl)
	shard.EXPECT().ShardID().Return(models.ShardID(1))
	shard.EXPECT().Database().Return(database)

	intervalDataSegment := NewMockIntervalDataSegment(ctrl)

	s := &segment{
		shard: shard,
		name:  "test",
		rollupTargets: map[timeutil.Interval]IntervalDataSegment{
			timeutil.Interval(commonTimeutil.OneSecond): intervalDataSegment,
		},
		logger: logger.GetLogger("TSDB", "Segment"),
	}

	intervalDataSegment.EXPECT().TTL().Return(fmt.Errorf("err"))
	assert.NoError(t, s.TTL())

	intervalDataSegment.EXPECT().TTL().Return(nil)
	assert.NoError(t, s.TTL())
}

func TestSegment_EvictSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	intervalDataSegment := NewMockIntervalDataSegment(ctrl)
	s := &segment{
		rollupTargets: map[timeutil.Interval]IntervalDataSegment{
			timeutil.Interval(commonTimeutil.OneSecond): intervalDataSegment,
		},
	}
	intervalDataSegment.EXPECT().EvictSegment()
	s.EvictSegment()
}
