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

package index

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/lindb/common/pkg/logger"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	v1 "github.com/lindb/lindb/index/v1"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

func TestMetricIndexDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	name := "./metric_index_database"
	defer func() {
		_ = os.RemoveAll(name)
		ctrl.Finish()
	}()

	metaDB, err := NewMetricMetaDatabase("test", path.Join(name, "meta"))
	assert.NoError(t, err)
	db, err := NewMetricIndexDatabase(path.Join(name, "index"), metaDB)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db1 := db.(*metricIndexDatabase)
	m := &protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		Tags:      []*protoMetricsV1.KeyValue{{Key: "key1", Value: "value1"}},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	}
	row1 := protoToStorageRow(m)
	_, err = db.GenSeriesID(0, row1)
	assert.NoError(t, err)
	_, err = db.GenSeriesID(0, row1)
	assert.NoError(t, err)
	m.Tags = []*protoMetricsV1.KeyValue{{Key: "key2", Value: "value1"}}
	row2 := protoToStorageRow(m)
	_, err = db.GenSeriesID(0, row2)
	assert.NoError(t, err)
	m.Tags = []*protoMetricsV1.KeyValue{{Key: "key2", Value: "value0000"}}
	row2 = protoToStorageRow(m)
	limits := models.NewDefaultLimits()
	limits.MaxSeriesPerMetric = 1
	models.SetDatabaseLimits("test", limits)
	_, err = db.GenSeriesID(0, row2)
	assert.Equal(t, constants.ErrTooManySeries, err)
	models.SetDatabaseLimits("test", models.NewDefaultLimits())
	test := func() {
		// get series ids by metric
		seriesIDs, err0 := db.GetSeriesIDsForMetric(0)
		assert.NoError(t, err0)
		assert.Equal(t, []uint32{0, 1}, seriesIDs.ToArray())
		// metric not exist
		seriesIDs, err = db.GetSeriesIDsForMetric(100)
		assert.NoError(t, err)
		assert.Empty(t, seriesIDs.ToArray())
		// get series ids by tag key
		seriesIDs, err = db.GetSeriesIDsForTag(0)
		assert.NoError(t, err)
		assert.Equal(t, []uint32{0}, seriesIDs.ToArray())
		// tag key not exist
		seriesIDs, err = db.GetSeriesIDsForTag(1000)
		assert.NoError(t, err)
		assert.Empty(t, seriesIDs.ToArray())
		// get series ids by tag value
		seriesIDs, err = db.GetSeriesIDsByTagValueIDs(0, roaring.BitmapOf(0))
		assert.NoError(t, err)
		assert.Equal(t, []uint32{0}, seriesIDs.ToArray())
		// tag value not exist
		seriesIDs, err = db.GetSeriesIDsByTagValueIDs(0, roaring.BitmapOf(100))
		assert.NoError(t, err)
		assert.Empty(t, seriesIDs.ToArray())

		assert.NoError(t, db.GetGroupingContext(&flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{
				GroupByTagKeyIDs: []tag.KeyID{0},
			},
			SeriesIDsAfterFiltering: roaring.BitmapOf(0, 1, 2),
		}))
		assert.Equal(t, constants.ErrNotFound, db.GetGroupingContext(&flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{
				GroupByTagKeyIDs: []tag.KeyID{0},
			},
			SeriesIDsAfterFiltering: roaring.BitmapOf(100, 200), // series ids not found
		}))
	}
	// wait tag index build completed
	time.Sleep(100 * time.Millisecond)
	// get series ids from memory
	test()
	db.PrepareFlush()
	assert.NoError(t, db.Flush())
	db1.sequenceCache.Purge()
	// get series ids from kv
	test()
	m.Tags = []*protoMetricsV1.KeyValue{{Key: "key3", Value: "value1"}}
	row3 := protoToStorageRow(m)
	_, err = db.GenSeriesID(0, row3)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	// get series ids by metric
	seriesIDs, err := db.GetSeriesIDsForMetric(0)
	assert.NoError(t, err)
	assert.Equal(t, []uint32{0, 1, 2}, seriesIDs.ToArray())
	// flushing
	db1.flushing.Store(true)
	assert.NoError(t, db.Close())
}

