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
	"sync"
	"time"

	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
)

// BoundDeltaHistogram is a histogram which has been Bound to a certain metric
// with field-name and metrics, used for non-negative values
//
// a default created bucket will be automatically created,
// however, you can also specify your own buckets.
// Prometheus Histogram's buckets are cumulative where values in each buckets is cumulative,
type BoundDeltaHistogram struct {
	mu             sync.Mutex
	bkts           *histogramBuckets
	lastValues     []float64
	lastTotalCount float64
	lastTotalSum   float64
}

func NewHistogram() *BoundDeltaHistogram {
	h := &BoundDeltaHistogram{
		bkts: newHistogramBuckets(
			defaultMinBucketUpperBound,
			defaultMaxBucketUpperBound,
			defaultBucketCount,
			exponentBucket,
		),
	}
	h.afterResetBuckets()
	return h
}

func (h *BoundDeltaHistogram) afterResetBuckets() {
	h.lastValues = cloneFloat64Slice(h.bkts.values)
	h.lastTotalCount = h.bkts.totalCount
	h.lastTotalSum = h.bkts.totalSum
}

func (h *BoundDeltaHistogram) WithExponentBuckets(lower, upper time.Duration, count int) *BoundDeltaHistogram {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.bkts.reset(float64(lower.Nanoseconds()/1e6), float64(upper.Nanoseconds()/1e6), count, exponentBucket)
	h.afterResetBuckets()
	return h
}

func (h *BoundDeltaHistogram) WithLinearBuckets(lower, upper time.Duration, count int) *BoundDeltaHistogram {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.bkts.reset(float64(lower.Nanoseconds()/1e6), float64(upper.Nanoseconds()/1e6), count, linearBucket)
	h.afterResetBuckets()
	return h
}

func (h *BoundDeltaHistogram) UpdateDuration(d time.Duration) {
	h.UpdateMilliseconds(float64(d.Nanoseconds() / 1e6))
}

func (h *BoundDeltaHistogram) UpdateSince(start time.Time) {
	h.UpdateMilliseconds(float64(time.Since(start).Nanoseconds() / 1e6))
}

func (h *BoundDeltaHistogram) UpdateMilliseconds(s float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.bkts.Update(s)
}

func (h *BoundDeltaHistogram) UpdateSeconds(s float64) {
	h.UpdateMilliseconds(s * 1000)
}

func (h *BoundDeltaHistogram) Update(f func()) {
	start := time.Now()
	f()
	h.UpdateSince(start)
}

func (h *BoundDeltaHistogram) marshalToCompoundField() *protoMetricsV1.CompoundField {
	h.mu.Lock()
	defer h.mu.Unlock()

	deltas := cloneFloat64Slice(h.bkts.values)
	for idx := range deltas {
		deltas[idx] -= h.lastValues[idx]
	}

	f := &protoMetricsV1.CompoundField{
		Min:            h.bkts.min,
		Max:            h.bkts.max,
		Sum:            h.bkts.totalSum - h.lastTotalSum,
		Count:          h.bkts.totalCount - h.lastTotalCount,
		ExplicitBounds: cloneFloat64Slice(h.bkts.upperBounds),
		Values:         deltas,
	}
	// resets min and max
	h.bkts.min = 0
	h.bkts.max = 0
	// resets total and sum
	h.lastTotalCount = h.bkts.totalCount
	h.lastTotalSum = h.bkts.totalSum
	copy(h.lastValues, h.bkts.values)
	return f
}
