package metadb

import (
	"io"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=metadb
var metaLogger = logger.GetLogger("tsdb", "MetaDB")

// IDGenerator generates unique ID numbers for metric, tag and field.
type IDGenerator interface {
	// GenMetricID generates ID(uint32) from metricName
	GenMetricID(metricName string) uint32
	//FIXME need check max tag key id????
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
	// GetTagKeyIDs returns tag IDs([]uint32), return ErrNotFound if not exist
	GetTagKeyIDs(metricID uint32) (tagKeyIDs []uint32, err error)
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

// Metadata represents all metadata of tsdb, like metric/tag metadata
type Metadata interface {
	// MetadataDatabase returns the metric level metadata
	MetadataDatabase() MetadataDatabase
	// TagMetadata returns the tag metadata
	TagMetadata() TagMetadata
}

// MetadataDatabase represents the metadata storage includes namespace/metric metadata
type MetadataDatabase interface {
	io.Closer

	// SuggestNamespace suggests the namespace by namespace's prefix
	SuggestNamespace(prefix string, limit int) (namespaces []string, err error)
	// SuggestMetricName suggests the metric name by name's prefix
	SuggestMetricName(namespace, prefix string, limit int) (namespaces []string, err error)

	// GetMetricID gets the metric id by namespace and metric name, if not exist return series.ErrNotFound
	GetMetricID(namespace, metricName string) (metricID uint32, err error)
	// GetTagKeyID gets the tag key id by namespace/metric name/tag key key, if not exist return series.ErrNotFound
	GetTagKeyID(namespace, metricName, tagKey string) (tagKeyID uint32, err error)
	// GetAllTagKeys returns the all tag keys by namespace/metric name, if not exist return series.ErrNotFound
	GetAllTagKeys(namespace, metricName string) (tags []tag.Meta, err error)
	// GetField gets the field meta by namespace/metric name/field name, if not exist return series.ErrNotFound
	GetField(namespace, metricName, fieldName string) (field field.Meta, err error)
	// GetAllFields returns the  all fields by namespace/metric name, if not exist return series.ErrNotFound
	GetAllFields(namespace, metricName string) (fields []field.Meta, err error)

	// GenMetricID generates the metric id in the memory
	GenMetricID(namespace, metricName string) (metricID uint32, err error)
	// GenFieldID generates the field id in the memory
	GenFieldID(namespace, metricName, fieldName string, fieldType field.Type) (uint16, error)
	// GenTagKeyID generates the tag key id in the memory
	GenTagKeyID(namespace, metricName, tagKey string) (uint32, error)

	// Sync syncs the pending metadata update event
	Sync() error
}
