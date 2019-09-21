package memdb

import (
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//////////////////////////////////////////////////////
// Scanner methods
//////////////////////////////////////////////////////

// Scan scans metric store based on scan context
func (ms *metricStore) Scan(sCtx *series.ScanContext) {
	// first need check query's fields is match store's fields, if not return.
	fmList := ms.fieldsMetas.Load().(field.Metas)
	subList, ok := fmList.Intersects(sCtx.FieldIDs)
	if !ok {
		return
	}

	// scan tagIndex when version matches the idSet
	scanOnVersionMatch := func(idx tagIndexINTF) {
		if _, ok := sCtx.SeriesIDSet.Versions()[idx.Version()]; ok {
			ms.scan(sCtx, idx, subList)
		}
	}
	ms.mux.RLock()
	scanOnVersionMatch(ms.mutable)
	immutable := ms.immutable.Load()
	ms.mux.RUnlock()
	if immutable != nil {
		tagIndex := immutable.(tagIndexINTF)
		scanOnVersionMatch(tagIndex)
	}
}

// scan finds time series store from tag index by series ids
func (ms *metricStore) scan(sCtx *series.ScanContext, tagIndex tagIndexINTF, existedFieldMetas field.Metas) {
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
		tStore.Scan(sCtx, version, seriesID, existedFieldMetas)
	}
}
