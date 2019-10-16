package memdb

import "github.com/lindb/lindb/aggregation"

// scan scans segment store data based on query time range
func (fs *simpleFieldStore) scan(agg aggregation.SeriesAggregator, memScanCtx *memScanContext) {
	// check family time is in query time range
	segmentAgg, ok := agg.GetAggregator(fs.familyTime)
	if !ok {
		return
	}
	aggregators := segmentAgg.GetAllAggregators()
	fs.block.scan(fs.aggFunc, aggregators, memScanCtx)
}
