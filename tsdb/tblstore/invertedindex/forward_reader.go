package invertedindex

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source ./forward_reader.go -destination=./forward_reader_mock.go -package invertedindex

// ForwardReader represents read forward index data(series id=>tag value id)
type ForwardReader interface {
	series.Grouping
	// GetSeriesIDsForTagKeyID returns series ids for spec metric's tag key id
	GetSeriesIDsForTagKeyID(tagKeyID uint32) (*roaring.Bitmap, error)
}

// forwardReader implements ForwardReader
type forwardReader struct {
	readers []table.Reader
}

// NewForwardReader creates a Reader for reading forward index
func NewForwardReader(readers []table.Reader) ForwardReader {
	return &forwardReader{
		readers: readers,
	}
}

// GetSeriesIDsForTagKeyID get series ids for spec metric's tag key id
func (r *forwardReader) GetSeriesIDsForTagKeyID(tagKeyID uint32) (*roaring.Bitmap, error) {
	seriesIDs := roaring.New()
	if err := r.findReader(tagKeyID, func(reader TagForwardReader) {
		seriesIDs.Or(reader.getSeriesIDs())
	}); err != nil {
		return nil, err
	}
	return seriesIDs, nil
}

// GetGroupingScanner returns the grouping scanners based on tag key ids and series ids
func (r *forwardReader) GetGroupingScanner(tagKeyID uint32, seriesIDs *roaring.Bitmap) ([]series.GroupingScanner, error) {
	var scanners []series.GroupingScanner
	if err := r.findReader(tagKeyID, func(reader TagForwardReader) {
		// check reader if has series ids(after filtering)
		finalSeriesIDs := roaring.FastAnd(seriesIDs, reader.getSeriesIDs())
		if finalSeriesIDs.IsEmpty() {
			// not found
			return
		}
		// found series ids in the sst file
		scanners = append(scanners, reader)
	}); err != nil {
		return nil, err
	}
	return scanners, nil
}

// findReader finds the tag forward reader by tag key id, if reader exist, will invoke callback function
func (r *forwardReader) findReader(tagKeyID uint32, callback func(reader TagForwardReader)) error {
	for _, reader := range r.readers {
		value, ok := reader.Get(tagKeyID)
		if !ok {
			continue
		}
		indexReader, err := newTagForwardReader(value)
		if err != nil {
			return err
		}

		callback(indexReader)
	}
	return nil
}

// TagForwardReader represents the forward index inverterReader for one tag(series id=>tag value id)
type TagForwardReader interface {
	series.GroupingScanner
	// getSeriesIDs gets all series ids under this tag key
	getSeriesIDs() *roaring.Bitmap
}

// tagForwardReader implements TagForwardReader interface
type tagForwardReader struct {
	baseReader
	tagValueIDs *encoding.DeltaBitPackingDecoder
}

// newTagForwardReader creates a forward index inverterReader
func newTagForwardReader(buf []byte) (TagForwardReader, error) {
	r := &tagForwardReader{baseReader: baseReader{
		buf: buf,
	}}
	if err := r.initReader(); err != nil {
		return nil, err
	}
	return r, nil
}

// GetSeriesAndTagValue returns group by container and tag value ids
func (r *tagForwardReader) GetSeriesAndTagValue(highKey uint16) (roaring.Container, []uint32) {
	index := r.keys.GetContainerIndex(highKey)
	if index < 0 {
		// data not found
		return nil, nil
	}
	if r.tagValueIDs == nil {
		r.tagValueIDs = encoding.NewDeltaBitPackingDecoder(r.buf[r.offsets.Get(index):])
	} else {
		r.tagValueIDs.Reset(r.buf[r.offsets.Get(index):])
	}
	container := r.keys.GetContainerAtIndex(index)
	tagValueIDsCount := container.GetCardinality()
	tagValueIDs := make([]uint32, tagValueIDsCount)
	i := 0
	for r.tagValueIDs.HasNext() {
		tagValueIDs[i] = uint32(r.tagValueIDs.Next())
		i++
	}
	return container, tagValueIDs
}

// getSeriesIDs gets all series ids under this tag key
func (r *tagForwardReader) getSeriesIDs() *roaring.Bitmap {
	return r.keys
}

// tagForwardScanner represents the tag forward index scanner which scans the index data when merge operation
type tagForwardScanner struct {
	reader          *tagForwardReader
	container       roaring.Container
	tagValueOffsets *encoding.FixedOffsetDecoder
	highKeys        []uint16
	highKey         uint16
	keyPos          int
	tagValueIDs     *encoding.DeltaBitPackingDecoder
}

// newTagForwardScanner creates a tag forward index scanner
func newTagForwardScanner(reader TagForwardReader) *tagForwardScanner {
	forwardReader := reader.(*tagForwardReader)
	s := &tagForwardScanner{
		reader:   forwardReader,
		highKeys: forwardReader.keys.GetHighKeys(),
	}
	s.nextContainer()
	return s
}

// nextContainer goes next container context for scanner
func (s *tagForwardScanner) nextContainer() {
	s.highKey = s.highKeys[s.keyPos]
	s.container = s.reader.keys.GetContainerAtIndex(s.keyPos)
	s.tagValueOffsets = encoding.NewFixedOffsetDecoder(s.reader.buf[s.reader.offsets.Get(s.keyPos):])
	if s.tagValueIDs == nil {
		s.tagValueIDs = encoding.NewDeltaBitPackingDecoder(s.reader.buf[s.reader.offsets.Get(s.keyPos):])
	} else {
		s.tagValueIDs.Reset(s.reader.buf[s.reader.offsets.Get(s.keyPos):])
	}
	s.keyPos++
}

// scan scans the data then merges the tag value ids into target tag value ids
func (s *tagForwardScanner) scan(highKey, lowSeriesID uint16, tagValueIDs []uint32) []uint32 {
	if s.highKey < highKey {
		if s.keyPos >= len(s.highKeys) {
			// current tag inverted no data can read
			return tagValueIDs
		}
		s.nextContainer()
	}
	if highKey != s.highKey {
		// high key not match, return it
		return tagValueIDs
	}
	// find data by low tag value id
	if s.container.Contains(lowSeriesID) {
		tagValueIDs = append(tagValueIDs, uint32(s.tagValueIDs.Next()))
	}
	return tagValueIDs
}
