package tsdb

import "github.com/lindb/lindb/pkg/concurrent"

type ExecutePool struct {
	Scan  concurrent.Pool
	Merge concurrent.Pool
}
