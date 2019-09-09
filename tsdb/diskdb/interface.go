package diskdb

import (
	"github.com/lindb/lindb/tsdb/field"
	"github.com/lindb/lindb/tsdb/series"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=diskdb

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
type IDSequencer interface {
	// Recover loads metric-names and metricIDs from the index file to build the tree
	Recover() error
	IDGenerator
	IDGetter
	// SuggestMetrics returns suggestions from a given prefix of metricName
	SuggestMetrics(metricPrefix string, limit int) []string
	// SuggestTagKeys returns suggestions from given metricName and prefix of tagKey
	SuggestTagKeys(metricName, tagKeyPrefix string, limit int) []string
	// FlushNameIDs flushes metricName and metricID to family
	FlushNameIDs() error
	// FlushMetricsMeta flushes tagKey, tagKeyId, fieldName, fieldID to family
	FlushMetricsMeta() error
}

// IndexDatabase represents a database of index files, it is shard-level
// it provides the abilities to filter seriesID from the index.
// See `tsdb/doc` for index file layout.
type IndexDatabase interface {
	series.MetaGetter
	series.Filter
	// SuggestTagValues returns suggestions from given metricName, tagKey and prefix of tagValue
	SuggestTagValues(metricName, tagKey, tagValuePrefix string, limit int) []string
}
