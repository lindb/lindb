package memdb

import (
	"sync"

	"github.com/RoaringBitmap/roaring"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series"
)

const scanBufSize = 4096

// define series IDs pool which reuses for scanning data
var seriesIDsPool = sync.Pool{
	New: func() interface{} {
		return make([]uint32, scanBufSize)
	},
}

// define time series store pool which reuses for scanning data
var storePool = sync.Pool{
	New: func() interface{} {
		return make([]tStoreINTF, scanBufSize)
	},
}

// getSeriesIDs gets series IDs from pool
func getSeriesIDs() []uint32 {
	IDs := seriesIDsPool.Get()
	return IDs.([]uint32)
}

// putSeriesIDs puts back series IDs to pool
func putSeriesIDs(seriesIDs interface{}) {
	seriesIDsPool.Put(seriesIDs)
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
	stores     []tStoreINTF
	seriesIDs  []uint32
	version    series.Version
	sCtx       *series.ScanContext
	length     int
	aggregates aggregation.FieldAggregates
}

// newScanEvent creates a new metric scan event
func newScanEvent(length int, stores []tStoreINTF, seriesIDs []uint32, version series.Version,
	sCtx *series.ScanContext) *metricScanEvent {
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
	return e.aggregates
}

// SeriesIDs returns the found series IDs
func (e *metricScanEvent) SeriesIDs() *roaring.Bitmap {
	return roaring.BitmapOf(e.seriesIDs[:e.length]...)
}

// Release releases the scan resource for reusing
func (e *metricScanEvent) Release() {
	if e.aggregates != nil {
		e.aggregates.Reset()
		e.sCtx.Release(e.aggregates)
	}
}

// release releases the memory metric store scan's resource
func (e *metricScanEvent) release() {
	for idx := range e.stores {
		e.stores[idx] = nil
	}
	putStores(e.stores)
	if e.seriesIDs != nil {
		putSeriesIDs(e.seriesIDs)
	}
}

// Scan scans the memory database, then aggregates the data
func (e *metricScanEvent) Scan() bool {
	defer e.release()
	//FIXME add lock?????
	aggregates, ok := e.sCtx.GetAggregator().(aggregation.FieldAggregates)
	if !ok {
		return false
	}
	e.aggregates = aggregates
	memScanCtx := &memScanContext{
		fieldIDs:   e.sCtx.FieldIDs,
		aggregates: aggregates,
		tsd:        encoding.GetTSDDecoder(),
		fieldCount: len(e.sCtx.FieldIDs),
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
	fieldIDs   []uint16
	aggregates aggregation.FieldAggregates
	tsd        *encoding.TSDDecoder

	fieldCount int
}
