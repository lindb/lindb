package memdb

import (
	"github.com/lindb/lindb/series"
)

//////////////////////////////////////////////////////
// Scanner methods
//////////////////////////////////////////////////////

// findFieldMetas returns if query's fields are in store, if all query fields found returns true, else returns false
func (ms *metricStore) findFieldMetas(fieldIDs []uint16) (map[uint16]*fieldMeta, bool) {
	ms.mutex4Fields.RLock()
	defer ms.mutex4Fields.RUnlock()
	result := make(map[uint16]*fieldMeta)
	for _, fieldID := range fieldIDs {
		result[fieldID] = &fieldMeta{}
	}

	found := 0
	for _, field := range ms.fieldsMetas {
		fieldMeta, ok := result[field.fieldID]
		if ok {
			fieldMeta.fieldID = field.fieldID
			fieldMeta.fieldType = field.fieldType
			fieldMeta.fieldName = field.fieldName
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
		if _, ok := sCtx.SeriesIDSet.Versions()[idx.getVersion()]; ok {
			ms.scan(sCtx, idx, fieldMetas)
		}
	}
	// find from immutable store
	ms.mutex4Immutable.RLock()
	for _, idx := range ms.immutable {
		collectOnVersionMatch(idx)
	}
	ms.mutex4Immutable.RUnlock()

	// find from mutable store
	ms.mutex4Mutable.RLock()
	collectOnVersionMatch(ms.mutable)
	ms.mutex4Mutable.RUnlock()
}

// scan finds time series store from tag index by series ids
func (ms *metricStore) scan(sCtx *series.ScanContext, tagIndex tagIndexINTF, fieldMetas map[uint16]*fieldMeta) {
	// support multi-version
	version := tagIndex.getVersion()
	seriesIDs := sCtx.SeriesIDSet.Versions()[version]
	seriesIDIt := seriesIDs.Iterator()

	for seriesIDIt.HasNext() {
		seriesID := seriesIDIt.Next()
		tStore, ok := tagIndex.getTStoreBySeriesID(seriesID)
		// if not found or no data
		if !ok || tStore.isNoData() {
			continue
		}

		// scan time series store
		tStore.scan(sCtx, version, seriesID, fieldMetas)
	}
}
