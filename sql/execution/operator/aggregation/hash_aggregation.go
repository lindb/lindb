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
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"

	operator "github.com/lindb/lindb/sql/execution/operator"
	seriesmetric "github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/sql/planner/plan"
)

// HashAggregationOperator computes aggregation results at the broker layer.
// For histogram functions (histogram_quantile, histogram_avg, etc.) it derives the
// final scalar or TimeSeries result from the physical bucket/stat columns sent by
// storage.  For standard aggregation functions (sum, count, min, max, first, last)
// the pre-aggregated values from storage are copied through unchanged.
type HashAggregationOperator struct {
	node    *plan.AggregationNode
	child   operator.Operator
	inbound *operator.Queue
}

// NewHashAggregationOperator creates a HashAggregationOperator for the given plan node.
func NewHashAggregationOperator(node *plan.AggregationNode, child operator.Operator) operator.Operator {
	return &HashAggregationOperator{
		child:   child,
		node:    node,
		inbound: operator.NewQueue(make(chan arrow.RecordBatch)),
	}
}

// Run collects all incoming batches from the inbound queue, merges rows with the
// same grouping key across all batches (using the hash table), and emits a single
// consolidated output batch.  This is necessary because multiple storage nodes or
// multiple shards can each return partial results for the same group key.
func (h *HashAggregationOperator) Run(ctx context.Context, output chan<- arrow.RecordBatch) {
	ht := newHashTable()
	for {
		source, ok := h.inbound.Consume(ctx)
		if !ok {
			break
		}
		if source == nil || source.NumRows() == 0 {
			continue
		}
		h.accumulateBatch(ht, source)
	}

	merged := ht.buildMergedBatch()
	if merged == nil || merged.NumRows() == 0 {
		return
	}
	output <- h.computeAggregations(merged)
}

// accumulateBatch adds a single batch from storage into the hash table.
// It resolves the grouping key column indices and delegates to hashTable.accumulate.
func (h *HashAggregationOperator) accumulateBatch(ht *hashTable, batch arrow.RecordBatch) {
	inputIdx := make(map[string]int, batch.NumCols())
	for i, f := range batch.Schema().Fields() {
		inputIdx[f.Name] = i
	}
	isTS := h.inputIsTimeSeries(batch, inputIdx)

	// Collect batch column indices for non-timestamp grouping key symbols.
	var keyColIdxs []int
	var keyColNames []string
	if h.node.GroupingSets != nil {
		for _, sym := range h.node.GroupingSets.GroupingKeys {
			if h.isTimestampKey(sym) {
				continue // timestamp is not a physical batch column
			}
			if ci, ok := inputIdx[sym.Name]; ok {
				keyColIdxs = append(keyColIdxs, ci)
				keyColNames = append(keyColNames, sym.Name)
			}
		}
	}

	ht.accumulate(batch, h.node, keyColIdxs, keyColNames, inputIdx, isTS)
}

// computeAggregations rebuilds the output batch by processing every output symbol.
// It handles grouping-key columns, timestamp columns, histogram aggregations, and
// standard aggregations uniformly through the columnAggregator dispatch mechanism.
//
// The input batch comes from buildTableScanOutputSymbols and may contain:
//   - Grouping-key columns (String): tags such as "grpc_service"
//   - Timestamp column: absent from the batch (reducer skips it; time is embedded in TimeSeries)
//   - Physical histogram columns: *larray.Aggregation (instant query) or *larray.TimeSeries (range)
//   - Standard aggregation columns: *larray.Aggregation or *larray.TimeSeries
//
// The output batch has one column per symbol in node.GetOutputSymbols(), in the same
// order, so result_set_output's position-based column mapping works correctly.
func (h *HashAggregationOperator) computeAggregations(batch arrow.RecordBatch) arrow.RecordBatch {
	outputSymbols := h.node.GetOutputSymbols()
	numRows := int(batch.NumRows())

	// Build column-name → batch-column-index lookup.
	inputIdx := make(map[string]int, batch.NumCols())
	for i, f := range batch.Schema().Fields() {
		inputIdx[f.Name] = i
	}

	// Detect whether aggregation columns carry TimeSeries (range query) or scalar values.
	isTimeSeries := h.inputIsTimeSeries(batch, inputIdx)

	// Build output Arrow schema.
	// For aggregation result columns in range queries, override the plan-declared type
	// with TimeSeries so the downstream result rendering can expand the values correctly.
	outFields := make([]arrow.Field, len(outputSymbols))
	for i, sym := range outputSymbols {
		dt := sym.DataType
		if isTimeSeries && !h.isGroupingKey(sym) && !h.isTimestampKey(sym) {
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
			(&nullColumnAgg{}).fill(b, batch, numRows)
			continue
		}

		// Non-timestamp grouping key (tag column): copy from input.
		if h.isGroupingKey(sym) {
			colIdx, ok := inputIdx[sym.Name]
			if !ok {
				(&nullColumnAgg{}).fill(b, batch, numRows)
				continue
			}
			col := batch.Column(colIdx)
			for rowIdx := 0; rowIdx < numRows; rowIdx++ {
				appendColumnValue(b, col, rowIdx)
			}
			continue
		}

		// Aggregation result column: dispatch by function name via buildColumnAgg.
		aggAssign := h.findAggForSymbol(sym)
		if aggAssign == nil {
			(&nullColumnAgg{}).fill(b, batch, numRows)
			continue
		}
		buildColumnAgg(aggAssign, inputIdx, isTimeSeries).fill(b, batch, numRows)
	}

	return rb.NewRecordBatch()
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

// inputIsTimeSeries reports whether the aggregation columns in the batch carry
// TimeSeries values (range query) rather than scalars (instant query).
// It checks both standard aggregation columns and histogram bucket columns.
//
// NOTE: stat columns (.__sum, .__count) may remain scalar even in range queries,
// so we intentionally skip them and only inspect bucket columns for histogram detection.
func (h *HashAggregationOperator) inputIsTimeSeries(batch arrow.RecordBatch, inputIdx map[string]int) bool {
	// Check standard aggregation output columns first.
	for _, agg := range h.node.Aggregations {
		if ci, ok := inputIdx[agg.Symbol.Name]; ok {
			if _, ok := batch.Column(ci).(*larray.TimeSeries); ok {
				return true
			}
		}
	}
	// Fall back to checking physical histogram bucket columns.
	// Stat columns (.__sum, .__count) are intentionally excluded — they may not be
	// promoted to TimeSeries even in range queries, so only bucket columns give a
	// reliable signal.  We scan all bucket columns rather than returning early on the
	// first non-bucket hit, because Go map iteration order is non-deterministic.
	for colName, ci := range inputIdx {
		if seriesmetric.IsBucketField(colName) {
			if _, ok := batch.Column(ci).(*larray.TimeSeries); ok {
				return true
			}
		}
	}
	return false
}

func (h *HashAggregationOperator) GetLayout() []*plan.Symbol {
	return h.node.GetOutputSymbols()
}

func (h *HashAggregationOperator) Children() []operator.Operator {
	return []operator.Operator{h.child}
}

func (h *HashAggregationOperator) GetInbounds() []chan arrow.RecordBatch {
	return []chan arrow.RecordBatch{h.inbound.GetInbound()}
}

func (h *HashAggregationOperator) String() string {
	return "HashAggregationOperator"
}
