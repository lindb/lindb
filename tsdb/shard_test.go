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
	"testing"

	"github.com/lindb/lindb/pkg/timeutil"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"go.uber.org/mock/gomock"

	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"

	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb/memdb"
)

func TestShard_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		mkDirIfNotExist = fileutil.MkDirIfNotExist
		newIndexDBFunc = index.NewMetricIndexDatabase
		newMemoryDBFunc = memdb.NewMemoryDatabase
		newIntervalDataSegmentFunc = newIntervalDataSegment

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
			name: "create interval dataSegment err",
			prepare: func() {
				newIntervalDataSegmentFunc = func(shard Shard, interval option.Interval) (segment IntervalDataSegment, err error) {
					return nil, fmt.Errorf("err")
				}
				db.EXPECT().GetOption().Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: 1}}})
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
				seq := NewMockIntervalDataSegment(ctrl)
				seq.EXPECT().Close().AnyTimes()

				newIntervalDataSegmentFunc = func(shard Shard,
					interval option.Interval,
				) (IntervalDataSegment, error) {
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

func TestShard_Database(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockDatabase(ctrl)
	s := &shard{db: db}
	assert.Equal(t, db, s.Database())
}

func TestShard_ShardID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := &shard{id: models.ShardID(1)}
	assert.Equal(t, models.ShardID(1), s.ShardID())
}

func TestShard_CurrentInterval(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	opt := &option.DatabaseOption{Intervals: option.Intervals{option.Interval{
		Interval: 1,
	}}}
	s := &shard{
		option: opt,
	}
	assert.Equal(t, opt.Intervals[0].Interval, s.CurrentInterval())
}

func TestShard_BufferManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bufferMgr := memdb.NewMockBufferManager(ctrl)
	s := &shard{
		bufferMgr: bufferMgr,
	}
	assert.Equal(t, bufferMgr, s.BufferManager())
}

func TestShard_GetIndexDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	indexDB := index.NewMockMetricIndexDatabase(ctrl)

	seg := NewMockSegment(ctrl)
	seg.EXPECT().IndexDB().Return(indexDB)

	sp := NewMockSegmentPartition(ctrl)

	s := &shard{
		segmentPartition: sp,
	}

	sp.EXPECT().GetOrCreateSegment(gomock.Any()).Return(seg, nil)
	assert.Equal(t, indexDB, s.GetIndexDB(1))
}

func TestShard_GetOrCreateDataFamily(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sp := NewMockSegmentPartition(ctrl)
	mockDataFamily := NewMockDataFamily(ctrl)
	segment := NewMockSegment(ctrl)
	s := &shard{
		segmentPartition: sp,
	}

	sp.EXPECT().GetOrCreateSegment(gomock.Any()).Return(nil, fmt.Errorf("err"))
	dataFamily, err := s.GetOrCreateDataFamily(1)
	assert.Error(t, err)
	assert.Nil(t, dataFamily)

	sp.EXPECT().GetOrCreateSegment(gomock.Any()).Return(segment, nil)
	segment.EXPECT().GetOrCrateDataFamily(gomock.Any()).Return(mockDataFamily, nil)
	dataFamily, err = s.GetOrCreateDataFamily(1)
	assert.NoError(t, err)
	assert.Equal(t, mockDataFamily, dataFamily)
}

func TestShard_GetDataFamilies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sp := NewMockSegmentPartition(ctrl)
	mockDataFamily := NewMockDataFamily(ctrl)
	segment := NewMockSegment(ctrl)
	s := &shard{
		segmentPartition: sp,
	}

	sp.EXPECT().GetSegments().Return([]Segment{segment})
	segment.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return([]DataFamily{mockDataFamily})
	dataFamilies := s.GetDataFamilies(timeutil.Day, timeutil.TimeRange{})
	assert.Len(t, dataFamilies, 1)
}

func TestShard_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	statistics := metrics.NewShardStatistics("test", "1")
	bufferMgr := memdb.NewMockBufferManager(ctrl)
	sp := NewMockSegmentPartition(ctrl)

	bufferMgr.EXPECT().Cleanup().AnyTimes()
	sp.EXPECT().WaitFlushIndexCompleted().AnyTimes()

	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()

	s := &shard{
		id:               models.ShardID(1),
		db:               db,
		statistics:       statistics,
		bufferMgr:        bufferMgr,
		segmentPartition: sp,
		logger:           logger.GetLogger("TSDB", "Shard"),
	}

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "flush index error",
			prepare: func() {
				sp.EXPECT().FlushIndex().Return(fmt.Errorf("err"))
				sp.EXPECT().Close().Return(fmt.Errorf("err"))
			},
		},
		{
			name: "close index error",
			prepare: func() {
				sp.EXPECT().FlushIndex().Return(nil)
				sp.EXPECT().Close().Return(fmt.Errorf("err"))
			},
		},
		{
			name: "successfully",
			prepare: func() {
				sp.EXPECT().FlushIndex().Return(nil)
				sp.EXPECT().Close().Return(nil)
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

func TestShard_FlushIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	statistics := metrics.NewShardStatistics("test", "1")
	bufferMgr := memdb.NewMockBufferManager(ctrl)
	sp := NewMockSegmentPartition(ctrl)

	bufferMgr.EXPECT().Cleanup().AnyTimes()
	sp.EXPECT().WaitFlushIndexCompleted().AnyTimes()

	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()

	s := &shard{
		id:               models.ShardID(1),
		statistics:       statistics,
		db:               db,
		bufferMgr:        bufferMgr,
		segmentPartition: sp,
		logger:           logger.GetLogger("TSDB", "Shard"),
	}

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "flush index error",
			prepare: func() {
				sp.EXPECT().FlushIndex().Return(fmt.Errorf("err"))
			},
		},
		{
			name: "successfully",
			prepare: func() {
				sp.EXPECT().FlushIndex().Return(nil)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

	sp := NewMockSegmentPartition(ctrl)
	sp.EXPECT().WaitFlushIndexCompleted()
	s := shard{
		segmentPartition: sp,
	}
	s.WaitFlushIndexCompleted()
}

func TestShard_TTL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()

	segmentPartition := NewMockSegmentPartition(ctrl)
	s := &shard{
		id:               models.ShardID(1),
		db:               db,
		segmentPartition: segmentPartition,
		logger:           logger.GetLogger("TSDB", "Test"),
	}

	segmentPartition.EXPECT().TTL().Return(fmt.Errorf("err"))
	s.TTL()

	segmentPartition.EXPECT().TTL().Return(nil)
	s.TTL()
}

func TestShard_EvictSegment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	segmentPartition := NewMockSegmentPartition(ctrl)
	s := &shard{
		segmentPartition: segmentPartition,
	}
	segmentPartition.EXPECT().EvictSegment()
	s.EvictSegment()
}

func TestShard_notifyLimitsChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := &shard{
		limitsChanged: atomic.Bool{},
	}
	s.notifyLimitsChange()
	assert.Equal(t, true, s.limitsChanged.Load())
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