func TestMetricIndexDatabase_New_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	oldMgr := kv.GetStoreManager()
	defer func() {
		ctrl.Finish()
		newSequence = NewSequence
		kv.InitStoreManager(oldMgr)
	}()
	kvStore := kv.NewMockStore(ctrl)
	storeMgr := kv.NewMockStoreManager(ctrl)
	kv.InitStoreManager(storeMgr)

	cases := []struct {
		prepare func()
		name    string
	}{
		{
			name: "create kv store error",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
		},
		{
			name: "create series family error",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(kvStore, nil)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
		},
		{
			name: "create inverted index family error",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(kvStore, nil)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
		},
		{
			name: "create metric series family error",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(kvStore, nil)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil).MaxTimes(2)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
		},
		{
			name: "create forward index family error",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(kvStore, nil)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil).MaxTimes(3)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
		},
	}

	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()

			_, err := NewMetricIndexDatabase("./dir", nil)
			assert.Error(t, err)
		})
	}
}

func TestMetricIndexDatabase_Flush_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	kvFamily := kv.NewMockFamily(ctrl)
	flush := kv.NewMockFlusher(ctrl)
	flush.EXPECT().Release().AnyTimes()
	kvFamily.EXPECT().NewFlusher().Return(flush).AnyTimes()
	indexStore := NewMockIndexKVStore(ctrl)
	indexStore.EXPECT().PrepareFlush().AnyTimes()
	forwardFlusher := v1.NewMockForwardIndexFlusher(ctrl)
	forwardFlusher.EXPECT().Close().Return(nil).AnyTimes()
	db := &metricIndexDatabase{
		metricInverted: newInvertedIndex(kvFamily),
		forward:        newForwardIndex(kvFamily),
		inverted:       newInvertedIndex(kvFamily),
		series:         indexStore,
	}
	cases := []struct {
		prepare func()
		name    string
	}{
		{
			name: "flush metric index error",
			prepare: func() {
				db.metricInverted.immutable = imap.NewIntMap[*roaring.Bitmap]()
				db.metricInverted.immutable.Put(10, roaring.New())
				newInvertedIndexFlusher = func(kvFlusher kv.Flusher) (v1.InvertedIndexFlusher, error) {
					return nil, fmt.Errorf("err")
				}
			},
		},
		{
			name: "flush forward index error",
			prepare: func() {
				db.metricInverted.immutable = nil
				db.forward.immutable = imap.NewIntMap[*imap.IntMap[uint32]]()
				db.forward.immutable.Put(1, nil)
				newForwardIndexFlusher = func(kvFlusher kv.Flusher) (v1.ForwardIndexFlusher, error) {
					return nil, fmt.Errorf("err")
				}
			},
		},
		{
			name: "flush inverted error",
			prepare: func() {
				db.metricInverted.immutable = nil
				db.forward.immutable = nil
				db.inverted.immutable = imap.NewIntMap[*roaring.Bitmap]()
				db.inverted.immutable.Put(10, roaring.New())
				newInvertedIndexFlusher = func(kvFlusher kv.Flusher) (v1.InvertedIndexFlusher, error) {
					return nil, fmt.Errorf("err")
				}
			},
		},
		{
			name: "flush series error",
			prepare: func() {
				db.metricInverted.immutable = nil
				db.forward.immutable = nil
				db.inverted.immutable = nil
				indexStore.EXPECT().Flush().Return(fmt.Errorf("err"))
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newInvertedIndexFlusher = v1.NewInvertedIndexFlusher
				newForwardIndexFlusher = v1.NewForwardIndexFlusher
			}()
			db.flushing.Store(false)
			tt.prepare()
			db.PrepareFlush()
			assert.Error(t, db.Flush())
		})
	}
}

