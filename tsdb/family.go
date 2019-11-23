package tsdb

import (
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/timeutil"
)

//go:generate mockgen -source=./family.go -destination=./family_mock.go -package=tsdb

// DataFamily represents a storage unit for time series data, support multi-version.
type DataFamily interface {
	// flow.Scanner scans files under kv store based on query condition
	//flow.Scanner
	// Interval returns the interval data family's interval
	Interval() int64
	// TimeRange returns the data family's base time range
	TimeRange() timeutil.TimeRange
	// Family returns the raw kv family
	Family() kv.Family
}

// dataFamily represents a wrapper of kv's family with basic info
type dataFamily struct {
	interval  timeutil.Interval
	timeRange timeutil.TimeRange
	family    kv.Family
}

// newDataFamily creates a data family storage unit
func newDataFamily(
	interval timeutil.Interval,
	timeRange timeutil.TimeRange,
	family kv.Family,
) DataFamily {
	return &dataFamily{
		interval:  interval,
		timeRange: timeRange,
		family:    family,
	}
}

// Scan scans time series data based on query condition
//func (f *dataFamily) Scan(qCtx *flow.StorageQueryContext) {
//	snapShot := f.family.GetSnapshot()
//	defer snapShot.Close()
//
//	readers, err := snapShot.FindReaders(qCtx.MetricID)
//	if err != nil {
//		return
//	}
//	scanner := metricsdata.NewScanner(readers)
//	scanner.Scan(qCtx)
//}

// Interval returns the data family's interval
func (f *dataFamily) Interval() int64 {
	return f.interval.Int64()
}

// TimeRange returns the data family's base time range
func (f *dataFamily) TimeRange() timeutil.TimeRange {
	return f.timeRange
}

func (f *dataFamily) Family() kv.Family {
	return f.family
}
