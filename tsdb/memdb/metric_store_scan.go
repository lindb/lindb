package memdb

import (
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

// Scan scans metric store based on scan context
func (ms *metricStore) Scan(sCtx *series.ScanContext) {
	// first need check query's fields is match store's fields, if not return.
	fmList := ms.fieldsMetas.Load().(field.Metas)
	_, ok := fmList.Intersects(sCtx.FieldIDs)
	if !ok {
		return
	}
	// scan tagIndex when version matches the idSet
	scanOnVersionMatch := func(idx tagIndexINTF) {
		if _, ok := sCtx.SeriesIDSet.Versions()[idx.Version()]; ok {
			idx.scan(sCtx)
		}
	}
	ms.mux.RLock()
	scanOnVersionMatch(ms.mutable)
	immutable := ms.atomicGetImmutable()
	ms.mux.RUnlock()
	if immutable != nil {
		scanOnVersionMatch(immutable)
	}
}
