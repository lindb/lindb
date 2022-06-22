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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/linmetrics"
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
	dataFamily := newDataFamily(shard, timeutil.Interval(timeutil.OneSecond*10), timeRange, 10, family)
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
	now := timeutil.Now()
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
			prepare: func(f *dataFamily) {
				snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "get metric reader data failure",
			prepare: func(f *dataFamily) {
				snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{reader}, nil)
				reader.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			wantErr: false,
			len:     0,
		},
		{
			name: "new metric reader failure",
			prepare: func(f *dataFamily) {
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
			prepare: func(f *dataFamily) {
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
			prepare: func(f *dataFamily) {
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
				familyTime: now,
				family:     family,
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			rs, err := f.Filter(&flow.ShardExecuteContext{
				StorageExecuteCtx: &flow.StorageExecuteContext{
					MetricID: 1,
					Query: &stmtpkg.Query{
						StorageInterval: timeutil.Interval(timeutil.OneMinute),
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
				memDB.EXPECT().Size().Return(0)
			},
			needFlush: false,
		},
		{
			name: "trigger time threshold",
			prepare: func(f *dataFamily) {
				cfg := config.NewDefaultStorageBase()
				cfg.TSDB.MutableMemDBTTL = ltoml.Duration(time.Second)
				config.SetGlobalStorageConfig(cfg)
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.mutableMemDB = memDB
				memDB.EXPECT().Size().Return(10)
				memDB.EXPECT().Uptime().Return(time.Duration(timeutil.Now() - timeutil.OneMinute)).MaxTimes(2)
			},
			needFlush: true,
		},
		{
			name: "trigger size threshold",
			prepare: func(f *dataFamily) {
				cfg := config.NewDefaultStorageBase()
				cfg.TSDB.MutableMemDBTTL = ltoml.Duration(time.Hour)
				cfg.TSDB.MaxMemDBSize = 10
				config.SetGlobalStorageConfig(cfg)
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.mutableMemDB = memDB
				memDB.EXPECT().Size().Return(10)
				memDB.EXPECT().Uptime().Return(time.Duration(timeutil.Now() - timeutil.OneMinute)).MaxTimes(2)
				memDB.EXPECT().MemSize().Return(int64(1000)).MaxTimes(2)
			},
			needFlush: true,
		},
		{
			name: "no trigger any threshold",
			prepare: func(f *dataFamily) {
				cfg := config.NewDefaultStorageBase()
				cfg.TSDB.MutableMemDBTTL = ltoml.Duration(time.Hour)
				cfg.TSDB.MaxMemDBSize = 10000
				config.SetGlobalStorageConfig(cfg)
				memDB := memdb.NewMockMemoryDatabase(ctrl)
				f.mutableMemDB = memDB
				memDB.EXPECT().Size().Return(10)
				memDB.EXPECT().Uptime().Return(time.Duration(timeutil.Now() - timeutil.OneMinute)).MaxTimes(2)
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
				memDB.EXPECT().Size().Return(100)
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
				memDB.EXPECT().Size().Return(100)
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
				memDB.EXPECT().Size().Return(100)
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
				memDB.EXPECT().Size().Return(100)
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
	shard.EXPECT().Database().Return(db).AnyTimes()
	db.EXPECT().Name().Return("db").AnyTimes()
	shard.EXPECT().BufferManager().Return(memdb.NewMockBufferManager(ctrl)).AnyTimes()

	f := &dataFamily{
		shard:      shard,
		statistics: metrics.NewFamilyStatistics("data", "1"),
	}
	newMemoryDBFunc = func(cfg memdb.MemoryDatabaseCfg) (memdb.MemoryDatabase, error) {
		return nil, fmt.Errorf("err")
	}
	memDB, err := f.GetOrCreateMemoryDatabase(1)
	assert.Error(t, err)
	assert.Nil(t, memDB)
	memDB2 := memdb.NewMockMemoryDatabase(ctrl)
	newMemoryDBFunc = func(cfg memdb.MemoryDatabaseCfg) (memdb.MemoryDatabase, error) {
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
		seq:       make(map[int32]atomic.Int64),
		callbacks: make(map[int32][]func(seq int64)),
	}
	f.CommitSequence(1, 10)
	assert.True(t, f.ValidateSequence(2, 10))
	assert.False(t, f.ValidateSequence(1, 5))
	c := 0
	f.AckSequence(2, func(seq int64) {
		c++
	})
	assert.Equal(t, 0, c)
	f.AckSequence(1, func(seq int64) {
		c++
	})
	assert.Equal(t, 1, c)
}

func TestDataFamily_WriteRows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	memDB := memdb.NewMockMemoryDatabase(ctrl)
	memDB.EXPECT().WithLock().Return(func() {}).AnyTimes()
	memDB.EXPECT().CompleteWrite().AnyTimes()
	memDB.EXPECT().AcquireWrite().AnyTimes()
	memDB.EXPECT().MemSize().Return(int64(10)).AnyTimes()
	shard := NewMockShard(ctrl)
	db := NewMockDatabase(ctrl)
	shard.EXPECT().Database().Return(db).AnyTimes()
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
				newMemoryDBFunc = func(cfg memdb.MemoryDatabaseCfg) (memdb.MemoryDatabase, error) {
					return nil, fmt.Errorf("err")
				}
				return mockBatchRows(&protoMetricsV1.Metric{
					Name:      "test",
					Timestamp: timeutil.Now(),
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
			name: "metric is not writable",
			prepare: func() []metric.StorageRow {
				return mockBatchRows(&protoMetricsV1.Metric{
					Name:      "test",
					Timestamp: timeutil.Now(),
					SimpleFields: []*protoMetricsV1.SimpleField{{
						Name:  "f1",
						Value: 1.0,
						Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
					}},
				})
			},
			wantErr: false,
		},
		{
			name: "write metric failure",
			prepare: func() []metric.StorageRow {
				memDB.EXPECT().WriteRow(gomock.Any()).Return(0, fmt.Errorf("err"))
				rows := mockBatchRows(&protoMetricsV1.Metric{
					Name:      "test",
					Timestamp: timeutil.Now(),
					SimpleFields: []*protoMetricsV1.SimpleField{{
						Name:  "f1",
						Value: 1.0,
						Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
					}},
				})
				rows[0].Writable = true
				return rows
			},
			wantErr: false,
		},
		{
			name: "write metric successfully",
			prepare: func() []metric.StorageRow {
				memDB.EXPECT().WriteRow(gomock.Any()).Return(10, nil)
				rows := mockBatchRows(&protoMetricsV1.Metric{
					Name:      "test",
					Timestamp: timeutil.Now(),
					SimpleFields: []*protoMetricsV1.SimpleField{{
						Name:  "f1",
						Value: 1.0,
						Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
					}},
				})
				rows[0].Writable = true
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
				interval:   timeutil.Interval(10 * timeutil.OneSecond),
				statistics: metrics.NewFamilyStatistics("data", "1"),
				logger:     logger.GetLogger("TSDB", "Test"),
			}
			f.intervalCalc = f.interval.Calculator()
			newMemoryDBFunc = func(cfg memdb.MemoryDatabaseCfg) (memdb.MemoryDatabase, error) {
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
