package invertedindex

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
)

//go:generate mockgen -source ./inverted_reader.go -destination=./inverted_reader_mock.go -package invertedindex

// InvertedReader reads seriesID bitmap from series-index-table
type InvertedReader interface {
	// GetSeriesIDsByTagValueIDs finds series ids by tag key id and tag value ids
	GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error)
}

// inverterReader implements InvertedReader
type inverterReader struct {
	readers []table.Reader
}

// NewInvertedReader creates a InvertedReader for reading inverted index
func NewInvertedReader(readers []table.Reader) InvertedReader {
	return &inverterReader{
		readers: readers,
	}
}

// GetSeriesIDsByTagValueIDs finds series ids by tag key id and tag value ids
func (r *inverterReader) GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	if tagValueIDs == nil || tagValueIDs.IsEmpty() {
		return roaring.New(), nil
	}
	fn := func(indexReader *tagInvertedReader) (*roaring.Bitmap, error) {
		return indexReader.getSeriesIDsByTagValueIDs(tagValueIDs)
	}
	return r.loadSeriesIDs(tagKeyID, fn)
}

// loadSeriesIDs loads the series ids by tag key id, function need implement condition
func (r *inverterReader) loadSeriesIDs(tagKeyID uint32, fn func(indexReader *tagInvertedReader) (*roaring.Bitmap, error)) (*roaring.Bitmap, error) {
	seriesIDs := roaring.New()
	for _, reader := range r.readers {
		value, ok := reader.Get(tagKeyID)
		if !ok {
			continue
		}
		indexReader := newTagInvertedReader(value)
		ids, err := fn(indexReader)
		if err != nil {
			return nil, err
		}
		seriesIDs.Or(ids)
	}
	return seriesIDs, nil
}

// tagInvertedReader represents the inverted index inverterReader for one tag(tag value ids=>series ids)
type tagInvertedReader struct {
	baseReader
	init bool
}

// newTagInvertedReader creates an inverted index tagInvertedReader
func newTagInvertedReader(buf []byte) *tagInvertedReader {
	r := &tagInvertedReader{
		baseReader: baseReader{buf: buf},
	}
	return r
}

// getSeriesIDsByTagValueIDs finds series ids by tag value ids under this tag key
func (r *tagInvertedReader) getSeriesIDsByTagValueIDs(tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	if !r.init {
		// if not init, init read
		if err := r.initReader(); err != nil {
			return nil, err
		}
		r.init = true
	}
	result := roaring.New()
	// get final tag value ids need to load
	finalTagValueIDs := roaring.And(tagValueIDs, r.keys)
	highKeys := finalTagValueIDs.GetHighKeys()
	for idx, highKey := range highKeys {
		loadLowContainer := finalTagValueIDs.GetContainerAtIndex(idx)
		lowContainerIdx := r.keys.GetContainerIndex(highKey)
		lowContainer := r.keys.GetContainerAtIndex(lowContainerIdx)
		seriesOffsets := encoding.NewFixedOffsetDecoder(r.buf[r.offsets.Get(lowContainerIdx):])
		it := loadLowContainer.PeekableIterator()
		for it.HasNext() {
			lowTagValueID := it.Next()
			// get the index of low tag value id in container
			lowIdx := lowContainer.Rank(lowTagValueID)
			seriesPos := seriesOffsets.Get(lowIdx - 1)
			// unmarshal series ids
			seriesIDs := roaring.New()
			if err := encoding.BitmapUnmarshal(seriesIDs, r.buf[seriesPos:]); err != nil {
				return nil, err
			}
			result.Or(seriesIDs)
		}
	}
	return result, nil
}
