package memdb

import (
	"github.com/lindb/lindb/tsdb/field"
	"github.com/lindb/lindb/tsdb/series"
)

//////////////////////////////////////////////////////
// primitiveIterator implements PrimitiveIterator
//////////////////////////////////////////////////////
type primitiveIterator struct {
	fStore fStoreINTF
	sStore sStoreINTF
	sCtx   series.ScanContext
	data   []byte
}

func newPrimitiveIterator(sCtx series.ScanContext) *primitiveIterator {
	return &primitiveIterator{sCtx: sCtx}
}

func (pi *primitiveIterator) reset(fStore fStoreINTF) {
	pi.fStore = fStore
	pi.sStore = nil
	pi.data = pi.data[:0]
}

func (pi *primitiveIterator) FieldID() uint16 {
	if pi.fStore == nil {
		return 0
	}
	return pi.fStore.GetFieldID()
}

func (pi *primitiveIterator) HasNext() bool {
	if pi.fStore == nil {
		return false
	}
	if pi.sStore == nil {
		var ok bool
		if pi.sStore, ok = pi.fStore.GetSStore(pi.sCtx.FamilyTime); !ok {
			return false
		}
	}
	if len(pi.data) == 0 {
		data, startSlot, endSlot, err := pi.sStore.bytes()
		if err != nil {
			return false
		}
		// todo: fixme
		_ = data
		_ = startSlot
		_ = endSlot
	}

	return false
}

func (pi *primitiveIterator) Next() (timeSlot int, value float64) {
	return 0, 0
}

//////////////////////////////////////////////////////
// fStoreIterator implements FieldIterator
//////////////////////////////////////////////////////
type fStoreIterator struct {
	tStore       tStoreINTF
	fieldName    string
	fieldID      uint16
	fieldType    field.Type
	sCtx         series.ScanContext
	fieldIDs     []uint16 // on iterating
	fieldMetas   fieldsMetas
	primitiveItr *primitiveIterator
}

func newFStoreIterator(metas fieldsMetas, sCtx series.ScanContext) *fStoreIterator {
	itr := &fStoreIterator{
		sCtx:         sCtx,
		fieldMetas:   metas,
		primitiveItr: newPrimitiveIterator(sCtx)}
	return itr
}

func (fsi *fStoreIterator) reset(tStore tStoreINTF) {
	fsi.tStore = tStore
	fsi.fieldIDs = append(fsi.fieldIDs[:0], fsi.sCtx.FieldIDs...)
}
func (fsi *fStoreIterator) FieldID() uint16       { return fsi.fieldID }
func (fsi *fStoreIterator) FieldName() string     { return fsi.fieldName }
func (fsi *fStoreIterator) FieldType() field.Type { return fsi.fieldType }

func (fsi *fStoreIterator) findID(fieldID uint16) bool {
	for _, meta := range fsi.fieldMetas {
		if fieldID == meta.fieldID {
			fsi.fieldID = meta.fieldID
			fsi.fieldName = meta.fieldName
			fsi.fieldType = meta.fieldType
			return true
		}
	}
	return false
}

func (fsi *fStoreIterator) HasNext() bool {
Loop:
	{
		if len(fsi.fieldIDs) == 0 {
			fsi.fieldID = 0
			fsi.fieldName = ""
			fsi.fieldType = 0
			return false
		}
		thisFieldID := fsi.fieldIDs[0]
		fsi.fieldIDs = fsi.fieldIDs[1:]
		if !fsi.findID(thisFieldID) {
			goto Loop
		}
		fStore, ok := fsi.tStore.getFStore(thisFieldID)
		if !ok {
			goto Loop
		}
		timeRange, ok := fStore.TimeRange(fsi.sCtx.TimeInterval)
		if !ok {
			goto Loop
		}
		if !timeRange.Overlap(&fsi.sCtx.TimeRange) {
			goto Loop
		}
		fsi.primitiveItr.reset(fStore)
		return true
	}
}

func (fsi *fStoreIterator) Next() series.PrimitiveIterator {
	return fsi.primitiveItr
}