func TestMetricIndexDatabase_ForwardIndex_Flush_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newForwardIndexFlusher = v1.NewForwardIndexFlusher
		ctrl.Finish()
	}()

	kvFamily := kv.NewMockFamily(ctrl)
	kvFlusher := kv.NewMockFlusher(ctrl)
	kvFamily.EXPECT().NewFlusher().Return(kvFlusher).AnyTimes()
	kvFlusher.EXPECT().Release().AnyTimes()
	flusher := v1.NewMockForwardIndexFlusher(ctrl)
	newForwardIndexFlusher = func(kvFlusher kv.Flusher) (v1.ForwardIndexFlusher, error) {
		return flusher, nil
	}
	index := newForwardIndex(kvFamily)

	cases := []struct {
		prepare func()
		name    string
		wantErr bool
	}{
		{
			name: "flusher close error",
			prepare: func() {
				index.immutable = imap.NewIntMap[*imap.IntMap[uint32]]()
				index.immutable.Put(10, imap.NewIntMap[uint32]())
				flusher.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "write series ids error",
			prepare: func() {
				flusher.EXPECT().Prepare(gomock.Any())
				flusher.EXPECT().WriteSeriesIDs(gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "write tag value ids error",
			prepare: func() {
				flusher.EXPECT().Prepare(gomock.Any())
				flusher.EXPECT().WriteSeriesIDs(gomock.Any()).Return(nil)
				flusher.EXPECT().WriteTagValueIDs(gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "flusher commit error",
			prepare: func() {
				flusher.EXPECT().Prepare(gomock.Any())
				flusher.EXPECT().WriteSeriesIDs(gomock.Any()).Return(nil)
				flusher.EXPECT().WriteTagValueIDs(gomock.Any()).Return(nil)
				flusher.EXPECT().Commit().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
	}

	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			index.immutable = imap.NewIntMap[*imap.IntMap[uint32]]()
			tags := imap.NewIntMap[uint32]()
			tags.Put(10, 10)
			index.immutable.Put(100, tags)
			tt.prepare()
			err := index.flush()
			if (err != nil) != tt.wantErr {
				t.Fatal(tt.name)
			}
		})
	}
}

func TestMetricIndexDatabase_Inverted_Flush_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newInvertedIndexFlusher = v1.NewInvertedIndexFlusher
		ctrl.Finish()
	}()

	kvFamily := kv.NewMockFamily(ctrl)
	kvFlusher := kv.NewMockFlusher(ctrl)
	kvFamily.EXPECT().NewFlusher().Return(kvFlusher).AnyTimes()
	kvFlusher.EXPECT().Release().AnyTimes()
	flusher := v1.NewMockInvertedIndexFlusher(ctrl)
	flusher.EXPECT().Prepare(gomock.Any()).AnyTimes()
	newInvertedIndexFlusher = func(kvFlusher kv.Flusher) (v1.InvertedIndexFlusher, error) {
		return flusher, nil
	}
	index := newInvertedIndex(kvFamily)
	cases := []struct {
		prepare func()
		name    string
		wantErr bool
	}{
		{
			name: "flusher close error",
			prepare: func() {
				index.immutable = imap.NewIntMap[*roaring.Bitmap]()
				index.immutable.Put(100, roaring.BitmapOf())
				flusher.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "write series error",
			prepare: func() {
				flusher.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "flusher commit error",
			prepare: func() {
				flusher.EXPECT().Write(gomock.Any()).Return(nil)
				flusher.EXPECT().Commit().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
	}
	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			index.immutable = imap.NewIntMap[*roaring.Bitmap]()
			index.immutable.Put(100, roaring.BitmapOf(1, 2, 3))
			tt.prepare()
			err := index.flush()
			if (err != nil) != tt.wantErr {
				t.Fatal(tt.name)
			}
		})
	}
}

func TestMetricIndexDatabase_Inverted_Read_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		bitmapUnmarshal = encoding.BitmapUnmarshal
		ctrl.Finish()
	}()

	kvFamily := kv.NewMockFamily(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	kvFamily.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	snapshot.EXPECT().Close().AnyTimes()
	index := newInvertedIndex(kvFamily)

	t.Run("unmarshal bitmap error when get series ids", func(t *testing.T) {
		bitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) (int64, error) {
			return 0, fmt.Errorf("err")
		}
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).DoAndReturn(func(_ uint32, loader func(data []byte) error) error {
			return loader(nil)
		})
		ids, err := index.getSeriesIDs(10)
		assert.Error(t, err)
		assert.Nil(t, ids)
	})

	t.Run("unmarshal bitmap error when find series ids", func(t *testing.T) {
		bitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) (int64, error) {
			return 0, fmt.Errorf("err")
		}
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).DoAndReturn(func(_ uint32, loader func(data []byte) error) error {
			return loader(nil)
		})
		ids, err := index.findSeriesIDsByKeys(roaring.BitmapOf(1, 2, 3))
		assert.Error(t, err)
		assert.Nil(t, ids)
	})
	t.Run("unmarshal bitmap error when find series ids", func(t *testing.T) {
		snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
		ids, err := index.findSeriesIDsByKeys(roaring.BitmapOf(1, 2, 3))
		assert.Error(t, err)
		assert.Nil(t, ids)
	})
}

