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

package aggregation

import (
	"math"
	"sort"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	larray "github.com/lindb/arrow/pkg/arrow/array"

	seriesmetric "github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

// histogramScalarAgg computes one scalar histogram result per row (instant queries).
type histogramScalarAgg struct {
	params histoParams
}

func (h *histogramScalarAgg) fill(b array.Builder, batch arrow.RecordBatch, numRows int) {
	appendHistoScalars(b, batch, numRows, h.params)
}

// histogramTimeSeriesAgg computes a TimeSeries histogram result per row (range queries).
type histogramTimeSeriesAgg struct {
	params histoParams
}

func (h *histogramTimeSeriesAgg) fill(b array.Builder, batch arrow.RecordBatch, numRows int) {
	appendHistoTimeSeries(b, batch, numRows, h.params)
}

// histoParams captures all parameters needed to compute one histogram aggregation.
type histoParams struct {
	fn       tree.FuncName
	phi      float64
	buckets  []bucketCol
	sumIdx   int
	countIdx int
	hasSum   bool
	hasCount bool
}

// bucketCol pairs a bucket upper bound with the batch column index holding its delta count.
type bucketCol struct {
	bound  float64
	colIdx int
}

// rowFunc builds the aggregation parameter struct.
func rowFunc(fn tree.FuncName, phi float64,
	buckets []bucketCol,
	sumIdx, countIdx int, hasSum, hasCount bool,
) histoParams {
	return histoParams{fn: fn, phi: phi, buckets: buckets,
		sumIdx: sumIdx, countIdx: countIdx, hasSum: hasSum, hasCount: hasCount}
}

// appendHistoScalars computes one scalar value per row (instant queries).
func appendHistoScalars(b array.Builder, batch arrow.RecordBatch, numRows int, p histoParams) {
	for rowIdx := 0; rowIdx < numRows; rowIdx++ {
		appendFloat64(b, computeHistoValue(p, func(ci int) float64 {
			return readFloat64(batch.Column(ci), rowIdx)
		}))
	}
}

// appendHistoTimeSeries computes a TimeSeries per row (range queries).
// Each physical bucket column is a TimeSeries; for every time slot the per-bucket
// delta counts are collected and the aggregation is applied.
func appendHistoTimeSeries(b array.Builder, batch arrow.RecordBatch, numRows int, p histoParams) {
	eb, ok := b.(*array.ExtensionBuilder)
	if !ok {
		// Schema mismatch — output type is not a TimeSeries extension builder.
		// Degrade gracefully rather than panicking.
		(&nullColumnAgg{}).fill(b, batch, numRows)
		return
	}
	tsB := larray.NewTimeSeriesBuilder(eb)

	for rowIdx := 0; rowIdx < numRows; rowIdx++ {
		// Get time metadata from the first bucket column (all share the same range/interval).
		var meta *tsData
		for _, bkt := range p.buckets {
			if ts := readTimeSeries(batch.Column(bkt.colIdx), rowIdx); ts != nil {
				meta = ts
				break
			}
		}
		if meta == nil {
			tsB.AppendNull()
			continue
		}

		n := len(meta.values)
		results := make([]float64, n)
		for slot := 0; slot < n; slot++ {
			results[slot] = computeHistoValue(p, func(ci int) float64 {
				return readTimeSeriesSlot(batch.Column(ci), rowIdx, slot)
			})
		}
		tsB.Append(meta.start, meta.end, meta.interval, results)
	}
}

// computeHistoValue applies the histogram function for a single point in time.
// readCol(colIdx) returns the value for the given column at this time slot.
func computeHistoValue(p histoParams, readCol func(int) float64) float64 {
	switch p.fn {
	case tree.HistogramQuantile:
		bounds := make([]float64, len(p.buckets))
		deltas := make([]float64, len(p.buckets))
		for i, bkt := range p.buckets {
			bounds[i] = bkt.bound
			deltas[i] = readCol(bkt.colIdx)
		}
		return histoQuantile(p.phi, bounds, deltas)
	case tree.HistogramAvg:
		if p.hasSum && p.hasCount {
			sum := readCol(p.sumIdx)
			count := readCol(p.countIdx)
			if count > 0 {
				return sum / count
			}
		}
	case tree.HistogramSum:
		if p.hasSum {
			return readCol(p.sumIdx)
		}
	case tree.HistogramCount:
		if p.hasCount {
			return readCol(p.countIdx)
		}
	}
	return 0
}

// histoQuantile computes the Prometheus-compatible linear-interpolation quantile
// from per-bucket delta counts and their upper bounds (sorted ascending, +Inf last).
//
// This mirrors the algorithm in spi/table/metric/histogram_column.go:computeQuantile
// but operates on plain float64 slices (for Arrow batch processing at the broker).
func histoQuantile(phi float64, bounds, deltaCounts []float64) float64 {
	n := len(bounds)
	if n == 0 {
		return math.NaN()
	}

	// Reject out-of-range quantile values per Prometheus convention.
	if phi < 0 || phi > 1 {
		return math.NaN()
	}

	// Build cumulative counts.
	cumCounts := make([]float64, n)
	var total float64
	for i, c := range deltaCounts {
		total += c
		cumCounts[i] = total
	}
	if total == 0 {
		return math.NaN()
	}

	target := phi * total

	// Binary search: first bucket whose cumulative count >= target.
	idx := sort.Search(n, func(i int) bool { return cumCounts[i] >= target })
	if idx >= n {
		idx = n - 1
	}

	upperBound := bounds[idx]
	var lowerBound float64
	if idx > 0 {
		lowerBound = bounds[idx-1]
	}

	// +Inf bucket: clamp to last finite upper bound.
	if math.IsInf(upperBound, 1) {
		return lowerBound
	}

	bucketCount := deltaCounts[idx]
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

// extractHistoName returns the logical histogram column name from an
// AggregationAssignment.  It first tries the original FunctionCall AST
// (where the column argument is a plain Identifier like "sent_duration"),
// then falls back to the mangled symbol-reference name.
func extractHistoName(agg *plan.AggregationAssignment) string {
	if call, ok := agg.ASTExpression.(*tree.FunctionCall); ok && len(call.Arguments) >= 2 {
		switch a := call.Arguments[1].(type) {
		case *tree.Identifier:
			return a.Value
		case *tree.SymbolReference:
			if name := seriesmetric.HistoNameFromField(a.Name); name != "" {
				return name
			}
			return a.Name
		}
	}
	for _, arg := range agg.Aggregation.Arguments {
		if sr, ok := arg.(*tree.SymbolReference); ok {
			if name := seriesmetric.HistoNameFromField(sr.Name); name != "" {
				return name
			}
			return sr.Name
		}
	}
	return ""
}

// extractHistoPhi returns the quantile phi from the aggregation's literal arguments.
// Defaults to 0.5 (median) when no literal is found.
func extractHistoPhi(args []tree.Expression) float64 {
	for _, arg := range args {
		switch a := arg.(type) {
		case *tree.FloatLiteral:
			return a.Value
		case *tree.LongLiteral:
			return float64(a.Value)
		}
	}
	return 0.5
}
