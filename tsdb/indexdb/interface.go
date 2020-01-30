package indexdb

import (
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=indexdb

var indexLogger = logger.GetLogger("tsdb", "IndexDB")

// IndexDatabase represents a database of index files, it is shard-level
// it provides the abilities to filter seriesID from the index.
// See `tsdb/doc` for index file layout.
type IndexDatabase interface {
	series.Filter
	series.TagValueSuggester
}

type MemoryIndexDatabase interface {
	GetTimeSeriesID(metricName string, tags map[string]string, tagsHash uint64) (metricID, seriesID uint32)

	// FlushInvertedIndexTo flushes the series data to a inverted-index file.
	FlushInvertedIndexTo(flusher invertedindex.Flusher) (err error)
}
