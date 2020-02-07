package indexdb

import (
	"github.com/lindb/roaring"
)

// TagIndex represents the tag inverted index
type TagIndex interface {
	// buildInvertedIndex builds inverted index for tag value id
	buildInvertedIndex(tagValueID uint32, seriesID uint32)
	// getSeriesIDsByTagValueIDs returns series ids by tag value ids
	getSeriesIDsByTagValueIDs(tagValueIDs *roaring.Bitmap) *roaring.Bitmap
	// getValues returns the all tag values and series ids
	getValues() *TagStore
	// getAllSeriesIDs returns all series ids
	getAllSeriesIDs() *roaring.Bitmap
}

// tagIndex is a inverted mapping relation of tag-value and seriesID group.
type tagIndex struct {
	seriesIDs *roaring.Bitmap // store all series ids of teg level
	values    *TagStore       // store all tag value id=>series ids of tag level
}

// newTagKVEntrySet returns a new tagKVEntrySet
func newTagIndex() TagIndex {
	return &tagIndex{
		values:    NewTagStore(),
		seriesIDs: roaring.New(),
	}
}

// buildInvertedIndex builds inverted index for tag value id
func (index *tagIndex) buildInvertedIndex(tagValueID uint32, seriesID uint32) {
	seriesIDs, ok := index.values.Get(tagValueID)
	if !ok {
		// create new series ids for new tag value
		seriesIDs = roaring.NewBitmap()
		index.values.Put(tagValueID, seriesIDs)
	}
	seriesIDs.Add(seriesID)
	index.seriesIDs.Add(seriesID)
}

// getSeriesIDsByTagValueIDs returns series ids by tag value ids
func (index *tagIndex) getSeriesIDsByTagValueIDs(tagValueIDs *roaring.Bitmap) *roaring.Bitmap {
	result := roaring.New()
	values := index.values.Values()
	keys := index.values.Keys()
	// get final tag value ids need to load
	finalTagValueIDs := roaring.And(tagValueIDs, keys)
	highKeys := finalTagValueIDs.GetHighKeys()
	for idx, highKey := range highKeys {
		loadLowContainer := finalTagValueIDs.GetContainerAtIndex(idx)
		lowContainerIdx := keys.GetContainerIndex(highKey)
		lowContainer := keys.GetContainerAtIndex(lowContainerIdx)
		it := loadLowContainer.PeekableIterator()
		for it.HasNext() {
			lowTagValueID := it.Next()
			// get the index of low tag value id in container
			lowIdx := lowContainer.Rank(lowTagValueID)
			result.Or(values[lowContainerIdx][lowIdx-1])
		}
	}
	return result
}

// getAllSeriesIDs returns all series ids
func (index *tagIndex) getAllSeriesIDs() *roaring.Bitmap {
	return index.seriesIDs.Clone()
}

// getValues returns the all tag values and series ids
func (index *tagIndex) getValues() *TagStore {
	return index.values
}
