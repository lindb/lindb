package memdb

import "github.com/lindb/lindb/aggregation"

// scan scans segment store data based on query time range for simple field store
func (fs *simpleFieldStore) scan(agg aggregation.SeriesAggregator, memScanCtx *memScanContext) {
	// check family time is in query time range
	segmentAgg, ok := agg.GetAggregator(fs.familyTime)
	if !ok {
		return
	}
	aggregators := segmentAgg.GetAllAggregators()
	fs.block.scan(fs.aggFunc, aggregators, memScanCtx)
}

// scan scans segment store data based on query time range for complex field store
func (fs *complexFieldStore) scan(agg aggregation.SeriesAggregator, memScanCtx *memScanContext) {
	// check family time is in query time range
	segmentAgg, ok := agg.GetAggregator(fs.familyTime)
	if !ok {
		return
	}
	for _, a1 := range segmentAgg.GetAllAggregators() {
		pFieldID := a1.FieldID()
		block := fs.blocks[pFieldID]
		if block != nil {
			block.scan(fs.schema.GetAggFunc(pFieldID), []aggregation.PrimitiveAggregator{a1}, memScanCtx)
		}
	}
}
