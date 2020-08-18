package flow

import (
	"io"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./filtering.go -destination=./filtering_mock.go -package=flow

// DataFilter represents the filter ability over memory database and files under data family.
type DataFilter interface {
	// Filter filters the data based on metricIDs/fieldIDs/seriesIDs/timeRange,
	// if finds data then returns filter result set, else returns nil.
	Filter(metricID uint32, fieldIDs []field.ID,
		seriesIDs *roaring.Bitmap, timeRange timeutil.TimeRange,
	) ([]FilterResultSet, error)
}

// FilterResultSet represents the filter result set, loads data and does down sampling need based on this interface.
type FilterResultSet interface {
	// Identifier identifies the source of result set(mem/kv etc.)
	Identifier() string
	// Load loads the data from storage, then returns the data scanner.
	Load(flow StorageQueryFlow, fieldIDs []field.ID, highKey uint16, seriesID roaring.Container) Scanner
	// SeriesIDs returns the series ids which matches with query series ids
	SeriesIDs() *roaring.Bitmap
}

// Scanner represents the scanner which scan metric data from storage.
type Scanner interface {
	io.Closer
	// Scan scans the metric data by given series id.
	Scan(lowSeriesID uint16)
}
