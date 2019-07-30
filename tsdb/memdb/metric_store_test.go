package memdb

import (
	"fmt"
	"testing"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/timeutil"
	pb "github.com/eleme/lindb/rpc/proto/field"

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
	assert.Zero(t, mStoreInterface.getTagsCount())
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
	mockTagIdx.EXPECT().len().Return(10).AnyTimes()

	mStore.mutable = mockTagIdx
	assert.NotNil(t, mStore.write(&pb.Metric{Name: "metric", Tags: "test"}, writeContext{}))
}

func Test_mStore_isFull(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)
	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().len().Return(10000000).AnyTimes()

	mStore.mutable = mockTagIdx
	assert.Equal(t, models.ErrTooManyTags,
		mStoreInterface.write(&pb.Metric{Name: "metric", Tags: "test"}, writeContext{}))
}

func Test_mStore_write_ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mStoreInterface := newMetricStore(100)
	mStore := mStoreInterface.(*metricStore)

	mockTStore := NewMocktStoreINTF(ctrl)
	mockTStore.EXPECT().write(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockTagIdx := NewMocktagIndexINTF(ctrl)
	mockTagIdx.EXPECT().len().Return(1).AnyTimes()
	mockTagIdx.EXPECT().getTStore(gomock.Any()).Return(nil, false).AnyTimes()
	mockTagIdx.EXPECT().getOrCreateTStore(gomock.Any()).Return(mockTStore, nil).AnyTimes()

	mStore.mutable = mockTagIdx
	assert.Nil(t, mStoreInterface.write(&pb.Metric{Name: "metric", Tags: "test"}, writeContext{}))
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
	mockTagIdx.EXPECT().allTStores().Return(map[uint32]tStoreNode{
		11: {tStoreINTF: mockTStore1, hash: 1},
		22: {tStoreINTF: mockTStore2, hash: 2},
		33: {tStoreINTF: mockTStore3, hash: 3},
		44: {tStoreINTF: mockTStore3, hash: 4},
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

	assert.Nil(t, mStoreInterface.flushMetricsTo(nil, flushContext{}))
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
	returnNotNil := mockTagIdx.EXPECT().findSeriesIDsByExpr(
		gomock.Any(), gomock.Any()).Return(roaring.New()).Times(2)
	returnNil := mockTagIdx.EXPECT().findSeriesIDsByExpr(
		gomock.Any(), gomock.Any()).Return(nil).Times(2)
	gomock.InOrder(returnNotNil, returnNil)
	// build mStore
	mStore.immutable = []tagIndexINTF{mockTagIdx}
	mStore.mutable = mockTagIdx
	// result assert
	set, err := mStoreInterface.findSeriesIDsByExpr(nil, timeutil.TimeRange{})
	assert.Nil(t, err)
	assert.NotNil(t, set)
	_, err2 := mStoreInterface.findSeriesIDsByExpr(nil, timeutil.TimeRange{})
	assert.Nil(t, err2)
	// mock getSeriesIDsForTag
	returnNotNil2 := mockTagIdx.EXPECT().getSeriesIDsForTag(
		gomock.Any(), gomock.Any()).Return(roaring.New()).Times(2)
	returnNil2 := mockTagIdx.EXPECT().getSeriesIDsForTag(
		gomock.Any(), gomock.Any()).Return(nil).Times(2)
	gomock.InOrder(returnNotNil2, returnNil2)
	mStoreInterface.getSeriesIDsForTag("", timeutil.TimeRange{})
	mStoreInterface.getSeriesIDsForTag("", timeutil.TimeRange{})

}
