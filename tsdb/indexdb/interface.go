package indexdb

import (
	"io"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
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
	series.TagValueSuggester
	series.Filter
	// GetOrCreateSeriesID gets series by tags hash, if not exist generate new series id in memory,
	// if generate a new series id returns isCreate is true
	// if generate fail return err
	GetOrCreateSeriesID(metricID uint32, tagsHash uint64) (seriesID uint32, isCreated bool, err error)
	// BuildInvertIndex builds the inverted index for tag value => series ids,
	// the tags is considered as a empty key-value pair while tags is nil.
	BuildInvertIndex(namespace, metricName string, tags map[string]string, seriesID uint32)
	// Flush flushes index data to disk
	Flush() error
}
