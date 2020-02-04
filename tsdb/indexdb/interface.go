package indexdb

import (
	"io"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=indexdb

var indexLogger = logger.GetLogger("tsdb", "IndexDB")

// FileIndexDatabase represents a database of index files, it is shard-level
// it provides the abilities to filter seriesID from the index.
// See `tsdb/doc` for index file layout.
type FileIndexDatabase interface {
	series.Filter
	series.TagValueSuggester
}

// IndexDatabase represents a index database includes memory/file storage, it is shard level.
// index database will generate series id if tags hash not exist in mapping storage, and
// builds inverted index for tags => series id
type IndexDatabase interface {
	io.Closer
	series.Filter
	series.TagValueSuggester
	// GetOrCreateSeriesID gets series by tags hash, if not exist generate new series id in memory, then
	// builds inverted index for tags => series id, if generate fail return err
	GetOrCreateSeriesID(metricID uint32, tags map[string]string, tagsHash uint64) (seriesID uint32, err error)

	// FlushInvertedIndexTo flushes the series data to a inverted-index file.
	FlushInvertedIndexTo(flusher invertedindex.TagFlusher) (err error)
}
