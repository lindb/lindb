package constants

import "math"

const (
	// use this limit of metric-level when maxSeriesIDsLimit is not set
	DefaultMaxSeriesIDsCount = 10000000
	//max tag keys limitation of a metric
	DefaultMaxTagKeysCount = 32
	// max fields limitation of a tsStore.
	DefaultMaxFieldsCount = math.MaxUint8
	// the max number of suggestions count
	MaxSuggestions = 10000

	// Check if the global memory usage is greater than the limit,
	// If so, engine will flush the biggest shard's memdb until we are down to the lower mark.
	MemoryHighWaterMark = 80
	MemoryLowWaterMark  = 60
	// Check if shard's memory usage is greater than this limit,
	// If so, engine will flush this shard to disk
	ShardMemoryUsedThreshold = 500 * 1024 * 1024
	// FlushConcurrency controls the concurrent number of flushers
	FlushConcurrency = 4

	// TagValueIDForTag represents tag value id placeholder for store all series ids under tag
	TagValueIDForTag = uint32(0)
	// DefaultNamespace represents default namespace if not set
	DefaultNamespace = "default-ns"
)
