package storage

import (
	"sync/atomic"

	meta "github.com/eleme/lindb/storage/version"
)

type Snapshot struct {
	version *meta.Version
	closed  *int32
}

// create snapshot for reading data
func newSnapshot(version *meta.Version) *Snapshot {
	var closed int32
	return &Snapshot{
		version: version,
		closed:  &closed,
	}
}

// release related resources
func (s *Snapshot) Close() {
	if atomic.CompareAndSwapInt32(s.closed, 0, 1) {
		s.version.Release()
	}
}
