// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package types

import "github.com/lindb/lindb/pkg/timeutil"

// TimeSeries represents time series data type.
type TimeSeries struct {
	Values      []float64          `json:"values,omitempty"`
	TimeRange   timeutil.TimeRange `json:"timeRange"`
	Interval    int64              `json:"interval,omitempty"`
	NumOfPoints int                `json:"numOfPoints,omitempty"`
	Value       float64            `json:"value,omitempty"`
	IsSingle    bool               `json:"isSingle,omitempty"`
}

// NewTimeSeries creates a time series with given time range and interval.
func NewTimeSeries(timeRange timeutil.TimeRange, interval timeutil.Interval) *TimeSeries {
	numOfPoints := (&timeRange).NumOfPoints(interval)
	return &TimeSeries{
		TimeRange:   timeRange,
		Interval:    interval.Int64(),
		NumOfPoints: numOfPoints,
		Values:      make([]float64, numOfPoints),
	}
}

// NewTimeSeriesWithValues creates a time series with given time range, interval and values.
func NewTimeSeriesWithValues(timeRange timeutil.TimeRange, interval timeutil.Interval, values []float64) *TimeSeries {
	return &TimeSeries{
		TimeRange:   timeRange,
		Interval:    interval.Int64(),
		NumOfPoints: len(values),
		Values:      values,
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

func (col *TimeSeries) GetValue() float64 {
	if col.IsSingle {
		return col.Value
	}
	if len(col.Values) > 0 {
		return col.Values[0]
	}
	return 0
}
