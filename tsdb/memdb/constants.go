package memdb

import "sync/atomic"

const (
	// buckets count for sharding metric-stores, 32
	shardingCountOfMStores = 2 << 4
	// mask for calculating sharding-index by AND
	shardingCountMask = shardingCountOfMStores - 1
	// use this limit of metric-store when maxTagsLimit is not set
	defaultMaxTagsLimit = 100000
	// max fields limitation of a tsStore.
	maxFieldsLimit = 1024
	// unit: millisecond, used to prevent resetting metric-store too frequently.
	minIntervalForResetMetricStore = 10 * 1000
)

// use var for mocking
var (
	// store will be purged if have not been used in this TTL, unit: milliseconds
	tagsIDTTL int64 = 300 * 1000
)

// getTagsIDTTL returns the tagsIDTTL
func getTagsIDTTL() int64 {
	return atomic.LoadInt64(&tagsIDTTL)
}

// setTagsIDTTL sets the tagsIDTTL
func setTagsIDTTL(ttl int64) {
	atomic.StoreInt64(&tagsIDTTL, ttl)
}
