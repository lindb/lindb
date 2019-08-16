package indextbl

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/kv"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_SeriesIndexFlusher_Commit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	indexFlusher := NewSeriesIndexFlusher(mockFlusher)
	assert.NotNil(t, indexFlusher)

	// mock commit error
	mockFlusher.EXPECT().Commit().Return(fmt.Errorf("commit error"))
	assert.NotNil(t, indexFlusher.Commit())

	// mock commit ok
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	err := indexFlusher.FlushTagKey(333)
	assert.Nil(t, err)
}

func Test_SeriesIndexFlusher_RS_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	indexFlusher := NewSeriesIndexFlusher(mockFlusher).(*seriesIndexFlusher)

	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	// mock trie tree
	mockTrie := NewMocktrieTreeBuilder(ctrl)
	mockTrie.EXPECT().Add(gomock.Any(), gomock.Any()).Return().AnyTimes()
	mockTrie.EXPECT().Reset().Return().AnyTimes()
	// mock rank select
	mockRS := NewMockRSINTF(ctrl)
	// mock binary return
	mockBin := &trieTreeData{
		trieTreeBlock: trieTreeBlock{LOUDS: mockRS, isPrefixKey: mockRS, labels: nil},
		values:        nil}

	mockTrie.EXPECT().MarshalBinary().Return(mockBin).AnyTimes()
	// replace trie with mock
	indexFlusher.trie = mockTrie

	// mock isPrefixKey error
	mockRS.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("error1"))
	assert.NotNil(t, indexFlusher.FlushTagKey(11))
	// mock LOUDS error
	mockRS.EXPECT().MarshalBinary().Return([]byte("12345"), nil)
	mockRS.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("error2"))
	assert.NotNil(t, indexFlusher.FlushTagKey(11))
}

func Test_SeriesIndexFlusher_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
	indexFlusher := NewSeriesIndexFlusher(mockFlusher)

	// flush versions of tagValue1
	indexFlusher.FlushVersion(1, 3, 5, roaring.New())
	indexFlusher.FlushVersion(2, 4, 6, roaring.New())
	indexFlusher.FlushVersion(3, 1, 2, roaring.New())
	// flush tagValue1
	indexFlusher.FlushTagValue("1")
	// flush versions of tagValue2
	indexFlusher.FlushVersion(1, 12, 15, roaring.New())
	indexFlusher.FlushVersion(2, 15, 20, roaring.New())
	indexFlusher.FlushVersion(3, 22, 24, roaring.New())
	// flush tagValue2
	indexFlusher.FlushTagValue("2")
	// flush tagKey
	assert.Nil(t, indexFlusher.FlushTagKey(0))
}
