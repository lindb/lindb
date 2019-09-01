package memdb

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/pkg/field"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/tsdb/diskdb"
	"github.com/lindb/lindb/tsdb/series"
	"github.com/lindb/lindb/tsdb/tblstore"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_mStore_getMetricId(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	assert.NotNil(t, mStoreInterface)
	assert.Equal(t, uint32(100), mStoreInterface.getMetricID())
	assert.True(t, mStoreInterface.isEmpty())
	assert.False(t, mStore.isFull())
	assert.Zero(t, mStoreInterface.getTagsUsed())
	assert.Zero(t, mStoreInterface.getTagsInUse())
}

func Test_mStore_setMaxTagsLimit(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	assert.NotZero(t, mStore.getMaxTagsLimit())
	mStoreInterface.setMaxTagsLimit(1000)
	assert.Equal(t, uint32(1000), mStore.getMaxTagsLimit())
}

func Test_mStore_write_getOrCreateTStore_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().getTStore(gomock.Any()).Return(nil, false).AnyTimes()
	mockTagIdx.EXPECT().getOrCreateTStore(gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()
	mockTagIdx.EXPECT().tagsUsed().Return(10).AnyTimes()

	mStore.mutable = mockTagIdx
	assert.NotNil(t, mStore.write(&pb.Metric{Name: "metric", Tags: map[string]string{"type": "test"}}, writeContext{}))
}

func Test_mStore_isFull(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().tagsUsed().Return(10000000).AnyTimes()

	mStore.mutable = mockTagIdx
	assert.Equal(t, series.ErrTooManyTags,
		mStoreInterface.write(&pb.Metric{Name: "metric", Tags: map[string]string{"type": "test"}}, writeContext{}))
}

func Test_mStore_write_ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	mockTStore := NewMocktStoreINTF(ctrl)
	mockTStore.EXPECT().write(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().tagsUsed().Return(1).AnyTimes()
	mockTagIdx.EXPECT().updateTime(gomock.Any()).Return().AnyTimes()
	mockTagIdx.EXPECT().getTStore(gomock.Any()).Return(nil, false).AnyTimes()
	mockTagIdx.EXPECT().getOrCreateTStore(gomock.Any()).Return(mockTStore, nil).AnyTimes()

	mStore.mutable = mockTagIdx
	assert.Nil(t, mStoreInterface.write(&pb.Metric{Name: "metric", Tags: map[string]string{"type": "test"}}, writeContext{}))
}

func Test_mStore_resetVersion(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	assert.Nil(t, mStore.immutable)
	assert.NotNil(t, mStoreInterface.resetVersion())

	tagIdx := mStore.mutable.(*tagIndex)
	tagIdx.version -= 3600 * 1000
	mStore.mutable = tagIdx
	assert.Nil(t, mStoreInterface.resetVersion())
	assert.NotNil(t, mStore.immutable)
}

func Test_mStore_evict(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	// evict on empty
	mStore.evict()
	assert.True(t, mStore.isEmpty())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock tStores
	mockTStore1 := NewMocktStoreINTF(ctrl)
	mockTStore1.EXPECT().isNoData().Return(true).AnyTimes()
	mockTStore1.EXPECT().isExpired().Return(false).AnyTimes()
	mockTStore2 := NewMocktStoreINTF(ctrl)
	mockTStore2.EXPECT().isNoData().Return(false).AnyTimes()
	mockTStore2.EXPECT().isExpired().Return(false).AnyTimes()
	mockTStore3 := NewMocktStoreINTF(ctrl)
	mockTStore3.EXPECT().isNoData().Return(true).AnyTimes()
	mockTStore3.EXPECT().isExpired().Return(true).AnyTimes()
	mockTStore4 := NewMocktStoreINTF(ctrl)
	mockTStore4.EXPECT().isNoData().Return(true).AnyTimes()
	mockTStore4.EXPECT().isExpired().Return(true).AnyTimes()
	// mock tagIndex
	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().allTStores().Return(map[uint32]tStoreINTF{
		11: mockTStore1,
		22: mockTStore2,
		33: mockTStore3,
		44: mockTStore3,
	})
	mockTagIdx.EXPECT().getTStoreBySeriesID(uint32(33)).Return(mockTStore3, true).AnyTimes()
	mockTagIdx.EXPECT().getTStoreBySeriesID(uint32(44)).Return(nil, false).AnyTimes()
	mockTagIdx.EXPECT().removeTStores(uint32(33)).Return().AnyTimes()

	mStore.mutable = mockTagIdx
	mStoreInterface.evict()
}