func TestMetricIndexDatabase_ForwardIndex_Read_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newForwardReader = v1.NewForwardReader
		ctrl.Finish()
	}()

	kvFamily := kv.NewMockFamily(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	kvFamily.EXPECT().GetSnapshot().Return(snapshot).AnyTimes()
	snapshot.EXPECT().Close().AnyTimes()
	index := newForwardIndex(kvFamily)
	reader := v1.NewMockForwardReader(ctrl)

	t.Run("find kv reader error", func(t *testing.T) {
		snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
		ids, err := index.findSeriesIDsForTag(10)
		assert.Error(t, err)
		assert.Nil(t, ids)
	})

	t.Run("read series ids error", func(t *testing.T) {
		snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{nil}, nil)
		newForwardReader = func(readers []table.Reader) v1.ForwardReader {
			return reader
		}
		reader.EXPECT().GetSeriesIDsForTagKeyID(gomock.Any()).Return(nil, fmt.Errorf("err"))
		ids, err := index.findSeriesIDsForTag(10)
		assert.Error(t, err)
		assert.Nil(t, ids)
	})

	t.Run("find kv reader error when grouping", func(t *testing.T) {
		snapshot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
		assert.Error(t, index.GetGroupingContext(&flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{
				GroupByTagKeyIDs: []tag.KeyID{0},
			},
			SeriesIDsAfterFiltering: roaring.BitmapOf(0, 1, 2),
		}))
	})
	t.Run("read group scaner error when grouping", func(t *testing.T) {
		snapshot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{nil}, nil)
		newForwardReader = func(readers []table.Reader) v1.ForwardReader {
			return reader
		}
		reader.EXPECT().GetGroupingScanner(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
		assert.Error(t, index.GetGroupingContext(&flow.ShardExecuteContext{
			StorageExecuteCtx: &flow.StorageExecuteContext{
				GroupByTagKeyIDs: []tag.KeyID{0},
			},
			SeriesIDsAfterFiltering: roaring.BitmapOf(0, 1, 2),
		}))
	})
}

func TestMetricIndexDatabase_createSeriesID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := kv.NewMockFamily(ctrl)
	idx := &metricIndexDatabase{
		sequenceCache:  expirable.NewLRU[metric.ID, uint32](100000, nil, time.Hour),
		metricInverted: newInvertedIndex(family),
	}
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close()
	family.EXPECT().GetSnapshot().Return(snapshot)
	snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	seriesID := idx.createSeriesID(10)
	assert.Zero(t, seriesID)
}

func TestMetricIndexDatabase_buildInvertIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metaDB := NewMockMetricMetaDatabase(ctrl)
	m := &protoMetricsV1.Metric{
		Name:      "test1",
		Namespace: "ns",
		Tags:      []*protoMetricsV1.KeyValue{{Key: "key1", Value: "value1"}},
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
		},
	}
	row1 := protoToStorageRow(m)
	idx := &metricIndexDatabase{
		metaDB:     metaDB,
		statistics: metrics.NewIndexDBStatistics("test"),
		logger:     logger.GetLogger("test", "test"),
	}
	metaDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(tag.KeyID(0), fmt.Errorf("err"))
	idx.buildInvertIndex(1, row1.NewKeyValueIterator(), 1)
	metaDB.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(tag.KeyID(0), nil)
	metaDB.EXPECT().GenTagValueID(gomock.Any(), gomock.Any()).Return(uint32(0), fmt.Errorf("err"))
	idx.buildInvertIndex(1, row1.NewKeyValueIterator(), 1)
}

func protoToStorageRow(m *protoMetricsV1.Metric) *metric.StorageRow {
	var ml protoMetricsV1.MetricList
	ml.Metrics = append(ml.Metrics, m)
	var buf bytes.Buffer
	converter := metric.NewProtoConverter(models.NewDefaultLimits())
	_, err := converter.MarshalProtoMetricListV1To(ml, &buf)
	if err != nil {
		panic(err)
	}

	var br metric.StorageBatchRows
	br.UnmarshalRows(buf.Bytes())
	return br.Rows()[0]
}
