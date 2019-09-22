package series

import (
	"testing"

	"github.com/RoaringBitmap/roaring"
	"github.com/stretchr/testify/assert"
)

func TestSeriesIDsIterator(t *testing.T) {
	buf := make([]uint32, 2)
	it := NewIDsIterator(roaring.BitmapOf(1, 2, 3, 4, 5), buf)
	n, b := it.Next()
	assert.Equal(t, n, 2)
	assert.Equal(t, buf, b)
	n, b = it.Next()
	assert.Equal(t, n, 2)
	assert.Equal(t, buf, b)
	n, b = it.Next()
	assert.Equal(t, n, 1)
	assert.Equal(t, buf, b)
}
