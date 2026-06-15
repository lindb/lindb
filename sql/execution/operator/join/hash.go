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

package join

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/samber/lo"

	operator "github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
)

// batchFlushThreshold is the maximum number of rows accumulated in the output
// RecordBuilder before flushing to the output channel.  Emitting in chunks
// avoids materialising the entire join result in memory for large inputs.
const batchFlushThreshold = 4096

// buildRowID uniquely identifies a row within the left (build) side buffer.
// batchIdx is the index into leftBatches; rowIdx is the row within that batch.
type buildRowID struct {
	batchIdx, rowIdx int
}

// HashJoinOperator implements a classic two-phase hash join:
//
//	Phase 1 – Buffer: drain both inbound channels concurrently.
//	Phase 2 – Build:  create a hash map from left (build) rows keyed on equi-join columns.
//	Phase 3 – Probe:  for each right row look up matching build rows, apply any residual filter,
//	                  and emit joined output rows.
//	Phase 4 – Emit unmatched: for LEFT / FULL / RIGHT joins, output rows that had no match
//	                           with NULL-padding on the absent side.
type HashJoinOperator struct {
	node         *plan.JoinNode
	left, right  operator.Operator
	filter       expression.Expression // compiled non-equality ON filter; nil if none

	leftScope, rightScope []*plan.Symbol
	leftKeys, rightKeys   []int // column indices within the respective scope

	leftInbound, rightInbound *operator.Queue
}

// NewHashJoinOperator creates a HashJoinOperator for the given logical JoinNode.
func NewHashJoinOperator(node *plan.JoinNode, left, right operator.Operator) operator.Operator {
	return &HashJoinOperator{
		node:         node,
		left:         left,
		right:        right,
		leftScope:    node.Left.GetOutputSymbols(),
		rightScope:   node.Right.GetOutputSymbols(),
		leftInbound:  operator.NewQueue(make(chan arrow.RecordBatch)),
		rightInbound: operator.NewQueue(make(chan arrow.RecordBatch)),
	}
}

// ── operator.Operator interface ──────────────────────────────────────────────

// Children implements operator.Operator.
func (h *HashJoinOperator) Children() []operator.Operator {
	return []operator.Operator{h.left, h.right}
}

// GetInbounds implements operator.Operator.
func (h *HashJoinOperator) GetInbounds() []chan arrow.RecordBatch {
	return []chan arrow.RecordBatch{h.leftInbound.GetInbound(), h.rightInbound.GetInbound()}
}

// GetLayout implements operator.Operator.
func (h *HashJoinOperator) GetLayout() []*plan.Symbol {
	return h.node.GetOutputSymbols()
}

// String implements operator.Operator.
func (h *HashJoinOperator) String() string {
	return "HashJoinOperator"
}

// ── Run ──────────────────────────────────────────────────────────────────────

