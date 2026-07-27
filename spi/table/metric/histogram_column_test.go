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
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/tree"
)

// makeFloatResult creates a result[float64] with a single value pre-set at index 0.
func makeFloatResult(v float64) Result {
	r := NewResult[float64](1)
	r.(*result[float64]).array.SetValue(0, v)
	return r
}

// makePhysCol creates a minimal column[float64] that points to the given aggregator slot.
func makePhysCol(slotIndex int) *column[float64] {
	agg := &aggregator[float64]{
		index: slotIndex,
		aggFn: func(a, b float64) float64 { return a + b },
	}
	return &column[float64]{aggs: []*aggregator[float64]{agg}}
}

// ----- computeQuantile tests -----

// TestComputeQuantile_NoBuckets verifies that a histogram with no buckets returns NaN.
func TestComputeQuantile_NoBuckets(t *testing.T) {
	h := &histogramColumn{phi: 0.99}
	assert.True(t, math.IsNaN(h.computeQuantile(nil)))
}

// TestComputeQuantile_ZeroTotal verifies that all-zero counts return NaN.
func TestComputeQuantile_ZeroTotal(t *testing.T) {
	results := []Result{makeFloatResult(0), makeFloatResult(0), makeFloatResult(0)}
	h := &histogramColumn{
		phi: 0.5,
		buckets: []histogramEntry{
			{col: makePhysCol(0), bound: 1.0},
			{col: makePhysCol(1), bound: 5.0},
			{col: makePhysCol(2), bound: math.Inf(1)},
		},
	}
	assert.True(t, math.IsNaN(h.computeQuantile(results)))
}

// TestComputeQuantile_p50_Uniform verifies p50 with a uniform distribution across 4 buckets.
func TestComputeQuantile_p50_Uniform(t *testing.T) {
	// bounds: (0,1], (1,5], (5,10], (10,+Inf)  each with 25 observations
	results := []Result{
		makeFloatResult(25), // bucket 0: (0,1]
		makeFloatResult(25), // bucket 1: (1,5]
		makeFloatResult(25), // bucket 2: (5,10]
		makeFloatResult(25), // bucket 3: (10,+Inf)
	}
	h := &histogramColumn{
		phi: 0.5,
		buckets: []histogramEntry{
			{col: makePhysCol(0), bound: 1.0},
			{col: makePhysCol(1), bound: 5.0},
			{col: makePhysCol(2), bound: 10.0},
			{col: makePhysCol(3), bound: math.Inf(1)},
		},
	}
	got := h.computeQuantile(results)
	// cumCounts = [25, 50, 75, 100]; target = 50 → bucket idx=1 (bound=5)
	// lowerBound=1, prevCum=25, bucketCount=25
	// fraction = (50-25)/25 = 1.0 → result = 1 + 1.0*(5-1) = 5.0
	assert.InDelta(t, 5.0, got, 1e-9)
}

// TestComputeQuantile_p99 verifies p99 with most observations in the first bucket.
func TestComputeQuantile_p99(t *testing.T) {
	results := []Result{
		makeFloatResult(90), // (0,1]
		makeFloatResult(9),  // (1,10]
		makeFloatResult(1),  // (10,+Inf)
	}
	h := &histogramColumn{
		phi: 0.99,
		buckets: []histogramEntry{
			{col: makePhysCol(0), bound: 1.0},
			{col: makePhysCol(1), bound: 10.0},
			{col: makePhysCol(2), bound: math.Inf(1)},
		},
	}
	got := h.computeQuantile(results)
	// cumCounts = [90, 99, 100]; target = 99 → bucket idx=1 (bound=10)
	// lowerBound=1, prevCum=90, bucketCount=9
	// fraction = (99-90)/9 = 1.0 → result = 1 + 1.0*(10-1) = 10.0
	assert.InDelta(t, 10.0, got, 1e-9)
}

