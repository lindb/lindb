package memdb

import (
	"time"

	"go.uber.org/atomic"
)

// use var for mocking
var (
	// series will be purged if have not been used in this TTL
	seriesTTL = atomic.NewDuration(5 * time.Minute)
)
