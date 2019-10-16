package series

import (
	"sync"

	"github.com/RoaringBitmap/roaring"

	"github.com/lindb/lindb/pkg/interval"
)

//go:generate mockgen -source=./scanner.go -destination=./scanner_mock.go -package=series

// ScanContext is the context for scanning data
type ScanContext struct {
	// required
	MetricID   uint32
	FieldIDs   []uint16
	HasGroupBy bool

	Worker ScanWorker // scan worker which handles field event

	// optional, if SeriesIDSet is nil, just search metric level data
	SeriesIDSet *MultiVerSeriesIDSet

	// runtime, required for memory scan
	IntervalCalc interval.Calculator

	Aggregators sync.Pool
}

// ContainsFieldID checks if fieldID is in search
func (sCtx *ScanContext) ContainsFieldID(fieldID uint16) bool {
	for _, id := range sCtx.FieldIDs {
		if id == fieldID {
			return true
		}
	}
	return false
}

// GetAggregator gets aggregator from the pool of scanner context
func (sCtx *ScanContext) GetAggregator() interface{} {
	return sCtx.Aggregators.Get()
}

// Release puts back aggregator to the pool of scanner context
func (sCtx *ScanContext) Release(agg interface{}) {
	sCtx.Aggregators.Put(agg)
}

// Scanner represents the scan ability over memory database and files under data family.
type Scanner interface {
	// Scan scans the data over memory or files
	Scan(sCtx *ScanContext)
}

// ScanEvent represents the scan event, includes scan context and result
type ScanEvent interface {
	// SeriesIDs returns the found series IDs
	SeriesIDs() *roaring.Bitmap
	// Release releases the scan resource for reusing
	Release()
	// ResultSet returns the result set of scanner
	ResultSet() interface{}
	// Scan scans the storage data, then aggregates the data
	Scan() bool
}

// ScanWorker represents the scan worker which handles the field event which scans result
type ScanWorker interface {
	// Emit emits the field event of one series,
	// make sure emit event in order based on series id.
	Emit(event ScanEvent)
	// Close closes scan worker, then releases the resources
	Close()
}
