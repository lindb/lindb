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
	fieldKey := buildFieldKey(familyID, fieldID, pField)
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
func (ts *timeSeriesStore) FlushSeriesTo(flusher metricsdata.Flusher, flushCtx flushContext) {
	for _, fStore := range ts.fStoreNodes {
		//FIXME need flush by family id
		fStore.FlushFieldTo(flusher, flushCtx)
	}
}

// scan scans the time series data based on key(family+field+primitive).
// NOTICE: field ids and fields aggregator must be in order.
func (ts *timeSeriesStore) scan(memScanCtx *memScanContext) {
	fieldLength := len(ts.fStoreNodes)
	fieldAggs := memScanCtx.fieldAggs
	// find small/equals family id index
	idx := sort.Search(fieldLength, func(i int) bool {
		return ts.fStoreNodes[i].GetFamilyID() >= fieldAggs[0].familyID
	})
	fieldCount := len(fieldAggs)
	j := 0
	for i := idx; i < fieldLength; i++ {
		fieldStore := ts.fStoreNodes[i]
		agg := fieldAggs[j]
		key := fieldStore.GetKey()
		switch {
		case key == agg.fieldKey:
			fieldStore.Load(agg.fieldMeta.Type, agg.aggregator, memScanCtx)
			j++ // goto next query field id
			// found all query fields return it
			if fieldCount == j {
				return
			}
		case key > agg.fieldKey:
			// store key > query key, return it
			return
		}
	}
}
