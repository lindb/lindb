package metricsmeta

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
)

func TestTagMerger_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newTagKVEntrySetFunc = newTagKVEntrySet
		ctrl.Finish()
	}()
	merger := NewTagMerger()
	merger.Init(nil)
	// case 1: merge success
	data, err := merger.Merge(20, mockMergeData())
	assert.NoError(t, err)
	tReader, err := newTagKVEntrySet(data)
	assert.NoError(t, err)
	assert.Equal(t, uint32(200), tReader.TagValueSeq())
	tagValueIDs := roaring.BitmapOf(1, 2, 3, 4, 6, 7, 8, 9)
	result, _ := tReader.TagValueIDs()
	assert.EqualValues(t, tagValueIDs.ToArray(), result.ToArray())
	trie, _ := tReader.TrieTree()
	it := tagValueIDs.Iterator()
	c := 0
	for it.HasNext() {
		tagValueID := it.Next()
		values := trie.FindOffsetsByEqual(fmt.Sprintf("192.168.1.%d", tagValueID))
		assert.Equal(t, tagValueID, tReader.GetTagValueID(values[0]))
		c++
	}
	assert.Equal(t, uint64(c), tagValueIDs.GetCardinality())
	// case 2: new tag reader err
	data, err = merger.Merge(20, [][]byte{{1, 2, 3}})
	assert.Error(t, err)
	assert.Nil(t, data)
	// case 3: get trie err
	reader := NewMockTagKVEntrySetINTF(ctrl)
	newTagKVEntrySetFunc = func(block []byte) (intf TagKVEntrySetINTF, err error) {
		return reader, nil
	}
	reader.EXPECT().TagValueSeq().Return(uint32(199)).AnyTimes()
	reader.EXPECT().TrieTree().Return(nil, fmt.Errorf("err"))
	data, err = merger.Merge(20, [][]byte{{1, 2, 3}})
	assert.Error(t, err)
	assert.Nil(t, data)
	// case 4: flush data err
	newTagKVEntrySetFunc = newTagKVEntrySet
	m := merger.(*tagMerger)
	flusher := NewMockTagFlusher(ctrl)
	m.tagFlusher = flusher
	flusher.EXPECT().FlushTagValue(gomock.Any(), gomock.Any()).AnyTimes()
	flusher.EXPECT().FlushTagKeyID(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	data, err = merger.Merge(20, mockMergeData())
	assert.Error(t, err)
	assert.Nil(t, data)
}

func mockMergeData() (data [][]byte) {
	nopKVFlusher := kv.NewNopFlusher()
	tagFlusher := NewTagFlusher(nopKVFlusher)
	tagFlusher.FlushTagValue("192.168.1.1", 1)
	tagFlusher.FlushTagValue("192.168.1.2", 2)
	tagFlusher.FlushTagValue("192.168.1.3", 3)
	tagFlusher.FlushTagValue("192.168.1.4", 4)
	_ = tagFlusher.FlushTagKeyID(20, 20)
	data = append(data, nopKVFlusher.Bytes())

	nopKVFlusher = kv.NewNopFlusher()
	tagFlusher = NewTagFlusher(nopKVFlusher)
	tagFlusher.FlushTagValue("192.168.1.9", 9)
	_ = tagFlusher.FlushTagKeyID(20, 200)
	data = append(data, nopKVFlusher.Bytes())
	nopKVFlusher = kv.NewNopFlusher()
	tagFlusher = NewTagFlusher(nopKVFlusher)
	tagFlusher.FlushTagValue("192.168.1.7", 7)
	tagFlusher.FlushTagValue("192.168.1.6", 6)
	tagFlusher.FlushTagValue("192.168.1.8", 8)
	_ = tagFlusher.FlushTagKeyID(20, 100)
	data = append(data, nopKVFlusher.Bytes())
	return data
}
