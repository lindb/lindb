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

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"go.uber.org/mock/gomock"

	"github.com/lindb/common/pkg/fasttime"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/ltoml"
	commontimeutil "github.com/lindb/common/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/index"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/metric"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestDataFamily_BaseTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := kv.NewMockFamily(ctrl)
	timeRange := timeutil.TimeRange{
		Start: 10,
		End:   50,
	}
	database := NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test").AnyTimes()
	snapshot := version.NewMockSnapshot(ctrl)
	v := version.NewMockVersion(ctrl)
	v.EXPECT().GetSequences().Return(map[int32]int64{1: 10})
	snapshot.EXPECT().GetCurrent().Return(v)
	snapshot.EXPECT().Close()
	family.EXPECT().GetSnapshot().Return(snapshot)
	shard := NewMockShard(ctrl)
	shard.EXPECT().Database().Return(database)
	shard.EXPECT().ShardID().Return(models.ShardID(1))
	dataFamily := newDataFamily(shard, nil, timeutil.Interval(commontimeutil.OneSecond*10), timeRange, 10, family)
	assert.Equal(t, timeRange, dataFamily.TimeRange())
	assert.Equal(t, timeutil.Interval(10000), dataFamily.Interval())
	assert.NotNil(t, dataFamily.Family())
	assert.Equal(t, shard, dataFamily.Shard())
	assert.Equal(t, int64(10), dataFamily.FamilyTime())

	err := dataFamily.Close()
	assert.NoError(t, err)
}

func TestDataFamily_Filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := kv.NewMockFamily(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close().AnyTimes()
	family.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	reader := table.NewMockReader(ctrl)
	reader.EXPECT().Path().Return("test").AnyTimes()
	now := commontimeutil.Now()
	cases := []struct {
		name    string
		prepare func(f *dataFamily)
		len     int
		wantErr bool
	}{
		{
			name: "filter memory database failure",
			prepare: func(f *dataFamily) {
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.mutableMemDB = memDB
				memDB.EXPECT().Filter(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "filter immutable memory database failure",
			prepare: func(f *dataFamily) {
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.immutableMemDB = memDB
				memDB.EXPECT().Filter(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "filter memory database successfully",
			prepare: func(f *dataFamily) {
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.immutableMemDB = memDB
				memDB.EXPECT().Filter(gomock.Any()).Return([]flow.FilterResultSet{nil}, nil)
				snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
			},
			wantErr: false,
			len:     1,
		},
		{
			name: "get file reader failure",
			prepare: func(_ *dataFamily) {
				snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "get metric reader data failure",
			prepare: func(_ *dataFamily) {
				snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{reader}, nil)
				reader.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: false,
			len:     0,
		},
		{
			name: "new metric reader failure",
			prepare: func(_ *dataFamily) {
				snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{reader}, nil)
				reader.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3}, nil)
				newReaderFunc = func(path string, metricBlock []byte) (metricsdata.MetricReader, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "time range not match",
			prepare: func(_ *dataFamily) {
				snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{reader}, nil)
				reader.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3}, nil)
				mReader := metricsdata.NewMockMetricReader(ctrl)
				newReaderFunc = func(path string, metricBlock []byte) (metricsdata.MetricReader, error) {
					return mReader, nil
				}
				mReader.EXPECT().GetTimeRange().Return(timeutil.SlotRange{Start: 1000, End: 1000})
			},
			wantErr: false,
			len:     0,
		},
		{
			name: "find data",
			prepare: func(_ *dataFamily) {
				snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{reader}, nil)
				reader.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3}, nil)
				mReader := metricsdata.NewMockMetricReader(ctrl)
				newReaderFunc = func(path string, metricBlock []byte) (metricsdata.MetricReader, error) {
					return mReader, nil
				}
				mReader.EXPECT().GetTimeRange().Return(timeutil.SlotRange{Start: 0, End: 1000})
				filter := metricsdata.NewMockFilter(ctrl)
				newFilterFunc = func(familyTime int64, snapshot version.Snapshot,
					readers []metricsdata.MetricReader) metricsdata.Filter {
					return filter
				}
				filter.EXPECT().Filter(gomock.Any(), gomock.Any()).Return([]flow.FilterResultSet{nil}, nil)
			},
			wantErr: false,
			len:     1,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newReaderFunc = metricsdata.NewReader
				newFilterFunc = metricsdata.NewFilter
			}()
			f := &dataFamily{
				familyTime:   now,
				family:       family,
				lastReadTime: atomic.NewInt64(fasttime.UnixMilliseconds()),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			rs, err := f.Filter(&flow.ShardExecuteContext{
				StorageExecuteCtx: &flow.StorageExecuteContext{
					MetricID: 1,
					Query: &stmtpkg.Query{
						StorageInterval: timeutil.Interval(commontimeutil.OneMinute),
						TimeRange:       timeutil.TimeRange{Start: now, End: now + 60000},
					},
				},
			})
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, rs)
			} else {
				assert.NoError(t, err)
				assert.Len(t, rs, tt.len)
			}
		})
	}
}

