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

package metric

import "github.com/lindb/lindb/series/field"

type rollups []*rollup

type rollup struct {
	timeseries *TimeSeries

	window        int64
	currTimestamp int64
	nextTimestamp int64
	currValue     float64
}

func newRollup(capacity int, window int64) *rollup {
	return &rollup{
		timeseries: newTimeSeries(capacity),
		window:     window,
	}
}

func (r *rollup) doRollup(aggType field.AggType, timestamp int64, value float64) {
	if r.currTimestamp == 0 {
		r.nextWindow(timestamp, value)
	} else if timestamp < r.nextTimestamp {
		r.currValue = aggType.Aggregate(r.currValue, value)
	} else {
		r.timeseries.Append(timestamp, value)

		r.nextWindow(timestamp, value)
	}
}

func (r *rollup) getTimeSeries() *TimeSeries {
	if len(r.timeseries.timestamps) == 0 || r.currTimestamp != r.timeseries.timestamps[len(r.timeseries.timestamps)-1] {
		// check last timestamp/value if append time series block
		r.timeseries.Append(r.currTimestamp, r.currValue)
	}
	return r.timeseries
}

func (r *rollup) nextWindow(timestamp int64, value float64) {
	r.currTimestamp = timestamp
	r.currValue = value
	r.nextTimestamp = timestamp + r.window
}

func (r *rollup) reset() {
	r.window = 0
	r.currTimestamp = 0

	r.timeseries.Reset()
}
