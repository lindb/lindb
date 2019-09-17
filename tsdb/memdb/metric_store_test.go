package memdb

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/diskdb"
	"github.com/lindb/lindb/tsdb/tblstore"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_mStore_GetMetricID(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	assert.NotNil(t, mStoreInterface)
	assert.Equal(t, uint32(100), mStoreInterface.GetMetricID())
	assert.True(t, mStoreInterface.IsEmpty())
	assert.False(t, mStore.isFull())
	assert.Zero(t, mStoreInterface.GetTagsUsed())
	assert.Zero(t, mStoreInterface.GetTagsInUse())
}

func Test_mStore_setMaxTagsLimit(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	assert.NotZero(t, mStore.getMaxTagsLimit())
	mStoreInterface.SetMaxTagsLimit(1000)
	assert.Equal(t, uint32(1000), mStore.getMaxTagsLimit())
}

func Test_mStore_write_getOrCreateTStore_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().GetTStore(gomock.Any()).Return(nil, false).AnyTimes()
	mockTagIdx.EXPECT().GetOrCreateTStore(gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()
	mockTagIdx.EXPECT().TagsUsed().Return(10).AnyTimes()

	mStore.mutable = mockTagIdx
	assert.NotNil(t, mStore.Write(&pb.Metric{Name: "metric", Tags: map[string]string{"type": "test"}}, writeContext{}))
}

func Test_mStore_isFull(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().TagsUsed().Return(10000000).AnyTimes()

	mStore.mutable = mockTagIdx
	assert.Equal(t, series.ErrTooManyTags,
		mStoreInterface.Write(&pb.Metric{Name: "metric", Tags: map[string]string{"type": "test"}}, writeContext{}))
}

func Test_mStore_write_ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	mockTStore := NewMocktStoreINTF(ctrl)
	mockTStore.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().TagsUsed().Return(1).AnyTimes()
	mockTagIdx.EXPECT().UpdateIndexTimeRange(gomock.Any()).Return().AnyTimes()
	mockTagIdx.EXPECT().GetTStore(gomock.Any()).Return(nil, false).AnyTimes()
	mockTagIdx.EXPECT().GetOrCreateTStore(gomock.Any()).Return(mockTStore, nil).AnyTimes()

	mStore.mutable = mockTagIdx
	assert.Nil(t, mStoreInterface.Write(&pb.Metric{Name: "metric", Tags: map[string]string{"type": "test"}}, writeContext{}))
}

func Test_mStore_resetVersion(t *testing.T) {
	mStoreInterface := newMetricStore(100)

	assert.Nil(t, mStoreInterface.ResetVersion())
	assert.NotNil(t, mStoreInterface.ResetVersion())
	assert.NotNil(t, mStoreInterface.ResetVersion())
}

func Test_mStore_evict(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	// evict on empty
	mStore.Evict()
	assert.True(t, mStore.IsEmpty())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock tStores
	mockTStore1 := NewMocktStoreINTF(ctrl)
	mockTStore1.EXPECT().IsNoData().Return(true).AnyTimes()
	mockTStore1.EXPECT().IsExpired().Return(false).AnyTimes()
	mockTStore2 := NewMocktStoreINTF(ctrl)
	mockTStore2.EXPECT().IsNoData().Return(false).AnyTimes()
	mockTStore2.EXPECT().IsExpired().Return(false).AnyTimes()
	mockTStore3 := NewMocktStoreINTF(ctrl)
	mockTStore3.EXPECT().IsNoData().Return(true).AnyTimes()
	mockTStore3.EXPECT().IsExpired().Return(true).AnyTimes()
	mockTStore4 := NewMocktStoreINTF(ctrl)
	mockTStore4.EXPECT().IsNoData().Return(true).AnyTimes()
	mockTStore4.EXPECT().IsExpired().Return(true).AnyTimes()
	// mock tagIndex
	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().AllTStores().Return(map[uint32]tStoreINTF{
		11: mockTStore1,
		22: mockTStore2,
		33: mockTStore3,
		44: mockTStore3,
	})
	mockTagIdx.EXPECT().GetTStoreBySeriesID(uint32(33)).Return(mockTStore3, true).AnyTimes()
	mockTagIdx.EXPECT().GetTStoreBySeriesID(uint32(44)).Return(nil, false).AnyTimes()
	mockTagIdx.EXPECT().RemoveTStores(uint32(33)).Return().AnyTimes()

	mStore.mutable = mockTagIdx
	mStoreInterface.Evict()
}

