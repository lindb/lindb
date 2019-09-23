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

	ts.sl.Lock()
	for _, fm := range existedFieldMetas {
		fStore, ok := ts.GetFStore(fm.ID)
		if !ok {
			continue
		}
		fStore.Scan(sCtx, version, seriesID, fm, ts)
	}
	ts.sl.Unlock()

	// send msg to notify current series scan completed
	worker.Complete(seriesID)
}
