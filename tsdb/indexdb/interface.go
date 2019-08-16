package indexdb

import (
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/tsdb/series"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=indexdb

// IDGenerator generates unique ID numbers for metric, tag and field.
type IDGenerator interface {
	// GenMetricID generates ID(uint32) from metricName
	GenMetricID(metricName string) uint32
	// GenTagID generates ID(uint32) from metricID + tagKey
	GenTagID(metricID uint32, tagKey string) uint32
	// GenFieldID generates ID(uint32) from metricID and fieldName
	GenFieldID(metricID uint32, fieldName string, fieldType field.Type) (uint16, error)
}

// IDGetter represents the query ability for metric level, such as metric id, field meta etc.
type IDGetter interface {
	// GetMetricID returns metric ID(uint32), if not exist return ErrNotFound error
	GetMetricID(metricName string) (uint32, error)
	// GetTagID returns tag ID(uint32), return ErrNotFound if not exist
	GetTagID(metricID uint32, tagKey string) (tagID uint32, err error)
	// GetFieldID returns field id and type by given metricID and field name,
	// if not exist return ErrNotFound error
	GetFieldID(metricID uint32, fieldName string) (fieldID uint16, fieldType field.Type, err error)
}

// IndexDatabase represents a database of index files,
// it provides the abilities of generate id and getting meta data from the index.
// See `tsdb/doc` for index file layout.
type IndexDatabase interface {
	IDGenerator
	IDGetter
	series.MetaGetter
	series.Filter
	// FlushNameIDsTo flushes metricName and metricID to flusher
	FlushNameIDsTo(flusher kv.Flusher) error
	// FlushMetricsMetaTo flushes tagKey, tagKeyId, fieldName, fieldID to flusher
	FlushMetricsMetaTo(flusher kv.Flusher) error
}
