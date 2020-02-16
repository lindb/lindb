package tsdb

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source=./family.go -destination=./family_mock.go -package=tsdb

// for testing
var (
	newReaderFunc = metricsdata.NewReader
	newFilterFunc = metricsdata.NewFilter
)

// DataFamily represents a storage unit for time series data, support multi-version.
type DataFamily interface {
	// Interval returns the interval data family's interval
	Interval() int64
	// TimeRange returns the data family's base time range
	TimeRange() timeutil.TimeRange
	// Family returns the raw kv family
	Family() kv.Family

	// flow.DataFilter filters data under data family based on query condition
	flow.DataFilter
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

// Interval returns the data family's interval
func (f *dataFamily) Interval() int64 {
	return f.interval.Int64()
}

// TimeRange returns the data family's base time range
func (f *dataFamily) TimeRange() timeutil.TimeRange {
	return f.timeRange
}

// Family returns the kv store's family
func (f *dataFamily) Family() kv.Family {
	return f.family
}

// Filter filters the data based on metric/version/seriesIDs,
// if finds data then returns the FilterResultSet, else returns nil
func (f *dataFamily) Filter(metricID uint32, fieldIDs []field.ID,
	seriesIDs *roaring.Bitmap, timeRange timeutil.TimeRange,
) (resultSet []flow.FilterResultSet, err error) {
	snapShot := f.family.GetSnapshot()
	defer func() {
		if err != nil || len(resultSet) == 0 {
			// if not find metrics data or has err, close snapshot directly
			snapShot.Close()
		}
	}()
	readers, err := snapShot.FindReaders(metricID)
	if err != nil {
		engineLogger.Error("filter data family error", logger.Error(err))
		return
	}
	var metricReaders []metricsdata.Reader
	for _, reader := range readers {
		value, ok := reader.Get(metricID)
		// metric data not found
		if !ok {
			continue
		}
		r, err := newReaderFunc(value)
		if err != nil {
			return nil, err
		}
		metricReaders = append(metricReaders, r)
	}
	if len(metricReaders) == 0 {
		return
	}
	filter := newFilterFunc(f.timeRange.Start, snapShot, metricReaders)
	return filter.Filter(fieldIDs, seriesIDs)
}
