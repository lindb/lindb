package memdb

import (
	"container/list"
	"time"

	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/lockers"
)

// timeSeriesStore holds a mapping relation of field and fieldStore.
type timeSeriesStore struct {
	//TODO: tsID
	// tsID             uint32 // tsId identifier
	fields         map[string]*fieldStore // key: fieldName
	tagsID         string                 // tags identifier
	lastAccessedAt int64                  // nanosecond
	element        *list.Element          // element in LRU-list
	sl             lockers.SpinLock       // spin-lock
}

// newTimeSeriesStore returns a new timeSeriesStore from tagsID.
func newTimeSeriesStore(tagsID string) *timeSeriesStore {
	return &timeSeriesStore{
		tagsID:         tagsID,
		lastAccessedAt: time.Now().UnixNano(),
		fields:         make(map[string]*fieldStore)}
}

// getFieldStore mustGet a fieldStore by fieldName.
func (ts *timeSeriesStore) getFieldStore(fieldName string) *fieldStore {
	ts.lastAccessedAt = time.Now().UnixNano()

	ts.sl.Lock()
	store, exist := ts.fields[fieldName]
	if !exist {
		store = newFieldStore(field.SumField)
		ts.fields[fieldName] = store
	}
	ts.sl.Unlock()
	return store
}

// shouldBeEvicted detects if thisStore has been accessed for tagsIDTTL.
func (ts *timeSeriesStore) shouldBeEvicted() bool {
	return ts.lastAccessedAt+(time.Duration(getTagsIDTTL())*time.Millisecond).Nanoseconds() < time.Now().UnixNano()
}

// getFieldsCount returns the count of fields thread-safely.
func (ts *timeSeriesStore) getFieldsCount() int {
	ts.sl.Lock()
	length := len(ts.fields)
	ts.sl.Unlock()
	return length
}