// Run executes the four-phase hash join and writes output batches to output.
func (h *HashJoinOperator) Run(ctx context.Context, output chan<- arrow.RecordBatch) {
	h.prepare(ctx)

	// ── Phase 1: buffer both sides concurrently ──────────────────────────────
	// The pipeline uses unbuffered channels; both sides MUST be consumed in
	// parallel to prevent deadlock when the upstream operator tries to send.
	// Each goroutine writes exclusively to its own slice — no mutex needed.
	var (
		leftBatches  []arrow.RecordBatch
		rightBatches []arrow.RecordBatch
		wg           sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			b, ok := h.leftInbound.Consume(ctx)
			if !ok {
				break
			}
			if b != nil && b.NumRows() > 0 {
				leftBatches = append(leftBatches, b)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			b, ok := h.rightInbound.Consume(ctx)
			if !ok {
				break
			}
			if b != nil && b.NumRows() > 0 {
				rightBatches = append(rightBatches, b)
			}
		}
	}()
	wg.Wait()

	// ── Phase 2: build hash map from left (build) side ───────────────────────
	isCross := h.node.IsCrossJoin()
	// A theta join has no equi-criteria but carries a non-equality ON filter.
	// Like CROSS JOIN, every left row is a candidate for every right row; the
	// filter is applied per pair to accept or reject rows.
	isThetaJoin := !isCross && len(h.node.Criteria) == 0 && h.node.Filter != nil
	isLeft := h.node.Type == plan.Left || h.node.Type == plan.Full
	isRight := h.node.Type == plan.Right || h.node.Type == plan.Full

	hashMap := make(map[string][]buildRowID)
	if !isCross && !isThetaJoin {
		for bIdx, batch := range leftBatches {
			for row := range int(batch.NumRows()) {
				key := buildJoinKey(batch, h.leftKeys, row)
				if key == "" {
					continue // NULL key never matches (SQL NULL ≠ NULL)
				}
				hashMap[key] = append(hashMap[key], buildRowID{bIdx, row})
			}
		}
	}

	// ── Phase 3: probe right side, emit joined rows ──────────────────────────
	outputSchema := h.buildOutputSchema()
	rb := array.NewRecordBuilder(memory.NewGoAllocator(), outputSchema)
	defer rb.Release()

	// matched tracks build-side rows that joined with at least one probe row (LEFT/FULL).
	matched := make(map[buildRowID]bool)

	// For CROSS or theta joins every build-side row is a candidate for every probe
	// row.  Pre-compute the candidate list once to avoid O(M×N) slice allocations.
	var allLeftCandidates []buildRowID
	if isCross || isThetaJoin {
		for bIdx, lb := range leftBatches {
			for lRow := range int(lb.NumRows()) {
				allLeftCandidates = append(allLeftCandidates, buildRowID{bIdx, lRow})
			}
		}
	}

	for _, rightBatch := range rightBatches {
		for rRow := range int(rightBatch.NumRows()) {
			var candidates []buildRowID

			if isCross || isThetaJoin {
				// CROSS / theta join: every left row pairs with every right row.
				candidates = allLeftCandidates
			} else {
				key := buildJoinKey(rightBatch, h.rightKeys, rRow)
				if key != "" {
					candidates = hashMap[key]
				}
			}

			foundMatch := false
			for _, id := range candidates {
				lb := leftBatches[id.batchIdx]
				lRow := id.rowIdx

				// Apply residual non-equality ON filter on the candidate joined row.
				if h.filter != nil {
					oneRow := h.buildOneRow(outputSchema, lb, lRow, rightBatch, rRow)
					passes, err := h.applyFilter(oneRow)
					oneRow.Release() // release single-row batch built for filter evaluation
					if err != nil {
						panic(fmt.Sprintf("HashJoinOperator: filter eval: %v", err))
					}
					if !passes {
						continue
					}
				}
				h.appendJoinedRow(rb, lb, lRow, rightBatch, rRow)
				matched[id] = true
				foundMatch = true
			}

			// RIGHT or FULL: emit probe row padded with NULLs for the build side.
			if isRight && !foundMatch {
				h.appendRightUnmatched(rb, rightBatch, rRow)
			}

			// Flush to output channel periodically to bound memory usage for large joins.
			if outputSchema.NumFields() > 0 && rb.Field(0).Len() >= batchFlushThreshold {
				output <- rb.NewRecordBatch()
			}
		}
	}

	// ── Phase 4: emit unmatched build (left) rows for LEFT / FULL joins ──────
	if isLeft {
		for bIdx, lb := range leftBatches {
			for lRow := range int(lb.NumRows()) {
				if !matched[buildRowID{bIdx, lRow}] {
					h.appendLeftUnmatched(rb, lb, lRow)
					// Flush periodically during Phase 4 as well.
					if outputSchema.NumFields() > 0 && rb.Field(0).Len() >= batchFlushThreshold {
						output <- rb.NewRecordBatch()
					}
				}
			}
		}
	}

	// Flush any rows remaining below the threshold.
	result := rb.NewRecordBatch()
	if result.NumRows() > 0 {
		output <- result
	}
}

