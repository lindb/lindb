package memdb

import (
	"testing"

	"github.com/cespare/xxhash"
	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func Test_tagIndex_tStore_get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGenerator := metadb.NewMockIDGenerator(ctrl)
	mockGenerator.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1)).AnyTimes()

	tagIdxInterface := newTagIndex()
	// test get empty map
	tStore30, createdSize := tagIdxInterface.GetOrCreateTStore(uint32(10))
	assert.True(t, createdSize > 0)
	assert.NotNil(t, tStore30)

	tagIdx := tagIdxInterface.(*tagIndex)
	// version
	assert.NotZero(t, tagIdxInterface.Version())
	// get empty key value tStore
	tStore30, createdSize = tagIdxInterface.GetOrCreateTStore(uint32(10))
	assert.True(t, createdSize == 0)
	assert.NotNil(t, tStore30)
	assert.Equal(t, 1, tagIdx.seriesID2TStore.size())
}

func Test_tagIndex_filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGenerator := metadb.NewMockIDGenerator(ctrl)
	mockGenerator.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1)).AnyTimes()

	tagIdxInterface := newTagIndex()
	tagIdx := tagIdxInterface.(*tagIndex)
	// too many tag keys
	for i := 0; i < 1000; i++ {
		_, _ = tagIdx.GetOrCreateTStore(uint32(i))
	}
	assert.True(t, tagIdx.filter(roaring.BitmapOf(1)))
	assert.False(t, tagIdx.filter(roaring.BitmapOf(3000)))

	tagIdx.loadData(nil, nil, 0, nil)
}

func Test_tagIndex_flushMetricTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagIdxInterface := newTagIndex()
	tagIdx := tagIdxInterface.(*tagIndex)

	mockTF := metricsdata.NewMockFlusher(ctrl)
	mockTF.EXPECT().FlushMetric(gomock.Any()).Return(nil).MaxTimes(2)
	mockTF.EXPECT().FlushSeriesBucket()
	mockTF.EXPECT().FlushVersion(gomock.Any(), gomock.Any()).Return().AnyTimes()

	// no data flushed, tStores is empty
	tagIdxInterface.FlushVersionDataTo(mockTF, flushContext{})

	// tStore is not empty
	mockTStore1 := NewMocktStoreINTF(ctrl)
	mockTStore1.EXPECT().FlushSeriesTo(gomock.Any(), gomock.Any()).AnyTimes()
	tagIdx.seriesID2TStore = newMetricMap()
	tagIdx.seriesID2TStore.put(1, mockTStore1)
	tagIdx.seriesID2TStore.put(2, mockTStore1)
	// data flushed
	tagIdxInterface.FlushVersionDataTo(mockTF, flushContext{})
}

type mockTagKey struct {
}

func (mockTagKey) TagKey() string {
	return "host"
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
