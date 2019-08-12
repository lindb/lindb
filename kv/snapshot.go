package kv

import (
	"sync/atomic"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
)

//go:generate mockgen -source ./snapshot.go -destination=./snapshot_mock.go -package kv

// Snapshot represents a current family version by given key, for reading data.
type Snapshot interface {
	// Readers returns store reader that match query condition
	Readers() []table.Reader
	// Close releases related resources
	Close()
}

// snapshot implements Snapshot interface
type snapshot struct {
	readers []table.Reader

	version *version.Version
	closed  int32
}

// newSnapshot new snapshot instance
func newSnapshot(version *version.Version, readers []table.Reader) Snapshot {
	return &snapshot{
		version: version,
		readers: readers,
	}
}

// Readers returns store reader that match query condition
func (s *snapshot) Readers() []table.Reader {
	return s.readers
}

// Close releases related resources
func (s *snapshot) Close() {
	// atomic set closed status, make sure only release once
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		s.version.Release()
	}
}
