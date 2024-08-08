package types

import "github.com/lindb/lindb/pkg/timeutil"

type TimeSeries struct {
	Values    []float64          `json:"values"`
	TimeRange timeutil.TimeRange `json:"timeRange"`
	Interval  timeutil.Interval  `json:"interval"`
}

func NewTimeSeries(timeRange timeutil.TimeRange, interval timeutil.Interval) *TimeSeries {
	return &TimeSeries{
		TimeRange: timeRange,
		Interval:  interval,
		Values:    make([]float64, (&timeRange).NumOfPoints(interval)),
	}
}

func (col *TimeSeries) Put(timestamp int, value float64) {
	col.Values[timestamp] = value
}
