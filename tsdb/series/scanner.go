package series

import (
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./scanner.go -destination=./scanner_mock.go -package=series

// ScanContext is the context for scanning data
type ScanContext struct {
	// required
	MetricID  uint32
	FieldIDs  []uint16
	TimeRange timeutil.TimeRange
	// optional, if SeriesIDSet is nil, just search metric level data
	SeriesIDSet *MultiVerSeriesIDSet
	// for context usage
	TimeInterval int64 // database interval in seconds
	FamilyTime   int64 // family time
}

// DataFamilyScanner represents the scan ability over family data.
type DataFamilyScanner interface {
	Scan(sCtx ScanContext) VersionIterator
}
