package memdb

import (
	"strconv"
	"testing"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/diskdb"
	"github.com/lindb/lindb/tsdb/tblstore"

	"github.com/cespare/xxhash"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_tagIndex_tStore_get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGenerator := diskdb.NewMockIDGenerator(ctrl)
	mockGenerator.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1)).AnyTimes()

	tagIdxInterface := newTagIndex()
	// test get empty map
	tStore30, ok := tagIdxInterface.GetTStoreBySeriesID(uint32(10))
	assert.False(t, ok)
	assert.Nil(t, tStore30)

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
	assert.NotZero(t, tagIdx.seriesID2TStore.size())
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
	mockTStore1.EXPECT().GetSeriesID().Return(uint32(1)).AnyTimes()
	mockTStore1.EXPECT().FlushSeriesTo(gomock.Any(), gomock.Any(), gomock.Any()).Return(false).AnyTimes()
	mockTStore2 := NewMocktStoreINTF(ctrl)
	mockTStore2.EXPECT().FlushSeriesTo(gomock.Any(), gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	mockTStore1.EXPECT().GetSeriesID().Return(uint32(2)).AnyTimes()
	tagIdx.seriesID2TStore = newMetricMap()
	tagIdx.seriesID2TStore.put(1, mockTStore1)
	tagIdx.seriesID2TStore.put(2, mockTStore1)
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
	it := tagIdx.seriesID2TStore.iterator()
	for it.hasNext() {
		seriesID, tStore := it.next()
		mockTStore := NewMocktStoreINTF(ctrl)
		mockTStore.EXPECT().GetSeriesID().Return(tStore.GetSeriesID()).AnyTimes()
		newMap[seriesID] = mockTStore
	}

	tagIdx.seriesID2TStore = newMetricMap()
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

var _testHashString = "abcdefghijklmnopqrstuvwxzy1234567890"

func Benchmark_Fnv1a(b *testing.B) {
	const (
		// FNV-1a
		offset64 = uint64(14695981039346656037)
		prime64  = uint64(1099511628211)

		// Init64 is what 64 bits hash values should be initialized with.
		Init64 = offset64
	)
	AddString64 := func(h uint64, s string) uint64 {
		i := 0
		n := (len(s) / 8) * 8

		for i != n {
			h = (h ^ uint64(s[i])) * prime64
			h = (h ^ uint64(s[i+1])) * prime64
			h = (h ^ uint64(s[i+2])) * prime64
			h = (h ^ uint64(s[i+3])) * prime64
			h = (h ^ uint64(s[i+4])) * prime64
			h = (h ^ uint64(s[i+5])) * prime64
			h = (h ^ uint64(s[i+6])) * prime64
			h = (h ^ uint64(s[i+7])) * prime64
			i += 8
		}
		for _, c := range s[i:] {
			h = (h ^ uint64(c)) * prime64
		}
		return h
	}
	// HashString64 returns the hash of s.
	HashString64 := func(s string) uint64 {
		return AddString64(Init64, s)
	}

	for i := 0; i < b.N; i++ {
		HashString64(_testHashString)
	}
}

func Benchmark_xxhash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xxhash.Sum64String(_testHashString)
	}
}