func TestDataFamily_NeedFlush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := NewMockDatabase(ctrl)
	shard := NewMockShard(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()

	cases := []struct {
		name      string
		prepare   func(f *dataFamily)
		needFlush bool
	}{
		{
			name: "flush job is doing",
			prepare: func(f *dataFamily) {
				f.isFlushing.Store(true)
			},
			needFlush: false,
		},
		{
			name: "immutable memory database is not empty",
			prepare: func(f *dataFamily) {
				f.immutableMemDB = memdb.NewMockMemoryDatabase(ctrl)
			},
			needFlush: false,
		},
		{
			name:      "memory database is nil",
			needFlush: false,
		},
		{
			name: "memory database size is zero",
			prepare: func(f *dataFamily) {
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.mutableMemDB = memDB
				memDB.EXPECT().NumOfMetrics().Return(0)
			},
			needFlush: false,
		},
		{
			name: "trigger time threshold",
			prepare: func(f *dataFamily) {
				cfg := config.NewDefaultStorageBase()
				cfg.TSDB.MutableMemDBTTL = ltoml.Duration(time.Second)
				db.EXPECT().GetOption().Return(&option.DatabaseOption{
					Intervals: option.Intervals{{Interval: timeutil.Interval(commontimeutil.OneSecond)}},
				})
				config.SetGlobalStorageConfig(cfg)
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.mutableMemDB = memDB
				memDB.EXPECT().NumOfMetrics().Return(10)
				memDB.EXPECT().Uptime().Return(time.Minute)
				memDB.EXPECT().MemSize().Return(int64(10))
			},
			needFlush: true,
		},
		{
			name: "trigger size threshold",
			prepare: func(f *dataFamily) {
				cfg := config.NewDefaultStorageBase()
				cfg.TSDB.MutableMemDBTTL = ltoml.Duration(time.Hour)
				db.EXPECT().GetOption().Return(&option.DatabaseOption{
					Intervals: option.Intervals{
						{Interval: timeutil.Interval(commontimeutil.OneSecond)},
						{Interval: timeutil.Interval(commontimeutil.OneMinute * 5)},
					},
				})
				cfg.TSDB.MaxMemDBSize = 10
				config.SetGlobalStorageConfig(cfg)
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.mutableMemDB = memDB
				memDB.EXPECT().NumOfMetrics().Return(10)
				memDB.EXPECT().Uptime().Return(time.Minute)
				memDB.EXPECT().MemSize().Return(int64(1000))
			},
			needFlush: true,
		},
		{
			name: "no trigger any threshold",
			prepare: func(f *dataFamily) {
				cfg := config.NewDefaultStorageBase()
				cfg.TSDB.MutableMemDBTTL = ltoml.Duration(time.Hour)
				cfg.TSDB.MaxMemDBSize = 10000
				db.EXPECT().GetOption().Return(&option.DatabaseOption{
					Intervals: option.Intervals{
						{Interval: timeutil.Interval(commontimeutil.OneSecond)},
						{Interval: timeutil.Interval(commontimeutil.OneMinute * 5)},
					},
				})
				config.SetGlobalStorageConfig(cfg)
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.mutableMemDB = memDB
				memDB.EXPECT().NumOfMetrics().Return(10)
				memDB.EXPECT().Uptime().Return(time.Minute)
				memDB.EXPECT().MemSize().Return(int64(10))
			},
			needFlush: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				config.SetGlobalStorageConfig(config.NewDefaultStorageBase())
			}()
			f := &dataFamily{
				shard:  shard,
				logger: logger.GetLogger("TSDB", "Test"),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			assert.Equal(t, tt.needFlush, f.NeedFlush())
		})
	}
}

