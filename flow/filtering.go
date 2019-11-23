package flow

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source=./filtering.go -destination=./filtering_mock.go -package=flow

// DataFilter represents the filter ability over memory database and files under data family.
type DataFilter interface {
	// Filter filters the data based on metricIDs/fieldIDs/version/seriesIDs,
	// if finds data then returns filter result set, else returns nil.
	Filter(metricID uint32, fieldIDs []uint16, version series.Version, seriesIDs *roaring.Bitmap) []FilterResultSet
}

// FilterResultSet represents the filter result set, loads data and does down sampling need based on this interface.
type FilterResultSet interface {
	// Load loads the data from storage, then does down sampling, finally reduces the down sampling results.
	Load(flow StorageQueryFlow, fieldIDs []uint16, highKey uint16, groupedSeries map[string][]uint16)
}
