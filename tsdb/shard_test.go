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
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	commontimeutil "github.com/lindb/common/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"

	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb/memdb"
)

func TestShard_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		mkDirIfNotExist = fileutil.MkDirIfNotExist
		newIndexDBFunc = index.NewMetricIndexDatabase
		newMemoryDBFunc = memdb.NewMemoryDatabase
		newIntervalSegmentFunc = newIntervalSegment

		ctrl.Finish()
	}()

	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("db").AnyTimes()
	db.EXPECT().MetaDB().Return(nil).AnyTimes()
	db.EXPECT().GetLimits().Return(models.NewDefaultLimits()).AnyTimes()

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "create shard path err",
			prepare: func() {
				mkDirIfNotExist = func(path string) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create interval segment err",
			prepare: func() {
				newIntervalSegmentFunc = func(shard Shard, interval option.Interval) (segment IntervalSegment, err error) {
					return nil, fmt.Errorf("err")
				}
				db.EXPECT().GetOption().Return(&option.DatabaseOption{Intervals: option.Intervals{{}}})
			},
			wantErr: true,
		},
		{
			name: "create index db err",
			prepare: func() {
				indexDB := index.NewMockMetricIndexDatabase(ctrl)
				newIndexDBFunc = func(_ string, _ index.MetricMetaDatabase) (index.MetricIndexDatabase, error) {
					return indexDB, fmt.Errorf("err")
				}
				indexDB.EXPECT().Notify(gomock.Any()).DoAndReturn(func(notifier index.Notifier) {
					mn := notifier.(*index.FlushNotifier)
					mn.Callback(fmt.Errorf("err"))
				})
				gomock.InOrder(
					db.EXPECT().GetOption().Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}),
				)
			},
			wantErr: true,
		},
		{
			name: "create shard successfully",
			prepare: func() {
				newIndexDBFunc = func(_ string, _ index.MetricMetaDatabase) (index.MetricIndexDatabase, error) {
					return nil, nil
				}
				gomock.InOrder(
					db.EXPECT().GetOption().Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}),
				)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				// mock mkdir
				mkDirIfNotExist = func(path string) error {
					return nil
				}
				newIndexDBFunc = index.NewMetricIndexDatabase
				seq := NewMockIntervalSegment(ctrl)
				seq.EXPECT().Close().AnyTimes()

				newIntervalSegmentFunc = func(shard Shard,
					interval option.Interval,
				) (IntervalSegment, error) {
					return seq, nil
				}
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			s, err := newShard(db, 1)
			if ((err != nil) != tt.wantErr && s == nil) || (!tt.wantErr && s == nil) {
				t.Errorf("newShard() error = %v, wantErr %v", err, tt.wantErr)
			}
			if s != nil {
				assert.Equal(t, db, s.Database())
				assert.Equal(t, models.ShardID(1), s.ShardID())
				assert.NotEmpty(t, s.Indicator())
				assert.NotNil(t, s.BufferManager())
				assert.True(t, s.CurrentInterval().Int64() > 0)
				s.notifyLimitsChange()
			}
		})
	}
}

