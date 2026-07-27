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
	"fmt"
	"math"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larray "github.com/lindb/arrow/pkg/arrow/array"

	seriesmetric "github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/sql/function"
	_ "github.com/lindb/lindb/sql/function/builtin" // register built-in merge functions
	"github.com/lindb/lindb/sql/planner/plan"
)

// mergeFunc combines an existing accumulated value with an incoming partial value.
type mergeFunc func(existing, incoming float64) float64

var (
	// mergeFuncAdd is the default for histogram bucket/stat columns (additive partial results).
	mergeFuncAdd mergeFunc = func(a, b float64) float64 { return a + b }
	// mergeFuncMin / mergeFuncMax are used only for histogram .__min / .__max stat columns.
	mergeFuncMin mergeFunc = math.Min
	mergeFuncMax mergeFunc = math.Max
)

// colMergeFor resolves the merge function for a given input column by name.
// Standard aggregation output columns delegate to the function registry;
// histogram physical columns use additive merge except for __min/__max stat columns.
func colMergeFor(colName string, node *plan.AggregationNode) mergeFunc {
	// Standard aggregation output symbol → look up merge semantics from registry.
	for _, agg := range node.Aggregations {
		if agg.Symbol.Name == colName {
			return function.DefaultRegistry.MergeFunc(agg.Aggregation.Function)
		}
	}
	// Histogram stat columns: __min/__max use directional merge; __sum/__count are additive.
	histoName := seriesmetric.HistoNameFromField(colName)
	if histoName != "" {
		switch colName {
		case seriesmetric.HistoStatFieldName(histoName, seriesmetric.HistoStatMin):
			return mergeFuncMin
		case seriesmetric.HistoStatFieldName(histoName, seriesmetric.HistoStatMax):
			return mergeFuncMax
		}
	}
	// Bucket columns and __sum/__count stats → additive.
	return mergeFuncAdd
}

// ── per-group accumulation state ──────────────────────────────────────────────

// colState holds the accumulated values for one non-key column in one group.
type colState struct {
	scalar float64   // instant query or scalar stat column in range query
	vector []float64 // range query: per-time-slot values; nil means instant query
	tsMeta *tsData   // time metadata (start, end, interval) from the first TimeSeries seen
	init   bool      // whether at least one value has been accumulated
}

// groupState holds the accumulated state for one unique combination of grouping key values.
type groupState struct {
	keyValues []string             // grouping key column values, in the same order as keyColNames
	keyIsNull []bool               // true when the corresponding keyValues entry is NULL
	cols      map[string]*colState // non-key column name → accumulated state
}

// ── hash table ──────────────────────────────────────────────────────────────────

// hashTable groups incoming rows by their grouping key values and merges aggregation columns.
// It is designed to accumulate rows from multiple Arrow RecordBatches (e.g. results from
// multiple storage nodes) before producing a single consolidated output batch.
type hashTable struct {
	order       []string // insertion-order group keys for deterministic output
	groups      map[string]*groupState
	inputSchema *arrow.Schema // schema captured from the first batch seen (nil until first batch)
	keyColNames []string      // column names of the grouping key columns
	isTS        bool          // whether input has TimeSeries columns (range query)
}

func newHashTable() *hashTable {
	return &hashTable{groups: make(map[string]*groupState)}
}

// buildGroupKey serialises the grouping-column values at rowIdx into a stable string.
// A NUL byte (0x00) separates fields; a SOH byte (0x01) represents a null field value.
func buildGroupKey(batch arrow.RecordBatch, keyColIdxs []int, rowIdx int) string {
	if len(keyColIdxs) == 0 {
		return ""
	}
	parts := make([]string, len(keyColIdxs))
	for i, ci := range keyColIdxs {
		col := batch.Column(ci)
		if col.IsNull(rowIdx) {
			parts[i] = "\x01"
		} else if s, ok := col.(*array.String); ok {
			parts[i] = s.Value(rowIdx)
		} else {
			parts[i] = fmt.Sprintf("?%d", ci)
		}
	}
	return strings.Join(parts, "\x00")
}

// accumulate merges all rows of batch into the hash table.
//
//   - keyColIdxs: batch column indices for the grouping key symbols, in symbol order.
//   - keyColNames: column names corresponding to keyColIdxs.
//   - inputIdx: full column-name → batch-column-index map for all columns.
//   - isTS: true when input columns are TimeSeries (range query), false for scalar (instant query).
func (t *hashTable) accumulate(
	batch arrow.RecordBatch,
	node *plan.AggregationNode,
	keyColIdxs []int,
	keyColNames []string,
	inputIdx map[string]int,
	isTS bool,
) {
	// Capture schema and key metadata from the first batch.
	if t.inputSchema == nil {
		t.inputSchema = batch.Schema()
		t.keyColNames = keyColNames
		t.isTS = isTS
	}

	numRows := int(batch.NumRows())
	for rowIdx := 0; rowIdx < numRows; rowIdx++ {
		gk := buildGroupKey(batch, keyColIdxs, rowIdx)

		state, exists := t.groups[gk]
		if !exists {
			keyValues := make([]string, len(keyColIdxs))
			keyIsNull := make([]bool, len(keyColIdxs))
			for i, ci := range keyColIdxs {
				col := batch.Column(ci)
				if col.IsNull(rowIdx) {
					keyIsNull[i] = true
				} else if s, ok := col.(*array.String); ok {
					keyValues[i] = s.Value(rowIdx)
				}
			}
			state = &groupState{
				keyValues: keyValues,
				keyIsNull: keyIsNull,
				cols:      make(map[string]*colState),
			}
			t.order = append(t.order, gk)
			t.groups[gk] = state
		}

		// Accumulate non-key columns using the appropriate merge function.
		for colName, ci := range inputIdx {
			if t.isKeyCol(colName) {
				continue
			}
			cs := state.cols[colName]
			if cs == nil {
				cs = &colState{}
				state.cols[colName] = cs
			}
			mergeFn := colMergeFor(colName, node)
			col := batch.Column(ci)
			if isTS {
				t.mergeTS(cs, col, rowIdx, mergeFn)
			} else {
				t.mergeScalar(cs, col, rowIdx, mergeFn)
			}
		}
	}
}

