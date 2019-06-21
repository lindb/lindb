package memdb

import (
	"github.com/eleme/lindb/pkg/lockers"
)

// fieldStore holds the relation of segmentTime and segmentStore.
type fieldStore struct {
	segments map[int64]*segmentStore
	lockers.SpinLock
}

// newFieldStore returns a new fieldStore.
func newFieldStore() *fieldStore {
	return &fieldStore{segments: make(map[int64]*segmentStore)}
}

// getSegmentStore returns a new segmentStore
func (fs *fieldStore) getSegmentStore(segmentTime int64) *segmentStore {
	fs.Lock()
	store, exist := fs.segments[segmentTime]
	if !exist {
		store = newSegmentStore(segmentTime)
		fs.segments[segmentTime] = store
	}
	fs.Unlock()
	return store
}
