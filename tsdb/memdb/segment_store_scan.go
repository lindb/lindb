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
	//FIXME check slot range is match???
	start, end := fs.SlotRange()
	fs.load(agg.GetFieldType(), start, end, aggregators, memScanCtx)
}

// scan scans segment store data based on query time range for complex field store
func (fs *complexFieldStore) scan(agg aggregation.SeriesAggregator, memScanCtx *memScanContext) {
	// check family time is in query time range
	segmentAgg, ok := agg.GetAggregator(fs.familyTime)
	if !ok {
		return
	}
	start, end := fs.SlotRange()
	aggregators := segmentAgg.GetAllAggregators()
	fieldType := agg.GetFieldType()
	for _, a1 := range aggregators {
		pFieldID := a1.FieldID()
		block := fs.blocks[pFieldID]
		if block != nil {
			fs.load(fieldType, start, end, aggregators, memScanCtx)
		}
	}
}
