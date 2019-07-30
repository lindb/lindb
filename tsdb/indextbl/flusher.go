package indextbl

import "github.com/RoaringBitmap/roaring"

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package indextbl

// SeriesTableFlusher is a wrapper of kv.Builder, provides the ability to build a versioned series-id index table.
// The layout is available in `tsdb/doc.go`
type SeriesTableFlusher interface {
	// FlushVersion writes a version of the tagValues and related bitmap to the index table.
	FlushVersion(version int64, entrySet map[string]*roaring.Bitmap)
	// FlushTagKey marks versions above belongs to this key and metric.
	FlushTagKey(tagKey string, metricID uint32)
	// Commit closes the writer, this will be called after writing all tagKeys and metrics.
	Commit() error
}