// TestComputeQuantile_Phi0 verifies that phi=0 returns the lower bound of the first non-empty bucket.
func TestComputeQuantile_Phi0(t *testing.T) {
	results := []Result{makeFloatResult(10), makeFloatResult(10), makeFloatResult(10)}
	h := &histogramColumn{
		phi: 0.0,
		buckets: []histogramEntry{
			{col: makePhysCol(0), bound: 5.0},
			{col: makePhysCol(1), bound: 10.0},
			{col: makePhysCol(2), bound: math.Inf(1)},
		},
	}
	got := h.computeQuantile(results)
	// target = 0; first bucket where cum >= 0 is idx=0
	// lowerBound=0 (no previous bucket), fraction = 0 → result = 0
	assert.InDelta(t, 0.0, got, 1e-9)
}

// TestComputeQuantile_AllInInfBucket verifies that when all observations fall in +Inf, return lowerBound.
func TestComputeQuantile_AllInInfBucket(t *testing.T) {
	results := []Result{makeFloatResult(0), makeFloatResult(10)}
	h := &histogramColumn{
		phi: 0.99,
		buckets: []histogramEntry{
			{col: makePhysCol(0), bound: 1.0},
			{col: makePhysCol(1), bound: math.Inf(1)},
		},
	}
	got := h.computeQuantile(results)
	// cumCounts = [0, 10]; target=9.9 → bucket idx=1 (+Inf)
	// return lowerBound = 1.0
	assert.InDelta(t, 1.0, got, 1e-9)
}

// ----- readScalar tests -----

func TestReadScalar_NilCol(t *testing.T) {
	h := &histogramColumn{}
	assert.Equal(t, float64(0), h.readScalar(nil, nil))
}

func TestReadScalar_Normal(t *testing.T) {
	results := []Result{nil, makeFloatResult(42.5)}
	h := &histogramColumn{}
	col := makePhysCol(1)
	got := h.readScalar(results, col)
	assert.InDelta(t, 42.5, got, 1e-9)
}

// ----- helper function tests -----

func TestIsBucketField(t *testing.T) {
	assert.True(t, isBucketField("latency.__bucket_1"))
	assert.True(t, isBucketField("req.__bucket_+Inf"))
	assert.False(t, isBucketField("latency_sum"))
	assert.False(t, isBucketField("__bucket_1"))    // no dot-prefix histoName
	assert.False(t, isBucketField("latency.other")) // suffix not starting with __bucket_
	assert.False(t, isBucketField(""))
}

func TestParseBucketBound(t *testing.T) {
	bound, ok := parseBucketBound("latency.__bucket_1.5")
	assert.True(t, ok)
	assert.InDelta(t, 1.5, bound, 1e-9)

	bound, ok = parseBucketBound("req.__bucket_+Inf")
	assert.True(t, ok)
	assert.True(t, math.IsInf(bound, 1))

	_, ok = parseBucketBound("latency_sum")
	assert.False(t, ok)

	_, ok = parseBucketBound("latency.other")
	assert.False(t, ok)
}

func TestHistoNameFromField(t *testing.T) {
	assert.Equal(t, "latency", histoNameFromField("latency.__bucket_1"))
	assert.Equal(t, "latency", histoNameFromField("latency.__sum"))
	assert.Equal(t, "latency", histoNameFromField("latency.__count"))
	assert.Equal(t, "latency", histoNameFromField("latency.__min"))
	assert.Equal(t, "latency", histoNameFromField("latency.__max"))
	assert.Equal(t, "", histoNameFromField("plain_field"))
}

func TestHistogramColumnBuilderSeal(t *testing.T) {
	b := &histogramColumnBuilder{
		aggFunc:      tree.HistogramQuantile,
		phi:          0.99,
		outputOffset: 0,
		buckets: []histogramEntry{
			{col: makePhysCol(2), bound: 10.0},
			{col: makePhysCol(0), bound: 1.0}, // intentionally out of order
			{col: makePhysCol(1), bound: 5.0},
		},
	}
	hCol := b.seal()
	assert.Equal(t, tree.HistogramQuantile, hCol.aggFunc)
	assert.Equal(t, 0.99, hCol.phi)
	// Verify sort: bounds should be ascending
	assert.Equal(t, 1.0, hCol.buckets[0].bound)
	assert.Equal(t, 5.0, hCol.buckets[1].bound)
	assert.Equal(t, 10.0, hCol.buckets[2].bound)
}
