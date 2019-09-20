package memdb

import (
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//////////////////////////////////////////////////////
// Scanner methods
//////////////////////////////////////////////////////

// findFieldMetas returns if query's fields are in store, if all query fields found returns true, else returns false
func (ms *metricStore) findFieldMetas(fieldIDs []uint16) (map[uint16]*field.Meta, bool) {
	fmList := ms.fieldsMetas.Load().(*fieldsMetas)
	result := make(map[uint16]*field.Meta)
	for _, fieldID := range fieldIDs {
		result[fieldID] = &field.Meta{}
	}

	found := 0
	for _, fm := range *fmList {
		fieldMeta, ok := result[fm.ID]
		if ok {
			*fieldMeta = fm
			found++
		}
	}
	return result, found == len(fieldIDs)
}

// Scan scans metric store based on scan context
func (ms *metricStore) Scan(sCtx *series.ScanContext) {
	// first need check query's fields is match store's fields, if not return.
	fieldMetas, ok := ms.findFieldMetas(sCtx.FieldIDs)
	if !ok {
		return
	}

	// collect all tagIndexes whose version matches the idSet
	collectOnVersionMatch := func(idx tagIndexINTF) {
		if _, ok := sCtx.SeriesIDSet.Versions()[idx.Version()]; ok {
			ms.scan(sCtx, idx, fieldMetas)
		}
	}
	ms.mux.RLock()
	collectOnVersionMatch(ms.mutable)
	immutable := ms.immutable.Load()
	ms.mux.RUnlock()
	if immutable != nil {
		tagIndex := immutable.(tagIndexINTF)
		collectOnVersionMatch(tagIndex)
	}
}

// scan finds time series store from tag index by series ids
func (ms *metricStore) scan(sCtx *series.ScanContext, tagIndex tagIndexINTF, fieldMetas map[uint16]*field.Meta) {
	// support multi-version
	version := tagIndex.Version()
	seriesIDs := sCtx.SeriesIDSet.Versions()[version]
	seriesIDIt := seriesIDs.Iterator()

	for seriesIDIt.HasNext() {
		seriesID := seriesIDIt.Next()
		tStore, ok := tagIndex.GetTStoreBySeriesID(seriesID)
		// if not found or no data
		if !ok || tStore.IsNoData() {
			continue
		}

		// scan time series store
		tStore.Scan(sCtx, version, seriesID, fieldMetas)
	}
}
