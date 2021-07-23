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
)

type bucketAllocateStrategy int

const (
	exponentBucket bucketAllocateStrategy = iota + 1
	linearBucket
)

const (
	defaultMinBucketUpperBound float64 = 0.1
	defaultMaxBucketUpperBound float64 = 30 * 1000
	defaultBucketCount         int     = 50
)

func assertBucketParams(start, to float64, count int) {
	if start >= to || start <= 0 {
		panic("valid exponent range is 0< start < to")
	}
	if count <= 2 {
		panic("bucket count must > 2")
	}
}

func makeExponentBuckets(start, to float64, count int) []float64 {
	assertBucketParams(start, to, count)

	multiplier := math.Pow(to/start, 1/float64(count-2))
	var buckets = make([]float64, count)
	for i := range buckets {
		if i == 0 {
			buckets[i] = start
		} else {
			buckets[i] = buckets[i-1] * multiplier
		}
	}
	buckets[len(buckets)-1] = math.Inf(1) + 1
	return buckets
}

func makeLinearBuckets(start, to float64, count int) []float64 {
	assertBucketParams(start, to, count)

	delta := (to - start) / float64(count-2)
	var buckets = make([]float64, count)
	for i := range buckets {
		if i == 0 {
			buckets[i] = start
		} else {
			buckets[i] = buckets[i-1] + delta
		}
	}
	buckets[len(buckets)-1] = math.Inf(1) + 1
	return buckets
}

// not thread-safe
type histogramBuckets struct {
	min         float64                // min updated value
	max         float64                // max updated value
	totalCount  float64                // total update count
	totalSum    float64                // total sum of updated values
	lower       float64                // left boundary, static
	upper       float64                // right boundary, static
	strategy    bucketAllocateStrategy // exponent or linear, static
	values      []float64              // count in different sections
	upperBounds []float64              // upper bounds
	gradient    float64                // multiplier or delta
}

func newHistogramBuckets(start, to float64, count int, strategy bucketAllocateStrategy) *histogramBuckets {
	bkt := &histogramBuckets{}
	bkt.reset(start, to, count, strategy)
	return bkt
}

// reset resets the counts
func (bkt *histogramBuckets) reset(start, to float64, count int, strategy bucketAllocateStrategy) {
	var upperBounds []float64
	switch strategy {
	case exponentBucket:
		upperBounds = makeExponentBuckets(start, to, count)
		bkt.gradient = math.Pow(to/start, 1/float64(count-2))
	case linearBucket:
		upperBounds = makeLinearBuckets(start, to, count)
		bkt.gradient = (to - start) / float64(count-2)
	default:
		panic("unrecognized strategy")
	}
	bkt.min = 0
	bkt.max = 0
	bkt.totalCount = 0
	bkt.totalSum = 0
	bkt.lower = start
	bkt.upper = to
	bkt.strategy = strategy
	bkt.values = make([]float64, count)
	bkt.upperBounds = upperBounds
}

func (bkt *histogramBuckets) Update(v float64) {
	if math.IsNaN(v) || v < 0 || math.IsInf(v, 1) {
		return
	}
	var bktIdx int
	if bkt.strategy == exponentBucket {
		bktIdx = int(math.Ceil(math.Log10(v/bkt.lower) / math.Log10(bkt.gradient)))
	} else {
		bktIdx = int(math.Ceil((v - bkt.lower) / bkt.gradient))
	}
	switch {
	case bktIdx <= 0:
		bkt.values[0]++
	case bktIdx >= len(bkt.values):
		bkt.values[len(bkt.values)-1]++
	default:
		bkt.values[bktIdx]++
	}

	bkt.totalCount++
	bkt.totalSum += v
	if bkt.min == 0 {
		bkt.min = v
	} else if v < bkt.min {
		bkt.min = v
	}
	if v > bkt.max {
		bkt.max = v
	}
}

func cloneFloat64Slice(src []float64) []float64 {
	var dst = make([]float64, len(src))
	copy(dst, src)
	return dst
}
