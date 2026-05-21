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
	"sort"

	seriesmetric "github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/sql/tree"
)

// noopColumn is a Column placeholder that performs no work.
// It is placed at the positions of "absorbed" histogram physical columns
// (i.e., bucket/stat columns whose data is consumed by histogramColumn instead).
type noopColumn struct{}

func (n *noopColumn) isExemplar() bool                      { return false }
func (n *noopColumn) createStream(_ int, _ int)             {}
func (n *noopColumn) downsampling(_ []*loader)              {}
func (n *noopColumn) aggregate(_ []Result)                  {}
func (n *noopColumn) reset()                                {}

// histogramEntry holds a single bucket field together with its parsed upper-bound value.
// The bound is cached so we avoid re-parsing the field name on every query.
type histogramEntry struct {
	col   *column[float64]
	bound float64 // upper bound parsed from field name; +Inf for the last bucket
}

// histogramColumn is a virtual Column that aggregates multiple physical bucket/stat fields
// into a single output value (quantile, avg, sum, or count).
//
// It is used for histogram_quantile and histogram_avg. For histogram_sum / histogram_count
// a plain column[float64] is sufficient.
//
// The histogramColumn occupies ONE slot in the shared aggregation Result[] slice (offset).
// All sub-columns (bucket, sum, count) use their own Result slots to accumulate data;
// histogramColumn.aggregate reads those and computes the final scalar at offset.
type histogramColumn struct {
	// aggFunc describes what to compute: quantile (with phi), avg, sum, or count
	aggFunc tree.FuncName
	phi     float64 // only meaningful for histogram_quantile

	// offset is the output slot index in the shared Result[] slice
	offset int

	// buckets is the sorted (by upper-bound) list of bucket columns
	buckets []histogramEntry

	// sumCol / countCol are the _sum and _count physical columns
	sumCol   *column[float64]
	countCol *column[float64]
}

// isExemplar always returns false; histograms are never exemplar columns.
func (h *histogramColumn) isExemplar() bool { return false }

// createStream delegates stream creation to every physical sub-column.
func (h *histogramColumn) createStream(numOfFamilies, numOfPoints int) {
	for i := range h.buckets {
		h.buckets[i].col.createStream(numOfFamilies, numOfPoints)
	}
	if h.sumCol != nil {
		h.sumCol.createStream(numOfFamilies, numOfPoints)
	}
	if h.countCol != nil {
		h.countCol.createStream(numOfFamilies, numOfPoints)
	}
}

// downsampling delegates time-windowed rollup to every physical sub-column.
func (h *histogramColumn) downsampling(familyLoaders []*loader) {
	for i := range h.buckets {
		h.buckets[i].col.downsampling(familyLoaders)
	}
	if h.sumCol != nil {
		h.sumCol.downsampling(familyLoaders)
	}
	if h.countCol != nil {
		h.countCol.downsampling(familyLoaders)
	}
}

// aggregate reads rolled-up values from each physical sub-column, computes the
// histogram aggregation result, and writes one value into aggregator[h.offset].
func (h *histogramColumn) aggregate(aggregator []Result) {
	// Let every sub-column deposit its values into their own slots first.
	for i := range h.buckets {
		h.buckets[i].col.aggregate(aggregator)
	}
	if h.sumCol != nil {
		h.sumCol.aggregate(aggregator)
	}
	if h.countCol != nil {
		h.countCol.aggregate(aggregator)
	}

	// Compute the final scalar.
	var scalar float64
	switch h.aggFunc {
	case tree.HistogramSum:
		scalar = h.readScalar(aggregator, h.sumCol)
	case tree.HistogramCount:
		scalar = h.readScalar(aggregator, h.countCol)
	case tree.HistogramAvg:
		sum := h.readScalar(aggregator, h.sumCol)
		cnt := h.readScalar(aggregator, h.countCol)
		if cnt == 0 {
			scalar = math.NaN()
		} else {
			scalar = sum / cnt
		}
	case tree.HistogramQuantile:
		scalar = h.computeQuantile(aggregator)
	default:
		scalar = math.NaN()
	}

	// Write the final scalar into the designated output slot.
	if aggregator[h.offset] == nil {
		aggregator[h.offset] = NewResult[float64](1)
	}
	dst := aggregator[h.offset].(*result[float64])
	dst.array.SetValue(0, scalar)
}

// reset clears all physical sub-column rollup state for reuse across series.
func (h *histogramColumn) reset() {
	for i := range h.buckets {
		h.buckets[i].col.reset()
	}
	if h.sumCol != nil {
		h.sumCol.reset()
	}
	if h.countCol != nil {
		h.countCol.reset()
	}
}

