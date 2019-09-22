package memdb

import (
	"strconv"
	"testing"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/diskdb"
	"github.com/lindb/lindb/tsdb/tblstore"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_tagIndex_tStore_get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGenerator := diskdb.NewMockIDGenerator(ctrl)
	mockGenerator.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1)).AnyTimes()

	tagIdxInterface := newTagIndex()
	tagIdx := tagIdxInterface.(*tagIndex)
	// version
	assert.NotZero(t, tagIdxInterface.Version())
	// get empty key value tStore
	tStore0, err := tagIdxInterface.GetOrCreateTStore(nil, writeContext{generator: mockGenerator})
	assert.NotNil(t, tStore0)
	assert.Nil(t, err)
	// get not exist tStore
	tStore1, ok := tagIdxInterface.GetTStore(map[string]string{"host": "adca", "ip": "1.1.1.1"})
	assert.Nil(t, tStore1)
	assert.False(t, ok)
	// get or create
	tStore2, err := tagIdxInterface.GetOrCreateTStore(
		map[string]string{"host": "adca", "ip": "1.1.1.1"},
		writeContext{generator: mockGenerator})
	assert.NotNil(t, tStore2)
	assert.Nil(t, err)
	_, _ = tagIdxInterface.GetOrCreateTStore(
		map[string]string{"host": "adca", "ip": "1.1.1.1"},
		writeContext{generator: mockGenerator})
	// get existed
	tStore3, ok := tagIdxInterface.GetTStore(
		map[string]string{"host": "adca", "ip": "1.1.1.1"})
	assert.NotNil(t, tStore3)
	assert.True(t, ok)
	// get tStore by seriesID
	assert.NotZero(t, len(tagIdx.seriesID2TStore))
	tStore4, ok := tagIdxInterface.GetTStoreBySeriesID(1)
	assert.NotNil(t, tStore4)
	assert.True(t, ok)
	// getOrInsertTagKeyEntry, present in the slice
	_, _ = tagIdxInterface.GetOrCreateTStore(map[string]string{"g": "32"}, writeContext{generator: mockGenerator})
	_, _ = tagIdxInterface.GetOrCreateTStore(map[string]string{"g": "33"}, writeContext{generator: mockGenerator})
	_, _ = tagIdxInterface.GetOrCreateTStore(map[string]string{"h": "32"}, writeContext{generator: mockGenerator})

	// getTagKVEntrySet test
	assert.NotNil(t, tagIdxInterface.GetTagKVEntrySets())
}

func Test_tagIndex_tStore_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGenerator := diskdb.NewMockIDGenerator(ctrl)
	mockGenerator.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1)).AnyTimes()

	tagIdxInterface := newTagIndex()
	tagIdx := tagIdxInterface.(*tagIndex)
	// too many tag keys
	for i := 0; i < 1000; i++ {
		_, _ = tagIdx.GetOrCreateTStore(
			map[string]string{strconv.Itoa(i): strconv.Itoa(i)}, writeContext{generator: mockGenerator})
	}
	assert.Equal(t, 512, tagIdx.TagsUsed())
	_, err := tagIdxInterface.GetOrCreateTStore(
		map[string]string{"zone": "nj"},
		writeContext{generator: mockGenerator})
	assert.Equal(t, series.ErrTooManyTagKeys, err)
	assert.Equal(t, 512, tagIdx.TagsUsed())
	// remove tStores
	tagIdx.RemoveTStores()
	tagIdx.RemoveTStores(1, 2, 3, 4, 1003)
	// used tags won't change
	assert.Equal(t, 512, tagIdx.TagsUsed())
	// in use tags was removed
	assert.Equal(t, 508, tagIdx.TagsInUse())
	// allTStores
	assert.NotNil(t, tagIdxInterface.AllTStores())
}

