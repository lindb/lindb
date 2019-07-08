package memdb

import "sync/atomic"

const (
	// buckets count for sharding metric-stores, 32
	shardingCountOfMStores = 2 << 4
	// mask for calculating sharding-index by AND
	shardingCountMask = shardingCountOfMStores - 1
	// use this limit of metric-store when maxTagsLimit is not set
	defaultMaxTagsLimit = 10000
)

// use var for mocking
var (
	// evictor evicts the stores in this interval, unit: milliseconds
	evictInterval int64 = 1000
	// store will be purged if have not been used in this TTL, unit: milliseconds
	tagsIDTTL int64 = 300 * 1000
)

// getEvictInterval returns the evictInterval
func getEvictInterval() int64 {
	return atomic.LoadInt64(&evictInterval)
}

// setEvictInterval sets the evictInterval
func setEvictInterval(interval int64) {
	atomic.StoreInt64(&evictInterval, interval)
}

// getTagsIDTTL returns the tagsIDTTL
func getTagsIDTTL() int64 {
	return atomic.LoadInt64(&tagsIDTTL)
}

// setTagsIDTTL sets the tagsIDTTL
func setTagsIDTTL(ttl int64) {
	atomic.StoreInt64(&tagsIDTTL, ttl)
}