func Test_mStore_FlushMetricsDataTo_withImmutable(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	flusher := tblstore.NewMockMetricsDataFlusher(ctrl)
	flusher.EXPECT().FlushMetric(gomock.Any()).Return(nil).AnyTimes()
	// mock tagIndex
	mStore.mutable = newTagIndex()
	_ = mStore.ResetVersion()
	assert.Nil(t, mStoreInterface.FlushMetricsDataTo(flusher, flushContext{}))
}

func Test_mStore_FlushMetricsDataTo_OK(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock tagIndex
	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().Version().Return(series.Version(1)).AnyTimes()
	mockTagIdx.EXPECT().FlushVersionDataTo(gomock.Any(), gomock.Any()).Return().AnyTimes()
	mStore.mutable = mockTagIdx

	assert.Nil(t, mStore.atomicGetImmutable())
	// mock flush field meta
	mockTF := tblstore.NewMockMetricsDataFlusher(ctrl)
	mockTF.EXPECT().FlushFieldMeta(gomock.Any(), gomock.Any()).AnyTimes()
	mockTF.EXPECT().FlushMetric(gomock.Any()).Return(nil).AnyTimes()
	mStore.fieldsMetas.Store(&fieldsMetas{fieldMeta{}, fieldMeta{}})

	assert.Nil(t, mStoreInterface.FlushMetricsDataTo(mockTF, flushContext{}))
	assert.Nil(t, mStore.atomicGetImmutable())
}

func Test_mStore_findSeriesIDsByExpr_getSeriesIDsForTag(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagIdx := NewMocktagIndexINTF(ctrl)
	count := int64(1)
	mockTagIdx.EXPECT().Version().DoAndReturn(func() series.Version {
		count++
		return series.Version(count)
	}).AnyTimes()

	// mock FindSeriesIDsByExpr
	returnNotNil := mockTagIdx.EXPECT().FindSeriesIDsByExpr(gomock.Any()).Return(roaring.New()).Times(2)
	returnNil := mockTagIdx.EXPECT().FindSeriesIDsByExpr(gomock.Any()).Return(nil).Times(2)
	gomock.InOrder(returnNotNil, returnNil)
	// build mStore
	mStore.immutable.Store(mockTagIdx)
	mStore.mutable = mockTagIdx
	// result assert
	set, err := mStoreInterface.FindSeriesIDsByExpr(nil)
	assert.Nil(t, err)
	assert.NotNil(t, set)
	_, err2 := mStoreInterface.FindSeriesIDsByExpr(nil)
	assert.Nil(t, err2)
	// mock GetSeriesIDsForTag
	returnNotNil2 := mockTagIdx.EXPECT().GetSeriesIDsForTag(gomock.Any()).Return(roaring.New()).Times(2)
	returnNil2 := mockTagIdx.EXPECT().GetSeriesIDsForTag(gomock.Any()).Return(nil).Times(2)
	gomock.InOrder(returnNotNil2, returnNil2)
	_, _ = mStoreInterface.GetSeriesIDsForTag("")
	_, _ = mStoreInterface.GetSeriesIDsForTag("")
}

