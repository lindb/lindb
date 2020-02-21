package indexdb

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

//go:generate mockgen -source ./tag_index.go -destination=./tag_index_mock.go -package=indexdb

// TagIndex represents the tag inverted index
type TagIndex interface {
	// GetGroupingScanner returns the grouping scanners based on series ids
	GetGroupingScanner(seriesIDs *roaring.Bitmap) ([]series.GroupingScanner, error)
	// buildInvertedIndex builds inverted index for tag value id
	buildInvertedIndex(tagValueID uint32, seriesID uint32)
	// getSeriesIDsByTagValueIDs returns series ids by tag value ids
	getSeriesIDsByTagValueIDs(tagValueIDs *roaring.Bitmap) *roaring.Bitmap
	// getValues returns the all tag values and series ids
	getValues() *InvertedStore
	// getAllSeriesIDs returns all series ids
	getAllSeriesIDs() *roaring.Bitmap
	// flush flushes tag index under spec tag key,
	// write series ids of tag key level with constants.TagValueIDForTag
	flush(tagKeyID uint32, forward invertedindex.ForwardFlusher, inverted invertedindex.InvertedFlusher) error
}

// memGroupingScanner implements series.GroupingScanner for memory tag index
type memGroupingScanner struct {
	forward *ForwardStore
}

// GetSeriesAndTagValue returns group by container and tag value ids
func (g *memGroupingScanner) GetSeriesAndTagValue(highKey uint16) (roaring.Container, []uint32) {
	index := g.forward.keys.GetContainerIndex(highKey)
	if index < 0 {
		// data not found
		return nil, nil
	}
	return g.forward.keys.GetContainerAtIndex(index), g.forward.values[index]
}

// tagIndex is a inverted mapping relation of tag-value and seriesID group.
type tagIndex struct {
	forward  *ForwardStore  // store forward index, series id=>tag value id, maybe have same tag value id
	inverted *InvertedStore // store all tag value id=>series ids of tag level
}

// newTagKVEntrySet returns a new tagKVEntrySet
func newTagIndex() TagIndex {
	return &tagIndex{
		inverted: NewInvertedStore(),
		forward:  NewForwardStore(),
	}
}

// GetGroupingScanner returns the grouping scanners based on series ids
func (index *tagIndex) GetGroupingScanner(seriesIDs *roaring.Bitmap) ([]series.GroupingScanner, error) {
	// check reader if has series ids(after filtering)
	finalSeriesIDs := roaring.FastAnd(seriesIDs, index.forward.Keys())
	if finalSeriesIDs.IsEmpty() {
		// not found
		return nil, nil
	}
	return []series.GroupingScanner{&memGroupingScanner{index.forward}}, nil
}

// buildInvertedIndex builds inverted index for tag value id
func (index *tagIndex) buildInvertedIndex(tagValueID uint32, seriesID uint32) {
	seriesIDs, ok := index.inverted.Get(tagValueID)
	if !ok {
		// create new series ids for new tag value
		seriesIDs = roaring.NewBitmap()
		index.inverted.Put(tagValueID, seriesIDs)
	}
	seriesIDs.Add(seriesID)

	// build forward index, because series id is an unique id, so just put into forward index
	index.forward.Put(seriesID, tagValueID)
}

// getSeriesIDsByTagValueIDs returns series ids by tag value ids
func (index *tagIndex) getSeriesIDsByTagValueIDs(tagValueIDs *roaring.Bitmap) *roaring.Bitmap {
	result := roaring.New()
	values := index.inverted.Values()
	keys := index.inverted.Keys()
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
	return index.forward.keys.Clone()
}

// getValues returns the all tag values and series ids
func (index *tagIndex) getValues() *InvertedStore {
	return index.inverted
}

// flush flushes tag index under spec tag key,
// write series ids of tag key level with constants.TagValueIDForTag
func (index *tagIndex) flush(tagKeyID uint32,
	forward invertedindex.ForwardFlusher, inverted invertedindex.InvertedFlusher,
) error {
	for _, tagValueIDs := range index.forward.values {
		forward.FlushForwardIndex(tagValueIDs)
	}
	if err := forward.FlushTagKeyID(tagKeyID, index.forward.keys); err != nil {
		return err
	}
	// write each tag value series ids
	if err := index.inverted.WalkEntry(inverted.FlushInvertedIndex); err != nil {
		return err
	}
	if err := inverted.FlushTagKeyID(tagKeyID); err != nil {
		return err
	}
	return nil
}
