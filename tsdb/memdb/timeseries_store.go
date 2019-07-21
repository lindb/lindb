package memdb

import (
	"sync/atomic"
	"time"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/lockers"
	"github.com/eleme/lindb/tsdb/index"
	"github.com/eleme/lindb/tsdb/metrictbl"
)

// timeSeriesStore holds a mapping relation of field and fieldStore.
type timeSeriesStore struct {
	tsID           uint32                 // tsId identifier
	fields         map[string]*fieldStore // key: fieldName
	lastAccessedAt int64                  // nanoseconds
	sl             lockers.SpinLock       // spin-lock
}

// newTimeSeriesStore returns a new timeSeriesStore.
func newTimeSeriesStore() *timeSeriesStore {
	return &timeSeriesStore{
		lastAccessedAt: time.Now().UnixNano(),
		fields:         make(map[string]*fieldStore)}
}

// mustGetTSID returns tsID, if unset, generate a new one.
func (ts *timeSeriesStore) mustGetTSID(generator index.IDGenerator, metricID uint32,
	tags string, version int64) uint32 {

	tsID := atomic.LoadUint32(&ts.tsID)
	if tsID > 0 {
		return tsID
	}
	atomic.CompareAndSwapUint32(&ts.tsID, 0, generator.GenTSID(metricID, tags, version))
	return atomic.LoadUint32(&ts.tsID)
}

// generateFieldsID generates fieldID.
func (ts *timeSeriesStore) generateFieldsID(metricID uint32, generator index.IDGenerator) {
	ts.sl.Lock()
	for fieldName, fStore := range ts.fields {
		fStore.mustGetFieldID(generator, metricID, fieldName)
	}
	ts.sl.Unlock()
}

// getOrCreateFStore mustGet a fieldStore by fieldName.
func (ts *timeSeriesStore) getOrCreateFStore(fieldName string, fieldType field.Type) (*fieldStore, error) {
	atomic.StoreInt64(&ts.lastAccessedAt, time.Now().UnixNano())

	ts.sl.Lock()
	store, exist := ts.fields[fieldName]
	if exist {
		if store.getFieldType() != fieldType {
			ts.sl.Unlock()
			return nil, models.ErrWrongFieldType
		}
	} else {
		store = newFieldStore(fieldType)
		ts.fields[fieldName] = store
	}
	ts.sl.Unlock()
	return store, nil
}

// shouldBeEvicted detects if thisStore has not been accessed for tagsIDTTL.
func (ts *timeSeriesStore) shouldBeEvicted() bool {
	// validate ttl
	expired := ts.lastAccessedAt+getTagsIDTTL()*int64(time.Millisecond) < time.Now().UnixNano()
	if !expired {
		return false
	}
	ts.sl.Lock()
	// make sure that family data has been flushed
	for _, fStore := range ts.fields {
		if fStore.getFamiliesCount() != 0 {
			ts.sl.Unlock()
			return false
		}
	}
	ts.sl.Unlock()
	return true
}

// isFull detects if the fields has too many fields.
func (ts *timeSeriesStore) isFull() bool {
	return ts.getFieldsCount() >= maxFieldsLimit
}

// getFieldsCount returns the count of fields thread-safely.
func (ts *timeSeriesStore) getFieldsCount() int {
	ts.sl.Lock()
	length := len(ts.fields)
	ts.sl.Unlock()
	return length
}

// flushTSEntryTo flushes the tsEntry data segment.
func (ts *timeSeriesStore) flushTSEntryTo(writer metrictbl.TableWriter, familyTime int64,
	generator index.IDGenerator, metricID uint32, tags string, version int64) {

	tsID := ts.mustGetTSID(generator, metricID, tags, version)
	ts.sl.Lock()
	for fieldName, fStore := range ts.fields {
		fStore.flushFieldTo(writer, familyTime, generator, metricID, fieldName)
	}
	writer.WriteTSEntry(tsID)
	ts.sl.Unlock()
}
