package memdb

import (
	"strconv"
	"testing"
	"time"

	"github.com/eleme/lindb/pkg/field"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_sortedTSStores(t *testing.T) {
	vm := newVersionedTSMap()
	vm.tsMap["1"] = newTimeSeriesStore()
	vm.tsMap["2"] = newTimeSeriesStore()
	vm.tsMap["3"] = newTimeSeriesStore()

	tagsList, tsStoreList, release := vm.allTSStores()
	defer release()
	assert.Len(t, *tagsList, 3)
	assert.Len(t, *tsStoreList, 3)
}

func Test_newMetricStore(t *testing.T) {
	mStore := newMetricStore()
	assert.NotNil(t, mStore)
	assert.NotNil(t, mStore.mutable)
	assert.NotNil(t, mStore.mutable.tsMap)
	assert.NotZero(t, mStore.maxTagsLimit)
}

func Test_metricStore_isEmpty_isFull(t *testing.T) {
	mStore := newMetricStore()
	assert.True(t, mStore.isEmpty())
	assert.False(t, mStore.isFull())

	for i := 0; i < 100; i++ {
		mStore.mutable.tsMap[strconv.Itoa(i)] = nil
	}
	assert.False(t, mStore.isFull())
	assert.False(t, mStore.isEmpty())

	for i := 0; i < defaultMaxTagsLimit; i++ {
		mStore.mutable.tsMap[strconv.Itoa(i)] = nil
	}
	assert.True(t, mStore.isFull())
	assert.False(t, mStore.isEmpty())
}

func Test_metricStore_getTimeSeries(t *testing.T) {
	mStore := newMetricStore()

	assert.NotNil(t, mStore.getOrCreateTSStore("host=alpha-1"))
	assert.Equal(t, mStore.getOrCreateTSStore("host=alpha-2"), mStore.getOrCreateTSStore("host=alpha-2"))
	assert.Equal(t, mStore.getTagsCount(), 2)
}

func Test_metricStore_evict(t *testing.T) {
	mStore := newMetricStore()
	mStore.evict()
	assert.True(t, mStore.isEmpty())
	// has not been purged
	for i := 0; i < 2000; i++ {
		mStore.getOrCreateTSStore(strconv.Itoa(i)).getOrCreateFStore("t", field.MaxField)
	}
	setTagsIDTTL(60 * 1000) // 1 minute
	assert.Equal(t, 2000, mStore.getTagsCount())
	mStore.evict()
	assert.Equal(t, 2000, mStore.getTagsCount())

	// purge half
	time.Sleep(time.Millisecond * 20)
	setTagsIDTTL(20) // 20 ms
	for i := 0; i < 1000; i++ {
		mStore.getOrCreateTSStore(strconv.Itoa(i)).getOrCreateFStore("t", field.MaxField)
	}
	mStore.evict()
	assert.Equal(t, 1000, mStore.getTagsCount())
	// purge all
	time.Sleep(time.Millisecond * 20)
	setTagsIDTTL(20) // 20 ms
	mStore.evict()
	assert.True(t, mStore.isEmpty())
}

func Test_mustGetMetricID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gen := makeMockIDGenerator(ctrl)
	mStore := newMetricStore()
	assert.NotZero(t, mStore.mustGetMetricID(gen, "m1"))
	assert.NotZero(t, mStore.mustGetMetricID(gen, "m2"))
}

func Test_unionFamilyTimesTo(t *testing.T) {
	vm := newVersionedTSMap()
	segments := map[int64]struct{}{1: {}, 2: {}, 3: {}}
	vm.familyTimes = map[int64]struct{}{2: {}, 4: {}}
	vm.unionFamilyTimesTo(segments)
	assert.Equal(t, 4, len(segments))

	ms := newMetricStore()
	ms.mutable = vm
	ms.immutable = append(ms.immutable, vm, vm)

	ms.unionFamilyTimesTo(segments)
	assert.Equal(t, 4, len(segments))
}

func Test_assignNewVersion(t *testing.T) {
	ms := newMetricStore()
	ms.mutable = newVersionedTSMap()

	assert.NotNil(t, ms.assignNewVersion())

	ms.mutable.version -= minIntervalForResetMetricStore * int64(time.Millisecond)
	assert.Nil(t, ms.assignNewVersion())
}
