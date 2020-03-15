package constants

import "math"

const (
	// DefaultMaxSeriesIDsCount represents series count limit, uses this limit of metric-level when maxSeriesIDsLimit is not set
	DefaultMaxSeriesIDsCount = 10000000
	// DefaultMaxTagKeysCount represents tag key count limit, uses this limit of max tag keys of a metric
	DefaultMaxTagKeysCount = 32
	// DefaultMaxFieldsCount represents field count limit, uses this limit of max fields of a metric
	DefaultMaxFieldsCount = math.MaxUint8
	// MaxSuggestions represents the max number of suggestions count
	MaxSuggestions = 10000

	// MemoryHighWaterMark checkes if the global memory usage is greater than the limit,
	// If so, engine will flush the biggest shard's memdb until we are down to the lower mark.
	MemoryHighWaterMark = 80
	// MemoryLowWaterMark checks if the global memory usage is low water mark
	MemoryLowWaterMark = 60
	// ShardMemoryUsedThreshold checks if shard's memory usage is greater than this limit,
	// If so, engine will flush this shard to disk
	ShardMemoryUsedThreshold = 500 * 1024 * 1024
	// FlushConcurrency controls the concurrent number of flushers
	FlushConcurrency = 4

	// TagValueIDForTag represents tag value id placeholder for store all series ids under tag
	TagValueIDForTag = uint32(0)
	// DefaultNamespace represents default namespace if not set
	DefaultNamespace = "default-ns"
	// SeriesIDWithoutTags represents the series ids under spec metric, but without nothing tags
	SeriesIDWithoutTags = uint32(0)
)
