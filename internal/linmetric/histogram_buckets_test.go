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

package linmetric

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Histogram_buckets_maker(t *testing.T) {
	buckets := makeExponentBuckets(0.1, 30, 100)
	assert.Len(t, buckets, 100)

	assert.True(t, math.IsInf(buckets[99], 1))

	assert.Panics(t, func() {
		makeExponentBuckets(0.1, 0.1, 100)
	})
	assert.Panics(t, func() {
		makeLinearBuckets(0.1, 0.1, 100)
	})
	assert.Panics(t, func() {
		makeExponentBuckets(-1, 100, 100)
	})
	assert.Panics(t, func() {
		makeLinearBuckets(-1, 100, 100)
	})
	assert.Panics(t, func() {
		makeExponentBuckets(1, 100, 2)
	})
	assert.Panics(t, func() {
		makeLinearBuckets(1, 100, 2)
	})

	// min count
	buckets = makeExponentBuckets(1, 100, 3)
	assert.InDeltaSlice(t, []float64{1, 100, math.Inf(1)}, buckets, 0.001)
	// normal case
	buckets = makeExponentBuckets(1, 100, 10)
	assert.InDeltaSlice(t,
		[]float64{1, 1.778, 3.162, 5.623, 10.000, 17.782, 31.622, 56.234, 100, math.Inf(1)},
		buckets,
		0.01,
	)

	// min count
	buckets = makeLinearBuckets(1, 100, 3)
	assert.InDeltaSlice(t, []float64{1, 100, math.Inf(1)}, buckets, 0.001)
	// normal case
	buckets = makeLinearBuckets(1, 100, 10)
	assert.InDeltaSlice(t,
		[]float64{1, 13.375, 25.75, 38.125, 50.5, 62.875, 75.25, 87.625, 100, math.Inf(1)},
		buckets,
		0.01,
	)
}

func Test_histogramBuckets_reset(t *testing.T) {
	bkt1 := newHistogramBuckets(
		1, 10*1000, 100, exponentBucket)
	bkt2 := newHistogramBuckets(
		2, 20*1000, 3, linearBucket)
	for i := 0; i < 10; i++ {
		bkt2.Update(float64(i))
	}

	bkt2.reset(1, 10*1000, 100, exponentBucket)
	assert.Equal(t, bkt1, bkt2)

	assert.Panics(t, func() {
		bkt2.reset(1, 2, 100, bucketAllocateStrategy(30))
	})
}

func Test_histogramBuckets_Update_Exponent(t *testing.T) {
	bkt1 := newHistogramBuckets(
		1, 10*1000, 100, exponentBucket)

	bkt1.Update(-2) // drop
	bkt1.Update(1000)
	bkt1.Update(100)
	bkt1.Update(200)
	assert.Equal(t, float64(3), bkt1.totalCount)
	assert.Equal(t, float64(1000), bkt1.max)
	assert.Equal(t, float64(1300), bkt1.totalSum)
	bkt1.Update(math.Inf(1) + 1) // drop
	bkt1.Update(10*1000 + 1)
	assert.Equal(t, float64(1), bkt1.values[99])
	bkt1.Update(0.3)
	assert.Equal(t, float64(5), bkt1.totalCount)
	assert.Equal(t, 0.3, bkt1.min)
	assert.Equal(t, float64(1), bkt1.values[0])
	assert.Equal(t, float64(10*1000+1), bkt1.max)

	assert.Equal(t, float64(0), bkt1.values[98])
	assert.Equal(t, float64(1), bkt1.values[99])
	bkt1.Update(10*1000 - 1)
	assert.Equal(t, float64(1), bkt1.values[99])
	bkt1.Update(10 * 1000)
	assert.Equal(t, float64(2), bkt1.values[98])

	bkt1.Update(1.2)
	assert.Equal(t, float64(1), bkt1.values[0])
	assert.Equal(t, float64(1), bkt1.values[2])

	bkt1.Update(100 * 1000)
	assert.Equal(t, float64(2), bkt1.values[99])
}

func Test_histogramBuckets_Update_Linear(t *testing.T) {
	bkt2 := newHistogramBuckets(
		1, 10*1000, 100, linearBucket)
	bkt2.Update(-2) // drop
	bkt2.Update(1000)
	bkt2.Update(100)
	bkt2.Update(200)
	assert.Equal(t, float64(3), bkt2.totalCount)
	assert.Equal(t, float64(1000), bkt2.max)
	assert.Equal(t, float64(1300), bkt2.totalSum)
	bkt2.Update(math.Inf(1) + 1) // drop
	bkt2.Update(10*1000 + 1)
	assert.Equal(t, float64(1), bkt2.values[99])
	bkt2.Update(0.3)
	assert.Equal(t, float64(5), bkt2.totalCount)
	assert.Equal(t, 0.3, bkt2.min)
	assert.Equal(t, float64(1), bkt2.values[0])
	assert.Equal(t, float64(10*1000+1), bkt2.max)

	assert.Equal(t, float64(0), bkt2.values[98])
	assert.Equal(t, float64(1), bkt2.values[99])
	bkt2.Update(10*1000 - 1)
	assert.Equal(t, float64(1), bkt2.values[99])
	bkt2.Update(10 * 1000)
	assert.Equal(t, float64(2), bkt2.values[98])

	bkt2.Update(1.2)
	assert.Equal(t, float64(1), bkt2.values[0])
	assert.Equal(t, float64(2), bkt2.values[1])

	bkt2.Update(100 * 1000)
	assert.Equal(t, float64(2), bkt2.values[99])
}
