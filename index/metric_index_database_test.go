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
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	v1 "github.com/lindb/lindb/index/v1"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
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

	metaDB, err := NewMetricMetaDatabase(path.Join(name, "meta"))
	assert.NoError(t, err)
	db, err := NewMetricIndexDatabase(path.Join(name, "index"), metaDB)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db1 := db.(*metricIndexDatabase)
	ch := make(chan struct{})
	db.Notify(nil)
	db.Notify(&MetaNotifier{
		Namespace:  "system",
		MetricName: "cpu",
		TagHash:    120,
		Tags: tag.Tags{
			{Key: []byte("key1"), Value: []byte("value1")},
			{Key: []byte("key2"), Value: []byte("value2")},
		},
		Callback: func(_ uint32, _ error) {},
	})
	db.Notify(&MetaNotifier{
		Namespace:  "system",
		MetricName: "cpu",
		TagHash:    120,
		Tags: tag.Tags{
			{Key: []byte("key1"), Value: []byte("value1")},
			{Key: []byte("key2"), Value: []byte("value2")},
		},
		Callback: func(_ uint32, _ error) {},
	})
	db.Notify(&MetaNotifier{
		Namespace:  "system",
		MetricName: "cpu",
		TagHash:    100,
		Tags: tag.Tags{
			{Key: []byte("key3"), Value: []byte("value1")},
			{Key: []byte("key4"), Value: []byte("value2")},
		},
		Callback: func(_ uint32, _ error) {
			ch <- struct{}{}
		},
	})
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

	<-ch

	// wait tag index build completed
	time.Sleep(100 * time.Millisecond)
	// get series ids from memory
	test()

	db.Notify(&FlushNotifier{
		Callback: func(err error) {
			assert.NoError(t, err)
			ch <- struct{}{}
		},
	})
	<-ch
	db1.sequenceCache.Purge()
	// get series ids from kv
	test()

	db.Notify(&MetaNotifier{
		Namespace:  "system",
		MetricName: "cpu",
		TagHash:    1000,
		Tags: tag.Tags{
			{Key: []byte("key3"), Value: []byte("value1")},
			{Key: []byte("key4"), Value: []byte("value2")},
		},
		Callback: func(_ uint32, _ error) {
			ch <- struct{}{}
		},
	})
	<-ch
	time.Sleep(100 * time.Millisecond)
	// get series ids by metric
	seriesIDs, err := db.GetSeriesIDsForMetric(0)
	assert.NoError(t, err)
	assert.Equal(t, []uint32{0, 1, 2}, seriesIDs.ToArray())
	// flushing
	db1.flushing.Store(true)
	db.Notify(&FlushNotifier{
		Callback: func(err error) {
			assert.NoError(t, err)
			ch <- struct{}{}
		},
	})
	<-ch
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
		name    string
		prepare func()
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

func TestMetricIndexDatabase_Notify_Error(t *testing.T) {
	name := "./notify_error"
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = os.RemoveAll(name)
	}()
	metaDB := NewMockMetricMetaDatabase(ctrl)
	db, err := NewMetricIndexDatabase(name, metaDB)
	assert.NoError(t, err)
	metaDB.EXPECT().Notify(gomock.Any()).Do(func(n Notifier) {
		mn := n.(*MetaNotifier)
		mn.Callback(0, fmt.Errorf("err"))
	}).AnyTimes()
	c := 0
	db.Notify(&MetaNotifier{
		Callback: func(_ uint32, err error) {
			assert.Error(t, err)
			c++
		},
	})
	assert.Equal(t, 1, c)
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
	invertedFlusher := v1.NewMockInvertedIndexFlusher(ctrl)
	forwardFlusher := v1.NewMockForwardIndexFlusher(ctrl)
	forwardFlusher.EXPECT().Close().Return(nil).AnyTimes()
	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "flush metric index error",
			prepare: func() {
				newInvertedIndexFlusher = func(kvFlusher kv.Flusher) (v1.InvertedIndexFlusher, error) {
					return nil, fmt.Errorf("err")
				}
			},
		},
		{
			name: "flush forward index error",
			prepare: func() {
				newInvertedIndexFlusher = func(kvFlusher kv.Flusher) (v1.InvertedIndexFlusher, error) {
					return invertedFlusher, nil
				}
				newForwardIndexFlusher = func(kvFlusher kv.Flusher) (v1.ForwardIndexFlusher, error) {
					return nil, fmt.Errorf("err")
				}
				invertedFlusher.EXPECT().Close().Return(nil)
			},
		},
		{
			name: "flush inverted error",
			prepare: func() {
				newInvertedIndexFlusher = func(kvFlusher kv.Flusher) (v1.InvertedIndexFlusher, error) {
					return invertedFlusher, nil
				}
				newForwardIndexFlusher = func(kvFlusher kv.Flusher) (v1.ForwardIndexFlusher, error) {
					return forwardFlusher, nil
				}
				invertedFlusher.EXPECT().Close().Return(nil)
				invertedFlusher.EXPECT().Close().Return(fmt.Errorf("err"))
			},
		},
		{
			name: "flush series error",
			prepare: func() {
				newInvertedIndexFlusher = func(kvFlusher kv.Flusher) (v1.InvertedIndexFlusher, error) {
					return invertedFlusher, nil
				}
				newForwardIndexFlusher = func(kvFlusher kv.Flusher) (v1.ForwardIndexFlusher, error) {
					return forwardFlusher, nil
				}
				invertedFlusher.EXPECT().Close().Return(nil).MaxTimes(2)
				indexStore.EXPECT().Flush().Return(fmt.Errorf("err"))
			},
		},
	}
	db := &metricIndexDatabase{
		metricInverted: newInvertedIndex(kvFamily),
		forward:        newForwardIndex(kvFamily),
		inverted:       newInvertedIndex(kvFamily),
		series:         indexStore,
	}
	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newInvertedIndexFlusher = v1.NewInvertedIndexFlusher
				newForwardIndexFlusher = v1.NewForwardIndexFlusher
			}()
			tt.prepare()
			db.PrepareFlush()
			assert.Error(t, db.Flush())
		})
	}
}

func TestMetricIndexDatabase_GenSeries_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	kvFamily := kv.NewMockFamily(ctrl)
	snapshot := version.NewMockSnapshot(ctrl)
	snapshot.EXPECT().Close()
	kvFamily.EXPECT().GetSnapshot().Return(snapshot)
	snapshot.EXPECT().Load(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	indexStore := NewMockIndexKVStore(ctrl)
	db := &metricIndexDatabase{
		series:         indexStore,
		sequenceCache:  expirable.NewLRU[metric.ID, uint32](100, nil, time.Minute),
		metricInverted: newInvertedIndex(kvFamily),
	}
	indexStore.EXPECT().GetOrCreateValue(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ uint32, _ []byte, createFn func() uint32) (uint32, error) {
			createFn()
			return 0, nil
		})
	c := 0
	db.handle(&MetaNotifier{
		Namespace:  "system",
		MetricName: "cpu",
		TagHash:    10,
		Callback: func(_ uint32, err error) {
			assert.Error(t, err)
			c++
		},
	})
	assert.Equal(t, 1, c)
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
		name    string
		prepare func()
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
		name    string
		prepare func()
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