func TestDataFamily_Flush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := kv.NewMockFamily(ctrl)
	flusher := kv.NewMockFlusher(ctrl)
	family.EXPECT().NewFlusher().Return(flusher).AnyTimes()
	flusher.EXPECT().Release().AnyTimes()
	flusher.EXPECT().Sequence(gomock.Any(), gomock.Any()).AnyTimes()
	cases := []struct {
		name    string
		prepare func(f *dataFamily)
		wantErr bool
	}{
		{
			name: "flush job doing",
			prepare: func(f *dataFamily) {
				f.isFlushing.Store(true)
			},
			wantErr: false,
		},
		{
			name:    "no data need flush",
			wantErr: false,
		},
		{
			name: "create data flusher failure",
			prepare: func(f *dataFamily) {
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				memDB.EXPECT().NumOfMetrics().Return(100)
				memDB.EXPECT().MarkReadOnly()
				f.mutableMemDB = memDB
				newMetricDataFlusher = func(kvFlusher kv.Flusher) (metricsdata.Flusher, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "flush successfully",
			prepare: func(f *dataFamily) {
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				memDB.EXPECT().NumOfMetrics().Return(100)
				memDB.EXPECT().MarkReadOnly()
				memDB.EXPECT().FlushFamilyTo(gomock.Any()).Return(nil)
				memDB.EXPECT().Close().Return(nil)
				memDB.EXPECT().MemSize().MaxTimes(2)
				f.mutableMemDB = memDB
				dataFlusher := metricsdata.NewMockFlusher(ctrl)
				newMetricDataFlusher = func(kvFlusher kv.Flusher) (metricsdata.Flusher, error) {
					return dataFlusher, nil
				}
			},
			wantErr: false,
		},
		{
			name: "flush metric data failure",
			prepare: func(f *dataFamily) {
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				memDB.EXPECT().NumOfMetrics().Return(100)
				memDB.EXPECT().MarkReadOnly()
				memDB.EXPECT().FlushFamilyTo(gomock.Any()).Return(fmt.Errorf("err"))
				memDB.EXPECT().MemSize()
				f.mutableMemDB = memDB
				dataFlusher := metricsdata.NewMockFlusher(ctrl)
				newMetricDataFlusher = func(kvFlusher kv.Flusher) (metricsdata.Flusher, error) {
					return dataFlusher, nil
				}
			},
			wantErr: true,
		},
		{
			name: "flush successfully, but close memory database failure",
			prepare: func(f *dataFamily) {
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				memDB.EXPECT().NumOfMetrics().Return(100)
				memDB.EXPECT().MarkReadOnly()
				memDB.EXPECT().FlushFamilyTo(gomock.Any()).Return(nil)
				memDB.EXPECT().Close().Return(fmt.Errorf("err"))
				memDB.EXPECT().MemSize().MaxTimes(3)
				f.mutableMemDB = memDB
				dataFlusher := metricsdata.NewMockFlusher(ctrl)
				newMetricDataFlusher = func(kvFlusher kv.Flusher) (metricsdata.Flusher, error) {
					return dataFlusher, nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newMetricDataFlusher = metricsdata.NewFlusher
			}()
			f := &dataFamily{
				family: family,
				seq: map[int32]atomic.Int64{
					1: *atomic.NewInt64(10),
				},
				persistSeq: make(map[int32]atomic.Int64),
				callbacks: map[int32][]func(seq int64){
					1: {func(seq int64) {}},
				},
				statistics: metrics.NewFamilyStatistics("data", "1"),
				logger:     logger.GetLogger("TSDB", "Test"),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			err := f.Flush()
			if (err != nil) != tt.wantErr {
				t.Errorf("Flush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDataFamily_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := kv.NewMockFamily(ctrl)
	flusher := kv.NewMockFlusher(ctrl)
	family.EXPECT().NewFlusher().Return(flusher).AnyTimes()
	flusher.EXPECT().Release().AnyTimes()
	flusher.EXPECT().Sequence(gomock.Any(), gomock.Any()).AnyTimes()
	cases := []struct {
		name    string
		prepare func(f *dataFamily)
		wantErr bool
	}{
		{
			name:    "no data need flush",
			wantErr: false,
		},
		{
			name: "flush immutable mem data failure",
			prepare: func(f *dataFamily) {
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.immutableMemDB = memDB
				newMetricDataFlusher = func(kvFlusher kv.Flusher) (metricsdata.Flusher, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "flush mem data failure",
			prepare: func(f *dataFamily) {
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.mutableMemDB = memDB
				newMetricDataFlusher = func(kvFlusher kv.Flusher) (metricsdata.Flusher, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newMetricDataFlusher = metricsdata.NewFlusher
			}()
			f := &dataFamily{
				family: family,
				seq: map[int32]atomic.Int64{
					1: *atomic.NewInt64(10),
				},
				callbacks: map[int32][]func(seq int64){
					1: {func(seq int64) {}},
				},
				statistics: metrics.NewFamilyStatistics("data", "1"),
				logger:     logger.GetLogger("TSDB", "Test"),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}

			err := f.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDataFamily_MemDBSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := &dataFamily{}
	assert.Equal(t, int64(0), f.MemDBSize())

	memDB := memdb.NewMockMemoryDatabase(ctrl)
	memDB.EXPECT().MemSize().Return(int64(1000))
	f.mutableMemDB = memDB
	assert.Equal(t, int64(1000), f.MemDBSize())
}

func TestDataFamily_GetOrCreateMemoryDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newMemoryDBFunc = memdb.NewMemoryDatabase
		ctrl.Finish()
	}()
	shard := NewMockShard(ctrl)
	db := NewMockDatabase(ctrl)
	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	shard.EXPECT().IndexDB().Return(indexDB).AnyTimes()
	shard.EXPECT().Database().Return(db).AnyTimes()
	db.EXPECT().Name().Return("db").AnyTimes()
	metaDB := index.NewMockMetricMetaDatabase(ctrl)
	db.EXPECT().MetaDB().Return(metaDB).AnyTimes()
	shard.EXPECT().BufferManager().Return(memdb.NewMockBufferManager(ctrl)).AnyTimes()

	f := &dataFamily{
		shard:      shard,
		statistics: metrics.NewFamilyStatistics("data", "1"),
	}
	newMemoryDBFunc = func(cfg *memdb.MemoryDatabaseCfg) (memdb.MemoryDatabase, error) {
		return nil, fmt.Errorf("err")
	}
	memDB, err := f.GetOrCreateMemoryDatabase(1)
	assert.Error(t, err)
	assert.Nil(t, memDB)
	memDB2 := memdb.NewMockMemoryDatabase(ctrl)
	newMemoryDBFunc = func(cfg *memdb.MemoryDatabaseCfg) (memdb.MemoryDatabase, error) {
		return memDB2, nil
	}

	memDB, err = f.GetOrCreateMemoryDatabase(1)
	assert.NoError(t, err)
	assert.Equal(t, memDB2, memDB)

	f.mutableMemDB = memDB
	memDB, err = f.GetOrCreateMemoryDatabase(1)
	assert.NoError(t, err)
	assert.Equal(t, memDB2, memDB)
}

func TestDataFamily_Sequence(t *testing.T) {
	f := &dataFamily{
		seq: make(map[int32]atomic.Int64),
		persistSeq: map[int32]atomic.Int64{
			1: *atomic.NewInt64(10),
		},
		callbacks: make(map[int32][]func(seq int64)),
		logger:    logger.GetLogger("TSDB", "Test"),
	}
	f.CommitSequence(1, 10)
	assert.True(t, f.ValidateSequence(2, 10))
	assert.False(t, f.ValidateSequence(1, 5))
	c := 0
	f.AckSequence(2, func(_ int64) {
		c++
	})
	assert.Equal(t, 0, c)
	f.AckSequence(1, func(_ int64) {
		c++
	})
	assert.Equal(t, 1, c)
}

func TestDataFamily_WriteRows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memDB := memdb.NewMockMemoryDatabase(ctrl)
	memDB.EXPECT().CompleteWrite().AnyTimes()
	memDB.EXPECT().AcquireWrite().AnyTimes()
	memDB.EXPECT().MemSize().Return(int64(10)).AnyTimes()
	shard := NewMockShard(ctrl)
	db := NewMockDatabase(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	indexDB := index.NewMockMetricIndexDatabase(ctrl)
	shard.EXPECT().IndexDB().Return(indexDB).AnyTimes()
	metaDB := index.NewMockMetricMetaDatabase(ctrl)
	db.EXPECT().MetaDB().Return(metaDB).AnyTimes()
	db.EXPECT().Name().Return("db").AnyTimes()
	shard.EXPECT().BufferManager().Return(memdb.NewMockBufferManager(ctrl)).AnyTimes()

	cases := []struct {
		name    string
		prepare func() []metric.StorageRow
		wantErr bool
	}{
		{
			name: "no rows",
			prepare: func() []metric.StorageRow {
				return nil
			},
			wantErr: false,
		},
		{
			name: "get memory database failure",
			prepare: func() []metric.StorageRow {
				newMemoryDBFunc = func(cfg *memdb.MemoryDatabaseCfg) (memdb.MemoryDatabase, error) {
					return nil, fmt.Errorf("err")
				}
				return mockBatchRows(&protoMetricsV1.Metric{
					Name:      "test",
					Timestamp: commontimeutil.Now(),
					SimpleFields: []*protoMetricsV1.SimpleField{{
						Name:  "f1",
						Value: 1.0,
						Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
					}},
				})
			},
			wantErr: true,
		},
		{
			name: "write metric failure",
			prepare: func() []metric.StorageRow {
				memDB.EXPECT().WriteRow(gomock.Any()).Return(fmt.Errorf("err"))
				rows := mockBatchRows(&protoMetricsV1.Metric{
					Name:      "test",
					Timestamp: commontimeutil.Now(),
					SimpleFields: []*protoMetricsV1.SimpleField{{
						Name:  "f1",
						Value: 1.0,
						Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
					}},
				})
				return rows
			},
			wantErr: false,
		},
		{
			name: "write metric successfully",
			prepare: func() []metric.StorageRow {
				memDB.EXPECT().WriteRow(gomock.Any()).Return(nil)
				rows := mockBatchRows(&protoMetricsV1.Metric{
					Name:      "test",
					Timestamp: commontimeutil.Now(),
					SimpleFields: []*protoMetricsV1.SimpleField{{
						Name:  "f1",
						Value: 1.0,
						Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
					}},
				})
				return rows
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newMemoryDBFunc = memdb.NewMemoryDatabase
			}()
			f := &dataFamily{
				shard:      shard,
				interval:   timeutil.Interval(10 * commontimeutil.OneSecond),
				statistics: metrics.NewFamilyStatistics("data", "1"),
				logger:     logger.GetLogger("TSDB", "Test"),
			}
			f.intervalCalc = f.interval.Calculator()
			newMemoryDBFunc = func(cfg *memdb.MemoryDatabaseCfg) (memdb.MemoryDatabase, error) {
				return memDB, nil
			}
			rows := tt.prepare()
			err := f.WriteRows(rows)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteRows() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDataFamily_GetState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard := NewMockShard(ctrl)
	shard.EXPECT().ShardID().Return(models.ShardID(1))
	db := memdb.NewMockMemoryDatabase(ctrl)
	db.EXPECT().NumOfMetrics().Return(10).MaxTimes(2)
	db.EXPECT().NumOfSeries().Return(100).MaxTimes(2)
	db.EXPECT().MemSize().Return(int64(10)).MaxTimes(2)
	db.EXPECT().Uptime().Return(time.Duration(10)).MaxTimes(2)
	now := commontimeutil.Now()
	f := &dataFamily{
		shard:          shard,
		familyTime:     now,
		mutableMemDB:   db,
		immutableMemDB: db,
		seq:            map[int32]atomic.Int64{10: *atomic.NewInt64(10)},
		persistSeq:     map[int32]atomic.Int64{10: *atomic.NewInt64(10)},
	}

	state := f.GetState()
	assert.Equal(t, models.DataFamilyState{
		ShardID:          models.ShardID(1),
		FamilyTime:       commontimeutil.FormatTimestamp(now, commontimeutil.DataTimeFormat2),
		AckSequences:     map[int32]int64{10: 10},
		ReplicaSequences: map[int32]int64{10: 10},
		MemoryDatabases: []models.MemoryDatabaseState{
			{
				State:        "immutable",
				Uptime:       time.Duration(10),
				MemSize:      10,
				NumOfMetrics: 10,
				NumOfSeries:  100,
			}, {
				State:        "mutable",
				Uptime:       time.Duration(10),
				MemSize:      10,
				NumOfMetrics: 10,
				NumOfSeries:  100,
			}},
	}, state)
}

func TestDataFamily_Compact(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	f := &dataFamily{
		lastFlushTime: fasttime.UnixMilliseconds(),
		mutableMemDB:  memdb.NewMockMemoryDatabase(ctrl),
	}
	f.Compact()

	f.mutableMemDB = nil
	f.lastFlushTime = fasttime.UnixMilliseconds() - 2*commontimeutil.OneHour - 5*commontimeutil.OneMinute
	kvFamily := kv.NewMockFamily(ctrl)
	f.family = kvFamily
	kvFamily.EXPECT().Compact()
	f.Compact()
}

func TestDataFamily_Evict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	shard := NewMockShard(ctrl)
	db := NewMockDatabase(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
	opt := &option.DatabaseOption{Ahead: "1h", Behind: "1h"}
	db.EXPECT().GetOption().Return(opt).AnyTimes()
	segment := NewMockSegment(ctrl)
	segment.EXPECT().EvictFamily(gomock.Any()).AnyTimes()

	cases := []struct {
		name    string
		prepare func(f *dataFamily)
	}{
		{
			name: "family write data",
			prepare: func(f *dataFamily) {
				f.Retain()
			},
		},
		{
			name: "family has mem database",
			prepare: func(f *dataFamily) {
				f.Retain()
				f.Release()
				f.mutableMemDB = memdb.NewMockMemoryDatabase(ctrl)
			},
		},
		{
			name: "family time in write time range",
			prepare: func(f *dataFamily) {
				f.familyTime = commontimeutil.Now()
				f.familyTime = f.familyTime - 6*commontimeutil.OneHour - commontimeutil.OneMinute
			},
		},
		{
			name: "family time expire",
			prepare: func(f *dataFamily) {
				f.familyTime = commontimeutil.Now()
				f.familyTime = f.familyTime - 7*commontimeutil.OneHour - commontimeutil.OneMinute
				f.lastReadTime.Store(commontimeutil.Now() - 3*commontimeutil.OneHour - commontimeutil.OneMinute)
			},
		},
		{
			name: "family time expire, but close family failure",
			prepare: func(f *dataFamily) {
				f.familyTime = commontimeutil.Now()
				f.familyTime = f.familyTime - 7*commontimeutil.OneHour - commontimeutil.OneMinute
				f.lastReadTime.Store(commontimeutil.Now() - 3*commontimeutil.OneHour - commontimeutil.OneMinute)
				closeFamilyFunc = func(_ *dataFamily) error {
					return fmt.Errorf("err")
				}
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				closeFamilyFunc = closeFamily
			}()
			f := &dataFamily{
				shard:        shard,
				segment:      segment,
				lastReadTime: atomic.NewInt64(fasttime.UnixMilliseconds()),
				statistics:   metrics.NewFamilyStatistics("data", "1"),
				logger:       logger.GetLogger("TSDB", "Test"),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			f.Evict()
		})
	}
}
