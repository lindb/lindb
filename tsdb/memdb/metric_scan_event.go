package memdb

import (
	"sync"

	"github.com/RoaringBitmap/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series"
)

// define time series store pool which reuses for scanning data
var storePool = sync.Pool{
	New: func() interface{} {
		return make([]tStoreINTF, series.ScanBufSize)
	},
}

// getStores gets time series store from pool
func getStores() []tStoreINTF {
	stores := storePool.Get()
	return stores.([]tStoreINTF)
}

// putStores puts back time series store to pool
func putStores(stores interface{}) {
	storePool.Put(stores)
}

// metricScanEvent represents the metric level scan event,includes found time series stores, IDs etc.
type metricScanEvent struct {
	stores      []tStoreINTF
	seriesIDs   []uint32
	version     series.Version
	sCtx        *series.ScanContext
	length      int
	aggregators aggregation.FieldAggregates
}

// newScanEvent creates a new metric scan event
func newScanEvent(
	length int,
	stores []tStoreINTF,
	seriesIDs []uint32,
	version series.Version,
	sCtx *series.ScanContext,
) *metricScanEvent {
	return &metricScanEvent{
		stores:    stores,
		seriesIDs: seriesIDs,
		version:   version,
		sCtx:      sCtx,
		length:    length,
	}
}

// ResultSet returns the result set of scanner
func (e *metricScanEvent) ResultSet() interface{} {
	return e.aggregators
}

// SeriesIDs returns the found series IDs
func (e *metricScanEvent) SeriesIDs() *roaring.Bitmap {
	return roaring.BitmapOf(e.seriesIDs[:e.length]...)
}

// Release releases the scan resource for reusing
func (e *metricScanEvent) Release() {
	if e.aggregators != nil {
		e.aggregators.Reset()
		e.sCtx.Release(e.aggregators)
	}
}

// release releases the memory metric store scan's resource
func (e *metricScanEvent) release() {
	for idx := range e.stores {
		e.stores[idx] = nil
	}
	putStores(e.stores)
	if e.seriesIDs != nil {
		series.Uint32Pool.Put(&e.seriesIDs)
	}
}

// Scan scans the memory database, then aggregates the data
func (e *metricScanEvent) Scan() bool {
	defer e.release()
	//FIXME add lock?????
	aggregators, ok := e.sCtx.GetAggregator().(aggregation.FieldAggregates)
	if !ok {
		return false
	}
	e.aggregators = aggregators
	memScanCtx := &memScanContext{
		fieldIDs:    e.sCtx.FieldIDs,
		aggregators: aggregators,
		tsd:         encoding.GetTSDDecoder(),
		fieldCount:  len(e.sCtx.FieldIDs),
	}

	for i := 0; i < e.length; i++ {
		//FIXME do group by and lock/using metric lock
		//seriesID := e.seriesIDs[i]
		store := e.stores[i]
		store.scan(memScanCtx)
	}
	encoding.ReleaseTSDDecoder(memScanCtx.tsd)
	return true
}

// memScanContext represents the memory metric store scan context
type memScanContext struct {
	fieldIDs    []uint16
	aggregators aggregation.FieldAggregates
	tsd         *encoding.TSDDecoder

	fieldCount int
}
