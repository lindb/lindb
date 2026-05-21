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

package operator

import (
	"context"
	"math"
	"sort"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"

	seriesmetric "github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type HashAggregationOperator struct {
	node  *plan.AggregationNode
	child Operator

	inbound *Queue
}

func NewHashAggregationOperator(node *plan.AggregationNode, child Operator) Operator {
	return &HashAggregationOperator{
		child:   child,
		node:    node,
		inbound: NewQueue(make(chan arrow.RecordBatch)),
	}
}

// Run processes each incoming batch from the storage layer.
// For histogram aggregation functions (histogram_quantile, histogram_avg, etc.) it
// computes the final scalar (instant query) or TimeSeries (range query) result from the
// physical bucket/stat columns.  All other aggregations are forwarded as-is.
func (h *HashAggregationOperator) Run(ctx context.Context, output chan<- arrow.RecordBatch) {
	for {
		source, ok := h.inbound.Consume(ctx)
		if !ok {
			return
		}
		output <- h.process(source)
	}
}

// process decides whether the batch needs histogram computation or can pass through.
func (h *HashAggregationOperator) process(batch arrow.RecordBatch) arrow.RecordBatch {
	if batch == nil || batch.NumRows() == 0 {
		return batch
	}
	for _, aggAssign := range h.node.Aggregations {
		if tree.IsHistogramFunc(aggAssign.Aggregation.Function) {
			return h.computeHistogramAggregations(batch)
		}
	}
	return batch // no histogram aggregations: pass through
}

// tsData holds the raw values extracted from a single TimeSeries column cell.
type tsData struct {
	start    int64
	end      int64
	interval int64
	values   []float64
}

