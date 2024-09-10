package types

import "github.com/lindb/lindb/pkg/timeutil"

type TimeSeries struct {
	Values      []float64          `json:"values"`
	TimeRange   timeutil.TimeRange `json:"timeRange"`
	Interval    timeutil.Interval  `json:"interval"`
	NumOfPoints int                `json:"numOfPoints"`
	Value       float64            `json:"value"`
	IsSingle    bool               `json:"isSingle"`
}

// NewTimeSeries creates a time series with given time range and interval.
func NewTimeSeries(timeRange timeutil.TimeRange, interval timeutil.Interval) *TimeSeries {
	numOfPoints := (&timeRange).NumOfPoints(interval)
	return &TimeSeries{
		TimeRange:   timeRange,
		Interval:    interval,
		NumOfPoints: numOfPoints,
		Values:      make([]float64, numOfPoints),
	}
}

// NewTimeSeriesWithSingleValue creates a time series with single value.
func NewTimeSeriesWithSingleValue(value float64) *TimeSeries {
	return &TimeSeries{
		Value:       value,
		IsSingle:    true,
		NumOfPoints: 1,
	}
}

// Put puts time series value for given timestamp offset.
func (col *TimeSeries) Put(tsOffset int, value float64) {
	if col.IsSingle {
		col.Value = value
	} else {
		col.Values[tsOffset] = value
	}
}

// Get returns time series value for given timestamp offset.
func (col *TimeSeries) Get(tsOffset int) float64 {
	if col.IsSingle {
		return col.Value
	}
	return col.Values[tsOffset]
}

// Size returns the number of time series points.
func (col *TimeSeries) Size() int {
	return col.NumOfPoints
}

// IsSingleValue returns whether the time series is single value.
func (col *TimeSeries) IsSingleValue() bool {
	return col.IsSingle
}
