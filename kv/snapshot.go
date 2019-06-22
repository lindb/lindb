package kv

import (
	"sync/atomic"

	"github.com/eleme/lindb/kv/table"
	"github.com/eleme/lindb/kv/version"
)

// Snapshot current family version, for reading data.
type Snapshot struct {
	readers []table.Reader

	version *version.Version
	closed  int32
}

// newSnapshot new snapshot instance
func newSnapshot(version *version.Version, readers []table.Reader) *Snapshot {
	return &Snapshot{
		version: version,
		readers: readers,
	}
}

// Readers returns store reader that match query condition
func (s *Snapshot) Readers() []table.Reader {
	return s.readers
}

// Close releases related resources
func (s *Snapshot) Close() {
	// atomic set closed status, make sure only release once
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		s.version.Release()
	}
}
