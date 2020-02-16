package invertedindex

import (
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./reader.go -destination=./reader_mock.go -package invertedindex

const (
	invertedIndexFooterSize = 4 + // tag value ids position
		4 + // high offsets position
		4 // crc32 checksum
)

// Reader reads versioned seriesID bitmap from series-index-table
type Reader interface {
	// GetSeriesIDsForTagKeyID get series ids for spec metric's tag key id
	GetSeriesIDsForTagKeyID(tagKeyID uint32) (*roaring.Bitmap, error)
	// FindSeriesIDsByTagValueIDs finds series ids by tag key id and tag value ids
	FindSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error)
}

// reader implements Reader
type reader struct {
	readers []table.Reader
}

// NewReader creates a Reader for reading inverted index
func NewReader(readers []table.Reader) Reader {
	return &reader{
		readers: readers,
	}
}

// GetSeriesIDsForTagKeyID get series ids for spec metric's tag key id
func (r *reader) GetSeriesIDsForTagKeyID(tagKeyID uint32) (*roaring.Bitmap, error) {
	fn := func(indexReader *invertedIndexReader) (*roaring.Bitmap, error) {
		return indexReader.getSeriesIDsForTagKeyID()
	}
	return r.loadSeriesIDs(tagKeyID, fn)
}

// FindSeriesIDsByTagValueIDs finds series ids by tag key id and tag value ids
func (r *reader) FindSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	if tagValueIDs == nil || tagValueIDs.IsEmpty() {
		return roaring.New(), nil
	}
	fn := func(indexReader *invertedIndexReader) (*roaring.Bitmap, error) {
		return indexReader.findSeriesIDsByTagValueIDs(tagValueIDs)
	}
	return r.loadSeriesIDs(tagKeyID, fn)
}

// loadSeriesIDs loads the series ids by tag key id, function need implement condition
func (r *reader) loadSeriesIDs(tagKeyID uint32, fn func(indexReader *invertedIndexReader) (*roaring.Bitmap, error)) (*roaring.Bitmap, error) {
	seriesIDs := roaring.New()
	for _, reader := range r.readers {
		value, ok := reader.Get(tagKeyID)
		if !ok {
			continue
		}
		indexReader := newInvertedIndexReader(value)
		ids, err := fn(indexReader)
		if err != nil {
			return nil, err
		}
		seriesIDs.Or(ids)
	}
	return seriesIDs, nil

}

// invertedIndexReader represents the inverted index reader for one tag(tag value ids=>series ids)
type invertedIndexReader struct {
	buf           []byte
	highOffsets   *encoding.FixedOffsetDecoder
	tagValueIDs   *roaring.Bitmap
	crc32CheckSum uint32
	init          bool
}

// newInvertedIndexReader creates an inverted index reader
func newInvertedIndexReader(buf []byte) *invertedIndexReader {
	r := &invertedIndexReader{buf: buf}
	return r
}

// getSeriesIDsForTagKeyID gets all series ids under this tag key
func (r *invertedIndexReader) getSeriesIDsForTagKeyID() (*roaring.Bitmap, error) {
	seriesIDs := roaring.New()
	// first value is tag level's series ids(tag value id is 0)
	if err := seriesIDs.UnmarshalBinary(r.buf); err != nil {
		return nil, err
	}
	return seriesIDs, nil
}

// findSeriesIDsByTagValueIDs finds series ids by tag value ids under this tag key
func (r *invertedIndexReader) findSeriesIDsByTagValueIDs(tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error) {
	if !r.init {
		// if not init, init read
		if err := r.initReader(); err != nil {
			return nil, err
		}
		r.init = true
	}
	result := roaring.New()
	// get final tag value ids need to load
	finalTagValueIDs := roaring.And(tagValueIDs, r.tagValueIDs)
	highKeys := finalTagValueIDs.GetHighKeys()
	for idx, highKey := range highKeys {
		loadLowContainer := finalTagValueIDs.GetContainerAtIndex(idx)
		lowContainerIdx := r.tagValueIDs.GetContainerIndex(highKey)
		lowContainer := r.tagValueIDs.GetContainerAtIndex(lowContainerIdx)
		seriesOffsets := encoding.NewFixedOffsetDecoder(r.buf[r.highOffsets.Get(lowContainerIdx):])
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

// initReader initializes the reader context includes tag value ids/high offsets
func (r *invertedIndexReader) initReader() error {
	if len(r.buf) <= invertedIndexFooterSize {
		return fmt.Errorf("block length no ok")
	}
	// read footer(4+4+4)
	footerPos := len(r.buf) - invertedIndexFooterSize
	tagValueIDsStartPos := int(stream.ReadUint32(r.buf, footerPos))
	highOffsetsPos := int(stream.ReadUint32(r.buf, footerPos+4))
	r.crc32CheckSum = stream.ReadUint32(r.buf, footerPos+8)
	// validate offsets
	if tagValueIDsStartPos > footerPos || highOffsetsPos > tagValueIDsStartPos {
		return fmt.Errorf("bad offsets")
	}
	// read tag value ids
	tagValueIDs := roaring.New()
	if err := encoding.BitmapUnmarshal(tagValueIDs, r.buf[tagValueIDsStartPos:]); err != nil {
		return err
	}
	r.tagValueIDs = tagValueIDs
	// read high offsets
	r.highOffsets = encoding.NewFixedOffsetDecoder(r.buf[highOffsetsPos:])
	return nil
}
