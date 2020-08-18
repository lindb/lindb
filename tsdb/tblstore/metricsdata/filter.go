package metricsdata

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./filter.go -destination=./filter_mock.go -package metricsdata

// Filter implements filtering metrics from sst files.
type Filter interface {
	// Filter filters data under each sst file based on query condition
	Filter(fieldIDs []field.ID, seriesIDs *roaring.Bitmap) ([]flow.FilterResultSet, error)
}

// metricsDataFilter represents the sst file data filter
type metricsDataFilter struct {
	familyTime int64
	snapshot   version.Snapshot //FIXME stone1100, need close version snapshot
	readers    []Reader
}

// NewFilter creates the sst file data filter
func NewFilter(familyTime int64, snapshot version.Snapshot, readers []Reader) Filter {
	return &metricsDataFilter{
		familyTime: familyTime,
		snapshot:   snapshot,
		readers:    readers,
	}
}

// Filter filters the data under each sst file based on metric/version/seriesIDs,
// if finds data then returns the FilterResultSet, else returns nil
func (f *metricsDataFilter) Filter(fieldIDs []field.ID,
	seriesIDs *roaring.Bitmap,
) (rs []flow.FilterResultSet, err error) {
	for _, reader := range f.readers {
		//FIXME add time range compare????
		fieldMetas, _ := reader.GetFields().Intersects(fieldIDs)
		if len(fieldMetas) == 0 {
			// field not found
			continue
		}
		// after and operator, query bitmap is sub of store bitmap
		matchSeriesIDs := roaring.FastAnd(seriesIDs, reader.GetSeriesIDs())
		if matchSeriesIDs.IsEmpty() {
			// series ids not found
			continue
		}
		rs = append(rs, newFileFilterResultSet(f.familyTime, fieldMetas, matchSeriesIDs, reader))
	}
	// not founds
	if len(rs) == 0 {
		return nil, constants.ErrNotFound
	}
	return
}

// fileFilterResultSet represents sst file reader for loading file data based on query condition
type fileFilterResultSet struct {
	reader     Reader
	familyTime int64
	fieldMetas field.Metas
	seriesIDs  *roaring.Bitmap
}

// newFileFilterResultSet creates the file filter result set
func newFileFilterResultSet(familyTime int64, fieldMetas field.Metas,
	seriesIDs *roaring.Bitmap, reader Reader,
) flow.FilterResultSet {
	return &fileFilterResultSet{
		familyTime: familyTime,
		reader:     reader,
		fieldMetas: fieldMetas,
		seriesIDs:  seriesIDs,
	}
}

// Identifier identifies the source of result set from kv store
func (f *fileFilterResultSet) Identifier() string {
	return f.reader.Path()
}

// SeriesIDs returns the series ids which matches with query series ids
func (f *fileFilterResultSet) SeriesIDs() *roaring.Bitmap {
	return f.seriesIDs
}

// Load reads data from sst files, then returns the data file scanner.
func (f *fileFilterResultSet) Load(flow flow.StorageQueryFlow, fieldIDs []field.ID,
	highKey uint16, seriesID roaring.Container,
) flow.Scanner {
	return f.reader.Load(flow, f.familyTime, fieldIDs, highKey, seriesID)
}
