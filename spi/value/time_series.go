package value

import "github.com/lindb/lindb/pkg/timeutil"

type TimeSeries struct {
	TimeRange timeutil.TimeRange `json:"timeRange"`
	Interval  timeutil.Interval  `json:"interval"`
	Values    []float64          `json:"values"`
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
