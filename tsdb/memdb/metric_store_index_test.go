package memdb

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metrictbl"
	"github.com/lindb/lindb/tsdb/series"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_tagIndex_tStore_get(t *testing.T) {
	tagIdxInterface := newTagIndex()
	tagIdx := tagIdxInterface.(*tagIndex)
	// version
	assert.NotZero(t, tagIdxInterface.getVersion())

	// get empty key value tStore
	tStore0, err := tagIdxInterface.getOrCreateTStore("")
	assert.NotNil(t, tStore0)
	assert.Nil(t, err)
	// get not exist tStore
	tStore1, ok := tagIdxInterface.getTStore("host=adca,ip=1.1.1.1")
	assert.Nil(t, tStore1)
	assert.False(t, ok)
	// get or create
	tStore2, err := tagIdxInterface.getOrCreateTStore("host=adca,ip=1.1.1.1")
	assert.NotNil(t, tStore2)
	assert.Nil(t, err)
	tagIdxInterface.getOrCreateTStore("host=adca,ip=1.1.1.1")
	// get existed
	tStore3, ok := tagIdxInterface.getTStore("host=adca,ip=1.1.1.1")
	assert.NotNil(t, tStore3)
	assert.True(t, ok)
	// get tStore by seriesID
	assert.NotZero(t, len(tagIdx.seriesID2TStore))
	tStore4, ok := tagIdxInterface.getTStoreBySeriesID(1)
	assert.NotNil(t, tStore4)
	assert.True(t, ok)
	// getOrInsertTagKeyEntry, present in the slice
	tagIdxInterface.getOrCreateTStore("g=32")
	tagIdxInterface.getOrCreateTStore("g=33")
	tagIdxInterface.getOrCreateTStore("h=33")

	// getTagKVEntrySet test
	assert.NotNil(t, tagIdxInterface.getTagKVEntrySets())
}

func Test_tagIndex_tStore_error(t *testing.T) {
	tagIdxInterface := newTagIndex()
	tagIdx := tagIdxInterface.(*tagIndex)
	// too many tag keys
	for i := 0; i < 1000; i++ {
		_, _ = tagIdx.getOrCreateTStore(fmt.Sprintf("%d=%d", i, i))
	}
	assert.Equal(t, 512, tagIdx.tagsUsed())
	_, err := tagIdxInterface.getOrCreateTStore("zone=nj")
	assert.Equal(t, series.ErrTooManyTagKeys, err)
	assert.Equal(t, 512, tagIdx.tagsUsed())
	// remove tStores
	tagIdx.removeTStores()
	tagIdx.removeTStores(1, 2, 3, 4, 1003)
	// used tags won't change
	assert.Equal(t, 512, tagIdx.tagsUsed())
	// in use tags was removed
	assert.Equal(t, 508, tagIdx.tagsInUse())
	// allTStores
	assert.NotNil(t, tagIdxInterface.allTStores())
}

func Test_tagIndex_flushMetricTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagIdxInterface := newTagIndex()
	tagIdx := tagIdxInterface.(*tagIndex)

	mockTF := metrictbl.NewMockTableFlusher(ctrl)
	mockTF.EXPECT().FlushMetric(gomock.Any()).Return(nil).MaxTimes(2)

	// tStores is empty
	assert.Nil(t, tagIdxInterface.flushMetricTo(mockTF, flushContext{}))

	// tStore is not empty
	mockTStore1 := NewMocktStoreINTF(ctrl)
	mockTStore1.EXPECT().getHash().Return(uint64(1)).AnyTimes()
	mockTStore1.EXPECT().flushSeriesTo(gomock.Any(), gomock.Any()).Return(false).AnyTimes()
	mockTStore2 := NewMocktStoreINTF(ctrl)
	mockTStore2.EXPECT().flushSeriesTo(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	mockTStore1.EXPECT().getHash().Return(uint64(2)).AnyTimes()
	tagIdx.seriesID2TStore = map[uint32]tStoreINTF{
		1: mockTStore1,
		2: mockTStore2,
	}
	// FlushMetric ok
	assert.Nil(t, tagIdxInterface.flushMetricTo(mockTF, flushContext{}))
}

func prepareTagIdx(ctrl *gomock.Controller) tagIndexINTF {
	tagIdxInterface := newTagIndex()
	tagIdx := tagIdxInterface.(*tagIndex)

	tagIdxInterface.getOrCreateTStore("host=a,zone=nj")   // seriesID: 1
	tagIdxInterface.getOrCreateTStore("host=abc,zone=sh") // 2
	tagIdxInterface.getOrCreateTStore("host=b,zone=nj")   // 3
	tagIdxInterface.getOrCreateTStore("host=c,zone=bj")   // 4
	tagIdxInterface.getOrCreateTStore("host=bc,zone=sz")  // 5
	tagIdxInterface.getOrCreateTStore("host=b21,zone=nj") // 6
	tagIdxInterface.getOrCreateTStore("host=b22,zone=sz") // 7
	tagIdxInterface.getOrCreateTStore("host=bcd,zone=sh") // 8

	newMap := make(map[uint32]tStoreINTF)
	for seriesID, tStore := range tagIdx.seriesID2TStore {
		mockTStore := NewMocktStoreINTF(ctrl)
		mockTStore.EXPECT().getHash().Return(tStore.getHash()).AnyTimes()
		newMap[seriesID] = mockTStore
	}

	tagIdx.seriesID2TStore = newMap
	return tagIdxInterface
}

