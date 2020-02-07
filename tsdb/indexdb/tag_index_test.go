package indexdb

import (
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func TestTagIndex_buildInvertedIndex(t *testing.T) {
	index := newTagIndex()
	index.buildInvertedIndex(2, 1)
	index.buildInvertedIndex(2, 3)
	index.buildInvertedIndex(2, 2)
	index.buildInvertedIndex(1, 1)
	index.buildInvertedIndex(1, 2)
	values := index.getValues()
	seriesIDs, ok := values.Get(1)
	assert.True(t, ok)
	assert.Equal(t, roaring.BitmapOf(1, 2), seriesIDs)
	seriesIDs, ok = values.Get(1)
	assert.True(t, ok)
	assert.Equal(t, roaring.BitmapOf(1, 2), seriesIDs)
	assert.Equal(t, roaring.BitmapOf(1, 2, 3), index.getAllSeriesIDs())
}

func TestTagIndex_getSeriesIDsByTagValueIDs(t *testing.T) {
	tagIndex := prepareTagIdx()
	// tag-value not exist
	assert.Equal(t, roaring.New(), tagIndex.getSeriesIDsByTagValueIDs(roaring.BitmapOf(40, 50, 30)))
	// tag-value exist
	assert.Equal(t, roaring.BitmapOf(4), tagIndex.getSeriesIDsByTagValueIDs(roaring.BitmapOf(4)))
}

func TestTagIndex_getAllSeriesIDs(t *testing.T) {
	tagIndex := prepareTagIdx()
	assert.Equal(t, roaring.BitmapOf(1, 2, 3, 4, 5, 6, 7, 8), tagIndex.getAllSeriesIDs())
}

func prepareTagIdx() TagIndex {
	tagIndex := newTagIndex()
	tagIndex.buildInvertedIndex(1, 1)
	tagIndex.buildInvertedIndex(2, 2)
	tagIndex.buildInvertedIndex(3, 3)
	tagIndex.buildInvertedIndex(4, 4)
	tagIndex.buildInvertedIndex(5, 5)
	tagIndex.buildInvertedIndex(6, 6)
	tagIndex.buildInvertedIndex(7, 7)
	tagIndex.buildInvertedIndex(8, 8)
	return tagIndex
}
