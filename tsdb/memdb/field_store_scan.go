package memdb

import "github.com/lindb/lindb/aggregation"

// scan scans the field store's data
func (fs *fieldStore) scan(agg aggregation.SeriesAggregator, memScanCtx *memScanContext) {
	for _, fsStore := range fs.sStoreNodes {
		fsStore.scan(agg, memScanCtx)
	}
}