func Test_tagIndex_flushMetricTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagIdxInterface := newTagIndex()
	tagIdx := tagIdxInterface.(*tagIndex)

	mockTF := tblstore.NewMockMetricsDataFlusher(ctrl)
	mockTF.EXPECT().FlushMetric(gomock.Any()).Return(nil).MaxTimes(2)
	mockTF.EXPECT().FlushVersion(gomock.Any()).Return().AnyTimes()

	// no data flushed, tStores is empty
	tagIdxInterface.FlushVersionDataTo(mockTF, flushContext{})

	// tStore is not empty
	mockTStore1 := NewMocktStoreINTF(ctrl)
	mockTStore1.EXPECT().GetHash().Return(uint64(1)).AnyTimes()
	mockTStore1.EXPECT().FlushSeriesTo(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()
	mockTStore2 := NewMocktStoreINTF(ctrl)
	mockTStore2.EXPECT().FlushSeriesTo(gomock.Any(), gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	mockTStore1.EXPECT().GetHash().Return(uint64(2)).AnyTimes()
	tagIdx.seriesID2TStore = map[uint32]tStoreINTF{
		1: mockTStore1,
		2: mockTStore2,
	}
	// data flushed
	tagIdxInterface.FlushVersionDataTo(mockTF, flushContext{})
}

func prepareTagIdx(ctrl *gomock.Controller) tagIndexINTF {
	tagIdxInterface := newTagIndex()
	tagIdx := tagIdxInterface.(*tagIndex)

	mockGenerator := diskdb.NewMockIDGenerator(ctrl)
	mockGenerator.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1)).AnyTimes()

	_, _ = tagIdxInterface.GetOrCreateTStore(
		map[string]string{"host": "a", "zone": "nj"},
		writeContext{generator: mockGenerator}) // seriesID: 1
	_, _ = tagIdxInterface.GetOrCreateTStore(
		map[string]string{"host": "abc", "zone": "sh"},
		writeContext{generator: mockGenerator}) // 2
	_, _ = tagIdxInterface.GetOrCreateTStore(
		map[string]string{"host": "b", "zone": "nj"},
		writeContext{generator: mockGenerator}) // 3
	_, _ = tagIdxInterface.GetOrCreateTStore(
		map[string]string{"host": "c", "zone": "bj"},
		writeContext{generator: mockGenerator}) // 4
	_, _ = tagIdxInterface.GetOrCreateTStore(
		map[string]string{"host": "bc", "zone": "sz"},
		writeContext{generator: mockGenerator}) // 5
	_, _ = tagIdxInterface.GetOrCreateTStore(
		map[string]string{"host": "b21", "zone": "nj"},
		writeContext{generator: mockGenerator}) // 6
	_, _ = tagIdxInterface.GetOrCreateTStore(
		map[string]string{"host": "b22", "zone": "sz"},
		writeContext{generator: mockGenerator}) // 7
	_, _ = tagIdxInterface.GetOrCreateTStore(
		map[string]string{"host": "bcd", "zone": "sh"},
		writeContext{generator: mockGenerator}) // 8

	newMap := make(map[uint32]tStoreINTF)
	for seriesID, tStore := range tagIdx.seriesID2TStore {
		mockTStore := NewMocktStoreINTF(ctrl)
		mockTStore.EXPECT().GetHash().Return(tStore.GetHash()).AnyTimes()
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
	bitmap := tagIdxInterface.FindSeriesIDsByExpr(&stmt.EqualsExpr{Key: "not-exist-key", Value: "alpha"})
	assert.Nil(t, bitmap)
	// tag-value not exist
	bitmap = tagIdxInterface.FindSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "alpha"})
	assert.Nil(t, bitmap)
	// tag-value exist
	bitmap = tagIdxInterface.FindSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "c"})
	assert.NotNil(t, bitmap)
	assert.Equal(t, uint64(1), bitmap.GetCardinality())
	// tag-value exist
	bitmap = tagIdxInterface.FindSeriesIDsByExpr(&stmt.EqualsExpr{Key: "host", Value: "bc"})
	assert.Equal(t, uint64(1), bitmap.GetCardinality())
}

func Test_tagIndex_findSeriesIDsByIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagIdxInterface := prepareTagIdx(ctrl)

	// tag-value exist
	bitmap := tagIdxInterface.FindSeriesIDsByExpr(&stmt.InExpr{Key: "host", Values: []string{"b", "bc", "bcd", "ahi"}})
	assert.Equal(t, uint64(3), bitmap.GetCardinality())
}

