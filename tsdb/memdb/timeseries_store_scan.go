package memdb

import (
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

// Scan scans time series data, then finds field store by field id
func (ts *timeSeriesStore) Scan(
	sCtx *series.ScanContext,
	version series.Version,
	seriesID uint32,
	existedFieldMetas field.Metas,
) {
	worker := sCtx.Worker
	for _, fm := range existedFieldMetas {
		ts.sl.Lock()
		fStore, ok := ts.GetFStore(fm.ID)
		ts.sl.Unlock()
		if !ok {
			continue
		}
		fStore.Scan(sCtx, version, seriesID, fm, ts)
	}

	// send msg to notify current series scan completed
	worker.Complete(seriesID)
}