// ── Preparation ──────────────────────────────────────────────────────────────

// prepare resolves join-key column indices in both scopes and compiles the
// residual ON filter (if the plan node carries one).
func (h *HashJoinOperator) prepare(ctx context.Context) {
	criteria := h.node.Criteria
	h.leftKeys = make([]int, len(criteria))
	h.rightKeys = make([]int, len(criteria))

	for idx, c := range criteria {
		_, ki, ok := lo.FindIndexOf(h.leftScope, func(s *plan.Symbol) bool { return s.Name == c.Left.Name })
		if !ok {
			panic(fmt.Sprintf("HashJoinOperator: left key symbol %q not found in left scope", c.Left.Name))
		}
		h.leftKeys[idx] = ki

		_, ki, ok = lo.FindIndexOf(h.rightScope, func(s *plan.Symbol) bool { return s.Name == c.Right.Name })
		if !ok {
			panic(fmt.Sprintf("HashJoinOperator: right key symbol %q not found in right scope", c.Right.Name))
		}
		h.rightKeys[idx] = ki
	}

	if h.node.Filter != nil {
		h.filter = expression.Rewrite(&expression.RewriteContext{
			SourceLayout: h.node.GetOutputSymbols(),
			EvalContext:  expression.NewEvalContext(ctx),
		}, h.node.Filter)
	}
}

// ── Key building ──────────────────────────────────────────────────────────────

// buildJoinKey serialises the key columns for one row into a canonical string.
// Returns "" when any key column is NULL (excluded from the hash map — SQL NULL ≠ NULL).
//
// Each column value is length-prefixed to prevent key collisions when string
// values contain the separator character or share a common prefix.
// Example encoding: columns ["abc", "d"] → "3:abc1:d" (unambiguous).
func buildJoinKey(batch arrow.RecordBatch, colIdxs []int, row int) string {
	if len(colIdxs) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, ci := range colIdxs {
		col := batch.Column(ci)
		if col.IsNull(row) {
			return "" // NULL key never participates in an equi-join
		}
		v := columnStringValue(col, row)
		// Length-prefix each column value: "<n>:<value>" — guarantees that no two
		// distinct multi-column key tuples produce the same encoded string.
		fmt.Fprintf(&sb, "%d:%s", len(v), v)
	}
	return sb.String()
}

// columnStringValue returns a string representation of col[row] for hash-map keying.
func columnStringValue(col arrow.Array, row int) string {
	switch c := col.(type) {
	case *array.Int64:
		return fmt.Sprintf("%d", c.Value(row))
	case *array.Float64:
		v := c.Value(row)
		if v == 0 {
			// Normalise IEEE-754 negative zero to "0" so that -0.0 and +0.0 join
			// to the same hash-map bucket (they are equal under SQL semantics).
			return "0"
		}
		return fmt.Sprintf("%g", v)
	case *array.String:
		return c.Value(row)
	case *array.Boolean:
		if c.Value(row) {
			return "1"
		}
		return "0"
	default:
		return fmt.Sprintf("%T@%d", col, row)
	}
}

// ── Output schema ─────────────────────────────────────────────────────────────

// buildOutputSchema returns the combined Arrow schema: left scope columns followed
// by right scope columns.  Symbols with no declared DataType default to Int64.
func (h *HashJoinOperator) buildOutputSchema() *arrow.Schema {
	fields := make([]arrow.Field, 0, len(h.leftScope)+len(h.rightScope))
	for _, sym := range h.leftScope {
		dt := sym.DataType
		if dt == nil {
			dt = arrow.PrimitiveTypes.Int64
		}
		fields = append(fields, arrow.Field{Name: sym.Name, Type: dt})
	}
	for _, sym := range h.rightScope {
		dt := sym.DataType
		if dt == nil {
			dt = arrow.PrimitiveTypes.Int64
		}
		fields = append(fields, arrow.Field{Name: sym.Name, Type: dt})
	}
	return arrow.NewSchema(fields, nil)
}