func TestShard_GetDataFamilies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segment := NewMockIntervalSegment(ctrl)
	rollupSeg := NewMockIntervalSegment(ctrl)
	s := &shard{
		interval: timeutil.Interval(10 * 1000), // 10s
		segment:  segment,
		rollupTargets: map[timeutil.Interval]IntervalSegment{
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
			name:         "match writable segment",
			intervalType: timeutil.Day,
			prepare: func() {
				segment.EXPECT().GetDataFamilies(gomock.Any()).Return([]DataFamily{nil})
			},
			assert: func(families []DataFamily) {
				assert.Len(t, families, 1)
			},
		},
		{
			name:         "match rollup segment",
			intervalType: timeutil.Month,
			prepare: func() {
				rollupSeg.EXPECT().GetDataFamilies(gomock.Any()).Return([]DataFamily{nil})
			},
			assert: func(families []DataFamily) {
				assert.Len(t, families, 1)
			},
		},
		{
			name:         "not match segment",
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
	s = &shard{
		interval: timeutil.Interval(10 * 1000), // 10s
		segment:  segment,
		rollupTargets: map[timeutil.Interval]IntervalSegment{
			timeutil.Interval(10 * 1000): rollupSeg, // 10s
		},
	}
	segment.EXPECT().GetDataFamilies(gomock.Any()).Return([]DataFamily{nil})
	families := s.GetDataFamilies(timeutil.Year, timeutil.TimeRange{})
	assert.Len(t, families, 1)
}

func TestShard_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	segment := NewMockIntervalSegment(ctrl)
	rollupSeg := NewMockIntervalSegment(ctrl)
	bufferMgr := memdb.NewMockBufferManager(ctrl)
	bufferMgr.EXPECT().Cleanup().AnyTimes()
	s := &shard{
		indexDB: indexDB,
		segment: segment,
		rollupTargets: map[timeutil.Interval]IntervalSegment{
			timeutil.Interval(10 * 60 * 1000): rollupSeg, // 10min
		},
		flushCondition: sync.NewCond(&sync.Mutex{}),
		bufferMgr:      bufferMgr,
	}
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "flush index db err",
			prepare: func() {
				indexDB.EXPECT().Notify(gomock.Any()).DoAndReturn(func(notifier index.Notifier) {
					mn := notifier.(*index.FlushNotifier)
					mn.Callback(fmt.Errorf("err"))
				})
			},
			wantErr: true,
		},
		{
			name: "close index db err",
			prepare: func() {
				indexDB.EXPECT().Notify(gomock.Any()).DoAndReturn(func(notifier index.Notifier) {
					mn := notifier.(*index.FlushNotifier)
					mn.Callback(nil)
				})
				indexDB.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "close segments",
			prepare: func() {
				gomock.InOrder(
					indexDB.EXPECT().Notify(gomock.Any()).DoAndReturn(func(notifier index.Notifier) {
						mn := notifier.(*index.FlushNotifier)
						mn.Callback(nil)
					}),
					indexDB.EXPECT().Close().Return(nil),
					segment.EXPECT().Close(),
					rollupSeg.EXPECT().Close(),
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
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShard_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	s := &shard{
		indexDB:        indexDB,
		db:             db,
		flushCondition: sync.NewCond(&sync.Mutex{}),
		statistics:     metrics.NewShardStatistics("data", "1"),
		logger:         logger.GetLogger("TSDB", "Test"),
	}
	assert.NotNil(t, s.IndexDB())
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
				indexDB.EXPECT().Notify(gomock.Any()).DoAndReturn(func(notifier index.Notifier) {
					mn := notifier.(*index.FlushNotifier)
					mn.Callback(fmt.Errorf("err"))
				})
			},
			wantErr: true,
		},
		{
			name: "flush successfully",
			prepare: func() {
				indexDB.EXPECT().Notify(gomock.Any()).DoAndReturn(func(notifier index.Notifier) {
					mn := notifier.(*index.FlushNotifier)
					mn.Callback(nil)
				})
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

func TestShard_WaitFlushIndexCompleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := commontimeutil.Now()

	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	s := &shard{
		indexDB:        indexDB,
		db:             db,
		flushCondition: sync.NewCond(&sync.Mutex{}),
		statistics:     metrics.NewShardStatistics("data", "1"),
		logger:         logger.GetLogger("TSDB", "Test"),
	}
	s.isFlushing.Store(false)
	indexDB.EXPECT().Notify(gomock.Any()).DoAndReturn(func(notifier index.Notifier) {
		time.Sleep(100 * time.Millisecond)
		mn := notifier.(*index.FlushNotifier)
		mn.Callback(nil)
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

func TestShard_TTL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	segment := NewMockIntervalSegment(ctrl)
	s := &shard{
		rollupTargets: map[timeutil.Interval]IntervalSegment{
			10: segment,
		},
		db:     db,
		logger: logger.GetLogger("TSDB", "Test"),
	}
	segment.EXPECT().TTL().Return(fmt.Errorf("err"))
	s.TTL()
}

func TestShard_EvictSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	segment := NewMockIntervalSegment(ctrl)
	s := &shard{
		rollupTargets: map[timeutil.Interval]IntervalSegment{
			10: segment,
		},
		db:     db,
		logger: logger.GetLogger("TSDB", "Test"),
	}
	segment.EXPECT().EvictSegment()
	s.EvictSegment()
}

func TestShard_GetOrCreateDataFamily(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	intervalSegment := NewMockIntervalSegment(ctrl)
	segment := NewMockSegment(ctrl)
	s := &shard{
		segment: intervalSegment,
	}
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "get or create segment error",
			prepare: func() {
				intervalSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "get or create data family error",
			prepare: func() {
				intervalSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(segment, nil)
				segment.EXPECT().GetOrCreateDataFamily(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "create rollup target segment error",
			prepare: func() {
				s.rollupTargets = map[timeutil.Interval]IntervalSegment{
					10: intervalSegment,
				}
				intervalSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(segment, nil)
				intervalSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "successfully",
			prepare: func() {
				intervalSegment.EXPECT().GetOrCreateSegment(gomock.Any()).Return(segment, nil)
				segment.EXPECT().GetOrCreateDataFamily(gomock.Any()).Return(nil, nil)
			},
		},
	}
	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				s.rollupTargets = nil
			}()
			tt.prepare()
			_, err := s.GetOrCrateDataFamily(time.Now().UnixMilli())
			if (err != nil) != tt.wantErr {
				t.Fatal(tt.name)
			}
		})
	}
}

func mockBatchRows(m *protoMetricsV1.Metric) []metric.StorageRow {
	var ml = protoMetricsV1.MetricList{Metrics: []*protoMetricsV1.Metric{m}}
	var buf bytes.Buffer
	converter := metric.NewProtoConverter(models.NewDefaultLimits())
	_, _ = converter.MarshalProtoMetricListV1To(ml, &buf)

	var br metric.StorageBatchRows
	br.UnmarshalRows(buf.Bytes())
	return br.Rows()
}