// computeHistogramAggregations builds a new output record batch by computing each
// histogram aggregation from the physical bucket/stat columns sent by storage.
//
// The input batch comes from buildTableScanOutputSymbols and may contain:
//   - Grouping-key columns (String): tags such as "grpc_service"
//   - Timestamp column: absent from the batch (reducer skips it; time is embedded in TimeSeries)
//   - Physical histogram columns: *larray.Aggregation (instant query) or *larray.TimeSeries (range)
//
// The output batch has one column per symbol in node.GetOutputSymbols(), in the same order,
// so result_set_output's position-based column mapping works correctly.
func (h *HashAggregationOperator) computeHistogramAggregations(batch arrow.RecordBatch) arrow.RecordBatch {
	outputSymbols := h.node.GetOutputSymbols()
	numRows := int(batch.NumRows())

	// Build column-name → batch-column-index lookup.
	inputIdx := make(map[string]int, batch.NumCols())
	for i, f := range batch.Schema().Fields() {
		inputIdx[f.Name] = i
	}

	// Detect whether physical histogram columns are TimeSeries (range query) or scalar.
	isTimeSeries := h.inputIsTimeSeries(batch, inputIdx)

	// Build output Arrow schema.  For histogram result columns, use TimeSeries type when
	// the input is a range query so the downstream result rendering expands the values.
	outFields := make([]arrow.Field, len(outputSymbols))
	for i, sym := range outputSymbols {
		dt := sym.DataType
		if isTimeSeries && !h.isGroupingKey(sym) && !h.isTimestampKey(sym) {
			// Override the plan-declared Sum type with TimeSeries for range results.
			if h.findAggForSymbol(sym) != nil {
				dt = larrow.ExtensionTypes.TimeSeries
			}
		}
		outFields[i] = arrow.Field{Name: sym.Name, Type: dt}
	}

	rb := array.NewRecordBuilder(memory.NewGoAllocator(), arrow.NewSchema(outFields, nil))
	defer rb.Release()

	for symIdx, sym := range outputSymbols {
		b := rb.Field(symIdx)

		// Timestamp grouping key: absent from input batch (embedded in TimeSeries structs).
		if h.isTimestampKey(sym) {
			for i := 0; i < numRows; i++ {
				b.AppendNull()
			}
			continue
		}

		// Non-timestamp grouping key (tag column): copy from input.
		if h.isGroupingKey(sym) {
			colIdx, ok := inputIdx[sym.Name]
			if !ok {
				for i := 0; i < numRows; i++ {
					b.AppendNull()
				}
				continue
			}
			for rowIdx := 0; rowIdx < numRows; rowIdx++ {
				appendColumnValue(b, batch.Column(colIdx), rowIdx)
			}
			continue
		}

		// Aggregation result column: find the matching aggregation definition.
		aggAssign := h.findAggForSymbol(sym)
		if aggAssign == nil || !tree.IsHistogramFunc(aggAssign.Aggregation.Function) {
			for i := 0; i < numRows; i++ {
				b.AppendNull()
			}
			continue
		}

		phi := extractHistoPhi(aggAssign.Aggregation.Arguments)
		histoName := extractHistoName(aggAssign)

		// Collect and sort bucket columns for this histogram.
		var buckets []bucketCol
		for colName, ci := range inputIdx {
			if seriesmetric.IsBucketField(colName) && seriesmetric.HistoNameFromField(colName) == histoName {
				bound, ok := seriesmetric.ParseBucketBound(colName)
				if ok {
					buckets = append(buckets, bucketCol{bound: bound, colIdx: ci})
				}
			}
		}
		sort.Slice(buckets, func(i, j int) bool {
			if math.IsInf(buckets[i].bound, 1) {
				return false
			}
			if math.IsInf(buckets[j].bound, 1) {
				return true
			}
			return buckets[i].bound < buckets[j].bound
		})

		sumIdx, hasSum := inputIdx[seriesmetric.HistoStatFieldName(histoName, seriesmetric.HistoStatSum)]
		countIdx, hasCount := inputIdx[seriesmetric.HistoStatFieldName(histoName, seriesmetric.HistoStatCount)]

		if isTimeSeries {
			h.appendHistoTimeSeries(b, batch, numRows, rowFunc(
				aggAssign.Aggregation.Function, phi,
				buckets, sumIdx, countIdx, hasSum, hasCount))
		} else {
			h.appendHistoScalars(b, batch, numRows, rowFunc(
				aggAssign.Aggregation.Function, phi,
				buckets, sumIdx, countIdx, hasSum, hasCount))
		}
	}

	return rb.NewRecordBatch()
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
func (h *HashAggregationOperator) appendHistoScalars(
	b array.Builder, batch arrow.RecordBatch, numRows int, p histoParams,
) {
	for rowIdx := 0; rowIdx < numRows; rowIdx++ {
		appendFloat64(b, computeHistoValue(p, func(ci int) float64 {
			return readFloat64(batch.Column(ci), rowIdx)
		}))
	}
}

// appendHistoTimeSeries computes a TimeSeries per row (range queries).
// Each physical bucket column is a TimeSeries; for every time slot the per-bucket
// delta counts are collected and the aggregation is applied.
func (h *HashAggregationOperator) appendHistoTimeSeries(
	b array.Builder, batch arrow.RecordBatch, numRows int, p histoParams,
) {
	tsB := larray.NewTimeSeriesBuilder(b.(*array.ExtensionBuilder))

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

// isGroupingKey reports whether sym is one of the AggregationNode's grouping keys.
func (h *HashAggregationOperator) isGroupingKey(sym *plan.Symbol) bool {
	if h.node.GroupingSets == nil {
		return false
	}
	for _, key := range h.node.GroupingSets.GroupingKeys {
		if key.Name == sym.Name {
			return true
		}
	}
	return false
}

// isTimestampKey reports whether sym is the timestamp grouping key.
func (h *HashAggregationOperator) isTimestampKey(sym *plan.Symbol) bool {
	return arrow.TypeEqual(sym.DataType, arrow.FixedWidthTypes.Timestamp_ns)
}

// findAggForSymbol returns the AggregationAssignment whose output Symbol matches sym.
func (h *HashAggregationOperator) findAggForSymbol(sym *plan.Symbol) *plan.AggregationAssignment {
	for _, agg := range h.node.Aggregations {
		if agg.Symbol.Name == sym.Name {
			return agg
		}
	}
	return nil
}

// inputIsTimeSeries checks whether the physical histogram columns in the input batch
// are TimeSeries (range query) vs Aggregation/scalar (instant query).
func (h *HashAggregationOperator) inputIsTimeSeries(batch arrow.RecordBatch, inputIdx map[string]int) bool {
	for colName := range inputIdx {
		if seriesmetric.IsHistogramRelatedField(colName) {
			col := batch.Column(inputIdx[colName])
			if _, ok := col.(*larray.TimeSeries); ok {
				return true
			}
			return false // first histogram column is scalar
		}
	}
	return false
}

// ── helpers ──────────────────────────────────────────────────────────────────

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

// readFloat64 reads the float64 value at rowIdx from col.
// Handles AggregationType (larray.Aggregation) and plain Float64 columns.
func readFloat64(col arrow.Array, rowIdx int) float64 {
	if col == nil || col.IsNull(rowIdx) {
		return 0
	}
	switch a := col.(type) {
	case *larray.Aggregation:
		return a.Value(rowIdx)
	case *array.Float64:
		return a.Value(rowIdx)
	}
	return 0
}

// readTimeSeries extracts the time metadata and values from a TimeSeries column cell.
func readTimeSeries(col arrow.Array, rowIdx int) *tsData {
	ts, ok := col.(*larray.TimeSeries)
	if !ok || ts.IsNull(rowIdx) {
		return nil
	}
	structArr := ts.Storage().(*array.Struct)
	start := structArr.Field(0).(*array.Int64).Value(rowIdx)
	end := structArr.Field(1).(*array.Int64).Value(rowIdx)
	interval := structArr.Field(2).(*array.Int64).Value(rowIdx)
	listArr := structArr.Field(3).(*array.List)
	offsets := listArr.Offsets()
	from, to := int(offsets[rowIdx]), int(offsets[rowIdx+1])
	floats := listArr.ListValues().(*array.Float64)
	values := make([]float64, to-from)
	for i := range values {
		values[i] = floats.Value(from + i)
	}
	return &tsData{start: start, end: end, interval: interval, values: values}
}

// readTimeSeriesSlot reads one time slot's value from a TimeSeries column cell.
// Returns 0 when the column is a scalar (Aggregation type) — used in range queries
// where physical fields may include non-TimeSeries stat columns.
func readTimeSeriesSlot(col arrow.Array, rowIdx, slot int) float64 {
	if col == nil || col.IsNull(rowIdx) {
		return 0
	}
	switch a := col.(type) {
	case *larray.TimeSeries:
		td := readTimeSeries(a, rowIdx)
		if td == nil || slot >= len(td.values) {
			return 0
		}
		return td.values[slot]
	case *larray.Aggregation:
		// Stat columns (sum/count) that didn't get promoted to TimeSeries: use scalar.
		return a.Value(rowIdx)
	case *array.Float64:
		return a.Value(rowIdx)
	}
	return 0
}

// appendFloat64 writes a float64 value to b, which may be a Float64Builder or
// an ExtensionBuilder (for AggregationType output columns).
func appendFloat64(b array.Builder, val float64) {
	switch eb := b.(type) {
	case *array.Float64Builder:
		eb.Append(val)
	case *array.ExtensionBuilder:
		larray.NewAggregationBuilder(eb).Append(val)
	default:
		b.AppendNull()
	}
}

// appendColumnValue copies the value at rowIdx from src into dst.
// Handles String, Float64, and AggregationType column types.
func appendColumnValue(dst array.Builder, src arrow.Array, rowIdx int) {
	if src.IsNull(rowIdx) {
		dst.AppendNull()
		return
	}
	switch a := src.(type) {
	case *array.String:
		dst.(*array.StringBuilder).Append(a.Value(rowIdx))
	case *array.Float64:
		appendFloat64(dst, a.Value(rowIdx))
	case *larray.Aggregation:
		appendFloat64(dst, a.Value(rowIdx))
	default:
		dst.AppendNull()
	}
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

func (h *HashAggregationOperator) GetLayout() []*plan.Symbol {
	return h.node.GetOutputSymbols()
}

func (h *HashAggregationOperator) Children() []Operator {
	return []Operator{h.child}
}

func (h *HashAggregationOperator) GetInbounds() []chan arrow.RecordBatch {
	return []chan arrow.RecordBatch{h.inbound.GetInbound()}
}

func (h *HashAggregationOperator) String() string {
	return "HashAggregationOperator"
}