// ── Row append helpers ─────────────────────────────────────────────────────────

// appendJoinedRow appends leftBatch[lRow] || rightBatch[rRow] into rb.
func (h *HashJoinOperator) appendJoinedRow(rb *array.RecordBuilder,
	leftBatch arrow.RecordBatch, lRow int,
	rightBatch arrow.RecordBatch, rRow int,
) {
	for i := range len(h.leftScope) {
		appendColValue(rb.Field(i), leftBatch.Column(i), lRow)
	}
	offset := len(h.leftScope)
	for i := range len(h.rightScope) {
		appendColValue(rb.Field(offset+i), rightBatch.Column(i), rRow)
	}
}

// appendLeftUnmatched appends leftBatch[lRow] into rb, NULL-padding the right side.
func (h *HashJoinOperator) appendLeftUnmatched(rb *array.RecordBuilder,
	leftBatch arrow.RecordBatch, lRow int,
) {
	for i := range len(h.leftScope) {
		appendColValue(rb.Field(i), leftBatch.Column(i), lRow)
	}
	offset := len(h.leftScope)
	for i := range len(h.rightScope) {
		rb.Field(offset + i).AppendNull()
	}
}

// appendRightUnmatched appends rightBatch[rRow] into rb, NULL-padding the left side.
func (h *HashJoinOperator) appendRightUnmatched(rb *array.RecordBuilder,
	rightBatch arrow.RecordBatch, rRow int,
) {
	for i := range len(h.leftScope) {
		rb.Field(i).AppendNull()
	}
	offset := len(h.leftScope)
	for i := range len(h.rightScope) {
		appendColValue(rb.Field(offset+i), rightBatch.Column(i), rRow)
	}
}

// appendColValue copies col[row] into the Arrow builder dst.
// Handles int64, float64, string, and boolean columns; unknown types append NULL.
func appendColValue(dst array.Builder, col arrow.Array, row int) {
	if col.IsNull(row) {
		dst.AppendNull()
		return
	}
	switch c := col.(type) {
	case *array.Int64:
		dst.(*array.Int64Builder).Append(c.Value(row))
	case *array.Float64:
		dst.(*array.Float64Builder).Append(c.Value(row))
	case *array.String:
		dst.(*array.StringBuilder).Append(c.Value(row))
	case *array.Boolean:
		dst.(*array.BooleanBuilder).Append(c.Value(row))
	default:
		dst.AppendNull()
	}
}

// ── Residual filter ───────────────────────────────────────────────────────────

// buildOneRow constructs a single-row RecordBatch from one candidate joined pair.
// This is used to evaluate the non-equality ON filter row-by-row.
// The caller is responsible for calling Release() on the returned batch.
func (h *HashJoinOperator) buildOneRow(schema *arrow.Schema,
	leftBatch arrow.RecordBatch, lRow int,
	rightBatch arrow.RecordBatch, rRow int,
) arrow.RecordBatch {
	rb := array.NewRecordBuilder(memory.NewGoAllocator(), schema)
	defer rb.Release()
	for i := range len(h.leftScope) {
		appendColValue(rb.Field(i), leftBatch.Column(i), lRow)
	}
	offset := len(h.leftScope)
	for i := range len(h.rightScope) {
		appendColValue(rb.Field(offset+i), rightBatch.Column(i), rRow)
	}
	return rb.NewRecordBatch()
}

// applyFilter evaluates the compiled residual ON filter against a single-row batch.
// Returns (true, nil) when the row passes, (false, nil) when it is rejected.
func (h *HashJoinOperator) applyFilter(record arrow.RecordBatch) (bool, error) {
	predArr, err := h.filter.Eval(record)
	if err != nil {
		return false, err
	}
	boolArr, ok := predArr.(*array.Boolean)
	if !ok {
		return false, fmt.Errorf("join filter must return boolean array, got %T", predArr)
	}
	if boolArr.Len() == 0 || boolArr.IsNull(0) {
		return false, nil
	}
	return boolArr.Value(0), nil
}
