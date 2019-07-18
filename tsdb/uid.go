package tsdb

import (
	"github.com/RoaringBitmap/roaring"

	"github.com/eleme/lindb/pkg/field"
)

//MetricUid represents metric name unique id under the database
type MetricUID interface {
	//GetOrCreateMetricID returns find the metric ID associated with a given name or create it.
	GetOrCreateMetricID(metricName string, create bool) (uint32, bool)
	//SuggestMetrics returns suggestions of metric names given a search prefix.
	SuggestMetrics(prefix string, limit int16) map[string]struct{}
	//Flush represents forces a flush of in-memory data, and clear it
	Flush() error
}

//TagsUID represents tags unique id under the metric name.
//Shard level sharing.
type TagsUID interface {
	//GetOrCreateTagsID returns find the tags ID associated with given tags or create it.
	GetOrCreateTagsID(metricID uint32, tags string) (uint32, error)
	//GetTagNames return get all tag names within the metric name
	GetTagNames(metricID uint32, limit int16) map[string]struct{}
	//GetTagValueBitmap returns find bitmap associated with a given tag value
	GetTagValueBitmap(metricID uint32, tagName string, tagValue string) *roaring.Bitmap
	//SuggestTagValues returns  suggestions of tag values given a search prefix
	SuggestTagValues(metricID uint32, tagName string, tagValuePrefix string, limit uint16) map[string]struct{}
	//Flush represents forces a flush of in-memory data, and clear it
	Flush() error
}

//FieldUID represents field unique under the metric name.
//Database level sharing.
type FieldUID interface {
	//GetOrCreateFieldID  returns find the ID associated with a given field name and field type or create it.
	GetOrCreateFieldID(metricID uint32, fieldName string, fieldType field.Type) (uint32, error)
	//GetFields returns get all fields within the metric name
	GetFields(metricID uint32, limit int16) map[string]struct{}
	//GetFieldID returns get fieldID by fieldName within the metric name
	GetFieldID(metricID uint32, fieldName string) uint32
	//Flush represents forces a flush of in-memory data, and clear it
	Flush() error
}
