package memdb

import (
	"sort"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./timeseries_store.go -destination=./timeseries_store_mock.go -package memdb

const emptyTimeSeriesStoreSize = 24 // fStores slice

// tStoreINTF abstracts a time-series store
type tStoreINTF interface {
	// GetFStore returns the fStore in field list by family/field
	GetFStore(familyID familyID, fieldID field.ID) (fStoreINTF, bool)
	// InsertFStore inserts a new fStore to field list.
	InsertFStore(fStore fStoreINTF) (createdSize int)
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
func (ts *timeSeriesStore) GetFStore(familyID familyID, fieldID field.ID) (fStoreINTF, bool) {
	fieldKey := buildFieldKey(familyID, fieldID)
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
func (ts *timeSeriesStore) InsertFStore(fStore fStoreINTF) (createdSize int) {
	createdSize = emptyFieldStoreSize + 8
	if ts.fStoreNodes == nil {
		ts.fStoreNodes = []fStoreINTF{fStore}
		return
	}
	ts.fStoreNodes = append(ts.fStoreNodes, fStore)
	sort.Sort(ts.fStoreNodes)
	return createdSize
}

// FlushSeriesTo flushes the series data segment.
func (ts *timeSeriesStore) FlushSeriesTo(flusher metricsdata.Flusher, flushCtx flushContext) {
	flushFamilyID := flushCtx.familyID
	var stores []fStoreINTF
	for _, fStore := range ts.fStoreNodes {
		// need flush by family id
		if flushFamilyID == fStore.GetFamilyID() {
			stores = append(stores, fStore)
		}
	}
	fStoreLen := len(stores)
	// if no field store under current data family
	if stores == nil || fStoreLen == 0 {
		return
	}
	fieldMetas := flusher.GetFieldMetas()
	idx := 0
	for _, fieldMeta := range fieldMetas {
		if idx < fStoreLen && fieldMeta.ID == stores[idx].GetFieldID() {
			// flush field data
			stores[idx].FlushFieldTo(flusher, fieldMeta, flushCtx)
			idx++
		} else {
			// must flush nil data for metric has mutli-field
			flusher.FlushField(nil)
		}
	}
}

// scan scans the time series data based on key(family+field).
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
			fieldStore.Load(agg.fieldMeta.Type, agg.aggregator.GetBlock(), memScanCtx)
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