func Test_getFieldIDOrGenerate(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGen := diskdb.NewMockIDGenerator(ctrl)
	// mock generate ok
	mockGen.EXPECT().GenFieldID(uint32(100), "sum", field.SumField).Return(uint16(1), nil).AnyTimes()
	fieldID, err := mStoreInterface.GetFieldIDOrGenerate("sum", field.SumField, mockGen)
	assert.Equal(t, uint16(1), fieldID)
	assert.Nil(t, err)
	// exist case
	_, err = mStoreInterface.GetFieldIDOrGenerate("sum", field.SumField, mockGen)
	// field not matches to the existed
	assert.Nil(t, err)
	_, err = mStoreInterface.GetFieldIDOrGenerate("sum", field.MinField, mockGen)
	assert.NotNil(t, err)
	// mock generate failure
	mockGen.EXPECT().GenFieldID(uint32(100), "gen-error", field.SumField).
		Return(uint16(1), series.ErrWrongFieldType)
	_, err = mStoreInterface.GetFieldIDOrGenerate("gen-error", field.SumField, mockGen)
	assert.NotNil(t, err)

	// mock too many fields
	var fieldsMetasList fieldsMetas
	for range [3000]struct{}{} {
		fieldsMetasList = append(fieldsMetasList, fieldMeta{})
	}
	mStore.fieldsMetas.Store(&fieldsMetasList)
	_, err = mStoreInterface.GetFieldIDOrGenerate("sum", field.SumField, mockGen)
	assert.NotNil(t, err)
}

func Test_getFieldIDOrGenerate_special_case(t *testing.T) {
	mStoreInterface := newMetricStore(100)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGen := diskdb.NewMockIDGenerator(ctrl)
	// fields meta sort
	mockGen.EXPECT().GenFieldID(uint32(100), "1", field.SumField).Return(uint16(1), nil).AnyTimes()
	mockGen.EXPECT().GenFieldID(uint32(100), "2", field.SumField).Return(uint16(2), nil).AnyTimes()
	mockGen.EXPECT().GenFieldID(uint32(100), "3", field.SumField).Return(uint16(3), nil).AnyTimes()
	_, _ = mStoreInterface.GetFieldIDOrGenerate("3", field.SumField, mockGen)
	_, _ = mStoreInterface.GetFieldIDOrGenerate("1", field.SumField, mockGen)
	_, _ = mStoreInterface.GetFieldIDOrGenerate("2", field.SumField, mockGen)
}

