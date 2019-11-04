package metadb

import (
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=metadb

// IDGenerator generates unique ID numbers for metric, tag and field.
type IDGenerator interface {
	// GenMetricID generates ID(uint32) from metricName
	GenMetricID(metricName string) uint32
	// GenTagKeyID generates ID(uint32) from metricID + tagKey
	GenTagKeyID(metricID uint32, tagKey string) uint32
	// GenFieldID generates ID(uint32) from metricID and fieldName
	GenFieldID(metricID uint32, fieldName string, fieldType field.Type) (uint16, error)
}

// IDGetter represents the query ability for metric level, such as metric id, field meta etc.
type IDGetter interface {
	// GetMetricID returns metric ID(uint32), if not exist return ErrNotFound error
	GetMetricID(metricName string) (uint32, error)
	// GetTagKeyID returns tag ID(uint32), return ErrNotFound if not exist
	GetTagKeyID(metricID uint32, tagKey string) (tagKeyID uint32, err error)
	// GetFieldID returns field id and type by given metricID and field name,
	// if not exist return ErrNotFound error
	GetFieldID(metricID uint32, fieldName string) (fieldID uint16, fieldType field.Type, err error)
}

// IDSequencer contains the abilities for querying and generating ID numbers.
// It is namespace level, and is used by all shards belong it.
type IDSequencer interface {
	// Recover loads metric-names and metricIDs from the index file to build the tree
	Recover() error
	IDGenerator
	IDGetter
	series.MetricMetaSuggester
	// FlushNameIDs flushes metricName and metricID to family
	FlushNameIDs() error
	// FlushMetricsMeta flushes tagKey, tagKeyId, fieldName, fieldID to family
	FlushMetricsMeta() error
}
