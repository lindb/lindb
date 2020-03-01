package metricsmeta

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
)

var bitmapUnmarshal = encoding.BitmapUnmarshal

func Test_newTagKVEntrySet_error_cases(t *testing.T) {
	// block length too short, 8 bytes
	_, err := newTagKVEntrySet([]byte{16, 86, 104, 89, 32, 63, 84, 101})
	assert.Error(t, err)
	// validate offsets failure
	_, err = newTagKVEntrySet([]byte{
		1, 1, 1, 1,
		2, 2, 2, 2,
		3, 3, 3, 3,
		4, 4, 4, 4,
		5})
	assert.Error(t, err)
}

func TestTagKVEntries_GetTagValueIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// case 1: get tag value ids
	es, err := newTagKVEntrySet(mockTagValues(10))
	assert.NoError(t, err)
	es2, err := newTagKVEntrySet(mockTagValues(101))
	assert.NoError(t, err)
	entries := TagKVEntries{es, es2}
	tagValueIDs, err := entries.GetTagValueIDs()
	assert.NoError(t, err)
	assert.EqualValues(t, roaring.BitmapOf(10, 101).ToArray(), tagValueIDs.ToArray())
	// case 2: get tag value ids err
	es3 := NewMockTagKVEntrySetINTF(ctrl)
	entries = TagKVEntries{es, es3}
	es3.EXPECT().TagValueIDs().Return(nil, fmt.Errorf("err"))
	tagValueIDs, err = entries.GetTagValueIDs()
	assert.Error(t, err)
	assert.Nil(t, tagValueIDs)
}

func TestTagKVEntrySet_TagValueSeq(t *testing.T) {
	es, err := newTagKVEntrySet(mockTagValues(1<<8 + 1))
	assert.NoError(t, err)
	assert.Equal(t, uint32(1<<8+1), es.TagValueSeq())
}

func TestTagKVEntrySet_CollectTagValues(t *testing.T) {
	defer func() {
		encoding.BitmapUnmarshal = bitmapUnmarshal
		trieTreeFunc = createTrieTree
	}()
	// mock data
	nopKVFlusher := kv.NewNopFlusher()
	seriesFlusher := NewTagFlusher(nopKVFlusher)
	for i := 0; i < 1000; i++ {
		seriesFlusher.FlushTagValue(fmt.Sprintf("test-%d", i), uint32(i))
	}
	_ = seriesFlusher.FlushTagKeyID(22, 30)
	data := nopKVFlusher.Bytes()
	es, err := newTagKVEntrySet(data)
	assert.NoError(t, err)
	// case 1: collect tag values
	result := make(map[uint32]string)
	tagValueIDs := roaring.BitmapOf(1, 3, 1002)
	err = es.CollectTagValues(tagValueIDs, result)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "test-1", result[1])
	assert.Equal(t, "test-3", result[3])
	assert.Equal(t, uint64(1), tagValueIDs.GetCardinality())
	assert.True(t, tagValueIDs.Contains(1002))
	// case 2: tag value ids not exist
	err = es.CollectTagValues(tagValueIDs, nil)
	assert.NoError(t, err)
	assert.True(t, tagValueIDs.Contains(1002))
	// case 3: get value ids err
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		return fmt.Errorf("err")
	}
	err = es.CollectTagValues(tagValueIDs, nil)
	assert.Error(t, err)
	encoding.BitmapUnmarshal = bitmapUnmarshal
	// case 4: create trie tree err
	trieTreeFunc = func(entrySet *tagKVEntrySet) (querier trieTreeQuerier, err error) {
		return nil, fmt.Errorf("err")
	}
	err = es.CollectTagValues(roaring.BitmapOf(1, 2, 3), nil)
	assert.Error(t, err)
}

func TestTagKVEntrySet_GetTagValueID(t *testing.T) {
	es, err := newTagKVEntrySet(mockTagValues(1<<8 + 1))
	assert.NoError(t, err)
	assert.Equal(t, uint32(1<<8+1), es.GetTagValueID(0))
	es, err = newTagKVEntrySet(mockTagValues(1<<16 + 1))
	assert.NoError(t, err)
	assert.Equal(t, uint32(1<<16+1), es.GetTagValueID(0))
	es, err = newTagKVEntrySet(mockTagValues(1<<24 + 1))
	assert.NoError(t, err)
	assert.Equal(t, uint32(1<<24+1), es.GetTagValueID(0))
}

func Test_tagKVEntrySet_TrieTree_error_cases(t *testing.T) {
	zoneBlock, _, _ := buildTagTrieBlock()
	entrySetIntf, _ := newTagKVEntrySet(zoneBlock)
	entrySet := entrySetIntf.(*tagKVEntrySet)
	// read stream eof
	entrySet.sr.Reset([]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 1, 1, 1, 1})
	// read stream eof
	_, err := entrySet.TrieTree()
	assert.Error(t, err)

	// failed validation of trie tree
	entrySet.sr.Reset([]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 1, 1, 1, 1, 1, 1})
	_, err = entrySet.TrieTree()
	assert.Error(t, err)

	// LOUDS block unmarshal failed
	entrySet.sr.Reset([]byte{1, 2, 3, 4, 5, 6, 7, 8, 6, 1, 1, 1, 1, 1, 1})
	_, err = entrySet.TrieTree()
	assert.Error(t, err)

	// isPrefixKey block unmarshal failed
	out, _ := NewRankSelect().MarshalBinary()
	badBLOCK := append([]byte{1, 2, 3, 4, 5, 6, 7, 8,
		18,   // trie tree length
		1, 1, // labels
		1, 1, // is prefix
		13}) // louds

	badBLOCK = append(badBLOCK, out...) // LOUDS block
	entrySet.sr.Reset(badBLOCK)
	_, err = entrySet.TrieTree()
	assert.Error(t, err)
}

func mockTagValues(tagValueID uint32) []byte {
	nopKVFlusher := kv.NewNopFlusher()
	seriesFlusher := NewTagFlusher(nopKVFlusher)
	seriesFlusher.FlushTagValue("test", tagValueID)
	_ = seriesFlusher.FlushTagKeyID(22, tagValueID)
	return nopKVFlusher.Bytes()
}
