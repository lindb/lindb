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
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	commonconstants "github.com/lindb/common/constants"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

func TestShard_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		mkDirIfNotExist = fileutil.MkDirIfNotExist
		newIndexDBFunc = indexdb.NewIndexDatabase
		newMemoryDBFunc = memdb.NewMemoryDatabase
		newIntervalSegmentFunc = newIntervalSegment

		kv.InitStoreManager(nil)
		ctrl.Finish()
	}()
	storeMgr := kv.NewMockStoreManager(ctrl)
	kv.InitStoreManager(storeMgr)

	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("db").AnyTimes()
	db.EXPECT().Metadata().Return(nil).AnyTimes()
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
			name: "create shard index store err",
			prepare: func() {
				gomock.InOrder(
					db.EXPECT().GetOption().Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}),
					storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).
						Return(nil, fmt.Errorf("err")),
				)
			},
			wantErr: true,
		},
		{
			name: "create forward index family err",
			prepare: func() {
				store := kv.NewMockStore(ctrl)
				gomock.InOrder(
					db.EXPECT().GetOption().Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}),
					storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).
						Return(store, nil),
					store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err")), // forward
					store.EXPECT().Name().Return("test"),
					storeMgr.EXPECT().CloseStore("test").Return(nil),
				)
			},
			wantErr: true,
		},
		{
			name: "create inverted index family err",
			prepare: func() {
				store := kv.NewMockStore(ctrl)
				gomock.InOrder(
					db.EXPECT().GetOption().Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}),
					storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).
						Return(store, nil),
					store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil),               // forward
					store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err")), // inverted
					store.EXPECT().Name().Return("test"),
					storeMgr.EXPECT().CloseStore("test").Return(nil),
				)
			},
			wantErr: true,
		},
		{
			name: "create index db err",
			prepare: func() {
				newIndexDBFunc = func(ctx context.Context, parent string, metadata metadb.Metadata,
					forwardFamily kv.Family, invertedFamily kv.Family) (indexdb.IndexDatabase, error) {
					return nil, fmt.Errorf("err")
				}
				store := kv.NewMockStore(ctrl)
				gomock.InOrder(
					db.EXPECT().GetOption().Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}),
					storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).
						Return(store, nil),
					store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil), // forward
					store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil), // inverted
					store.EXPECT().Name().Return("test"),
					storeMgr.EXPECT().CloseStore("test").Return(fmt.Errorf("err")),
				)
			},
			wantErr: true,
		},
		{
			name: "create shard successfully",
			prepare: func() {
				newIndexDBFunc = func(ctx context.Context, parent string, metadata metadb.Metadata,
					forwardFamily kv.Family, invertedFamily kv.Family) (indexdb.IndexDatabase, error) {
					return nil, nil
				}
				store := kv.NewMockStore(ctrl)
				gomock.InOrder(
					db.EXPECT().GetOption().Return(&option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}),
					storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).
						Return(store, nil),
					store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil), // forward
					store.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil), // inverted
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
				newIndexDBFunc = indexdb.NewIndexDatabase
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
		kv.InitStoreManager(nil)
		ctrl.Finish()
	}()
	storeMgr := kv.NewMockStoreManager(ctrl)
	kv.InitStoreManager(storeMgr)
	index := indexdb.NewMockIndexDatabase(ctrl)
	segment := NewMockIntervalSegment(ctrl)
	rollupSeg := NewMockIntervalSegment(ctrl)
	bufferMgr := memdb.NewMockBufferManager(ctrl)
	bufferMgr.EXPECT().Cleanup().AnyTimes()
	s := &shard{
		indexDB: index,
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
			name: "close index db err",
			prepare: func() {
				index.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "close store err",
			prepare: func() {
				store := kv.NewMockStore(ctrl)
				s.indexStore = store
				gomock.InOrder(
					index.EXPECT().Close().Return(nil),
					store.EXPECT().Name().Return("test"),
					storeMgr.EXPECT().CloseStore("test").Return(fmt.Errorf("err")),
				)
			},
			wantErr: true,
		},
		{
			name: "close segments",
			prepare: func() {
				gomock.InOrder(
					index.EXPECT().Close().Return(nil),
					segment.EXPECT().Close(),
					rollupSeg.EXPECT().Close(),
				)
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				s.indexStore = nil
			}()
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
	index := indexdb.NewMockIndexDatabase(ctrl)
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	s := &shard{
		indexDB:        index,
		db:             db,
		flushCondition: sync.NewCond(&sync.Mutex{}),
		statistics:     metrics.NewShardStatistics("data", "1"),
		logger:         logger.GetLogger("TSDB", "Test"),
	}
	assert.NotNil(t, s.IndexDatabase())
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
				index.EXPECT().Flush().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "flush successfully",
			prepare: func() {
				index.EXPECT().Flush().Return(nil)
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

func TestShard_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("tet").AnyTimes()
	s := &shard{
		indexDB:    indexDB,
		db:         db,
		metadata:   metadata,
		statistics: metrics.NewShardStatistics("data", "1"),
		logger:     logger.GetLogger("TSDB", "Test"),
	}
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "gen metric id err",
			prepare: func() {
				metadataDB.EXPECT().GenMetricID(commonconstants.DefaultNamespace, "test").
					Return(metric.ID(0), fmt.Errorf("err"))
			},
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			err := s.LookupRowMetricMeta(mockBatchRows(&protoMetricsV1.Metric{
				Name:      "test",
				Timestamp: timeutil.Now(),
				SimpleFields: []*protoMetricsV1.SimpleField{{
					Name:  "f1",
					Value: 1.0,
					Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
				}},
			}))
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteRows() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShard_lookupRowMeta(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("tet").AnyTimes()
	s := &shard{
		indexDB:    indexDB,
		db:         db,
		metadata:   metadata,
		statistics: metrics.NewShardStatistics("data", "1"),
		logger:     logger.GetLogger("TSDB", "Test"),
	}
	cases := []struct {
		name      string
		namespace string
		tags      []*protoMetricsV1.KeyValue
		prepare   func()
		wantErr   bool
	}{
		{
			name: "gen metric id err",
			prepare: func() {
				metadataDB.EXPECT().GenMetricID(commonconstants.DefaultNamespace, "test").Return(metric.ID(0), fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "gen series id err",
			tags: tag.KeyValuesFromMap(map[string]string{"ip": "1.1.1.1"}),
			prepare: func() {
				metadataDB.EXPECT().GenMetricID(commonconstants.DefaultNamespace, "test").Return(metric.ID(10), nil).AnyTimes()
				indexDB.EXPECT().GetOrCreateSeriesID(gomock.Any(), gomock.Any(),
					metric.ID(10), gomock.Any(), gomock.Any()).Return(uint32(0), false, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "get old series id",
			namespace: "ns",
			tags:      tag.KeyValuesFromMap(map[string]string{"ip": "1.1.1.1"}),
			prepare: func() {
				metadataDB.EXPECT().GenMetricID("ns", "test").Return(metric.ID(10), nil).AnyTimes()
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Return(field.ID(1), nil)
				indexDB.EXPECT().GetOrCreateSeriesID(gomock.Any(), gomock.Any(),
					metric.ID(10), gomock.Any(), gomock.Any()).Return(uint32(10), false, nil)
			},
		},
		{
			name: "empty tags",
			prepare: func() {
				metadataDB.EXPECT().GenMetricID(commonconstants.DefaultNamespace, "test").Return(metric.ID(10), nil).AnyTimes()
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Return(field.ID(1), nil)
			},
		},
		{
			name: "get new limits",
			tags: tag.KeyValuesFromMap(map[string]string{"ip": "1.1.1.1"}),
			prepare: func() {
				s.notifyLimitsChange()
				db.EXPECT().GetLimits().Return(models.NewDefaultLimits())
				metadataDB.EXPECT().GenMetricID(commonconstants.DefaultNamespace, "test").Return(metric.ID(10), nil).AnyTimes()
				indexDB.EXPECT().GetOrCreateSeriesID(gomock.Any(), gomock.Any(),
					metric.ID(10), gomock.Any(), gomock.Any()).Return(uint32(0), false, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "build inverted index, but gen field failure",
			tags: tag.KeyValuesFromMap(map[string]string{"ip": "1.1.1.1"}),
			prepare: func() {
				metadataDB.EXPECT().GenMetricID(commonconstants.DefaultNamespace, "test").Return(metric.ID(10), nil).AnyTimes()
				indexDB.EXPECT().GetOrCreateSeriesID(gomock.Any(), gomock.Any(),
					metric.ID(10), gomock.Any(), gomock.Any()).Return(uint32(1), true, nil)
				indexDB.EXPECT().BuildInvertIndex(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(field.ID(0), fmt.Errorf("err"))
			},
			wantErr: true,
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			err := s.lookupRowMeta(&(mockBatchRows(&protoMetricsV1.Metric{
				Name:      "test",
				Namespace: tt.namespace,
				Timestamp: timeutil.Now(),
				Tags:      tt.tags,
				SimpleFields: []*protoMetricsV1.SimpleField{{
					Name:  "f1",
					Value: 1.0,
					Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
				}},
			})[0]))
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteRows() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShard_lookup_histogram_fields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	indexDB := indexdb.NewMockIndexDatabase(ctrl)
	metadata := metadb.NewMockMetadata(ctrl)
	metadataDB := metadb.NewMockMetadataDatabase(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataDB).AnyTimes()
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("tet").AnyTimes()
	metadataDB.EXPECT().GenMetricID(commonconstants.DefaultNamespace, "test").Return(metric.ID(10), nil).AnyTimes()
	s := &shard{
		indexDB:    indexDB,
		db:         db,
		metadata:   metadata,
		statistics: metrics.NewShardStatistics("data", "1"),
		logger:     logger.GetLogger("TSDB", "Test"),
	}
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "gen min field failure",
			prepare: func() {
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(field.ID(0), fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "gen max field failure",
			prepare: func() {
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(field.ID(1), nil)
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(field.ID(0), fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "gen sum field failure",
			prepare: func() {
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(field.ID(1), nil).MaxTimes(2)
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(field.ID(0), fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "gen count field failure",
			prepare: func() {
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(field.ID(1), nil).MaxTimes(3)
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(field.ID(0), fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "gen bucket field failure",
			prepare: func() {
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(field.ID(1), nil).MaxTimes(4)
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(field.ID(0), fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "gen all fields successfully",
			prepare: func() {
				metadataDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(field.ID(1), nil).AnyTimes()
			},
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			err := s.lookupRowMeta(&(mockBatchRows(&protoMetricsV1.Metric{
				Name:      "test",
				Timestamp: timeutil.Now(),
				CompoundField: &protoMetricsV1.CompoundField{
					Min:            10,
					Max:            10,
					Sum:            10,
					Count:          10,
					ExplicitBounds: []float64{1, 1, 1, 1, 1, math.Inf(1) + 1},
					Values:         []float64{1, 1, 1, 1, 1, 1},
				},
			})[0]))
			if (err != nil) != tt.wantErr {
				t.Errorf("lookupRowMeta() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShard_WaitFlushIndexCompleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := timeutil.Now()

	index := indexdb.NewMockIndexDatabase(ctrl)
	db := NewMockDatabase(ctrl)
	db.EXPECT().Name().Return("test").AnyTimes()
	s := &shard{
		indexDB:        index,
		db:             db,
		flushCondition: sync.NewCond(&sync.Mutex{}),
		statistics:     metrics.NewShardStatistics("data", "1"),
		logger:         logger.GetLogger("TSDB", "Test"),
	}
	s.isFlushing.Store(false)
	index.EXPECT().Flush().DoAndReturn(func() error {
		time.Sleep(100 * time.Millisecond)
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
	assert.True(t, timeutil.Now()-now >= 90*time.Millisecond.Milliseconds())
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

func mockBatchRows(m *protoMetricsV1.Metric) []metric.StorageRow {
	var ml = protoMetricsV1.MetricList{Metrics: []*protoMetricsV1.Metric{m}}
	var buf bytes.Buffer
	converter := metric.NewProtoConverter()
	_, _ = converter.MarshalProtoMetricListV1To(ml, &buf)

	var br metric.StorageBatchRows
	br.UnmarshalRows(buf.Bytes())
	return br.Rows()
}