func prepareMockTagIndexes(ctrl *gomock.Controller) (*MocktagIndexINTF, *MocktagIndexINTF, *MocktagIndexINTF) {

	fakeKVEntrySet1 := []tagKVEntrySet{
		{key: "host", values: map[string]*roaring.Bitmap{
			"alpha": roaring.BitmapOf(1, 2, 3, 4, 5),
			"beta":  roaring.BitmapOf(6, 7, 8, 9, 10)}},
		{key: "zone", values: map[string]*roaring.Bitmap{
			"nj": roaring.BitmapOf(1, 2, 3, 4),
			"bj": roaring.BitmapOf(7, 8, 9, 10)}}}
	fakeKVEntrySet2 := []tagKVEntrySet{
		{key: "ip", values: map[string]*roaring.Bitmap{
			"1.1.1.1": roaring.BitmapOf(1, 2, 3, 4, 5),
			"2.2.2.2": roaring.BitmapOf(6, 7, 8, 9, 10)}},
		{key: "zone", values: map[string]*roaring.Bitmap{
			"sh": roaring.BitmapOf(1, 2, 3, 4, 5),
			"bj": roaring.BitmapOf(6, 7, 8, 9, 10)}}}
	fakeKVEntrySet3 := []tagKVEntrySet{
		{key: "usage", values: map[string]*roaring.Bitmap{
			"idle":   roaring.BitmapOf(1, 2, 3, 8, 9),
			"system": roaring.BitmapOf(4, 5, 6, 7, 10)}},
		{key: "zone", values: map[string]*roaring.Bitmap{
			"nj": roaring.BitmapOf(1, 2, 3, 4, 5),
			"nt": roaring.BitmapOf(6, 7, 8, 9, 10)}}}
	// mock tag index interface
	mockTagIdx1 := NewMocktagIndexINTF(ctrl)
	mockTagIdx1.EXPECT().GetTagKVEntrySets().Return(fakeKVEntrySet1).AnyTimes()
	mockTagIdx1.EXPECT().IndexTimeRange().Return(timeutil.TimeRange{Start: 1, End: 2}).AnyTimes()
	mockTagIdx1.EXPECT().Version().Return(series.Version(1)).AnyTimes()
	mockTagIdx1.EXPECT().GetTagKVEntrySet("host").Return(&fakeKVEntrySet1[0], true).AnyTimes()
	mockTagIdx1.EXPECT().GetTagKVEntrySet("zone").Return(&fakeKVEntrySet1[1], true).AnyTimes()
	mockTagIdx1.EXPECT().GetTagKVEntrySet("ip").Return(nil, false).AnyTimes()
	mockTagIdx1.EXPECT().GetTagKVEntrySet("usage").Return(nil, false).AnyTimes()

	mockTagIdx2 := NewMocktagIndexINTF(ctrl)
	mockTagIdx2.EXPECT().GetTagKVEntrySets().Return(fakeKVEntrySet2).AnyTimes()
	mockTagIdx2.EXPECT().IndexTimeRange().Return(timeutil.TimeRange{Start: 1, End: 2}).AnyTimes()
	mockTagIdx2.EXPECT().Version().Return(series.Version(2)).AnyTimes()
	mockTagIdx2.EXPECT().GetTagKVEntrySet("ip").Return(&fakeKVEntrySet2[0], true).AnyTimes()
	mockTagIdx2.EXPECT().GetTagKVEntrySet("host").Return(nil, false).AnyTimes()
	mockTagIdx2.EXPECT().GetTagKVEntrySet("usage").Return(nil, false).AnyTimes()
	mockTagIdx2.EXPECT().GetTagKVEntrySet("zone").Return(&fakeKVEntrySet2[1], true).AnyTimes()

	mockTagIdx3 := NewMocktagIndexINTF(ctrl)
	mockTagIdx3.EXPECT().GetTagKVEntrySets().Return(fakeKVEntrySet3).AnyTimes()
	mockTagIdx3.EXPECT().IndexTimeRange().Return(timeutil.TimeRange{Start: 1, End: 2}).AnyTimes()
	mockTagIdx3.EXPECT().Version().Return(series.Version(3)).AnyTimes()
	mockTagIdx3.EXPECT().GetTagKVEntrySet("usage").Return(&fakeKVEntrySet3[0], true).AnyTimes()
	mockTagIdx3.EXPECT().GetTagKVEntrySet("host").Return(nil, false).AnyTimes()
	mockTagIdx3.EXPECT().GetTagKVEntrySet("ip").Return(nil, false).AnyTimes()
	mockTagIdx3.EXPECT().GetTagKVEntrySet("zone").Return(&fakeKVEntrySet3[1], true).AnyTimes()

	return mockTagIdx1, mockTagIdx2, mockTagIdx3
}

func Test_mStore_flushInvertedIndexTo(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagIdx1, _, mockTagIdx3 := prepareMockTagIndexes(ctrl)

	// mock index-table series flusher
	mockTableFlusher := tblstore.NewMockInvertedIndexFlusher(ctrl)
	mockTableFlusher.EXPECT().FlushVersion(gomock.Any(), gomock.Any(), gomock.Any()).
		Return().AnyTimes()
	mockTableFlusher.EXPECT().FlushTagValue(gomock.Any()).Return().AnyTimes()

	//////////////////////////////////////////////
	// immutable part empty
	//////////////////////////////////////////////
	mStore.mutable = mockTagIdx1
	// flush ok
	mockTableFlusher.EXPECT().FlushTagKeyID(gomock.Any()).Return(nil).Times(2)
	assert.Nil(t, mStore.FlushInvertedIndexTo(mockTableFlusher, makeMockIDGenerator(ctrl)))
	// flush error
	mockTableFlusher.EXPECT().FlushTagKeyID(gomock.Any()).Return(fmt.Errorf("error")).Times(1)
	assert.NotNil(t, mStore.FlushInvertedIndexTo(mockTableFlusher, makeMockIDGenerator(ctrl)))

	//////////////////////////////////////////////
	// neither mutable nor immutable part is empty
	//////////////////////////////////////////////
	mStore.immutable.Store(mockTagIdx1)
	mStore.mutable = mockTagIdx3
	// flush error
	mockTableFlusher.EXPECT().FlushTagKeyID(gomock.Any()).Return(fmt.Errorf("error")).Times(1)
	assert.NotNil(t, mStore.FlushInvertedIndexTo(mockTableFlusher, makeMockIDGenerator(ctrl)))
	// flush ok
	mockTableFlusher.EXPECT().FlushTagKeyID(gomock.Any()).Return(nil).Times(3)
	assert.Nil(t, mStore.FlushInvertedIndexTo(mockTableFlusher, makeMockIDGenerator(ctrl)))
}

