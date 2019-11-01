package tsdb

import "github.com/lindb/lindb/pkg/concurrent"

type ExecutorPool struct {
	Scanners concurrent.Pool
	Mergers  concurrent.Pool
}
