package series

import (
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=series

// MetadataGetter represents the query ability for metric level metadata
type MetadataGetter interface {
	// GetTagValues returns tag values by tag keys and spec version for metric level
	GetTagValues(metricID uint32, tagKeys []string, version int64) (tagValues [][]string, err error)
}

// Filter represents the query ability for filtering seriesIDs by expr from an index of tags.
// to support multi-version based on timestamp, time range for filtering spec version is necessary
type Filter interface {
	// FindSeriesIDsByExpr finds series ids by tag filter expr for metric id
	FindSeriesIDsByExpr(metricID uint32, expr stmt.TagFilter,
		timeRange timeutil.TimeRange) (*MultiVerSeriesIDSet, error)
	// GetSeriesIDsForTag get series ids for spec metric's tag key
	GetSeriesIDsForTag(metricID uint32, tagKey string,
		timeRange timeutil.TimeRange) (*MultiVerSeriesIDSet, error)
}
