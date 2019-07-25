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
	// GenTSID generates ID(uint32) from metricID, sortedTags and version.
	GenTSID(metricID uint32, sortedTags string, version int64) uint32
	// GenFieldID generates ID(uint32) from metricID and fieldName
	GenFieldID(metricID uint32, fieldName string, fieldType field.Type) uint32
}

// Index represents an index of tags for searching series by expr
// support multi-version based on timestamp, so need time range for filter spec version
type Index interface {
	// FindSeriesIDsByExpr finds series ids by tag filter expr for metric id
	FindSeriesIDsByExpr(metricID uint32, expr stmt.TagFilter, timeRange *timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error)
	// GetSeriesIDsForTag get series ids for spec metric's tag key
	GetSeriesIDsForTag(metricID uint32, tagKey string, timeRange *timeutil.TimeRange) (*series.MultiVerSeriesIDSet, error)
}
