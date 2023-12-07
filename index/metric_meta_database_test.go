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
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	commonfileutil "github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/sql/stmt"
)

func TestMetricMetaDatabase(t *testing.T) {
	name := "./metric_meta_database"
	defer func() {
		_ = os.RemoveAll(name)
	}()

	db, err := NewMetricMetaDatabase(name)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	db.Notify(&MetaNotifier{
		Namespace:  "system",
		MetricName: "cpu",
		Callback:   func(id uint32, err error) {},
	})
	db.Notify(&FieldNotifier{
		Namespace:  "system",
		MetricName: "cpu",
		Field:      field.Meta{Name: "f1", Type: field.SumField},
		Callback:   func(fieldID field.ID, err error) {},
	})
	ch := make(chan struct{})
	db.Notify(&TagNotifier{
		metricID: 0,
		tags: tag.Tags{
			{Key: []byte("key1"), Value: []byte("value1")},
			{Key: []byte("key2"), Value: []byte("value2")},
		},
		buildIndex: func(_, _ uint32) {
			ch <- struct{}{}
		},
	})
	// two tags
	<-ch
	<-ch

	test := func() {
		mid, err := db.GetMetricID("system", "cpu")
		assert.Equal(t, metric.ID(0), mid)
		assert.NoError(t, err)
		schema, err := db.GetSchema(mid)
		assert.NoError(t, err)
		assert.NotNil(t, schema)
		assert.Len(t, schema.Fields, 1)
		assert.Equal(t, field.ID(0), schema.Fields[0].ID)
		assert.Equal(t, field.SumField, schema.Fields[0].Type)
		assert.Len(t, schema.TagKeys, 2)
		assert.Equal(t, "key1", schema.TagKeys[0].Key)
		assert.Equal(t, tag.KeyID(0), schema.TagKeys[0].ID)
		assert.Equal(t, "key2", schema.TagKeys[1].Key)
		assert.Equal(t, tag.KeyID(1), schema.TagKeys[1].ID)
		ids, err := db.FindTagValueDsByExpr(0, &stmt.EqualsExpr{Key: "key1", Value: "value1"})
		assert.NoError(t, err)
		assert.Equal(t, []uint32{0}, ids.ToArray())
		ids, err = db.FindTagValueIDsForTag(1)
		assert.NoError(t, err)
		assert.Equal(t, []uint32{1}, ids.ToArray())
		rs := make(map[uint32]string)
		err = db.CollectTagValues(0, roaring.BitmapOf(0), rs)
		assert.NoError(t, err)
		assert.Equal(t, map[uint32]string{0: "value1"}, rs)
		result, err := db.SuggestNamespace("sy", 1)
		assert.NoError(t, err)
		assert.Equal(t, []string{"system"}, result)
		result, err = db.SuggestNamespace("tt", 10)
		assert.NoError(t, err)
		assert.Empty(t, result)
		result, err = db.SuggestMetrics("system", "c", 10)
		assert.NoError(t, err)
		assert.Equal(t, []string{"cpu"}, result)
		result, err = db.SuggestMetrics("system1", "c", 10)
		assert.NoError(t, err)
		assert.Empty(t, result)
		result, err = db.SuggestMetrics("system", "tt", 10)
		assert.NoError(t, err)
		assert.Empty(t, result)
		result, err = db.SuggestTagValues(0, "val", 10)
		assert.NoError(t, err)
		assert.Equal(t, []string{"value1"}, result)
		result, err = db.SuggestTagValues(0, "tt", 10)
		assert.NoError(t, err)
		assert.Empty(t, result)
		result, err = db.SuggestTagValues(100, "val", 10)
		assert.NoError(t, err)
		assert.Empty(t, result)
	}

	// from memory
	test()

	db.Notify(&FlushNotifier{
		Callback: func(err error) {
			assert.NoError(t, err)
			ch <- struct{}{}
		},
	})

	<-ch

	// from kv store
	test()

	// flushing
	db1 := db.(*metricMetaDatabase)
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

func TestMetricMetaDatabase_Flush_Error(t *testing.T) {
	name := "./metric_meta_database_flush_error"
	ctrl := gomock.NewController(t)
	defer func() {
		syncFn = fileutil.Sync
		unmapFn = fileutil.Unmap
		_ = os.RemoveAll(name)
		ctrl.Finish()
	}()

	db, err := NewMetricMetaDatabase(name)
	assert.NoError(t, err)
	t.Run("sync sequence error", func(t *testing.T) {
		syncFn = func(_ []byte) error {
			return fmt.Errorf("err")
		}
		assert.Error(t, db.Flush())
	})
	syncFn = fileutil.Sync
	db1 := db.(*metricMetaDatabase)
	kvStore := NewMockIndexKVStore(ctrl)
	schemaStore := NewMockMetricSchemaStore(ctrl)
	db1.ns = kvStore
	db1.metric = kvStore
	db1.tagValue = kvStore
	db1.schemaStore = schemaStore
	t.Run("flush ns error", func(t *testing.T) {
		kvStore.EXPECT().Flush().Return(fmt.Errorf("err"))
		assert.Error(t, db.Flush())
	})

	t.Run("flush metric error", func(t *testing.T) {
		kvStore.EXPECT().Flush().Return(nil)
		kvStore.EXPECT().Flush().Return(fmt.Errorf("err"))
		assert.Error(t, db.Flush())
	})
	t.Run("flush metric schema error", func(t *testing.T) {
		kvStore.EXPECT().Flush().Return(nil).MaxTimes(2)
		schemaStore.EXPECT().Flush().Return(fmt.Errorf("err"))
		assert.Error(t, db.Flush())
	})
	t.Run("flush tag value error", func(t *testing.T) {
		kvStore.EXPECT().Flush().Return(nil).MaxTimes(2)
		schemaStore.EXPECT().Flush().Return(nil)
		kvStore.EXPECT().Flush().Return(fmt.Errorf("err"))
		assert.Error(t, db.Flush())
	})
	t.Run("close error", func(t *testing.T) {
		unmapFn = func(f *os.File, data []byte) error {
			return fmt.Errorf("err")
		}
		assert.Error(t, db.Close())
	})
}

func TestMetricMetaDatabase_New_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	oldMgr := kv.GetStoreManager()
	defer func() {
		ctrl.Finish()
		mkdir = commonfileutil.MkDirIfNotExist
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
			name: "mkdir error",
			prepare: func() {
				mkdir = func(path string) error {
					return fmt.Errorf("err")
				}
			},
		},
		{
			name: "new sequence error",
			prepare: func() {
				newSequence = func(fileName string) (*Sequence, error) {
					return nil, fmt.Errorf("err")
				}
			},
		},
		{
			name: "create kv store error",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
		},
		{
			name: "create ns family error",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(kvStore, nil)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
		},
		{
			name: "create metric family error",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(kvStore, nil)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
		},
		{
			name: "create tag value family error",
			prepare: func() {
				storeMgr.EXPECT().CreateStore(gomock.Any(), gomock.Any()).Return(kvStore, nil)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, nil).MaxTimes(2)
				kvStore.EXPECT().CreateFamily(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
		},
		{
			name: "create metric schema family error",
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
			mkdir = func(path string) error {
				return nil
			}
			newSequence = func(fileName string) (*Sequence, error) {
				return nil, nil
			}
			tt.prepare()

			_, err := NewMetricMetaDatabase("./dir")
			assert.Error(t, err)
		})
	}
}

