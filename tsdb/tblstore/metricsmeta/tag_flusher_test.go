package metricsmeta

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/kv"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_InvertedIndexFlusher_Commit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	indexFlusher := NewTagFlusher(mockFlusher)
	assert.NotNil(t, indexFlusher)

	// mock commit error
	mockFlusher.EXPECT().Commit().Return(fmt.Errorf("commit error"))
	assert.NotNil(t, indexFlusher.Commit())

	// mock commit ok
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	err := indexFlusher.FlushTagKeyID(333, 100)
	assert.Nil(t, err)
}

func Test_InvertedIndexFlusher_RS_error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	indexFlusher := NewTagFlusher(mockFlusher).(*tagFlusher)

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
	assert.NotNil(t, indexFlusher.FlushTagKeyID(11, 10))
	// mock LOUDS error
	mockRS.EXPECT().MarshalBinary().Return([]byte("12345"), nil)
	mockRS.EXPECT().MarshalBinary().Return(nil, fmt.Errorf("error2"))
	assert.NotNil(t, indexFlusher.FlushTagKeyID(11, 10))
}

func Test_SeriesIndexFlusher_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
	indexFlusher := NewTagFlusher(mockFlusher)

	// flush tagValue1
	indexFlusher.FlushTagValue("1", 1)
	// flush tagValue2
	indexFlusher.FlushTagValue("2", 2)
	// flush tagKeyID
	assert.Nil(t, indexFlusher.FlushTagKeyID(0, 10))
}