func Test_mStore_flushForwardIndexTo(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagIdx1, mockTagIdx2, mockTagIdx3 := prepareMockTagIndexes(ctrl)

	// mock index-table series flusher
	mockTableFlusher := tblstore.NewMockForwardIndexFlusher(ctrl)
	mockTableFlusher.EXPECT().FlushTagValue(gomock.Any(), gomock.Any()).Return().AnyTimes()
	mockTableFlusher.EXPECT().FlushTagKey(gomock.Any()).Return().AnyTimes()
	mockTableFlusher.EXPECT().FlushVersion(gomock.Any(), gomock.Any()).Return().AnyTimes()
	mockTableFlusher.EXPECT().FlushMetricID(gomock.Any()).Return(nil).AnyTimes()

	//////////////////////////////////////////////
	// immutable part empty
	//////////////////////////////////////////////
	mStore.mutable = mockTagIdx1
	assert.Nil(t, mStoreInterface.FlushForwardIndexTo(mockTableFlusher))
	//////////////////////////////////////////////
	// neither mutable nor immutable part is empty
	//////////////////////////////////////////////
	mStore.immutable.Store(mockTagIdx2)
	mStore.mutable = mockTagIdx3
	assert.Nil(t, mStoreInterface.FlushForwardIndexTo(mockTableFlusher))
}

func Test_mStore_getTagValues(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, mockTagIdx2, mockTagIdx3 := prepareMockTagIndexes(ctrl)

	//////////////////////////////////////////////
	// immutable part empty
	//////////////////////////////////////////////
	mStore.mutable = mockTagIdx3
	// host not exist
	mappings, err := mStoreInterface.GetTagValues(
		[]string{"host", "zone", "usage"}, 3, roaring.BitmapOf(3, 4, 5, 6, 11))
	assert.NotNil(t, err)
	assert.Nil(t, mappings)

	// zone, usage exist
	mappings, err = mStoreInterface.GetTagValues(
		[]string{"zone", "usage"}, 3, roaring.BitmapOf(3, 4, 5, 6, 11))
	assert.Nil(t, err)
	assert.Len(t, mappings, 5)
	assert.Equal(t, []string{"nj", "idle"}, mappings[3])
	assert.Equal(t, []string{"nj", "system"}, mappings[4])
	assert.Equal(t, []string{"nj", "system"}, mappings[5])
	assert.Equal(t, []string{"nt", "system"}, mappings[6])
	assert.Equal(t, []string{"", ""}, mappings[11])
	//////////////////////////////////////////////
	// immutable part not empty
	//////////////////////////////////////////////
	mStore.immutable.Store(mockTagIdx2)
	mStore.mutable = mockTagIdx3
	// version not match
	_, err = mStoreInterface.GetTagValues([]string{"ip"}, 4, roaring.BitmapOf(1, 2, 3))
	assert.NotNil(t, err)
	// version match, ip not exist
	_, err = mStoreInterface.GetTagValues([]string{"ip"}, 1, roaring.BitmapOf(1, 2, 3))
	assert.NotNil(t, err)
}

func Test_mStore_suggest(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagIdx1, _, mockTagIdx3 := prepareMockTagIndexes(ctrl)

	// invalid limit
	assert.Nil(t, mStoreInterface.SuggestTagValues("", "", 0))
	assert.Nil(t, mStoreInterface.SuggestTagKeys("", 0))

	mStore.immutable.Store(mockTagIdx1)
	mStore.mutable = mockTagIdx3

	assert.Len(t, mStoreInterface.SuggestTagKeys("host", 1), 1)
	assert.Len(t, mStoreInterface.SuggestTagKeys("host", 3), 1)
	assert.Len(t, mStoreInterface.SuggestTagValues("host", "a", 1), 1)
	assert.Len(t, mStoreInterface.SuggestTagValues("host", "a", 100000), 1)
}