func Test_tagIndex_findSeriesIDsByLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagIdxInterface := prepareTagIdx(ctrl)

	// tag-value exist
	bitmap := tagIdxInterface.FindSeriesIDsByExpr(&stmt.LikeExpr{Key: "host", Value: "bc"})
	assert.Equal(t, uint64(3), bitmap.GetCardinality())
	bitmap = tagIdxInterface.FindSeriesIDsByExpr(&stmt.LikeExpr{Key: "zone", Value: "s"})
	assert.Equal(t, uint64(4), bitmap.GetCardinality())
	// tag-value not exist
	bitmap = tagIdxInterface.FindSeriesIDsByExpr(&stmt.LikeExpr{Key: "zone", Value: "not-exist"})
	assert.Zero(t, bitmap.GetCardinality())

	// tag-value is empty
	bitmap = tagIdxInterface.FindSeriesIDsByExpr(&stmt.LikeExpr{Key: "host", Value: ""})
	assert.Equal(t, uint64(0), bitmap.GetCardinality())
	// tag-value is *
	bitmap = tagIdxInterface.FindSeriesIDsByExpr(&stmt.LikeExpr{Key: "host", Value: "*"})
	assert.Equal(t, uint64(8), bitmap.GetCardinality())
}

func Test_tagIndex_findSeriesIDsByRegex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagIdxInterface := prepareTagIdx(ctrl)

	// pattern not match
	bitmap := tagIdxInterface.FindSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: "bbbbbbbbbbb"})
	assert.Zero(t, bitmap.GetCardinality())
	// pattern error
	bitmap = tagIdxInterface.FindSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: "b.32*++++\n"})
	assert.Nil(t, bitmap)
	// tag-value exist
	bitmap = tagIdxInterface.FindSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: `b2[0-9]+`})
	assert.Equal(t, uint64(2), bitmap.GetCardinality())
	// literal prefix:22 not exist
	bitmap = tagIdxInterface.FindSeriesIDsByExpr(&stmt.RegexExpr{Key: "host", Regexp: `22+`})
	assert.Equal(t, uint64(0), bitmap.GetCardinality())

}

func Test_tagIndex_getSeriesIDsForTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tagIdxInterface := prepareTagIdx(ctrl)

	// not-exist
	bitmap := tagIdxInterface.GetSeriesIDsForTag("not-exist-key")
	assert.Nil(t, bitmap)
	// overlap
	bitmap = tagIdxInterface.GetSeriesIDsForTag("host")
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
	assert.Nil(t, tagIdxInterface.FindSeriesIDsByExpr(mockTagKey{}))
}

func Test_TagIndex_recreateEvictedTStores(t *testing.T) {
	tagIdxInterface := newTagIndex()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGenerator := diskdb.NewMockIDGenerator(ctrl)
	mockGenerator.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1)).AnyTimes()

	_, _ = tagIdxInterface.GetOrCreateTStore(map[string]string{"host": "a"}, writeContext{generator: mockGenerator})
	_, _ = tagIdxInterface.GetOrCreateTStore(map[string]string{"host": "a"}, writeContext{generator: mockGenerator})
	_, _ = tagIdxInterface.GetOrCreateTStore(map[string]string{"host": "b"}, writeContext{generator: mockGenerator})
	assert.Equal(t, 2, tagIdxInterface.TagsInUse())
	assert.Equal(t, 2, tagIdxInterface.TagsUsed())
	// remove seriesID = 1
	tagIdxInterface.RemoveTStores(0, 1)
	assert.Equal(t, 1, tagIdxInterface.TagsInUse())
	assert.Equal(t, 2, tagIdxInterface.TagsUsed())
	_, _ = tagIdxInterface.GetOrCreateTStore(map[string]string{"host": "a"}, writeContext{generator: mockGenerator})
	assert.Equal(t, 2, tagIdxInterface.TagsInUse())
	tagIdxInterface.RemoveTStores(1, 2)
	assert.Equal(t, 0, tagIdxInterface.TagsInUse())
	assert.Equal(t, 2, tagIdxInterface.TagsUsed())
}

func Test_TagIndex_timeRange(t *testing.T) {
	tagIdxInterface := newTagIndex()
	timeRange := tagIdxInterface.IndexTimeRange()
	tagIdxInterface.UpdateIndexTimeRange(timeutil.Now() - 10*1000)
	tagIdxInterface.UpdateIndexTimeRange(timeutil.Now() + 10*1000)

	newTimeRange := tagIdxInterface.IndexTimeRange()
	assert.True(t, newTimeRange.Start < timeRange.Start)
	assert.True(t, newTimeRange.End > timeRange.End)
}
