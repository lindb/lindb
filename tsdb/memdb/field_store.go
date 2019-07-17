package memdb

import (
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/lockers"
	"github.com/eleme/lindb/pkg/logger"
)

// fieldStore holds the relation of segmentTime and segmentStore.
type fieldStore struct {
	fieldType field.Type

	segments map[int64]segmentStore
	lockers.SpinLock
}

// newFieldStore returns a new fieldStore.
func newFieldStore(fieldType field.Type) *fieldStore {
	return &fieldStore{
		fieldType: fieldType,
		segments:  make(map[int64]segmentStore),
	}
}

// getFieldType returns field type for current field store
func (fs *fieldStore) getFieldType() field.Type {
	return fs.fieldType
}

// getSegmentStore returns a segmentStore, if segment store not exist reutrn nil
func (fs *fieldStore) getSegmentStore(familyStartTime int64) segmentStore {
	fs.Lock()
	store := fs.segments[familyStartTime]
	fs.Unlock()
	return store
}

func (fs *fieldStore) write(blockStore *blockStore, familyStartTime int64, slot int, f models.Field) {
	fs.Lock()
	if !f.IsComplex() {
		sf, ok := f.(models.SimpleField)
		if !ok {
			logger.GetLogger("mem/field/store").Warn("convert field to simple field error")
			return
		}
		store, exist := fs.segments[familyStartTime]
		if !exist {
			//TODO ???
			store = newSimpleFieldStore(familyStartTime, field.GetAggFunc(sf.AggType()))
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

	fs.Unlock()
}
