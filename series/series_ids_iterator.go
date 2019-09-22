package series

import "github.com/RoaringBitmap/roaring"

// IDsIterator represents series IDs iterator
type IDsIterator struct {
	seriesIt roaring.ManyIntIterable
	buf      []uint32
}

// NewIDsIterator creates a series IDs iterator
func NewIDsIterator(seriesIDs *roaring.Bitmap, buf []uint32) *IDsIterator {
	return &IDsIterator{
		seriesIt: seriesIDs.ManyIterator(),
		buf:      buf,
	}
}

// Next passes the series IDs in a buffer to fill up with values, returns how many values were returned
func (it *IDsIterator) Next() (n int, buf []uint32) {
	n = it.seriesIt.NextMany(it.buf)
	buf = it.buf
	return
}