// isKeyCol reports whether colName is one of the stored grouping key column names.
func (t *hashTable) isKeyCol(colName string) bool {
	for _, kn := range t.keyColNames {
		if kn == colName {
			return true
		}
	}
	return false
}

// mergeScalar accumulates a scalar column value into cs.
func (t *hashTable) mergeScalar(cs *colState, col arrow.Array, rowIdx int, fn mergeFunc) {
	val := readFloat64(col, rowIdx)
	if cs.init {
		cs.scalar = fn(cs.scalar, val)
	} else {
		cs.scalar = val
		cs.init = true
	}
}

// mergeTS accumulates a TimeSeries column value into cs element-wise.
// When the incoming vector is longer than the accumulated one, the extra slots are
// appended (instead of being silently dropped).  Scalar stat columns that were not
// promoted to TimeSeries are accumulated into cs.scalar.
func (t *hashTable) mergeTS(cs *colState, col arrow.Array, rowIdx int, fn mergeFunc) {
	switch a := col.(type) {
	case *larray.TimeSeries:
		td := readTimeSeries(a, rowIdx)
		if td == nil {
			return
		}
		if !cs.init {
			cs.tsMeta = td
			cs.vector = append([]float64(nil), td.values...)
			cs.init = true
			return
		}
		// Merge element-wise; extend cs.vector when the new batch has more slots.
		if len(td.values) > len(cs.vector) {
			extra := make([]float64, len(td.values)-len(cs.vector))
			cs.vector = append(cs.vector, extra...)
		}
		for i, v := range td.values {
			cs.vector[i] = fn(cs.vector[i], v)
		}
	default:
		// Scalar stat column (larray.Aggregation) or plain float that stayed scalar
		// even in a range query.  Accumulate into the scalar slot.
		val := readFloat64(col, rowIdx)
		if cs.init {
			cs.scalar = fn(cs.scalar, val)
		} else {
			cs.scalar = val
			cs.init = true
		}
	}
}

// buildMergedBatch creates a consolidated Arrow RecordBatch from all accumulated groups.
// The output schema matches the original input batches, making it suitable as direct
// input to computeAggregations.
func (t *hashTable) buildMergedBatch() arrow.RecordBatch {
	if len(t.order) == 0 || t.inputSchema == nil {
		return nil
	}

	rb := array.NewRecordBuilder(memory.NewGoAllocator(), t.inputSchema)
	defer rb.Release()

	for fieldIdx, field := range t.inputSchema.Fields() {
		b := rb.Field(fieldIdx)
		colName := field.Name

		if t.isKeyCol(colName) {
			// Grouping key column: preserve NULL vs empty-string distinction.
			sb := b.(*array.StringBuilder)
			keyIdx := t.keyColIdx(colName)
			for _, gk := range t.order {
				state := t.groups[gk]
				if keyIdx < 0 || keyIdx >= len(state.keyValues) || state.keyIsNull[keyIdx] {
					sb.AppendNull()
				} else {
					sb.Append(state.keyValues[keyIdx])
				}
			}
			continue
		}

		// Aggregation / histogram physical column: write accumulated value.
		for _, gk := range t.order {
			state := t.groups[gk]
			cs := state.cols[colName]
			if cs == nil || !cs.init {
				b.AppendNull()
				continue
			}
			if cs.vector != nil && cs.tsMeta != nil {
				// TimeSeries column (range query).
				eb, ok := b.(*array.ExtensionBuilder)
				if !ok {
					b.AppendNull()
					continue
				}
				tsB := larray.NewTimeSeriesBuilder(eb)
				tsB.Append(cs.tsMeta.start, cs.tsMeta.end, cs.tsMeta.interval, cs.vector)
			} else {
				// Scalar column (instant query or unpromoted stat column).
				appendFloat64(b, cs.scalar)
			}
		}
	}

	return rb.NewRecordBatch()
}

// keyColIdx returns the position of colName in t.keyColNames, or -1 if not found.
func (t *hashTable) keyColIdx(colName string) int {
	for i, kn := range t.keyColNames {
		if kn == colName {
			return i
		}
	}
	return -1
}
