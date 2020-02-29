package invertedindex

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
)

func TestForwardMerger_Merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	merge := NewForwardMerger()
	// case 1: merge data success
	data, err := merge.Merge(1, mockMergeForwardBlock())
	assert.NoError(t, err)
	reader, err := newTagForwardReader(data)
	assert.NoError(t, err)
	assert.EqualValues(t,
		roaring.BitmapOf(1, 2, 3, 4, 65535+10, 65535+20, 65535+30, 65535+40).ToArray(),
		reader.getSeriesIDs().ToArray())
	_, tagValueIDs := reader.GetSeriesAndTagValue(0)
	assert.Equal(t, []uint32{1, 2, 3, 4}, tagValueIDs)
	_, tagValueIDs = reader.GetSeriesAndTagValue(1)
	assert.Equal(t, []uint32{10, 20, 30, 40}, tagValueIDs)
	// case 2: new reader err
	data, err = merge.Merge(1, [][]byte{{1, 2, 3}})
	assert.Error(t, err)
	assert.Nil(t, data)
	// case 3: flush tag key data err
	flusher := NewMockForwardFlusher(ctrl)
	m := merge.(*forwardMerger)
	m.forwardFlusher = flusher
	flusher.EXPECT().FlushForwardIndex(gomock.Any()).AnyTimes()
	flusher.EXPECT().FlushTagKeyID(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	data, err = merge.Merge(1, mockMergeForwardBlock())
	assert.Error(t, err)
	assert.Nil(t, data)
}

func mockMergeForwardBlock() (block [][]byte) {
	nopKVFlusher := kv.NewNopFlusher()
	forwardFlusher := NewForwardFlusher(nopKVFlusher)
	forwardFlusher.FlushForwardIndex([]uint32{1, 3})
	forwardFlusher.FlushForwardIndex([]uint32{10, 20})
	_ = forwardFlusher.FlushTagKeyID(10, roaring.BitmapOf(1, 3, 65535+10, 65535+20))
	block = append(block, nopKVFlusher.Bytes())

	// create new nop flusher, because under nop flusher share buffer
	nopKVFlusher = kv.NewNopFlusher()
	forwardFlusher = NewForwardFlusher(nopKVFlusher)
	forwardFlusher.FlushForwardIndex([]uint32{2, 4})
	forwardFlusher.FlushForwardIndex([]uint32{30, 40})
	_ = forwardFlusher.FlushTagKeyID(10, roaring.BitmapOf(2, 4, 65535+30, 65535+40))
	block = append(block, nopKVFlusher.Bytes())
	return
}
