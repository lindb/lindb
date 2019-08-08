package indextbl

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/kv"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func prepareTestKVSets() []VersionedTagKVEntrySet {
	bitmap1 := roaring.New()
	bitmap1.AddRange(1, 100)
	bitmap2 := roaring.New()
	bitmap2.AddRange(300, 400)

	return []VersionedTagKVEntrySet{
		{Version: 1, EntrySet: map[string]*roaring.Bitmap{
			"nj": bitmap1,
			"bj": bitmap2}},
		{Version: 2, EntrySet: map[string]*roaring.Bitmap{
			"nt": roaring.New(),
			"sz": roaring.New()}},
	}
}

func Test_SeriesIndexFlusher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	indexFlusher := NewSeriesIndexFlusher(mockFlusher)
	assert.NotNil(t, indexFlusher)

	// mock commit error
	mockFlusher.EXPECT().Commit().Return(fmt.Errorf("commit error"))
	assert.NotNil(t, indexFlusher.Commit())

	// flush tag key, zone=nj, zone=bj
	testData := prepareTestKVSets()
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	err := indexFlusher.FlushTagKey(333, testData)
	assert.Nil(t, err)
}

func Test_SeriesIndexFlusher_RS_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	indexFlusher := NewSeriesIndexFlusher(mockFlusher).(*seriesIndexFlusher)

	testData := prepareTestKVSets()
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	// mock trie tree
	mockTrie := NewMocktrieTreeINTF(ctrl)
	mockTrie.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any()).Return().AnyTimes()
	mockTrie.EXPECT().Reset().Return().AnyTimes()
	// mock rank select
	mockRS := NewMockRSINTF(ctrl)
	// mock binary return
	mockBin := &seriesBinData{LOUDS: mockRS, isPrefixKey: mockRS, labels: nil, values: nil}
	mockTrie.EXPECT().MarshalBinary().Return(mockBin).AnyTimes()
	// replace trie with mock
	indexFlusher.trie = mockTrie

	// mock isPrefixKey error
	mockRS.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("error1"))
	assert.NotNil(t, indexFlusher.FlushTagKey(11, testData))
	// mock LOUDS error
	mockRS.EXPECT().MarshalBinary().Return([]byte("12345"), nil)
	mockRS.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("error2"))
	assert.NotNil(t, indexFlusher.FlushTagKey(11, testData))
}
