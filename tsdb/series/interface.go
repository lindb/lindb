package series

import (
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=series

// MetaGetter represents the query ability for metric level metadata
type MetaGetter interface {
	// GetTagValues returns tag values by tag keys and spec version for metric level
	GetTagValues(metricID uint32, tagKeys []string, version uint32) (tagValues [][]string, err error)
}

// Suggester represents the suggest ability for metricNames, tagKeys and tagValues.
// default max limit of suggestions is set in constants
type Suggester interface {
	// SuggestMetrics returns suggestions from a given prefix of metricName
	SuggestMetrics(metricPrefix string, limit int) []string
	// SuggestTagKeys returns suggestions from given metricName and prefix of tagKey
	SuggestTagKeys(metricName, tagKeyPrefix string, limit int) []string
	// SuggestTagValues returns suggestions from given metricName, tagKey and prefix of tagValue
	SuggestTagValues(metricName, tagKey, tagValuePrefix string, limit int) []string
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

// DataGetter represents the query ability for querying data of the seriesIDs
type DataGetter interface{}
