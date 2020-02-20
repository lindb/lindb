package invertedindex

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv/table"
)

//go:generate mockgen -source ./forward_reader.go -destination=./forward_reader_mock.go -package invertedindex

// ForwardReader represents read forward index data(series id=>tag value id)
type ForwardReader interface {
	// GetSeriesIDsForTagKeyID get series ids for spec metric's tag key id
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
	for _, reader := range r.readers {
		value, ok := reader.Get(tagKeyID)
		if !ok {
			continue
		}
		indexReader, err := newTagForwardReader(value)
		if err != nil {
			return nil, err
		}
		seriesIDs.Or(indexReader.getSeriesIDs())
	}
	return seriesIDs, nil
}

// tagForwardReader represents the forward index inverterReader for one tag(series id=>tag value id)
type tagForwardReader struct {
	baseReader
}

// newTagForwardReader creates a forward index inverterReader
func newTagForwardReader(buf []byte) (*tagForwardReader, error) {
	r := &tagForwardReader{baseReader: baseReader{
		buf: buf,
	}}
	if err := r.initReader(); err != nil {
		return nil, err
	}
	return r, nil
}

// getSeriesIDsForTagKeyID gets all series ids under this tag key
func (r *tagForwardReader) getSeriesIDs() *roaring.Bitmap {
	return r.keys
}
