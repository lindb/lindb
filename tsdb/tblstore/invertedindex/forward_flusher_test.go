package invertedindex

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
)

func TestForwardFlusher_Flusher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)
	indexFlusher := NewForwardFlusher(mockFlusher)
	assert.NotNil(t, indexFlusher)
	indexFlusher.FlushForwardIndex([]uint32{1, 2, 3, 4})
	indexFlusher.FlushForwardIndex([]uint32{1, 2, 3, 4})
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
	err := indexFlusher.FlushTagKeyID(3, roaring.BitmapOf(1, 2, 3))
	assert.NoError(t, err)
	mockFlusher.EXPECT().Commit().Return(nil)
	err = indexFlusher.Commit()
	assert.NoError(t, err)
}

func TestForwardFlusher_Flush_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.BitmapMarshal = bitMapMarshal
		ctrl.Finish()
	}()

	indexFlusher := NewForwardFlusher(nil)
	assert.NotNil(t, indexFlusher)
	encoding.BitmapMarshal = func(bitmap *roaring.Bitmap) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	err := indexFlusher.FlushTagKeyID(3, roaring.BitmapOf(1, 2, 3))
	assert.Error(t, err)
}
