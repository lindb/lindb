package tsdb

import (
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/tsdb/series"
)

// DataFamily represents a storage unit for time series data, support multi-version.
type DataFamily interface {
	// series.DataFamilyScanner returns a iterator which scans time series data based on query condition
	series.DataFamilyScanner
	// BaseTime returns the data family's base time
	BaseTime() int64
}

// dataFamily represents a wrapper of kv's family with basic info
type dataFamily struct {
	baseTime int64
	family   kv.Family
}

// newDataFamily creates a data family storage unit
func newDataFamily(baseTime int64, family kv.Family) DataFamily {
	return &dataFamily{
		baseTime: baseTime,
		family:   family,
	}
}

// Scan scans time series data based on query condition, returns scan iterator
func (f *dataFamily) Scan(scanContext series.ScanContext) series.VersionIterator {
	//TODO codingcrush
	return nil
}

// BaseTime returns the data family's base time
func (f *dataFamily) BaseTime() int64 {
	return f.baseTime
}
