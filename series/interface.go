package series

import (
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=series

// MetricMetaSuggester represents the suggest ability for metricNames and tagKeys.
// default max limit of suggestions is set in constants
type MetricMetaSuggester interface {
	// SuggestMetrics returns suggestions from a given prefix of metricName
	SuggestMetrics(metricPrefix string, limit int) []string
	// SuggestTagKeys returns suggestions from given metricName and prefix of tagKey
	SuggestTagKeys(metricName, tagKeyPrefix string, limit int) []string
}

// TagValueSuggester represents the suggest ability for tagValues.
// default max limit of suggestions is set in constants
type TagValueSuggester interface {
	// SuggestTagValues returns suggestions from given tag key id and prefix of tagValue
	SuggestTagValues(tagKeyID uint32, tagValuePrefix string, limit int) []string
}

// Filter represents the query ability for filtering seriesIDs by expr from an index of tags.
// to support multi-version based on timestamp, time range for filtering spec version is necessary
type Filter interface {
	// FindSeriesIDsByExpr finds series ids by tag filter expr for tag key id
	FindSeriesIDsByExpr(tagKeyID uint32, expr stmt.TagFilter, timeRange timeutil.TimeRange) (
		*MultiVerSeriesIDSet, error)
	// GetSeriesIDsForTag get series ids for spec metric's tag key
	GetSeriesIDsForTag(tagKeyID uint32, timeRange timeutil.TimeRange) (
		*MultiVerSeriesIDSet, error)
	//FIXME stone1100
	//// GetSeriesIDsForMetric get series ids for spec metric's id
	//GetSeriesIDsForMetric(tagKeyID uint32, timeRange timeutil.TimeRange) (
	//	*MultiVerSeriesIDSet, error)
	// GetGroupingContext returns the context of group by
	GetGroupingContext(tagKeyIDs []uint32, version Version) (GroupingContext, error)
}