// readScalar returns the single aggregated float64 from a physical sub-column's result slot.
func (h *histogramColumn) readScalar(aggregator []Result, col *column[float64]) float64 {
	if col == nil || len(col.aggs) == 0 {
		return 0
	}
	idx := col.aggs[0].index
	if idx >= len(aggregator) || aggregator[idx] == nil {
		return 0
	}
	r, ok := aggregator[idx].(*result[float64])
	if !ok {
		return 0
	}
	if !r.array.HasValue(0) {
		return 0
	}
	return r.array.GetValue(0)
}

// computeQuantile implements Prometheus-compatible linear interpolation within histogram buckets.
//
// Algorithm:
//  1. Read per-bucket delta counts from the aggregator.
//  2. Build cumulative counts.
//  3. Find the first bucket whose cumulative count >= phi * totalCount.
//  4. Linearly interpolate within that bucket's (lowerBound, upperBound] interval.
//
// Edge cases: empty histogram → NaN; phi=0 → lowerBound of first bucket;
// phi=1 or +Inf bucket → upperBound of last finite bucket.
func (h *histogramColumn) computeQuantile(aggregator []Result) float64 {
	n := len(h.buckets)
	if n == 0 {
		return math.NaN()
	}

	// Read per-bucket delta counts.
	bucketCounts := make([]float64, n)
	for i, entry := range h.buckets {
		if len(entry.col.aggs) == 0 {
			continue
		}
		idx := entry.col.aggs[0].index
		if idx < len(aggregator) && aggregator[idx] != nil {
			if r, ok := aggregator[idx].(*result[float64]); ok && r.array.HasValue(0) {
				bucketCounts[i] = r.array.GetValue(0)
			}
		}
	}

	// Build cumulative counts.
	cumCounts := make([]float64, n)
	var total float64
	for i, c := range bucketCounts {
		total += c
		cumCounts[i] = total
	}
	if total == 0 {
		return math.NaN()
	}

	target := h.phi * total

	// Binary search for the first bucket where cumCounts[i] >= target.
	idx := sort.Search(n, func(i int) bool { return cumCounts[i] >= target })
	if idx >= n {
		idx = n - 1
	}

	upperBound := h.buckets[idx].bound
	var lowerBound float64
	if idx > 0 {
		lowerBound = h.buckets[idx-1].bound
	}

	// For the +Inf bucket, clamp to the last finite upper bound.
	if math.IsInf(upperBound, 1) {
		return lowerBound
	}

	bucketCount := bucketCounts[idx]
	if bucketCount == 0 {
		return lowerBound
	}

	prevCum := float64(0)
	if idx > 0 {
		prevCum = cumCounts[idx-1]
	}

	// Linear interpolation within (lowerBound, upperBound].
	fraction := (target - prevCum) / bucketCount
	return lowerBound + fraction*(upperBound-lowerBound)
}

// histogramColumnBuilder accumulates the physical columns for a single histogram field
// while buildTableScan iterates over output columns.  When seal() is called it returns
// the final histogramColumn ready to replace the placeholder in tableScan.columns.
type histogramColumnBuilder struct {
	aggFunc      tree.FuncName
	phi          float64
	outputOffset int // index in tableScan.columns for the output slot

	buckets  []histogramEntry
	sumCol   *column[float64]
	countCol *column[float64]
}

func (b *histogramColumnBuilder) seal() *histogramColumn {
	// Sort buckets by ascending upper bound so computeQuantile works correctly.
	sort.Slice(b.buckets, func(i, j int) bool {
		// +Inf should always be last
		if math.IsInf(b.buckets[i].bound, 1) {
			return false
		}
		if math.IsInf(b.buckets[j].bound, 1) {
			return true
		}
		return b.buckets[i].bound < b.buckets[j].bound
	})
	return &histogramColumn{
		aggFunc:  b.aggFunc,
		phi:      b.phi,
		offset:   b.outputOffset,
		buckets:  b.buckets,
		sumCol:   b.sumCol,
		countCol: b.countCol,
	}
}

// isBucketField reports whether the field name matches the <histoName>.<__bucket_*> pattern.
// Delegates to seriesmetric.IsBucketField which is the canonical implementation.
func isBucketField(fieldName string) bool { return seriesmetric.IsBucketField(fieldName) }

// parseBucketBound extracts the upper-bound float64 from a <histoName>.<__bucket_bound> field name.
// Delegates to seriesmetric.ParseBucketBound which is the canonical implementation.
func parseBucketBound(fieldName string) (float64, bool) {
	return seriesmetric.ParseBucketBound(fieldName)
}

// histoNameFromField extracts the histogram logical name from a physical field name.
// Delegates to seriesmetric.HistoNameFromField which is the canonical implementation.
func histoNameFromField(fieldName string) string { return seriesmetric.HistoNameFromField(fieldName) }
