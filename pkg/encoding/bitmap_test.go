package encoding

import (
	"fmt"
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func TestBitmapMarshal(t *testing.T) {
	defer func() {
		BitmapMarshal = bitmapMarshal
		BitmapUnmarshal = bitmapUnmarshal
	}()
	data, err := BitmapMarshal(roaring.BitmapOf(1))
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)
	BitmapMarshal = func(bitmap *roaring.Bitmap) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	_, err = BitmapMarshal(roaring.BitmapOf(1))
	assert.Error(t, err)

	bitmap := roaring.New()
	err = BitmapUnmarshal(bitmap, data)
	assert.NoError(t, err)
	assert.EqualValues(t, roaring.BitmapOf(1).ToArray(), bitmap.ToArray())

	BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		return fmt.Errorf("err")
	}
	err = BitmapUnmarshal(bitmap, data)
	assert.Error(t, err)
}