func TestMetricMetaDatabase_Read_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	indexStore := NewMockIndexKVStore(ctrl)
	db := &metricMetaDatabase{
		tagValue: indexStore,
		ns:       indexStore,
		metric:   indexStore,
	}

	t.Run("read tag value error", func(t *testing.T) {
		indexStore.EXPECT().FindValuesByExpr(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
		ids, err := db.FindTagValueDsByExpr(10, nil)
		assert.Nil(t, ids)
		assert.Error(t, err)

		indexStore.EXPECT().GetValues(uint32(10)).Return(nil, fmt.Errorf("err"))
		ids, err = db.FindTagValueIDsForTag(10)
		assert.Nil(t, ids)
		assert.Error(t, err)
	})
	t.Run("read metric id error", func(t *testing.T) {
		indexStore.EXPECT().GetValue(gomock.Any(), gomock.Any()).Return(uint32(0), false, fmt.Errorf("err"))
		mid, err := db.GetMetricID("system", "cpu")
		assert.Error(t, err)
		assert.Equal(t, metric.ID(0), mid)
		// ns not exit
		indexStore.EXPECT().GetValue(gomock.Any(), gomock.Any()).Return(uint32(0), false, nil)
		mid, err = db.GetMetricID("system", "cpu")
		assert.Error(t, err)
		assert.Equal(t, metric.ID(0), mid)

		indexStore.EXPECT().GetValue(gomock.Any(), gomock.Any()).Return(uint32(0), true, nil)
		indexStore.EXPECT().GetValue(gomock.Any(), gomock.Any()).Return(uint32(0), false, fmt.Errorf("err"))
		mid, err = db.GetMetricID("system", "cpu")
		assert.Error(t, err)
		assert.Equal(t, metric.ID(0), mid)
		indexStore.EXPECT().GetValue(gomock.Any(), gomock.Any()).Return(uint32(0), true, nil)
		indexStore.EXPECT().GetValue(gomock.Any(), gomock.Any()).Return(uint32(0), false, nil)
		mid, err = db.GetMetricID("system", "cpu")
		assert.Error(t, err)
		assert.Equal(t, metric.ID(0), mid)
	})

	t.Run("gen metric id error", func(t *testing.T) {
		indexStore.EXPECT().GetOrCreateValue(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), fmt.Errorf("err"))
		_, err := db.genMetricID([]byte("system"), []byte("cpu"))
		assert.Error(t, err)

		indexStore.EXPECT().GetOrCreateValue(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), nil)
		indexStore.EXPECT().GetOrCreateValue(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), fmt.Errorf("err"))
		_, err = db.genMetricID([]byte("system"), []byte("cpu"))
		assert.Error(t, err)
	})
	t.Run("suggest metric name error", func(t *testing.T) {
		indexStore.EXPECT().GetValue(gomock.Any(), gomock.Any()).Return(uint32(0), false, fmt.Errorf("err"))
		rs, err := db.SuggestMetrics("system", "cpu", 10)
		assert.Error(t, err)
		assert.Empty(t, rs)
	})
	t.Run("suggest ns error", func(t *testing.T) {
		indexStore.EXPECT().Suggest(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
		rs, err := db.SuggestNamespace("system", 10)
		assert.Error(t, err)
		assert.Empty(t, rs)
	})
}

func TestMetricMetaDatabase_GenTag_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	schemaStore := NewMockMetricSchemaStore(ctrl)
	indexStore := NewMockIndexKVStore(ctrl)
	db := &metricMetaDatabase{
		tagValue:    indexStore,
		schemaStore: schemaStore,
		logger:      logger.GetLogger("Index", "Test"),
	}

	// gen tag key error
	schemaStore.EXPECT().genTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag.KeyID(0), fmt.Errorf("Err"))
	db.handle(&TagNotifier{
		metricID: 1,
		tags: tag.Tags{
			{Key: []byte("key1"), Value: []byte("value1")},
		},
	})
	// gen tag value error
	schemaStore.EXPECT().genTagKeyID(gomock.Any(), gomock.Any(), gomock.Any()).Return(tag.KeyID(0), nil)
	indexStore.EXPECT().GetOrCreateValue(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), fmt.Errorf("err"))
	db.handle(&TagNotifier{
		metricID: 1,
		tags: tag.Tags{
			{Key: []byte("key1"), Value: []byte("value1")},
		},
	})
}