func Test_mStore_flushMetricsTo_error(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock tagIndex
	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().flushMetricTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error")).AnyTimes()
	mStore.immutable = []tagIndexINTF{mockTagIdx}

	assert.NotNil(t, mStoreInterface.flushMetricsTo(nil, flushContext{}))
}

func Test_mStore_flushMetricsTo_OK(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock tagIndex
	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().flushMetricTo(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mStore.immutable = []tagIndexINTF{mockTagIdx}
	mStore.mutable = mockTagIdx

	// mock flush field meta
	mockTF := tblstore.NewMockMetricsDataFlusher(ctrl)
	mockTF.EXPECT().FlushFieldMeta(gomock.Any(), gomock.Any()).AnyTimes()
	mStore.fieldsMetas = append(mStore.fieldsMetas, fieldMeta{}, fieldMeta{})

	assert.Nil(t, mStoreInterface.flushMetricsTo(mockTF, flushContext{}))
	assert.Nil(t, mStore.immutable)
}

func Test_mStore_findSeriesIDsByExpr_getSeriesIDsForTag(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTagIdx := NewMocktagIndexINTF(ctrl)
	count := int64(0)
	mockTagIdx.EXPECT().getVersion().DoAndReturn(func() int64 {
		count++
		return count
	}).AnyTimes()

	// mock findSeriesIDsByExpr
	returnNotNil := mockTagIdx.EXPECT().findSeriesIDsByExpr(gomock.Any()).Return(roaring.New()).Times(2)
	returnNil := mockTagIdx.EXPECT().findSeriesIDsByExpr(gomock.Any()).Return(nil).Times(2)
	gomock.InOrder(returnNotNil, returnNil)
	// build mStore
	mStore.immutable = []tagIndexINTF{mockTagIdx}
	mStore.mutable = mockTagIdx
	// result assert
	set, err := mStoreInterface.findSeriesIDsByExpr(nil)
	assert.Nil(t, err)
	assert.NotNil(t, set)
	_, err2 := mStoreInterface.findSeriesIDsByExpr(nil)
	assert.Nil(t, err2)
	// mock getSeriesIDsForTag
	returnNotNil2 := mockTagIdx.EXPECT().getSeriesIDsForTag(gomock.Any()).Return(roaring.New()).Times(2)
	returnNil2 := mockTagIdx.EXPECT().getSeriesIDsForTag(gomock.Any()).Return(nil).Times(2)
	gomock.InOrder(returnNotNil2, returnNil2)
	mStoreInterface.getSeriesIDsForTag("")
	mStoreInterface.getSeriesIDsForTag("")
}

func Test_getFieldIDOrGenerate(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGen := diskdb.NewMockIDGenerator(ctrl)
	// mock generate ok
	mockGen.EXPECT().GenFieldID(uint32(100), "sum", field.SumField).Return(uint16(1), nil).AnyTimes()
	fieldID, err := mStoreInterface.getFieldIDOrGenerate("sum", field.SumField, mockGen)
	assert.Equal(t, uint16(1), fieldID)
	assert.Nil(t, err)
	// exist case
	_, err = mStoreInterface.getFieldIDOrGenerate("sum", field.SumField, mockGen)
	// field not matches to the existed
	assert.Nil(t, err)
	_, err = mStoreInterface.getFieldIDOrGenerate("sum", field.MinField, mockGen)
	assert.NotNil(t, err)
	// mock generate failure
	mockGen.EXPECT().GenFieldID(uint32(100), "gen-error", field.SumField).
		Return(uint16(1), series.ErrWrongFieldType)
	_, err = mStoreInterface.getFieldIDOrGenerate("gen-error", field.SumField, mockGen)
	assert.NotNil(t, err)

	// mock too many fields
	for range [3000]struct{}{} {
		mStore.fieldsMetas = append(mStore.fieldsMetas, fieldMeta{})
	}
	_, err = mStoreInterface.getFieldIDOrGenerate("sum", field.SumField, mockGen)
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
	mStoreInterface.getFieldIDOrGenerate("3", field.SumField, mockGen)
	mStoreInterface.getFieldIDOrGenerate("1", field.SumField, mockGen)
	mStoreInterface.getFieldIDOrGenerate("2", field.SumField, mockGen)
}

func prepareMockTagIndexes(ctrl *gomock.Controller) (*MocktagIndexINTF, *MocktagIndexINTF, *MocktagIndexINTF) {

	fakeKVEntrySet1 := []tagKVEntrySet{
		{key: "host", values: map[string]*roaring.Bitmap{"alpha": roaring.New(), "beta": roaring.New()}},
		{key: "zone", values: map[string]*roaring.Bitmap{"nj": roaring.New(), "bj": roaring.New()}},
	}
	fakeKVEntrySet2 := []tagKVEntrySet{
		{key: "ip", values: map[string]*roaring.Bitmap{"1.1.1.1": roaring.New(), "2.2.2.2": roaring.New()}},
		{key: "zone", values: map[string]*roaring.Bitmap{"sh": roaring.New(), "bj": roaring.New()}},
	}
	fakeKVEntrySet3 := []tagKVEntrySet{
		{key: "usage", values: map[string]*roaring.Bitmap{"idle": roaring.New(), "system": roaring.New()}},
		{key: "zone", values: map[string]*roaring.Bitmap{"nj": roaring.New(), "nt": roaring.New()}},
	}
	// mock tag index interface
	mockTagIdx1 := NewMocktagIndexINTF(ctrl)
	mockTagIdx1.EXPECT().getTagKVEntrySets().Return(fakeKVEntrySet1).AnyTimes()
	mockTagIdx1.EXPECT().getTimeRange().Return(uint32(1), uint32(2)).AnyTimes()
	mockTagIdx1.EXPECT().getVersion().Return(uint32(1)).AnyTimes()
	mockTagIdx1.EXPECT().getTagKVEntrySet("host").Return(&fakeKVEntrySet1[0], true).AnyTimes()
	mockTagIdx1.EXPECT().getTagKVEntrySet("zone").Return(&fakeKVEntrySet1[1], true).AnyTimes()
	mockTagIdx1.EXPECT().getTagKVEntrySet("ip").Return(nil, false).AnyTimes()
	mockTagIdx1.EXPECT().getTagKVEntrySet("usage").Return(nil, false).AnyTimes()

	mockTagIdx2 := NewMocktagIndexINTF(ctrl)
	mockTagIdx2.EXPECT().getTagKVEntrySets().Return(fakeKVEntrySet2).AnyTimes()
	mockTagIdx2.EXPECT().getTimeRange().Return(uint32(1), uint32(2)).AnyTimes()
	mockTagIdx2.EXPECT().getVersion().Return(uint32(2)).AnyTimes()
	mockTagIdx2.EXPECT().getTagKVEntrySet("ip").Return(&fakeKVEntrySet2[0], true).AnyTimes()
	mockTagIdx2.EXPECT().getTagKVEntrySet("host").Return(nil, false).AnyTimes()
	mockTagIdx2.EXPECT().getTagKVEntrySet("usage").Return(nil, false).AnyTimes()
	mockTagIdx2.EXPECT().getTagKVEntrySet("zone").Return(&fakeKVEntrySet2[1], true).AnyTimes()

	mockTagIdx3 := NewMocktagIndexINTF(ctrl)
	mockTagIdx3.EXPECT().getTagKVEntrySets().Return(fakeKVEntrySet3).AnyTimes()
	mockTagIdx3.EXPECT().getTimeRange().Return(uint32(1), uint32(2)).AnyTimes()
	mockTagIdx3.EXPECT().getVersion().Return(uint32(3)).AnyTimes()
	mockTagIdx3.EXPECT().getTagKVEntrySet("usage").Return(&fakeKVEntrySet3[0], true).AnyTimes()
	mockTagIdx3.EXPECT().getTagKVEntrySet("host").Return(nil, false).AnyTimes()
	mockTagIdx3.EXPECT().getTagKVEntrySet("ip").Return(nil, false).AnyTimes()
	mockTagIdx3.EXPECT().getTagKVEntrySet("zone").Return(&fakeKVEntrySet3[1], true).AnyTimes()

	return mockTagIdx1, mockTagIdx2, mockTagIdx3
}

func Test_mStore_flushInvertedIndexTo(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagIdx1, mockTagIdx2, mockTagIdx3 := prepareMockTagIndexes(ctrl)

	// mock index-table series flusher
	mockTableFlusher := tblstore.NewMockInvertedIndexFlusher(ctrl)
	mockTableFlusher.EXPECT().FlushVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return().AnyTimes()
	mockTableFlusher.EXPECT().FlushTagValue(gomock.Any()).Return().AnyTimes()

	//////////////////////////////////////////////
	// immutable part empty
	//////////////////////////////////////////////
	mStore.mutable = mockTagIdx1
	// flush ok
	mockTableFlusher.EXPECT().FlushTagID(gomock.Any()).Return(nil).Times(2)
	assert.Nil(t, mStore.flushInvertedIndexTo(mockTableFlusher, makeMockIDGenerator(ctrl)))
	// flush error
	mockTableFlusher.EXPECT().FlushTagID(gomock.Any()).Return(fmt.Errorf("error")).Times(1)
	assert.NotNil(t, mStore.flushInvertedIndexTo(mockTableFlusher, makeMockIDGenerator(ctrl)))

	//////////////////////////////////////////////
	// neither mutable nor immutable part is empty
	//////////////////////////////////////////////
	mStore.immutable = []tagIndexINTF{mockTagIdx1, mockTagIdx2}
	mStore.mutable = mockTagIdx3
	// flush error
	mockTableFlusher.EXPECT().FlushTagID(gomock.Any()).Return(fmt.Errorf("error")).Times(1)
	assert.NotNil(t, mStore.flushInvertedIndexTo(mockTableFlusher, makeMockIDGenerator(ctrl)))
	// flush ok
	mockTableFlusher.EXPECT().FlushTagID(gomock.Any()).Return(nil).Times(4)
	assert.Nil(t, mStore.flushInvertedIndexTo(mockTableFlusher, makeMockIDGenerator(ctrl)))
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
	mockTableFlusher.EXPECT().FlushVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return().AnyTimes()
	mockTableFlusher.EXPECT().FlushMetricID(gomock.Any()).Return(nil).AnyTimes()

	//////////////////////////////////////////////
	// immutable part empty
	//////////////////////////////////////////////
	mStore.mutable = mockTagIdx1
	assert.Nil(t, mStoreInterface.flushForwardIndexTo(mockTableFlusher))
	//////////////////////////////////////////////
	// neither mutable nor immutable part is empty
	//////////////////////////////////////////////
	mStore.immutable = []tagIndexINTF{mockTagIdx1, mockTagIdx2}
	mStore.mutable = mockTagIdx3
	assert.Nil(t, mStoreInterface.flushForwardIndexTo(mockTableFlusher))
}

func Test_mStore_getTagValues(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagIdx1, mockTagIdx2, mockTagIdx3 := prepareMockTagIndexes(ctrl)

	//////////////////////////////////////////////
	// immutable part empty
	//////////////////////////////////////////////
	mStore.mutable = mockTagIdx3
	tagValues, err := mStoreInterface.getTagValues([]string{"host", "zone", "ip", "usage"}, 3)
	assert.Nil(t, err)
	assert.Len(t, tagValues, 4)
	assert.Len(t, tagValues[0], 0)
	assert.Len(t, tagValues[1], 2)
	assert.Len(t, tagValues[2], 0)
	assert.Len(t, tagValues[3], 2)
	//////////////////////////////////////////////
	// immutable part not empty
	//////////////////////////////////////////////
	mStore.immutable = []tagIndexINTF{mockTagIdx1, mockTagIdx2}
	mStore.mutable = mockTagIdx3
	// version not match
	_, err = mStoreInterface.getTagValues([]string{"ip"}, 4)
	assert.NotNil(t, err)
	// version match, found
	_, err = mStoreInterface.getTagValues([]string{"ip"}, 1)
	assert.Nil(t, err)
}

func Test_mStore_suggest(t *testing.T) {
	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTagIdx1, mockTagIdx2, mockTagIdx3 := prepareMockTagIndexes(ctrl)

	// invalid limit
	assert.Nil(t, mStoreInterface.suggestTagValues("", "", 0))
	assert.Nil(t, mStoreInterface.suggestTagKeys("", 0))

	mStore.immutable = []tagIndexINTF{mockTagIdx1, mockTagIdx2}
	mStore.mutable = mockTagIdx3

	assert.Len(t, mStoreInterface.suggestTagKeys("host", 1), 1)
	assert.Len(t, mStoreInterface.suggestTagKeys("host", 3), 1)
	assert.Len(t, mStoreInterface.suggestTagValues("host", "a", 1), 1)
	assert.Len(t, mStoreInterface.suggestTagValues("host", "a", 100000), 1)
}
