package memdb

import (
	"sync/atomic"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/lockers"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/tsdb/index"
	"github.com/eleme/lindb/tsdb/metrictbl"
)

// fieldStore holds the relation of familyStartTime and segmentStore.
type fieldStore struct {
	fieldName *string                // field-name, nil after fieldID is generated
	fieldType field.Type             // sum, gauge, min, max
	fieldID   uint32                 // default 0
	segments  map[int64]segmentStore // familyTime
	sl        lockers.SpinLock       // spin-lock
}

// newFieldStore returns a new fieldStore.
func newFieldStore(fieldName string, fieldType field.Type) *fieldStore {
	return &fieldStore{
		fieldName: &fieldName,
		fieldType: fieldType,
		segments:  make(map[int64]segmentStore),
	}
}

// mustGetFieldID returns fieldID, if unset, generate a new one.
func (fs *fieldStore) mustGetFieldID(metricID uint32, generator index.IDGenerator) uint32 {
	fieldID := atomic.LoadUint32(&fs.fieldID)
	if fieldID > 0 {
		return fieldID
	}
	theFieldName := fs.fieldName
	if theFieldName != nil {
		atomic.CompareAndSwapUint32(&fs.fieldID, 0, generator.GenFieldID(metricID, *theFieldName, fs.fieldType))
		fs.fieldName = nil
	}
	return atomic.LoadUint32(&fs.fieldID)
}

// getFieldType returns field type for current field store
func (fs *fieldStore) getFieldType() field.Type {
	return fs.fieldType
}

// getSegmentsCount returns count of families.
func (fs *fieldStore) getFamiliesCount() int {
	fs.sl.Lock()
	length := len(fs.segments)
	fs.sl.Unlock()
	return length
}

// getSegmentStore returns a segmentStore, if segment store not exist returns nil
func (fs *fieldStore) getSegmentStore(familyStartTime int64) (segmentStore, bool) {
	fs.sl.Lock()
	store, ok := fs.segments[familyStartTime]
	fs.sl.Unlock()
	return store, ok
}

func (fs *fieldStore) write(blockStore *blockStore, familyStartTime int64, slot int, f models.Field) {
	fs.sl.Lock()
	if !f.IsComplex() {
		sf, ok := f.(models.SimpleField)
		if !ok {
			memDBLogger.Warn("convert field to simple field error")
			return
		}
		store, exist := fs.segments[familyStartTime]
		if !exist {
			//TODO ???
			store = newSimpleFieldStore(field.GetAggFunc(sf.AggType()))
			fs.segments[familyStartTime] = store
		}
		simpleStore, ok := store.(*simpleFieldStore)
		if ok {
			val := sf.Value()
			switch value := val.(type) {
			case int64:
				simpleStore.writeInt(blockStore, slot, value)
			case float64:
				//TODO handle float value
			}
		}
	}

	fs.sl.Unlock()
}

// flushFieldTo flushes segments' data to writer and reset the segments-map.
func (fs *fieldStore) flushFieldTo(writer metrictbl.TableWriter, metricID uint32,
	familyTime int64, generator index.IDGenerator) {
	fieldID := fs.mustGetFieldID(metricID, generator)

	fs.sl.Lock()
	defer fs.sl.Unlock()

	ss, ok := fs.segments[familyTime]
	if !ok {
		return
	}
	delete(fs.segments, familyTime)

	data, startSlot, endSlot, err := ss.bytes()

	if err != nil {
		memDBLogger.Error("read segment data error:", logger.Error(err))
		return
	}
	writer.WriteField(fieldID, data, startSlot, endSlot)
}
