package series

import (
	"testing"

	"github.com/RoaringBitmap/roaring"
	"github.com/stretchr/testify/assert"
)

func TestMultiVerSeriesIDSet_And(t *testing.T) {
	multiVer1 := NewMultiVerSeriesIDSet()
	multiVer1.Add(int64(12), roaring.BitmapOf(1, 2, 3, 4))
	// will ignore
	multiVer1.Add(int64(12), roaring.BitmapOf(1, 2, 3, 4, 5, 6))
	assert.Equal(t, *roaring.BitmapOf(1, 2, 3, 4), *(multiVer1.versions[int64(12)]))

	multiVer2 := NewMultiVerSeriesIDSet()
	multiVer2.Add(int64(12), roaring.BitmapOf(2, 3, 4))

	multiVer1.And(multiVer2)
	assert.Equal(t, *roaring.BitmapOf(2, 3, 4), *(multiVer1.versions[int64(12)]))

	multiVer3 := NewMultiVerSeriesIDSet()
	multiVer3.Add(int64(13), roaring.BitmapOf(2, 3, 4))
	multiVer1.And(multiVer3)
	assert.Equal(t, 0, len(multiVer1.versions))
}

func TestMultiVerSeriesIDSet_Or(t *testing.T) {
	multiVer1 := NewMultiVerSeriesIDSet()
	multiVer1.Add(int64(12), roaring.BitmapOf(1, 4, 5))
	assert.Equal(t, *roaring.BitmapOf(1, 4, 5), *(multiVer1.versions[int64(12)]))

	multiVer2 := NewMultiVerSeriesIDSet()
	multiVer2.Add(int64(12), roaring.BitmapOf(2, 3, 4))

	multiVer1.Or(multiVer2)
	assert.Equal(t, *roaring.BitmapOf(1, 2, 3, 4, 5), *(multiVer1.versions[int64(12)]))

	multiVer3 := NewMultiVerSeriesIDSet()
	multiVer3.Add(int64(13), roaring.BitmapOf(7, 8, 9))
	multiVer1.Or(multiVer3)
	assert.Equal(t, 2, len(multiVer1.versions))
	assert.Equal(t, *roaring.BitmapOf(1, 2, 3, 4, 5), *(multiVer1.versions[int64(12)]))
	assert.Equal(t, *roaring.BitmapOf(7, 8, 9), *(multiVer1.versions[int64(13)]))
}

func TestMultiVerSeriesIDSet_AndNot(t *testing.T) {
	multiVer1 := NewMultiVerSeriesIDSet()
	multiVer1.Add(int64(12), roaring.BitmapOf(1, 2, 4, 6, 7, 8, 9))
	multiVer1.Add(int64(13), roaring.BitmapOf(7, 8, 9))
	assert.Equal(t, *roaring.BitmapOf(1, 2, 4, 6, 7, 8, 9), *(multiVer1.versions[int64(12)]))

	multiVer2 := NewMultiVerSeriesIDSet()
	multiVer2.Add(int64(12), roaring.BitmapOf(2, 3, 4, 9))

	multiVer1.AndNot(multiVer2)
	assert.Equal(t, 2, len(multiVer1.versions))
	assert.Equal(t, *roaring.BitmapOf(1, 6, 7, 8), *(multiVer1.versions[int64(12)]))
	assert.Equal(t, *roaring.BitmapOf(7, 8, 9), *(multiVer1.versions[int64(13)]))

	multiVer3 := NewMultiVerSeriesIDSet()
	multiVer3.Add(int64(14), roaring.BitmapOf(7))
	multiVer1.AndNot(multiVer3)

	assert.Equal(t, 2, len(multiVer1.versions))
	assert.Equal(t, *roaring.BitmapOf(1, 6, 7, 8), *(multiVer1.versions[int64(12)]))
	assert.Equal(t, *roaring.BitmapOf(7, 8, 9), *(multiVer1.versions[int64(13)]))
}
