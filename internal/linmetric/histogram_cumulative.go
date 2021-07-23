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

// BoundCumulativeHistogram is a histogram which has been Bound to a certain metric
// with field-name and metrics, used for non-negative values
//
// a default created bucket will be automatically created,
// however, you can also specify your own buckets.
// Notice, cumulative means that value in each buckets is cumulative,
// value in different buckets are not cumulative.
type BoundCumulativeHistogram struct {
	mu   sync.Mutex
	bkts *histogramBuckets
}

func newCumulativeHistogram() *BoundCumulativeHistogram {
	return &BoundCumulativeHistogram{
		bkts: newHistogramBuckets(
			defaultMinBucketUpperBound,
			defaultMaxBucketUpperBound,
			defaultBucketCount,
			exponentBucket,
		),
	}
}

func (h *BoundCumulativeHistogram) WithExponentBuckets(lower, upper time.Duration, count int) *BoundCumulativeHistogram {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.bkts.reset(float64(lower.Nanoseconds()/1e6), float64(upper.Nanoseconds()/1e6), count, exponentBucket)
	return h
}

func (h *BoundCumulativeHistogram) WithLinearBuckets(lower, upper time.Duration, count int) *BoundCumulativeHistogram {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.bkts.reset(float64(lower.Nanoseconds()/1e6), float64(upper.Nanoseconds()/1e6), count, linearBucket)
	return h
}

func (h *BoundCumulativeHistogram) UpdateDuration(d time.Duration) {
	h.UpdateMilliseconds(float64(d.Nanoseconds() / 1e6))
}

func (h *BoundCumulativeHistogram) UpdateSince(start time.Time) {
	h.UpdateMilliseconds(float64(time.Since(start).Nanoseconds() / 1e6))
}

func (h *BoundCumulativeHistogram) UpdateMilliseconds(s float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.bkts.Update(s)
}

func (h *BoundCumulativeHistogram) UpdateSeconds(s float64) {
	h.UpdateMilliseconds(s * 1000)
}

func (h *BoundCumulativeHistogram) Update(f func()) {
	start := time.Now()
	f()
	h.UpdateSince(start)
}

func (h *BoundCumulativeHistogram) marshalToCompoundField() *protoMetricsV1.CompoundField {
	h.mu.Lock()
	defer h.mu.Unlock()

	f := &protoMetricsV1.CompoundField{
		Type:           protoMetricsV1.CompoundFieldType_CUMULATIVE_HISTOGRAM,
		Min:            h.bkts.min,
		Max:            h.bkts.max,
		Sum:            h.bkts.totalSum,
		Count:          h.bkts.totalCount,
		ExplicitBounds: cloneFloat64Slice(h.bkts.upperBounds),
		Values:         cloneFloat64Slice(h.bkts.values),
	}
	// resets min and max
	h.bkts.min = 0
	h.bkts.max = 0
	return f
}
