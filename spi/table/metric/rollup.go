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

import (
	"github.com/lindb/common/models"
)

type rollup[V float64 | *models.Exemplar] struct {
	timeseries *TimeSeries[V]

	window        int64
	currTimestamp int64
	nextTimestamp int64
	currValue     V
	agg           func(a, b V) V
}

func newRollup[V float64 | *models.Exemplar](capacity int, window int64, agg aggregateFunc[V]) *rollup[V] {
	return &rollup[V]{
		timeseries: newTimeSeries[V](capacity),
		window:     window,
		agg:        agg,
	}
}

func (r *rollup[V]) doRollup(timestamp int64, value V) {
	if r.currTimestamp == 0 {
		r.nextWindow(timestamp, value)
	} else if timestamp < r.nextTimestamp {
		r.currValue = r.agg(r.currValue, value)
	} else {
		// Flush the previous window's aggregated value before starting the new window.
		r.timeseries.Append(r.currTimestamp, r.currValue)
		r.nextWindow(timestamp, value)
	}
}

func (r *rollup[V]) getTimeSeries() *TimeSeries[V] {
	if len(r.timeseries.timestamps) == 0 || r.currTimestamp != r.timeseries.timestamps[len(r.timeseries.timestamps)-1] {
		// check last timestamp/value if append time series block
		r.timeseries.Append(r.currTimestamp, r.currValue)
	}
	return r.timeseries
}

func (r *rollup[V]) nextWindow(timestamp int64, value V) {
	r.currTimestamp = timestamp
	r.currValue = value
	r.nextTimestamp = timestamp + r.window
}

func (r *rollup[V]) reset() {
	r.currTimestamp = 0
	r.nextTimestamp = 0

	r.timeseries.Reset()
}
