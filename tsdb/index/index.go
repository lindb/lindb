package index

import (
	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/timeutil"
	"github.com/eleme/lindb/sql/stmt"
	"github.com/eleme/lindb/tsdb/series"
)

//go:generate mockgen -source ./index.go -destination=./index_mock.go -package=index

// IDGenerator generates unique ID numbers for metric, tag and field.
type IDGenerator interface {
	// GenMetricID generates ID(uint32) from metricName
	GenMetricID(metricName string) uint32
	// GenFieldID generates ID(uint32) from metricID and fieldName
	GenFieldID(metricID uint32, fieldName string, fieldType field.Type) (uint16, error)
}

// SeriesIDsFilter represents the query ability for filtering seriesIDs by expr from an index of tags.
// to support multi-version based on timestamp, time range for filtering spec version is necessary
type SeriesIDsFilter interface {
	// FindSeriesIDsByExpr finds series ids by tag filter expr for metric id
	FindSeriesIDsByExpr(metricID uint32, expr stmt.TagFilter,
		timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error)
	// GetSeriesIDsForTag get series ids for spec metric's tag key
	GetSeriesIDsForTag(metricID uint32, tagKey string,
		timeRange timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error)
}
