package tblstore

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/pkg/timeutil"

	"github.com/lindb/lindb/kv"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_InvertedIndexFlusher_Commit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	indexFlusher := NewInvertedIndexFlusher(mockFlusher)
	assert.NotNil(t, indexFlusher)

	// mock commit error
	mockFlusher.EXPECT().Commit().Return(fmt.Errorf("commit error"))
	assert.NotNil(t, indexFlusher.Commit())

	// mock commit ok
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	err := indexFlusher.FlushTagKeyID(333)
	assert.Nil(t, err)
}

func Test_InvertedIndexFlusher_RS_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	indexFlusher := NewInvertedIndexFlusher(mockFlusher).(*invertedIndexFlusher)

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
	assert.NotNil(t, indexFlusher.FlushTagKeyID(11))
	// mock LOUDS error
	mockRS.EXPECT().MarshalBinary().Return([]byte("12345"), nil)
	mockRS.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("error2"))
	assert.NotNil(t, indexFlusher.FlushTagKeyID(11))
}

func Test_SeriesIndexFlusher_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
	indexFlusher := NewInvertedIndexFlusher(mockFlusher)

	// flush versions of tagValue1
	indexFlusher.FlushVersion(1, timeutil.TimeRange{Start: 3, End: 5}, roaring.New())
	indexFlusher.FlushVersion(2, timeutil.TimeRange{Start: 4, End: 6}, roaring.New())
	indexFlusher.FlushVersion(3, timeutil.TimeRange{Start: 1, End: 2}, roaring.New())
	// flush tagValue1
	indexFlusher.FlushTagValue("1")
	// flush versions of tagValue2
	indexFlusher.FlushVersion(1, timeutil.TimeRange{Start: 12, End: 15}, roaring.New())
	indexFlusher.FlushVersion(2, timeutil.TimeRange{Start: 15, End: 20}, roaring.New())
	indexFlusher.FlushVersion(3, timeutil.TimeRange{Start: 22, End: 24}, roaring.New())
	// flush tagValue2
	indexFlusher.FlushTagValue("2")
	// flush tagKeyID
	assert.Nil(t, indexFlusher.FlushTagKeyID(0))
}
