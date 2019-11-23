package tsdb

import "github.com/lindb/lindb/pkg/concurrent"

// ExecutorPool represents the executor pool used by query flow for each storage engine
type ExecutorPool struct {
	Filtering concurrent.Pool
	Grouping  concurrent.Pool
	Scanner   concurrent.Pool
}
