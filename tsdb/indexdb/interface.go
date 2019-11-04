package indexdb

import (
	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=indexdb

// IndexDatabase represents a database of index files, it is shard-level
// it provides the abilities to filter seriesID from the index.
// See `tsdb/doc` for index file layout.
type IndexDatabase interface {
	series.MetaGetter
	series.Filter
	series.TagValueSuggester
}
