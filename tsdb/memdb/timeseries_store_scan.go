package memdb

// scan scans the time series data based on field ids.
// NOTICE: field ids and fields aggregator must be in order.
func (ts *timeSeriesStore) scan(memScanCtx *memScanContext) {
	//idx := 0
	//for _, fieldStore := range ts.fStoreNodes {
	//	fieldID := fieldStore.GetFieldID()
	//	switch {
	//	case fieldID == memScanCtx.fieldIDs[idx]:
	//		agg := memScanCtx.aggregators[idx]
	//		fieldStore.scan(agg, memScanCtx)
	//		idx++ // goto next query field id
	//		// found all query fields return it
	//		if memScanCtx.fieldCount == idx {
	//			return
	//		}
	//	case fieldID > memScanCtx.fieldIDs[idx]:
	//		// store field id > query field id, return it
	//		return
	//	}
	//}
}
