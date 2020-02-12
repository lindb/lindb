package memdb

import (
	"sort"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./timeseries_store.go -destination=./timeseries_store_mock.go -package memdb

const emptyTimeSeriesStoreSize = 24 // fStores

// tStoreINTF abstracts a time-series store
type tStoreINTF interface {
	// GetFStore returns the fStore in field list by family/field/primitive
	GetFStore(familyID familyID, fieldID field.ID, pField field.PrimitiveID) (fStoreINTF, bool)
	// InsertFStore inserts a new fStore to field list.
	InsertFStore(fStore fStoreINTF)
	// FlushSeriesTo flushes the series data segment.
	FlushSeriesTo(flusher metricsdata.Flusher, flushCtx flushContext)
	// scan scans the time series data based on field ids
	scan(memScanCtx *memScanContext)
}

// fStoreNodes implements sort.Interface
type fStoreNodes []fStoreINTF

func (f fStoreNodes) Len() int           { return len(f) }
func (f fStoreNodes) Less(i, j int) bool { return f[i].GetKey() < f[j].GetKey() }
func (f fStoreNodes) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

// timeSeriesStore holds a mapping relation of field and fieldStore.
type timeSeriesStore struct {
	fStoreNodes fStoreNodes // key: sorted fStore list by field-name, insert-only
}

// newTimeSeriesStore returns a new tStoreINTF.
func newTimeSeriesStore() tStoreINTF {
	return &timeSeriesStore{}
}

// GetFStore returns the fStore in this list from field-id.
func (ts *timeSeriesStore) GetFStore(familyID familyID, fieldID field.ID, pField field.PrimitiveID) (fStoreINTF, bool) {
	fieldKey := uint32(pField) | uint32(fieldID)<<8 | uint32(familyID)<<16
	fieldLength := len(ts.fStoreNodes)
	if fieldLength == 1 {
		if ts.fStoreNodes[0].GetKey() != fieldKey {
			return nil, false
		}
		return ts.fStoreNodes[0], true

	}
	idx := sort.Search(fieldLength, func(i int) bool {
		return ts.fStoreNodes[i].GetKey() >= fieldKey
	})
	if idx >= fieldLength || ts.fStoreNodes[idx].GetKey() != fieldKey {
		return nil, false
	}
	return ts.fStoreNodes[idx], true
}

// InsertFStore inserts a new fStore to field list.
func (ts *timeSeriesStore) InsertFStore(fStore fStoreINTF) {
	if ts.fStoreNodes == nil {
		ts.fStoreNodes = []fStoreINTF{fStore}
		return
	}
	ts.fStoreNodes = append(ts.fStoreNodes, fStore)
	sort.Sort(ts.fStoreNodes)
}

// FlushSeriesTo flushes the series data segment.
func (ts *timeSeriesStore) FlushSeriesTo(
	flusher metricsdata.Flusher,
	flushCtx flushContext,
) {
	// FIXME stone100
	//for _, fStore := range ts.fStoreNodes {
	//	fStore.FlushFieldTo(flusher, flushCtx)
	//}
	//flusher.FlushSeries()
}
