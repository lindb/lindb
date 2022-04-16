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

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestMetricStore_GetOrCreateTStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStoreInterface := newMetricStore()
	mStore := mStoreInterface.(*metricStore)
	tStore, created := mStore.GetOrCreateTStore(uint32(10))
	assert.NotNil(t, tStore)
	assert.True(t, created)
	tStore2, size := mStore.GetOrCreateTStore(uint32(10))
	assert.Zero(t, size)
	assert.Equal(t, tStore, tStore2)
}

func TestMetricStore_AddField(t *testing.T) {
	mStoreInterface := newMetricStore()
	mStore := mStoreInterface.(*metricStore)
	mStoreInterface.AddField(1, field.SumField)
	mStoreInterface.AddField(1, field.SumField)
	mStoreInterface.AddField(2, field.MinField)
	assert.Len(t, mStore.fields, 2)
	assert.Equal(t, field.Meta{ID: 1, Type: field.SumField}, mStore.fields[0])
	assert.Equal(t, field.Meta{ID: 2, Type: field.MinField}, mStore.fields[1])
}

func TestMetricStore_SetTimestamp(t *testing.T) {
	mStoreInterface := newMetricStore()
	mStoreInterface.SetSlot(10)
	slotRange := mStoreInterface.GetSlotRange()
	assert.Equal(t, uint16(10), slotRange.Start)
	assert.Equal(t, uint16(10), slotRange.End)
	mStoreInterface.SetSlot(5)
	slotRange = mStoreInterface.GetSlotRange()
	assert.Equal(t, uint16(5), slotRange.Start)
	assert.Equal(t, uint16(10), slotRange.End)
	mStoreInterface.SetSlot(50)
	slotRange = mStoreInterface.GetSlotRange()
	assert.Equal(t, uint16(5), slotRange.Start)
	assert.Equal(t, uint16(50), slotRange.End)
}

func TestMetricStore_FlushMetricsDataTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		flushFunc = flush
	}()

	flusher := metricsdata.NewMockFlusher(ctrl)

	mStoreInterface := newMetricStore()
	mStore := mStoreInterface.(*metricStore)
	tStore := NewMocktStoreINTF(ctrl)
	mStore.Put(10, tStore)

	// case 1: family time not exist
	err := mStoreInterface.FlushMetricsDataTo(flusher, &flushContext{})
	assert.NoError(t, err)
	// case 2: field not exist
	mStoreInterface.SetSlot(10)
	err = mStoreInterface.FlushMetricsDataTo(flusher, &flushContext{})
	assert.NoError(t, err)
	// case 3: flush success
	mStoreInterface.AddField(1, field.SumField)
	mStoreInterface.AddField(2, field.MinField)
	gomock.InOrder(
		flusher.EXPECT().PrepareMetric(gomock.Any(), gomock.Any()),
		tStore.EXPECT().FlushFieldsTo(gomock.Any(), gomock.Any()).Return(nil),
		flusher.EXPECT().FlushSeries(uint32(10)).Return(nil),
		flusher.EXPECT().CommitMetric(gomock.Any()).Return(nil),
	)
	err = mStoreInterface.FlushMetricsDataTo(flusher, &flushContext{})
	assert.NoError(t, err)
	// case 4: flush field err
	flusher.EXPECT().PrepareMetric(gomock.Any(), gomock.Any())
	tStore.EXPECT().FlushFieldsTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = mStoreInterface.FlushMetricsDataTo(flusher, &flushContext{})
	assert.Error(t, err)
	// case 5: flush err
	flushFunc = func(flusher metricsdata.Flusher, flushCtx *flushContext, key uint32, value tStoreINTF) error {
		return fmt.Errorf("err")
	}
	gomock.InOrder(
		flusher.EXPECT().PrepareMetric(gomock.Any(), gomock.Any()),
	)
	err = mStoreInterface.FlushMetricsDataTo(flusher, &flushContext{})
	assert.Error(t, err)
}

func Benchmark_MetricBucketStore_get(b *testing.B) {
	noOptimization := func(count int) func(b *testing.B) {
		m := NewMetricBucketStore()
		for i := 0; i < count; i += 2 {
			m.Put(uint32(i), nil)
		}

		return func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				m.Get(uint32(b.N % count))
			}
			b.StopTimer()
		}
	}

	withOptimization := func(count int) func(b *testing.B) {
		m := NewMetricBucketStore()
		for i := 0; i < count; i++ {
			m.Put(uint32(i), nil)
		}
		m.keys.RunOptimize()

		return func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				m.Get(uint32(b.N % count))
			}
			b.StopTimer()
		}
	}

	b.Run("10_without_optimize", noOptimization(10))
	b.Run("100_without_optimize", noOptimization(100))
	b.Run("500_without_optimize", noOptimization(500))
	b.Run("1000_without_optimize", noOptimization(1000))
	b.Run("5000_without_optimize", noOptimization(5000))
	b.Run("10000_without_optimize", noOptimization(10000))
	b.Run("50000_without_optimize", noOptimization(50000))
	b.Run("100000_without_optimize", noOptimization(100000))

	b.Run("100_with_optimize", withOptimization(100))
	b.Run("1000_with_optimize", withOptimization(1000))
	b.Run("5000_with_optimize", withOptimization(5000))
	b.Run("10000_with_optimize", withOptimization(10000))
	b.Run("50000_with_optimize", withOptimization(50000))
	b.Run("100000_with_optimize", withOptimization(100000))
}
