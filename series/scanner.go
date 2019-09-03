package series

import (
	"github.com/lindb/lindb/pkg/interval"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./scanner.go -destination=./scanner_mock.go -package=series

// ScanContext is the context for scanning data
type ScanContext struct {
	// required
	MetricID  uint32
	FieldIDs  []uint16
	TimeRange timeutil.TimeRange
	Worker    ScanWorker // scan worker which handles field event

	// optional, if SeriesIDSet is nil, just search metric level data
	SeriesIDSet *MultiVerSeriesIDSet

	// runtime, required for memory scan
	Interval     int64
	IntervalCalc interval.Calculator
}

// Scanner represents the scan ability over memory database and files under data family.
type Scanner interface {
	// Scan scans the data over memory or files
	Scan(sCtx *ScanContext)
}

// ScanWorker represents the scan worker which handles the field event which scans result
type ScanWorker interface {
	// Emit emits the field event of one series,
	// make sure emit event in order based on series id.
	Emit(event *FieldEvent)
	// Complete notifies current series scan completed
	Complete(seriesID uint32)
	// Close closes scan worker, then releases the resources
	Close()
}
