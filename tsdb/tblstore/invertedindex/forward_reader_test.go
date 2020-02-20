package invertedindex

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
)

func TestForwardReader_GetSeriesIDsForTagKeyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.BitmapUnmarshal = bitmapUnmarshal
		ctrl.Finish()
	}()

	reader := buildForwardReader(ctrl)
	// case 1: read not tagID key
	idSet, err := reader.GetSeriesIDsForTagKeyID(19)
	assert.NoError(t, err)
	assert.Equal(t, roaring.New(), idSet)
	// case 2: data is empty
	idSet, err = reader.GetSeriesIDsForTagKeyID(10)
	assert.Error(t, err)
	assert.Nil(t, idSet)
	// case 3: read series ids
	idSet, err = reader.GetSeriesIDsForTagKeyID(20)
	a := roaring.BitmapOf(1, 2, 3, 65535+10, 65535+20, 65535+30, 65535+40)
	assert.NoError(t, err)
	assert.EqualValues(t, a.ToArray(), idSet.ToArray())
	// case 4: unmarshal series ids err
	reader = buildForwardReader(ctrl)
	encoding.BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		return fmt.Errorf("err")
	}
	idSet, err = reader.GetSeriesIDsForTagKeyID(20)
	assert.Error(t, err)
	assert.Nil(t, idSet)
}

func TestForwardReader_offset_err(t *testing.T) {
	reader, err := newTagForwardReader([]byte{
		1, 1, 1, 1,
		2, 2, 2, 2,
		3, 3, 3, 3,
		4, 4, 4, 4,
		5})
	assert.Error(t, err)
	assert.Nil(t, reader)
}

func buildForwardReader(ctrl *gomock.Controller) ForwardReader {
	block := buildForwardBlock()
	// mock readers
	mockReader := table.NewMockReader(ctrl)
	mockReader.EXPECT().Get(uint32(10)).Return(nil, true).AnyTimes()
	mockReader.EXPECT().Get(uint32(19)).Return(nil, false).AnyTimes()
	mockReader.EXPECT().Get(uint32(20)).Return(block, true).AnyTimes()
	// build series index inverterReader
	return NewForwardReader([]table.Reader{mockReader})
}

func buildForwardBlock() (block []byte) {
	nopKVFlusher := kv.NewNopFlusher()
	forwardFlusher := NewForwardFlusher(nopKVFlusher)
	forwardFlusher.FlushForwardIndex([]uint32{1, 2, 3, 4})
	forwardFlusher.FlushForwardIndex([]uint32{10, 20, 30, 40})
	_ = forwardFlusher.FlushTagKeyID(10, roaring.BitmapOf(1, 2, 3, 65535+10, 65535+20, 65535+30, 65535+40))
	return nopKVFlusher.Bytes()
}
