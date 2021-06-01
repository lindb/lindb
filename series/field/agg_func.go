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

package field

import "math"

var (
	sumAggregator       = sumAgg{aggType: Sum}
	countAggregator     = sumAgg{aggType: Count}
	minAggregator       = minAgg{aggType: Min}
	maxAggregator       = maxAgg{aggType: Max}
	lastValueAggregator = lastValueAgg{aggType: LastValue}
)

// AggFunc returns aggregator function by given func type
func (t AggType) AggFunc() AggFunc {
	switch t {
	case Sum:
		return sumAggregator
	case Count:
		return countAggregator
	case Min:
		return minAggregator
	case Max:
		return maxAggregator
	case LastValue:
		return lastValueAggregator
	default:
		return nil
	}
}

// AggFunc represents field's aggregator function for int64 or float64 value
type AggFunc interface {
	// Aggregate aggregates two float64 values into one
	Aggregate(a, b float64) float64
	// AggType return aggregator type
	AggType() AggType
}

// sumAgg represents sum aggregator
type sumAgg struct {
	aggType AggType
}

func (s sumAgg) AggType() AggType               { return s.aggType }
func (s sumAgg) Aggregate(a, b float64) float64 { return a + b }

// minAgg represents min aggregator
type minAgg struct {
	aggType AggType
}

func (m minAgg) AggType() AggType               { return m.aggType }
func (m minAgg) Aggregate(a, b float64) float64 { return math.Min(a, b) }

// maxAgg represents max aggregator
type maxAgg struct {
	aggType AggType
}

func (m maxAgg) AggType() AggType               { return m.aggType }
func (m maxAgg) Aggregate(a, b float64) float64 { return math.Max(a, b) }

// lastValueAgg represents last value aggregator
type lastValueAgg struct {
	aggType AggType
}

func (m lastValueAgg) AggType() AggType               { return m.aggType }
func (m lastValueAgg) Aggregate(_, b float64) float64 { return b }
