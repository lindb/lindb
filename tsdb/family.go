package tsdb

import (
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source=./family.go -destination=./family_mock.go -package=tsdb

// DataFamily represents a storage unit for time series data, support multi-version.
type DataFamily interface {
	// series.Scanner scans files under kv store based on query condition
	series.Scanner
	// Interval returns the interval data family's interval
	Interval() int64
	// TimeRange returns the data family's base time range
	TimeRange() timeutil.TimeRange
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
func (f *dataFamily) Scan(sCtx *series.ScanContext) {
	snapShot := f.family.GetSnapshot()
	defer snapShot.Close()

	readers, err := snapShot.FindReaders(sCtx.MetricID)
	if err != nil {
		return
	}
	scanner := metricsdata.NewScanner(readers)
	scanner.Scan(sCtx)
}

// Interval returns the data family's interval
func (f *dataFamily) Interval() int64 {
	return f.interval.Int64()
}

// TimeRange returns the data family's base time range
func (f *dataFamily) TimeRange() timeutil.TimeRange {
	return f.timeRange
}
