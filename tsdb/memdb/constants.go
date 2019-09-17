package memdb

import (
	"time"

	"go.uber.org/atomic"
)

const (
	// buckets count for sharding metric-stores, 32
	shardingCountOfMStores = 2 << 4
	// mask for calculating sharding-index by AND
	shardingCountMask = shardingCountOfMStores - 1
)

// use var for mocking
var (
	// series will be purged if have not been used in this TTL
	seriesTTL = atomic.NewDuration(5 * time.Minute)
)
