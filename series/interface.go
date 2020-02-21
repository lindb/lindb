package series

import (
	"github.com/lindb/roaring"
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
type Filter interface {
	// GetSeriesIDsByTagValueIDs gets series ids by tag value ids for spec metric's tag key
	GetSeriesIDsByTagValueIDs(tagKeyID uint32, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error)
	// GetSeriesIDsForTag gets series ids for spec metric's tag key
	GetSeriesIDsForTag(tagKeyID uint32) (*roaring.Bitmap, error)
	// GetSeriesIDsForMetric gets series ids for spec metric name
	GetSeriesIDsForMetric(namespace, metricName string) (*roaring.Bitmap, error)
	// GetGroupingContext returns the context of group by
	GetGroupingContext(tagKeyIDs []uint32, seriesIDs *roaring.Bitmap) (GroupingContext, error)
}
