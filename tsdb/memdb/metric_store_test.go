package memdb

import (
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
	tStore, size := mStore.GetOrCreateTStore(uint32(10))
	assert.NotNil(t, tStore)
	assert.True(t, size > 0)
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
	assert.Equal(t, field.SumField, mStore.fields[1])
	assert.Equal(t, field.MinField, mStore.fields[2])
}

func TestMetricStore_SetTimestamp(t *testing.T) {
	mStoreInterface := newMetricStore()
	mStore := mStoreInterface.(*metricStore)
	mStoreInterface.SetTimestamp(1, 10)
	slotRange := mStore.families[1]
	assert.Equal(t, uint16(10), slotRange.start)
	assert.Equal(t, uint16(10), slotRange.end)
	mStoreInterface.SetTimestamp(1, 5)
	slotRange = mStore.families[1]
	assert.Equal(t, uint16(5), slotRange.start)
	assert.Equal(t, uint16(10), slotRange.end)
	mStoreInterface.SetTimestamp(1, 50)
	slotRange = mStore.families[1]
	assert.Equal(t, uint16(5), slotRange.start)
	assert.Equal(t, uint16(50), slotRange.end)
}

func TestMetricStore_FlushMetricsDataTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flusher := metricsdata.NewMockFlusher(ctrl)

	mStoreInterface := newMetricStore()
	mStore := mStoreInterface.(*metricStore)
	tStore := NewMocktStoreINTF(ctrl)
	mStore.Put(10, tStore)

	// case 1: family time not exist
	err := mStoreInterface.FlushMetricsDataTo(flusher, flushContext{familyID: 1})
	assert.NoError(t, err)
	// case 2: field not exist
	mStoreInterface.SetTimestamp(1, 10)
	err = mStoreInterface.FlushMetricsDataTo(flusher, flushContext{familyID: 1})
	assert.NoError(t, err)
	// case 3: flush err
	mStoreInterface.AddField(1, field.SumField)
	mStoreInterface.AddField(2, field.MinField)
	flusher.EXPECT().FlushFieldMetas(gomock.Any())
	tStore.EXPECT().FlushSeriesTo(gomock.Any(), gomock.Any())
	flusher.EXPECT().FlushMetric(gomock.Any()).Return(nil)
	err = mStoreInterface.FlushMetricsDataTo(flusher, flushContext{familyID: 1})
	assert.NoError(t, err)
}
