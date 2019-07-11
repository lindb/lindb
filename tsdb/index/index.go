package index

import "github.com/eleme/lindb/pkg/field"

//go:generate mockgen -source ./index.go -destination=./index_mock.go -package index

// IDGenerator generates unique ID numbers for metric, tag and field.
type IDGenerator interface {
	// GenMetricID generates ID(uint32) from metricName
	GenMetricID(metricName string) uint32
	// GenTSID generates ID(uint32) from metricID and sortedTags
	GenTSID(metricID uint32, sortedTags string) uint32
	// GenFieldID generates ID(uint32) from metricID and fieldName
	GenFieldID(metricID uint32, fieldName string, fieldType field.Type) uint32
}

type Index interface {
}
