package series

import (
	"fmt"
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func Test_AAA(t *testing.T) {
	r := roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8, 9)
	r.RemoveRange(0, 8)
	//r1:=r.RemoveRange(0, 8)
	//r1:=r.Flip(0,8)
	fmt.Println(r)
}

func TestMultiVerSeriesIDSet_IsEmpty(t *testing.T) {
	multiVer1 := NewMultiVerSeriesIDSet()
	assert.True(t, multiVer1.IsEmpty())
	multiVer1.Add(Version(12), roaring.BitmapOf())
	assert.True(t, multiVer1.IsEmpty())
	multiVer1 = NewMultiVerSeriesIDSet()
	multiVer1.Add(Version(12), roaring.BitmapOf(1, 2, 3, 4, 5))
	assert.False(t, multiVer1.IsEmpty())
	assert.Len(t, multiVer1.Versions(), 1)
	assert.True(t, multiVer1.Contains(Version(12)))
}

func TestMultiVerSeriesIDSet_And(t *testing.T) {
	multiVer1 := NewMultiVerSeriesIDSet()
	multiVer1.Add(Version(12), roaring.BitmapOf(1, 2, 3, 4))
	// will ignore
	multiVer1.Add(Version(12), roaring.BitmapOf(1, 2, 3, 4, 5, 6))
	assert.Equal(t, *roaring.BitmapOf(1, 2, 3, 4), *(multiVer1.versions[Version(12)]))

	multiVer2 := NewMultiVerSeriesIDSet()
	multiVer2.Add(Version(12), roaring.BitmapOf(2, 3, 4))

	multiVer1.And(multiVer2)
	assert.Equal(t, *roaring.BitmapOf(2, 3, 4), *(multiVer1.versions[Version(12)]))

	multiVer3 := NewMultiVerSeriesIDSet()
	multiVer3.Add(Version(13), roaring.BitmapOf(2, 3, 4))
	multiVer1.And(multiVer3)
	assert.Equal(t, 0, len(multiVer1.versions))
}

func TestMultiVerSeriesIDSet_Or(t *testing.T) {
	multiVer1 := NewMultiVerSeriesIDSet()
	multiVer1.Add(Version(12), roaring.BitmapOf(1, 4, 5))
	assert.Equal(t, *roaring.BitmapOf(1, 4, 5), *(multiVer1.versions[Version(12)]))

	multiVer2 := NewMultiVerSeriesIDSet()
	multiVer2.Add(Version(12), roaring.BitmapOf(2, 3, 4))

	multiVer1.Or(multiVer2)
	assert.Equal(t, *roaring.BitmapOf(1, 2, 3, 4, 5), *(multiVer1.versions[Version(12)]))

	multiVer3 := NewMultiVerSeriesIDSet()
	multiVer3.Add(Version(13), roaring.BitmapOf(7, 8, 9))
	multiVer1.Or(multiVer3)
	assert.Equal(t, 2, len(multiVer1.versions))
	assert.Equal(t, *roaring.BitmapOf(1, 2, 3, 4, 5), *(multiVer1.versions[Version(12)]))
	assert.Equal(t, *roaring.BitmapOf(7, 8, 9), *(multiVer1.versions[Version(13)]))
}

func TestMultiVerSeriesIDSet_AndNot(t *testing.T) {
	multiVer1 := NewMultiVerSeriesIDSet()
	multiVer1.Add(Version(12), roaring.BitmapOf(1, 2, 4, 6, 7, 8, 9))
	multiVer1.Add(Version(13), roaring.BitmapOf(7, 8, 9))
	assert.Equal(t, *roaring.BitmapOf(1, 2, 4, 6, 7, 8, 9), *(multiVer1.versions[Version(12)]))

	multiVer2 := NewMultiVerSeriesIDSet()
	multiVer2.Add(Version(12), roaring.BitmapOf(2, 3, 4, 9))

	multiVer1.AndNot(multiVer2)
	assert.Equal(t, 2, len(multiVer1.versions))
	assert.Equal(t, *roaring.BitmapOf(1, 6, 7, 8), *(multiVer1.versions[Version(12)]))
	assert.Equal(t, *roaring.BitmapOf(7, 8, 9), *(multiVer1.versions[Version(13)]))

	multiVer3 := NewMultiVerSeriesIDSet()
	multiVer3.Add(Version(14), roaring.BitmapOf(7))
	multiVer1.AndNot(multiVer3)

	assert.Equal(t, 2, len(multiVer1.versions))
	assert.Equal(t, *roaring.BitmapOf(1, 6, 7, 8), *(multiVer1.versions[Version(12)]))
	assert.Equal(t, *roaring.BitmapOf(7, 8, 9), *(multiVer1.versions[Version(13)]))
}
