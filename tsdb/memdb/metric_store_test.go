package memdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/metadb"
)

//func Test_mStore_write_getOrCreateTStore_error(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mStoreInterface := newMetricStore()
//	mStore := mStoreInterface.(*metricStore)
//
//	mockTagIdx := NewMocktagIndexINTF(ctrl)
//	mockTagIdx.EXPECT().GetTStore(gomock.Any()).Return(nil, false).AnyTimes()
//	mockTagIdx.EXPECT().GetOrCreateTStore(gomock.Any(), gomock.Any()).
//		Return(nil, 0, fmt.Errorf("error")).AnyTimes()
//	mockTagIdx.EXPECT().TagsUsed().Return(10).AnyTimes()
//
//	mStore.mutable = mockTagIdx
//	writtenSize, err := mStore.Write(&pb.Metric{Name: "metric", Tags: map[string]string{"type": "test"}}, writeContext{metricID: 1})
//	assert.Zero(t, writtenSize)
//	assert.NotNil(t, err)
//}
//
//
//func Test_mStore_write_ok(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mStoreInterface := newMetricStore()
//	mStore := mStoreInterface.(*metricStore)
//
//	mockTStore := NewMocktStoreINTF(ctrl)
//	mockTStore.EXPECT().Write(gomock.Any(), gomock.Any()).Return(0, nil).AnyTimes()
//
//	mockTagIdx := NewMocktagIndexINTF(ctrl)
//	mockTagIdx.EXPECT().TagsUsed().Return(1).AnyTimes()
//	mockTagIdx.EXPECT().UpdateIndexTimeRange(gomock.Any()).Return().AnyTimes()
//	mockTagIdx.EXPECT().GetTStore(gomock.Any()).Return(nil, false).AnyTimes()
//	mockTagIdx.EXPECT().GetOrCreateTStore(gomock.Any(), gomock.Any()).Return(mockTStore, 30, nil).AnyTimes()
//
//	mStore.mutable = mockTagIdx
//	writtenSize, err := mStoreInterface.Write(
//		&pb.Metric{Name: "metric", Tags: map[string]string{"type": "test"}}, writeContext{metricID: 1, slotIndex: 10})
//	assert.NoError(t, err)
//	assert.NotZero(t, writtenSize)
//	// test metric slot range
//	writtenSize, err = mStoreInterface.Write(
//		&pb.Metric{Name: "metric", Tags: map[string]string{"type": "test"}}, writeContext{metricID: 1, slotIndex: 9})
//	assert.NoError(t, err)
//	assert.NotZero(t, writtenSize)
//}
//
//func Test_mStore_resetVersion(t *testing.T) {
//	mStoreInterface := newMetricStore()
//	size1 := mStoreInterface.MemSize()
//	createdSize, err := mStoreInterface.ResetVersion()
//	assert.Nil(t, err)
//	assert.NotZero(t, createdSize)
//
//	createdSize, err = mStoreInterface.ResetVersion()
//	assert.NotNil(t, err)
//	assert.Zero(t, createdSize)
//
//	createdSize, err = mStoreInterface.ResetVersion()
//	assert.NotNil(t, err)
//	assert.Zero(t, createdSize)
//	size2 := mStoreInterface.MemSize()
//	assert.NotEqual(t, size1, size2)
//}
//
//func Test_mStore_evict(t *testing.T) {
//	mStoreInterface := newMetricStore()
//	mStore := mStoreInterface.(*metricStore)
//	// evict on empty
//	mStore.Evict()
//	assert.True(t, mStore.IsEmpty())
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// mock tStores
//	mockTStore1 := NewMocktStoreINTF(ctrl)
//	mockTStore1.EXPECT().IsNoData().Return(true).AnyTimes()
//	mockTStore1.EXPECT().IsExpired().Return(false).AnyTimes()
//	mockTStore2 := NewMocktStoreINTF(ctrl)
//	mockTStore2.EXPECT().IsNoData().Return(false).AnyTimes()
//	mockTStore2.EXPECT().IsExpired().Return(false).AnyTimes()
//	mockTStore3 := NewMocktStoreINTF(ctrl)
//	mockTStore3.EXPECT().IsNoData().Return(true).AnyTimes()
//	mockTStore3.EXPECT().IsExpired().Return(true).AnyTimes()
//	mockTStore4 := NewMocktStoreINTF(ctrl)
//	mockTStore4.EXPECT().IsNoData().Return(true).AnyTimes()
//	mockTStore4.EXPECT().IsExpired().Return(true).AnyTimes()
//	// mock tagIndex
//	mockTagIdx := NewMocktagIndexINTF(ctrl)
//	metricMap := newMetricMap()
//	metricMap.put(11, mockTStore1)
//	metricMap.put(22, mockTStore2)
//	metricMap.put(33, mockTStore3)
//	metricMap.put(44, mockTStore4)
//	mockTagIdx.EXPECT().AllTStores().Return(metricMap)
//	mockTagIdx.EXPECT().GetTStoreBySeriesID(uint32(33)).Return(mockTStore3, true).AnyTimes()
//	mockTagIdx.EXPECT().GetTStoreBySeriesID(uint32(44)).Return(nil, false).AnyTimes()
//	mockTagIdx.EXPECT().RemoveTStores(uint32(33)).Return([]tStoreINTF{mockTStore3}).AnyTimes()
//
//	mStore.mutable = mockTagIdx
//	mockTStore3.EXPECT().MemSize().Return(10)
//	mStoreInterface.Evict()
//}
//
//func Test_mStore_FlushMetricsDataTo_withImmutable(t *testing.T) {
//	mStoreInterface := newMetricStore()
//	mStore := mStoreInterface.(*metricStore)
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	flusher := metricsdata.NewMockFlusher(ctrl)
//	flusher.EXPECT().FlushMetric(gomock.Any()).Return(nil).AnyTimes()
//	flusher.EXPECT().FlushFieldMetas(gomock.Any()).Return().AnyTimes()
//	// mock tagIndex
//	mStore.mutable = newTagIndex()
//	_, _ = mStore.ResetVersion()
//	flushedSize, err := mStoreInterface.FlushMetricsDataTo(flusher, flushContext{})
//	assert.Nil(t, err)
//	assert.Zero(t, flushedSize)
//}
//
//func Test_mStore_FlushMetricsDataTo_OK(t *testing.T) {
//	mStoreInterface := newMetricStore()
//	mStore := mStoreInterface.(*metricStore)
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// mock tagIndex
//	mockTagIdx := NewMocktagIndexINTF(ctrl)
//	mockTagIdx.EXPECT().Version().Return(series.Version(1)).AnyTimes()
//	mockTagIdx.EXPECT().FlushVersionDataTo(gomock.Any(), gomock.Any()).Return(10).AnyTimes()
//	mStore.mutable = mockTagIdx
//
//	assert.Nil(t, mStore.atomicGetImmutable())
//	// mock flush field meta
//	mockTF := metricsdata.NewMockFlusher(ctrl)
//	mockTF.EXPECT().FlushFieldMetas(gomock.Any()).AnyTimes()
//	mockTF.EXPECT().FlushMetric(gomock.Any()).Return(nil).AnyTimes()
//	mStore.fieldsMetas.Store(field.Metas{field.Meta{}, field.Meta{}})
//
//	flushedSize, err := mStoreInterface.FlushMetricsDataTo(mockTF, flushContext{})
//	assert.NotZero(t, flushedSize)
//	assert.Nil(t, err)
//	assert.Nil(t, mStore.atomicGetImmutable())
//}
func Test_getFieldIDOrGenerate(t *testing.T) {
	mStoreInterface := newMetricStore()
	mStore := mStoreInterface.(*metricStore)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGen := metadb.NewMockIDGenerator(ctrl)
	// mock generate ok
	mockGen.EXPECT().GenFieldID(uint32(100), "sum", field.SumField).Return(uint16(1), nil).AnyTimes()
	fieldID, err := mStoreInterface.GetFieldIDOrGenerate(uint32(100), "sum", field.SumField, mockGen)
	assert.Equal(t, uint16(1), fieldID)
	assert.Nil(t, err)
	// exist case
	_, err = mStoreInterface.GetFieldIDOrGenerate(uint32(100), "sum", field.SumField, mockGen)
	// field not matches to the existed
	assert.Nil(t, err)
	_, err = mStoreInterface.GetFieldIDOrGenerate(uint32(100), "sum", field.MinField, mockGen)
	assert.NotNil(t, err)
	// mock generate failure
	mockGen.EXPECT().GenFieldID(uint32(100), "gen-error", field.SumField).
		Return(uint16(1), series.ErrWrongFieldType)
	_, err = mStoreInterface.GetFieldIDOrGenerate(uint32(100), "gen-error", field.SumField, mockGen)
	assert.NotNil(t, err)

	// mock too many fields
	var fieldsMetasList field.Metas
	for range [3000]struct{}{} {
		fieldsMetasList = append(fieldsMetasList, field.Meta{})
	}
	mStore.fieldsMetas.Store(fieldsMetasList)
	_, err = mStoreInterface.GetFieldIDOrGenerate(uint32(100), "sum", field.SumField, mockGen)
	assert.NotNil(t, err)
}

func Test_getFieldIDOrGenerate_special_case(t *testing.T) {
	mStoreInterface := newMetricStore()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGen := metadb.NewMockIDGenerator(ctrl)
	// fields meta sort
	mockGen.EXPECT().GenFieldID(uint32(100), "1", field.SumField).Return(uint16(1), nil).AnyTimes()
	mockGen.EXPECT().GenFieldID(uint32(100), "2", field.SumField).Return(uint16(2), nil).AnyTimes()
	mockGen.EXPECT().GenFieldID(uint32(100), "3", field.SumField).Return(uint16(3), nil).AnyTimes()
	_, _ = mStoreInterface.GetFieldIDOrGenerate(uint32(100), "3", field.SumField, mockGen)
	_, _ = mStoreInterface.GetFieldIDOrGenerate(uint32(100), "1", field.SumField, mockGen)
	_, _ = mStoreInterface.GetFieldIDOrGenerate(uint32(100), "2", field.SumField, mockGen)
}
