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

package memdb

// import (
// 	"fmt"
// 	"testing"
//
// 	"github.com/lindb/common/pkg/fasttime"
// 	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"
// 	"github.com/stretchr/testify/assert"
// 	"go.uber.org/mock/gomock"
//
// 	"github.com/lindb/lindb/index"
// 	"github.com/lindb/lindb/models"
// 	"github.com/lindb/lindb/pkg/imap"
// 	"github.com/lindb/lindb/series/field"
// 	"github.com/lindb/lindb/series/metric"
// )
//
// func TestMetadataDatabase_flush(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	metaDB := index.NewMockMetricMetaDatabase(ctrl)
// 	metaDB.EXPECT().PrepareFlush()
// 	metaDB.EXPECT().Flush().Return(nil)
// 	mdb := NewMetadataDatabase(&models.DatabaseConfig{}, metaDB)
// 	assert.NotNil(t, mdb.Config())
// 	ch := make(chan error)
// 	mdb.Notify(&FlushEvent{
// 		Callback: func(err error) {
// 			ch <- err
// 		},
// 	})
// 	err := <-ch
// 	assert.NoError(t, err)
// 	mdb.Close()
// }
//
// func TestMetadataDatabase_GetOrCreateMetricMeta(t *testing.T) {
// 	mdb := NewMetadataDatabase(&models.DatabaseConfig{}, nil)
//
// 	m := &protoMetricsV1.Metric{
// 		Name:      "test1",
// 		Namespace: "ns",
// 		Tags:      []*protoMetricsV1.KeyValue{{Key: "key1", Value: "value1"}},
// 		SimpleFields: []*protoMetricsV1.SimpleField{
// 			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
// 		},
// 	}
//
// 	row := protoToStorageRow(m)
// 	ms, ok := mdb.GetOrCreateMetricMeta(row)
// 	assert.NotNil(t, ms)
// 	assert.True(t, ok)
//
// 	ms, ok = mdb.GetOrCreateMetricMeta(row)
// 	assert.NotNil(t, ms)
// 	assert.False(t, ok)
//
// 	ms, ok = mdb.(*metadataDatabase).getOrCreateMetricMeta(row.NameHash())
// 	assert.NotNil(t, ms)
// 	assert.False(t, ok)
//
// 	ms, ok = mdb.GetMetricMeta(1000)
// 	assert.Nil(t, ms)
// 	assert.False(t, ok)
// }
//
// func TestMetadataDatabase_handleRow(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
//
// 	metaDB := index.NewMockMetricMetaDatabase(ctrl)
// 	mdb := NewMetadataDatabase(&models.DatabaseConfig{}, metaDB).(*metadataDatabase)
// 	m := &protoMetricsV1.Metric{
// 		Name:      "test1",
// 		Namespace: "ns",
// 		Tags:      []*protoMetricsV1.KeyValue{{Key: "key1", Value: "value1"}},
// 		SimpleFields: []*protoMetricsV1.SimpleField{
// 			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 10},
// 		},
// 	}
//
// 	row := protoToStorageRow(m)
// 	row.Add(1000)
// 	// gen metric id err
// 	metaDB.EXPECT().GenMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(0), fmt.Errorf("err"))
// 	mdb.handleRow(row)
//
// 	// no fields
// 	metaDB.EXPECT().GenMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(0), nil)
// 	mdb.handleRow(row)
// 	// metric meta not found
// 	row.Fields = field.Metas{{Name: "test"}}
// 	metaDB.EXPECT().GenMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(0), nil)
// 	mdb.handleRow(row)
// 	// gen field err
// 	metaDB.EXPECT().GenFieldID(gomock.Any(), gomock.Any()).Return(field.ID(0), fmt.Errorf("err"))
// 	metaDB.EXPECT().GenMetricID(gomock.Any(), gomock.Any()).Return(metric.ID(0), nil)
// 	_, _ = mdb.GetOrCreateMetricMeta(row)
// 	mdb.handleRow(row)
// }
//
// func TestMetadataDatabase_gc(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
//
// 	mStore := NewMockmStoreINTF(ctrl)
//
// 	t.Run("no data gc", func(t *testing.T) {
// 		mdb := &metadataDatabase{
// 			metricIndexStore: imap.NewIntMap[uint64](),
// 		}
// 		mdb.gc(fasttime.UnixMilliseconds())
// 	})
//
// 	t.Run("metric store active", func(t *testing.T) {
// 		mdb := &metadataDatabase{
// 			metricIndexStore: imap.NewIntMap[uint64](),
// 		}
// 		mdb.metricMetadatas.Store(uint64(100), mStore)
// 		mStore.EXPECT().IsActive(gomock.Any()).Return(true)
// 		mdb.gc(fasttime.UnixMilliseconds())
// 		_, ok := mdb.metricMetadatas.Load(uint64(100))
// 		assert.True(t, ok)
// 	})
//
// 	t.Run("gc metric store", func(t *testing.T) {
// 		mdb := &metadataDatabase{
// 			metricIndexStore: imap.NewIntMap[uint64](),
// 		}
// 		mdb.metricIndexStore.Put(10, 100)
// 		mdb.metricMetadatas.Store(uint64(100), mStore)
// 		mStore.EXPECT().IsActive(gomock.Any()).Return(false)
// 		mdb.gc(fasttime.UnixMilliseconds())
// 		_, ok := mdb.metricMetadatas.Load(uint64(100))
// 		assert.False(t, ok)
// 		assert.True(t, mdb.metricIndexStore.IsEmpty())
// 	})
//
// 	t.Run("not gc metric index", func(t *testing.T) {
// 		mdb := &metadataDatabase{
// 			metricIndexStore: imap.NewIntMap[uint64](),
// 		}
// 		mdb.metricIndexStore.Put(10, 100)
// 		mdb.metricIndexStore.Put(20, 200)
// 		mdb.metricMetadatas.Store(uint64(100), mStore)
// 		mdb.metricMetadatas.Store(uint64(200), mStore)
// 		mdb.metricMetadatas.Store(uint64(300), mStore)
// 		mStore.EXPECT().IsActive(gomock.Any()).Return(false)
// 		mStore.EXPECT().IsActive(gomock.Any()).Return(true).MaxTimes(2)
// 		mdb.gc(fasttime.UnixMilliseconds())
// 		assert.Equal(t, 2, mdb.metricIndexStore.Size())
// 	})
//
// 	t.Run("not gc metric index", func(t *testing.T) {
// 		mdb := &metadataDatabase{
// 			metricIndexStore: imap.NewIntMap[uint64](),
// 		}
// 		mdb.metricIndexStore.Put(10, 100)
// 		mdb.metricIndexStore.Put(20, 200)
// 		mdb.metricMetadatas.Store(uint64(100), mStore)
// 		mdb.metricMetadatas.Store(uint64(200), mStore)
// 		mStore.EXPECT().IsActive(gomock.Any()).Return(false)
// 		mStore.EXPECT().IsActive(gomock.Any()).Return(true)
// 		mdb.gc(fasttime.UnixMilliseconds())
// 		assert.Equal(t, 1, mdb.metricIndexStore.Size())
// 	})
// }