func Test_tagIndex_findSeriesIDsByEqual(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagIdxInterface := prepareTagIdx(ctrl)

	// tag-key not exist
	bitmap := tagIdxInterface.findSeriesIDsByExpr(&stmt.EqualsExpr{Key: "not-exist-key", Value: "alpha"})
	assert.Nil(t, bitmap)
	// tag-value not exist
	bitmap = tagIdxInterface.findSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "alpha"})
	assert.Nil(t, bitmap)
	// tag-value exist
	bitmap = tagIdxInterface.findSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "c"})
	assert.NotNil(t, bitmap)
	assert.Equal(t, uint64(1), bitmap.GetCardinality())
	// tag-value exist
	bitmap = tagIdxInterface.findSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "bc"})
	assert.Equal(t, uint64(1), bitmap.GetCardinality())
}

func Test_tagIndex_findSeriesIDsByIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagIdxInterface := prepareTagIdx(ctrl)

	// tag-value exist
	bitmap := tagIdxInterface.findSeriesIDsByExpr(&stmt.InExpr{Key: "host", Values: []string{"b", "bc", "bcd", "ahi"}})
	assert.Equal(t, uint64(3), bitmap.GetCardinality())
}

func Test_tagIndex_findSeriesIDsByLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagIdxInterface := prepareTagIdx(ctrl)

	// tag-value exist
	bitmap := tagIdxInterface.findSeriesIDsByExpr(&stmt.LikeExpr{Key: "host", Value: "bc"})
	assert.Equal(t, uint64(3), bitmap.GetCardinality())
	bitmap = tagIdxInterface.findSeriesIDsByExpr(&stmt.LikeExpr{Key: "zone", Value: "s"})
	assert.Equal(t, uint64(4), bitmap.GetCardinality())
	// tag-value not exist
	bitmap = tagIdxInterface.findSeriesIDsByExpr(&stmt.LikeExpr{Key: "zone", Value: "not-exist"})
	assert.Zero(t, bitmap.GetCardinality())
}

func Test_tagIndex_findSeriesIDsByRegex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagIdxInterface := prepareTagIdx(ctrl)

	// pattern not match
	bitmap := tagIdxInterface.findSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: "bbbbbbbbbbb"})
	assert.Zero(t, bitmap.GetCardinality())
	// pattern error
	bitmap = tagIdxInterface.findSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: "b.32*++++\n"})
	assert.Nil(t, bitmap)
	// tag-value exist
	bitmap = tagIdxInterface.findSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: `b2[0-9]+`})
	assert.Equal(t, uint64(2), bitmap.GetCardinality())
}

func Test_tagIndex_getSeriesIDsForTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagIdxInterface := prepareTagIdx(ctrl)

	// not-exist
	bitmap := tagIdxInterface.getSeriesIDsForTag("not-exist-key")
	assert.Nil(t, bitmap)
	// overlap
	bitmap = tagIdxInterface.getSeriesIDsForTag("host")
	assert.Equal(t, uint64(8), bitmap.GetCardinality())
}

type mockTagKey struct {
}

func (mockTagKey) TagKey() string {
	return "host"
}

func Test_tagIndex_special_case(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagIdxInterface := prepareTagIdx(ctrl)
	// test expr type assertion failure
	assert.Nil(t, tagIdxInterface.findSeriesIDsByExpr(mockTagKey{}))
}

func Test_TagIndex_recreateEvictedTStores(t *testing.T) {
	tagIdxInterface := newTagIndex()

	tagIdxInterface.getOrCreateTStore("host=a")
	tagIdxInterface.getOrCreateTStore("host=a")
	tagIdxInterface.getOrCreateTStore("host=b")
	assert.Equal(t, 2, tagIdxInterface.tagsInUse())
	assert.Equal(t, 2, tagIdxInterface.tagsUsed())
	// remove seriesID = 1
	tagIdxInterface.removeTStores(0, 1)
	assert.Equal(t, 1, tagIdxInterface.tagsInUse())
	assert.Equal(t, 2, tagIdxInterface.tagsUsed())
	tagIdxInterface.getOrCreateTStore("host=a")
	assert.Equal(t, 2, tagIdxInterface.tagsInUse())
	tagIdxInterface.removeTStores(1, 2)
	assert.Equal(t, 0, tagIdxInterface.tagsInUse())
	assert.Equal(t, 2, tagIdxInterface.tagsUsed())
}

func Test_TagIndex_timeRange(t *testing.T) {
	tagIdxInterface := newTagIndex()
	startTime, endTime := tagIdxInterface.getTimeRange()
	tagIdxInterface.updateTime(uint32(timeutil.Now()/1000 - 10))
	tagIdxInterface.updateTime(uint32(timeutil.Now()/1000 + 10))

	newStartTime, newEndTime := tagIdxInterface.getTimeRange()
	assert.True(t, newStartTime < startTime)
	assert.True(t, newEndTime > endTime)
}
