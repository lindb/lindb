package tagkeymeta

import (
	"io"
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func TestMerger_Merge_Success(t *testing.T) {
	merger := NewMerger()
	merger.Init(nil)

	data, err := merger.Merge(20, mockMergeData())
	assert.NoError(t, err)
	meta, err := newTagKeyMeta(data)
	assert.NoError(t, err)
	// validate TagValueIDSeq
	assert.Equal(t, uint32(200), meta.TagValueIDSeq())
	// validate tagValueIDs
	tagValueIDs := roaring.BitmapOf(1, 2, 3, 4, 6, 7, 8, 9)
	result, _ := meta.TagValueIDs()
	assert.EqualValues(t, tagValueIDs.ToArray(), result.ToArray())

	// validate trie tree
	itr, err := meta.PrefixIterator(nil)
	assert.NoError(t, err)
	var ips []string
	var ids []uint32
	for itr.Valid() {
		ips = append(ips, string(itr.Key()))
		ids = append(ids, encoding.ByteSlice2Uint32(itr.Value()))
		itr.Next()
	}
	assert.Equal(t, []uint32{1, 2, 3, 4, 6, 7, 8, 9}, ids)
	assert.Equal(t, []string{
		"192.168.1.1",
		"192.168.1.2",
		"192.168.1.3",
		"192.168.1.4",
		"192.168.1.6",
		"192.168.1.7",
		"192.168.1.8",
		"192.168.1.9",
	}, ips)
}

func Test_Merger_error(t *testing.T) {
	assert.Nil(t, cloneSlice(nil))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// case1: bad tagKeyMeta
	metaMerger1 := NewMerger()
	_, err := metaMerger1.Merge(20, append([][]byte{{1}}, mockMergeData()...))
	assert.Error(t, err)

	// case2: flush error
	metaMerger2 := NewMerger()
	mergerImpl2 := metaMerger2.(*merger)

	mockFlusher := NewMockFlusher(ctrl)
	mockFlusher.EXPECT().FlushTagValue(gomock.Any(), gomock.Any()).AnyTimes()
	mockFlusher.EXPECT().FlushTagKeyID(gomock.Any(), gomock.Any()).Return(io.ErrClosedPipe)
	mergerImpl2.flusher = mockFlusher

	data, err := mergerImpl2.Merge(20, mockMergeData())
	assert.Error(t, err)
	assert.Len(t, data, 0)
}

func mockMergeData() (data [][]byte) {
	nopKVFlusher1 := kv.NewNopFlusher()
	flusher1 := NewFlusher(nopKVFlusher1)
	flusher1.FlushTagValue([]byte("192.168.1.1"), 1)
	flusher1.FlushTagValue([]byte("192.168.1.2"), 2)
	flusher1.FlushTagValue([]byte("192.168.1.3"), 3)
	flusher1.FlushTagValue([]byte("192.168.1.4"), 4)
	_ = flusher1.FlushTagKeyID(20, 20)
	data = append(data, nopKVFlusher1.Bytes())

	nopKVFlusher2 := kv.NewNopFlusher()
	flusher2 := NewFlusher(nopKVFlusher2)
	flusher2.FlushTagValue([]byte("192.168.1.9"), 9)
	_ = flusher2.FlushTagKeyID(20, 200)
	data = append(data, nopKVFlusher2.Bytes())

	nopKVFlusher3 := kv.NewNopFlusher()
	flusher3 := NewFlusher(nopKVFlusher3)
	flusher3.FlushTagValue([]byte("192.168.1.7"), 7)
	flusher3.FlushTagValue([]byte("192.168.1.6"), 6)
	flusher3.FlushTagValue([]byte("192.168.1.8"), 8)
	_ = flusher3.FlushTagKeyID(20, 100)
	data = append(data, nopKVFlusher3.Bytes())
	return data
}
