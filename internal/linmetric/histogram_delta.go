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

	commonMetric "github.com/lindb/common/metric"

	"github.com/lindb/arrow/pkg/model"
	seriesmetric "github.com/lindb/lindb/series/metric"
)

// BoundHistogram is a histogram which has been Bound to a certain metric
// with field-name and metrics, used for non-negative values
//
// a default created bucket will be automatically created,
// however, you can also specify your own buckets.
// Prometheus Histogram's buckets are cumulative where values in each buckets is cumulative,
type BoundHistogram struct {
	bkts           *histogramBuckets
	lastValues     []float64
	lastTotalCount float64
	lastTotalSum   float64
	mu             sync.Mutex
	// name is the logical field name used when decomposing into storage fields.
	// Empty means "unnamed single histogram" — storage uses HistogramSum / __bucket_* convention.
	name string
}

func newBoundHistogram(name string) *BoundHistogram {
	h := &BoundHistogram{
		name: name,
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

func (h *BoundHistogram) afterResetBuckets() {
	h.lastValues = cloneFloat64Slice(h.bkts.values)
	h.lastTotalCount = h.bkts.totalCount
	h.lastTotalSum = h.bkts.totalSum
}

func (h *BoundHistogram) WithExponentBuckets(lower, upper time.Duration, count int) *BoundHistogram {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.bkts.reset(float64(lower.Nanoseconds()/1e6), float64(upper.Nanoseconds()/1e6), count, exponentBucket)
	h.afterResetBuckets()
	return h
}

func (h *BoundHistogram) WithLinearBuckets(lower, upper time.Duration, count int) *BoundHistogram {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.bkts.reset(float64(lower.Nanoseconds()/1e6), float64(upper.Nanoseconds()/1e6), count, linearBucket)
	h.afterResetBuckets()
	return h
}

func (h *BoundHistogram) UpdateDuration(d time.Duration) {
	h.UpdateMilliseconds(float64(d.Nanoseconds() / 1e6))
}

func (h *BoundHistogram) UpdateSince(start time.Time) {
	h.UpdateMilliseconds(float64(time.Since(start).Nanoseconds() / 1e6))
}

func (h *BoundHistogram) UpdateMilliseconds(s float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.bkts.Update(s)
}

func (h *BoundHistogram) UpdateSeconds(s float64) {
	h.UpdateMilliseconds(s * 1000)
}

func (h *BoundHistogram) Update(f func()) {
	start := time.Now()
	f()
	h.UpdateSince(start)
}

func (h *BoundHistogram) marshalToCompoundField(builder *commonMetric.RowBuilder) {
	h.mu.Lock()
	defer h.mu.Unlock()

	deltas := cloneFloat64Slice(h.bkts.values)
	for idx := range deltas {
		deltas[idx] -= h.lastValues[idx]
	}

	_ = builder.AddCompoundFieldMMSC(
		h.bkts.min, h.bkts.max,
		h.bkts.totalSum-h.lastTotalSum,
		h.bkts.totalCount-h.lastTotalCount,
	)
	_ = builder.AddCompoundFieldData(
		deltas,
		cloneFloat64Slice(h.bkts.upperBounds),
	)
	// resets min and max
	h.bkts.min = 0
	h.bkts.max = 0
	// resets total and sum
	h.lastTotalCount = h.bkts.totalCount
	h.lastTotalSum = h.bkts.totalSum
	copy(h.lastValues, h.bkts.values)
}

// AppendFields computes the per-field delta since the last gather and appends each
// histogram component as an individual model.Field into m.
//
// The physical field naming uses the ".__" separator to avoid collision with
// user-defined fields:
//   - <name>.__sum   (AggregationSum)
//   - <name>.__count (AggregationSum  — count is sum-aggregated across replicas)
//   - <name>.__min   (AggregationMin)
//   - <name>.__max   (AggregationMax)
//   - <name>.__bucket_<bound> (AggregationSum — bucket counts are summed across replicas;
//     storage identifies them as HistogramField via IsBucketField name detection)
func (h *BoundHistogram) AppendFields(m *model.Metric) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// compute per-bucket delta counts
	deltas := cloneFloat64Slice(h.bkts.values)
	for idx := range deltas {
		deltas[idx] -= h.lastValues[idx]
	}

	sumDelta := h.bkts.totalSum - h.lastTotalSum
	countDelta := h.bkts.totalCount - h.lastTotalCount

	// stat fields
	m.Fields = append(m.Fields,
		model.Field{Name: seriesmetric.HistoStatFieldName(h.name, seriesmetric.HistoStatSum), Kind: model.AggregationSum, Value: sumDelta},
		model.Field{Name: seriesmetric.HistoStatFieldName(h.name, seriesmetric.HistoStatCount), Kind: model.AggregationSum, Value: countDelta},
		model.Field{Name: seriesmetric.HistoStatFieldName(h.name, seriesmetric.HistoStatMin), Kind: model.AggregationMin, Value: h.bkts.min},
		model.Field{Name: seriesmetric.HistoStatFieldName(h.name, seriesmetric.HistoStatMax), Kind: model.AggregationMax, Value: h.bkts.max},
	)

	// bucket fields — AggregationSum because counts are simply summed across replicas.
	// Storage detects HistogramField type via IsBucketField(name) on the field name.
	for i, bound := range h.bkts.upperBounds {
		name := seriesmetric.BucketFieldName(h.name, bound)
		m.Fields = append(m.Fields, model.Field{
			Name:  name,
			Kind:  model.AggregationSum,
			Value: deltas[i],
		})
	}

	// advance delta state
	h.bkts.min = 0
	h.bkts.max = 0
	h.lastTotalCount = h.bkts.totalCount
	h.lastTotalSum = h.bkts.totalSum
	copy(h.lastValues, h.bkts.values)
}
