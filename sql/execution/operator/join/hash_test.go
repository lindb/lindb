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
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	operator "github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

// ── test stubs ───────────────────────────────────────────────────────────────

// stubPlanNode is a minimal plan.PlanNode stub that returns a fixed symbol list.
// It satisfies the plan.PlanNode interface without requiring real plan infrastructure.
type stubPlanNode struct {
	symbols []*plan.Symbol
}

func (s *stubPlanNode) GetNodeID() plan.PlanNodeID                      { return 0 }
func (s *stubPlanNode) GetSources() []plan.PlanNode                     { return nil }
func (s *stubPlanNode) GetOutputSymbols() []*plan.Symbol                { return s.symbols }
func (s *stubPlanNode) ReplaceChildren(_ []plan.PlanNode) plan.PlanNode { return s }
func (s *stubPlanNode) Accept(_ any, _ plan.Visitor) any                { return nil }

// ── column builders ──────────────────────────────────────────────────────────

func testInt64Col(vals []int64) arrow.Array {
	b := array.NewInt64Builder(memory.DefaultAllocator)
	defer b.Release()
	b.AppendValues(vals, nil)
	return b.NewArray()
}

func testFloat64Col(vals []float64) arrow.Array {
	b := array.NewFloat64Builder(memory.DefaultAllocator)
	defer b.Release()
	b.AppendValues(vals, nil)
	return b.NewArray()
}

func testStringCol(vals []string) arrow.Array {
	b := array.NewStringBuilder(memory.DefaultAllocator)
	defer b.Release()
	b.AppendValues(vals, nil)
	return b.NewArray()
}

// testNullableInt64Col builds an int64 column where valid[i]=false means NULL at row i.
func testNullableInt64Col(vals []int64, valid []bool) arrow.Array {
	b := array.NewInt64Builder(memory.DefaultAllocator)
	defer b.Release()
	b.AppendValues(vals, valid)
	return b.NewArray()
}

// ── record batch builder ──────────────────────────────────────────────────────

func makeBatch(schema *arrow.Schema, cols []arrow.Array) arrow.RecordBatch {
	return array.NewRecordBatch(schema, cols, int64(cols[0].Len()))
}

// ── symbol helpers ────────────────────────────────────────────────────────────

func symI64(name string) *plan.Symbol {
	return &plan.Symbol{Name: name, DataType: arrow.PrimitiveTypes.Int64}
}

func symF64(name string) *plan.Symbol {
	return &plan.Symbol{Name: name, DataType: arrow.PrimitiveTypes.Float64}
}

func symStr(name string) *plan.Symbol {
	return &plan.Symbol{Name: name, DataType: arrow.BinaryTypes.String}
}

// ── join node builder ─────────────────────────────────────────────────────────

// makeJoinNode creates a JoinNode with a single equi-join criterion between leftKey and rightKey.
func makeJoinNode(jt plan.JoinType, leftSymbols, rightSymbols []*plan.Symbol, leftKey, rightKey string) *plan.JoinNode {
	var leftSym, rightSym *plan.Symbol
	for _, s := range leftSymbols {
		if s.Name == leftKey {
			leftSym = s
		}
	}
	for _, s := range rightSymbols {
		if s.Name == rightKey {
			rightSym = s
		}
	}
	return &plan.JoinNode{
		Type:     jt,
		Left:     &stubPlanNode{symbols: leftSymbols},
		Right:    &stubPlanNode{symbols: rightSymbols},
		Criteria: []*plan.EqualJoinCriteria{{Left: leftSym, Right: rightSym}},
	}
}

// ── test runner ───────────────────────────────────────────────────────────────

// runJoin feeds batches into both inbound channels concurrently, then calls op.Run
// and collects all output batches.  The output channel is buffered so Run's final
// emit never blocks after the probe phase completes.
//
// Both sender goroutines close their channel on completion, which signals EOF to
// the Queue.Consume loop inside Run, allowing wg.Wait() to unblock.
func runJoin(op operator.Operator, leftBatches, rightBatches []arrow.RecordBatch) []arrow.RecordBatch {
	inbounds := op.GetInbounds()
	leftCh, rightCh := inbounds[0], inbounds[1]

	output := make(chan arrow.RecordBatch, 128)

	go func() {
		defer close(leftCh)
		for _, b := range leftBatches {
			leftCh <- b
		}
	}()
	go func() {
		defer close(rightCh)
		for _, b := range rightBatches {
			rightCh <- b
		}
	}()

	op.Run(context.Background(), output)
	close(output)

	var results []arrow.RecordBatch
	for b := range output {
		results = append(results, b)
	}
	return results
}

// totalRows sums NumRows across all result batches.
func totalRows(batches []arrow.RecordBatch) int64 {
	var n int64
	for _, b := range batches {
		n += b.NumRows()
	}
	return n
}

// ── tests ─────────────────────────────────────────────────────────────────────

// TestHashJoin_Inner_BasicMatch verifies a simple equi-join returns only the matching rows.
func TestHashJoin_Inner_BasicMatch(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}
	node := makeJoinNode(plan.Inner, lSyms, rSyms, "l_id", "r_id")

	lSchema := arrow.NewSchema([]arrow.Field{{Name: "l_id", Type: arrow.PrimitiveTypes.Int64}}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{{Name: "r_id", Type: arrow.PrimitiveTypes.Int64}}, nil)

	// Left: 1, 2, 3 — Right: 2, 3, 4 → matches on 2 and 3 only.
	leftBatches := []arrow.RecordBatch{makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{1, 2, 3})})}
	rightBatches := []arrow.RecordBatch{makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{2, 3, 4})})}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	assert.Equal(t, int64(2), totalRows(results))
}

// TestHashJoin_Inner_NoMatch verifies that no rows are emitted when nothing matches.
func TestHashJoin_Inner_NoMatch(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}
	node := makeJoinNode(plan.Inner, lSyms, rSyms, "l_id", "r_id")

	lSchema := arrow.NewSchema([]arrow.Field{{Name: "l_id", Type: arrow.PrimitiveTypes.Int64}}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{{Name: "r_id", Type: arrow.PrimitiveTypes.Int64}}, nil)

	leftBatches := []arrow.RecordBatch{makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{1, 2})})}
	rightBatches := []arrow.RecordBatch{makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{3, 4})})}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	assert.Equal(t, int64(0), totalRows(results))
}

// TestHashJoin_Inner_MultipleMatchesPerKey verifies that one build row matching many probe
// rows produces one output row per pair.
func TestHashJoin_Inner_MultipleMatchesPerKey(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id"), symI64("val")}
	rSyms := []*plan.Symbol{symI64("r_id"), symI64("score")}
	node := makeJoinNode(plan.Inner, lSyms, rSyms, "l_id", "r_id")

	lSchema := arrow.NewSchema([]arrow.Field{
		{Name: "l_id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "val", Type: arrow.PrimitiveTypes.Int64},
	}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{
		{Name: "r_id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "score", Type: arrow.PrimitiveTypes.Int64},
	}, nil)

	// Single left row (key=1) matches two right rows → 2 output rows.
	leftBatches := []arrow.RecordBatch{
		makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{1}), testInt64Col([]int64{100})}),
	}
	rightBatches := []arrow.RecordBatch{
		makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{1, 1}), testInt64Col([]int64{10, 20})}),
	}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	assert.Equal(t, int64(2), totalRows(results))
}

// TestHashJoin_Inner_MultipleColumns verifies that all columns from both sides appear in output.
func TestHashJoin_Inner_MultipleColumns(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id"), symF64("price"), symStr("name")}
	rSyms := []*plan.Symbol{symI64("r_id"), symStr("category")}
	node := makeJoinNode(plan.Inner, lSyms, rSyms, "l_id", "r_id")

	lSchema := arrow.NewSchema([]arrow.Field{
		{Name: "l_id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "price", Type: arrow.PrimitiveTypes.Float64},
		{Name: "name", Type: arrow.BinaryTypes.String},
	}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{
		{Name: "r_id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "category", Type: arrow.BinaryTypes.String},
	}, nil)

	leftBatches := []arrow.RecordBatch{
		makeBatch(lSchema, []arrow.Array{
			testInt64Col([]int64{1}),
			testFloat64Col([]float64{9.99}),
			testStringCol([]string{"widget"}),
		}),
	}
	rightBatches := []arrow.RecordBatch{
		makeBatch(rSchema, []arrow.Array{
			testInt64Col([]int64{1}),
			testStringCol([]string{"gadgets"}),
		}),
	}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	require.Equal(t, int64(1), totalRows(results))
	// Output schema: left(3) + right(2) = 5 columns.
	assert.Equal(t, int64(5), results[0].NumCols())
}

// TestHashJoin_Left_UnmatchedBuildRows verifies that LEFT JOIN pads unmatched left rows
// with NULLs on the right side.
func TestHashJoin_Left_UnmatchedBuildRows(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}
	node := makeJoinNode(plan.Left, lSyms, rSyms, "l_id", "r_id")

	lSchema := arrow.NewSchema([]arrow.Field{{Name: "l_id", Type: arrow.PrimitiveTypes.Int64}}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{{Name: "r_id", Type: arrow.PrimitiveTypes.Int64}}, nil)

	// Left: 1, 2, 3 — Right: only 1 matches → rows 2 and 3 have no right counterpart.
	leftBatches := []arrow.RecordBatch{makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{1, 2, 3})})}
	rightBatches := []arrow.RecordBatch{makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{1})})}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	require.Equal(t, int64(3), totalRows(results))

	// Two unmatched left rows should carry NULL in the right (r_id) column.
	nullCount := 0
	for _, b := range results {
		rCol := b.Column(1) // r_id is index 1 in output schema
		for i := range int(b.NumRows()) {
			if rCol.IsNull(i) {
				nullCount++
			}
		}
	}
	assert.Equal(t, 2, nullCount, "two unmatched left rows should have NULL on the right side")
}

// TestHashJoin_Right_UnmatchedProbeRows verifies that RIGHT JOIN pads unmatched right rows
// with NULLs on the left side.
func TestHashJoin_Right_UnmatchedProbeRows(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}
	node := makeJoinNode(plan.Right, lSyms, rSyms, "l_id", "r_id")

	lSchema := arrow.NewSchema([]arrow.Field{{Name: "l_id", Type: arrow.PrimitiveTypes.Int64}}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{{Name: "r_id", Type: arrow.PrimitiveTypes.Int64}}, nil)

	// Left: only 1 — Right: 1, 2, 3 → rows 2 and 3 have no left counterpart.
	leftBatches := []arrow.RecordBatch{makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{1})})}
	rightBatches := []arrow.RecordBatch{makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{1, 2, 3})})}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	require.Equal(t, int64(3), totalRows(results))

	// Two unmatched right rows should carry NULL in the left (l_id) column.
	nullCount := 0
	for _, b := range results {
		lCol := b.Column(0) // l_id is index 0 in output schema
		for i := range int(b.NumRows()) {
			if lCol.IsNull(i) {
				nullCount++
			}
		}
	}
	assert.Equal(t, 2, nullCount, "two unmatched right rows should have NULL on the left side")
}

// TestHashJoin_Full_BothUnmatched verifies that FULL JOIN preserves unmatched rows from
// both sides, each padded with NULLs on the absent side.
func TestHashJoin_Full_BothUnmatched(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}
	node := makeJoinNode(plan.Full, lSyms, rSyms, "l_id", "r_id")

	lSchema := arrow.NewSchema([]arrow.Field{{Name: "l_id", Type: arrow.PrimitiveTypes.Int64}}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{{Name: "r_id", Type: arrow.PrimitiveTypes.Int64}}, nil)

	// Left: 1, 2 — Right: 1, 3
	// 1 matched pair + 1 left-unmatched (l_id=2) + 1 right-unmatched (r_id=3) = 3 rows.
	leftBatches := []arrow.RecordBatch{makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{1, 2})})}
	rightBatches := []arrow.RecordBatch{makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{1, 3})})}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	assert.Equal(t, int64(3), totalRows(results))
}

// TestHashJoin_Cross_CartesianProduct verifies that a CROSS JOIN (no criteria) produces
// M × N output rows.
func TestHashJoin_Cross_CartesianProduct(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}
	// IsCrossJoin() == true when Type==Inner, Criteria==nil, Filter==nil.
	node := &plan.JoinNode{
		Type:  plan.Inner,
		Left:  &stubPlanNode{symbols: lSyms},
		Right: &stubPlanNode{symbols: rSyms},
	}

	lSchema := arrow.NewSchema([]arrow.Field{{Name: "l_id", Type: arrow.PrimitiveTypes.Int64}}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{{Name: "r_id", Type: arrow.PrimitiveTypes.Int64}}, nil)

	// 3 left rows × 2 right rows = 6 output rows.
	leftBatches := []arrow.RecordBatch{makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{1, 2, 3})})}
	rightBatches := []arrow.RecordBatch{makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{10, 20})})}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	assert.Equal(t, int64(6), totalRows(results))
}

// TestHashJoin_Inner_CompositeKey verifies that multi-column join criteria work correctly,
// requiring all key columns to match simultaneously.
func TestHashJoin_Inner_CompositeKey(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("a"), symI64("b")}
	rSyms := []*plan.Symbol{symI64("x"), symI64("y")}
	node := &plan.JoinNode{
		Type:  plan.Inner,
		Left:  &stubPlanNode{symbols: lSyms},
		Right: &stubPlanNode{symbols: rSyms},
		// Composite key: a=x AND b=y
		Criteria: []*plan.EqualJoinCriteria{
			{Left: lSyms[0], Right: rSyms[0]},
			{Left: lSyms[1], Right: rSyms[1]},
		},
	}

	lSchema := arrow.NewSchema([]arrow.Field{
		{Name: "a", Type: arrow.PrimitiveTypes.Int64},
		{Name: "b", Type: arrow.PrimitiveTypes.Int64},
	}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{
		{Name: "x", Type: arrow.PrimitiveTypes.Int64},
		{Name: "y", Type: arrow.PrimitiveTypes.Int64},
	}, nil)

	// Left: (1,1),(1,2),(2,1) — Right: (1,1),(1,3),(2,1)
	// Matches: (1,1)↔(1,1) and (2,1)↔(2,1) → 2 rows; (1,2) has no counterpart.
	leftBatches := []arrow.RecordBatch{
		makeBatch(lSchema, []arrow.Array{
			testInt64Col([]int64{1, 1, 2}),
			testInt64Col([]int64{1, 2, 1}),
		}),
	}
	rightBatches := []arrow.RecordBatch{
		makeBatch(rSchema, []arrow.Array{
			testInt64Col([]int64{1, 1, 2}),
			testInt64Col([]int64{1, 3, 1}),
		}),
	}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	assert.Equal(t, int64(2), totalRows(results))
}

// TestHashJoin_Inner_NullKey_Excluded verifies SQL NULL semantics: NULL key columns
// never participate in equi-join matches (NULL ≠ NULL).
func TestHashJoin_Inner_NullKey_Excluded(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}
	node := makeJoinNode(plan.Inner, lSyms, rSyms, "l_id", "r_id")

	lSchema := arrow.NewSchema([]arrow.Field{{Name: "l_id", Type: arrow.PrimitiveTypes.Int64}}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{{Name: "r_id", Type: arrow.PrimitiveTypes.Int64}}, nil)

	// Left: 1, NULL, 2 — Right: 1, NULL, 2
	// NULL keys are excluded from the hash map; only 1↔1 and 2↔2 match → 2 rows.
	leftBatches := []arrow.RecordBatch{
		makeBatch(lSchema, []arrow.Array{
			testNullableInt64Col([]int64{1, 0, 2}, []bool{true, false, true}),
		}),
	}
	rightBatches := []arrow.RecordBatch{
		makeBatch(rSchema, []arrow.Array{
			testNullableInt64Col([]int64{1, 0, 2}, []bool{true, false, true}),
		}),
	}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	assert.Equal(t, int64(2), totalRows(results))
}

// TestHashJoin_WithFilter_NonEquality verifies that a non-equality ON condition stored in
// JoinNode.Filter is applied as a per-row residual filter after the hash lookup.
func TestHashJoin_WithFilter_NonEquality(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("id"), symI64("price")}
	rSyms := []*plan.Symbol{symI64("ref_id"), symI64("max_price")}

	// Residual filter: price < max_price.
	// prepare() compiles this against the output layout [id(0), price(1), ref_id(2), max_price(3)].
	filter := &tree.ComparisonExpression{
		Operator: tree.ComparisonLT,
		Left:     &tree.SymbolReference{Name: "price"},
		Right:    &tree.SymbolReference{Name: "max_price"},
	}
	node := &plan.JoinNode{
		Type:  plan.Inner,
		Left:  &stubPlanNode{symbols: lSyms},
		Right: &stubPlanNode{symbols: rSyms},
		Criteria: []*plan.EqualJoinCriteria{
			{Left: lSyms[0], Right: rSyms[0]}, // id = ref_id
		},
		Filter: filter,
	}

	lSchema := arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "price", Type: arrow.PrimitiveTypes.Int64},
	}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{
		{Name: "ref_id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "max_price", Type: arrow.PrimitiveTypes.Int64},
	}, nil)

	// Left: (id=1,price=100),(id=2,price=200) — Right: (ref_id=1,max=150),(ref_id=1,max=90),(ref_id=2,max=250)
	// Equi-join on id=ref_id, then residual price < max_price:
	//   (1,100)×(1,150): 100<150 ✓
	//   (1,100)×(1,90):  100<90  ✗
	//   (2,200)×(2,250): 200<250 ✓
	// Expected: 2 rows.
	leftBatches := []arrow.RecordBatch{
		makeBatch(lSchema, []arrow.Array{
			testInt64Col([]int64{1, 2}),
			testInt64Col([]int64{100, 200}),
		}),
	}
	rightBatches := []arrow.RecordBatch{
		makeBatch(rSchema, []arrow.Array{
			testInt64Col([]int64{1, 1, 2}),
			testInt64Col([]int64{150, 90, 250}),
		}),
	}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	assert.Equal(t, int64(2), totalRows(results))
}

// TestHashJoin_MultiBatch_BuildSide verifies that build-side input spread across multiple
// RecordBatches is fully consumed and all rows are indexed into the hash map.
func TestHashJoin_MultiBatch_BuildSide(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}
	node := makeJoinNode(plan.Inner, lSyms, rSyms, "l_id", "r_id")

	lSchema := arrow.NewSchema([]arrow.Field{{Name: "l_id", Type: arrow.PrimitiveTypes.Int64}}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{{Name: "r_id", Type: arrow.PrimitiveTypes.Int64}}, nil)

	// Build side: batch1=[1,2], batch2=[3,4] — Probe: [2,3,5]
	// Matches: 2 (from batch1) and 3 (from batch2) → 2 rows.
	leftBatches := []arrow.RecordBatch{
		makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{1, 2})}),
		makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{3, 4})}),
	}
	rightBatches := []arrow.RecordBatch{
		makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{2, 3, 5})}),
	}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	assert.Equal(t, int64(2), totalRows(results))
}

// TestHashJoin_MultiBatch_ProbeSide verifies that probe-side input spread across multiple
// RecordBatches is fully consumed and each batch is probed against the hash map.
func TestHashJoin_MultiBatch_ProbeSide(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}
	node := makeJoinNode(plan.Inner, lSyms, rSyms, "l_id", "r_id")

	lSchema := arrow.NewSchema([]arrow.Field{{Name: "l_id", Type: arrow.PrimitiveTypes.Int64}}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{{Name: "r_id", Type: arrow.PrimitiveTypes.Int64}}, nil)

	// Build: [1,2,3] — Probe: batch1=[1,4], batch2=[2,3]
	// Matches: 1 (batch1), 2 and 3 (batch2) → 3 rows.
	leftBatches := []arrow.RecordBatch{
		makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{1, 2, 3})}),
	}
	rightBatches := []arrow.RecordBatch{
		makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{1, 4})}),
		makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{2, 3})}),
	}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	assert.Equal(t, int64(3), totalRows(results))
}

// TestHashJoin_String_ReturnsName verifies the operator's String() method.
func TestHashJoin_String_ReturnsName(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}
	node := makeJoinNode(plan.Inner, lSyms, rSyms, "l_id", "r_id")
	op := NewHashJoinOperator(node, nil, nil)
	assert.Equal(t, "HashJoinOperator", op.String())
}

// TestHashJoin_ThetaJoin_FilterOnly verifies that a join with no equi-criteria
// but a non-equality ON filter (theta join) correctly enumerates all left×right
// candidate pairs and applies the filter per pair.
func TestHashJoin_ThetaJoin_FilterOnly(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_val")}
	rSyms := []*plan.Symbol{symI64("r_val")}

	// Theta filter: l_val < r_val.
	// prepare() compiles against output layout [l_val(0), r_val(1)].
	filter := &tree.ComparisonExpression{
		Operator: tree.ComparisonLT,
		Left:     &tree.SymbolReference{Name: "l_val"},
		Right:    &tree.SymbolReference{Name: "r_val"},
	}
	node := &plan.JoinNode{
		Type:     plan.Inner,
		Left:     &stubPlanNode{symbols: lSyms},
		Right:    &stubPlanNode{symbols: rSyms},
		Criteria: nil,  // no equi-criteria → theta join
		Filter:   filter,
	}

	lSchema := arrow.NewSchema([]arrow.Field{{Name: "l_val", Type: arrow.PrimitiveTypes.Int64}}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{{Name: "r_val", Type: arrow.PrimitiveTypes.Int64}}, nil)

	// Left: [1, 3] — Right: [2, 4]
	// All 4 pairs: (1,2)✓  (1,4)✓  (3,2)✗  (3,4)✓ → 3 output rows.
	leftBatches := []arrow.RecordBatch{
		makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{1, 3})}),
	}
	rightBatches := []arrow.RecordBatch{
		makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{2, 4})}),
	}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	assert.Equal(t, int64(3), totalRows(results))
}

// TestHashJoin_EmptyBuildSide verifies that an empty left (build) side produces
// no matched rows for INNER JOIN, but a RIGHT JOIN still emits all right rows
// null-padded on the left.
func TestHashJoin_EmptyBuildSide(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}

	rSchema := arrow.NewSchema([]arrow.Field{{Name: "r_id", Type: arrow.PrimitiveTypes.Int64}}, nil)

	rightBatches := []arrow.RecordBatch{
		makeBatch(rSchema, []arrow.Array{testInt64Col([]int64{1, 2, 3})}),
	}

	t.Run("INNER_returns_no_rows", func(t *testing.T) {
		node := makeJoinNode(plan.Inner, lSyms, rSyms, "l_id", "r_id")
		op := NewHashJoinOperator(node, nil, nil)
		results := runJoin(op, nil, rightBatches)
		assert.Equal(t, int64(0), totalRows(results))
	})

	t.Run("RIGHT_returns_all_right_rows_with_null_left", func(t *testing.T) {
		node := makeJoinNode(plan.Right, lSyms, rSyms, "l_id", "r_id")
		op := NewHashJoinOperator(node, nil, nil)
		results := runJoin(op, nil, rightBatches)
		require.Equal(t, int64(3), totalRows(results))

		// All left-side (l_id) values must be NULL.
		for _, b := range results {
			for i := range int(b.NumRows()) {
				assert.True(t, b.Column(0).IsNull(i), "l_id must be NULL for unmatched right row")
			}
		}
	})
}

// TestHashJoin_EmptyProbeSide verifies that an empty right (probe) side produces
// no matched rows for INNER JOIN, but a LEFT JOIN still emits all left rows
// null-padded on the right.
func TestHashJoin_EmptyProbeSide(t *testing.T) {
	lSyms := []*plan.Symbol{symI64("l_id")}
	rSyms := []*plan.Symbol{symI64("r_id")}

	lSchema := arrow.NewSchema([]arrow.Field{{Name: "l_id", Type: arrow.PrimitiveTypes.Int64}}, nil)

	leftBatches := []arrow.RecordBatch{
		makeBatch(lSchema, []arrow.Array{testInt64Col([]int64{10, 20})}),
	}

	t.Run("INNER_returns_no_rows", func(t *testing.T) {
		node := makeJoinNode(plan.Inner, lSyms, rSyms, "l_id", "r_id")
		op := NewHashJoinOperator(node, nil, nil)
		results := runJoin(op, leftBatches, nil)
		assert.Equal(t, int64(0), totalRows(results))
	})

	t.Run("LEFT_returns_all_left_rows_with_null_right", func(t *testing.T) {
		node := makeJoinNode(plan.Left, lSyms, rSyms, "l_id", "r_id")
		op := NewHashJoinOperator(node, nil, nil)
		results := runJoin(op, leftBatches, nil)
		require.Equal(t, int64(2), totalRows(results))

		// All right-side (r_id) values must be NULL.
		for _, b := range results {
			for i := range int(b.NumRows()) {
				assert.True(t, b.Column(1).IsNull(i), "r_id must be NULL for unmatched left row")
			}
		}
	})
}

// TestHashJoin_KeyCollision_StringWithNullByte verifies that string column values
// containing the zero byte '\x00' do not cause false-positive key collisions
// after the length-prefix encoding fix.
func TestHashJoin_KeyCollision_StringWithNullByte(t *testing.T) {
	lSyms := []*plan.Symbol{symStr("a"), symStr("b")}
	rSyms := []*plan.Symbol{symStr("x"), symStr("y")}
	node := &plan.JoinNode{
		Type:  plan.Inner,
		Left:  &stubPlanNode{symbols: lSyms},
		Right: &stubPlanNode{symbols: rSyms},
		Criteria: []*plan.EqualJoinCriteria{
			{Left: lSyms[0], Right: rSyms[0]},
			{Left: lSyms[1], Right: rSyms[1]},
		},
	}

	lSchema := arrow.NewSchema([]arrow.Field{
		{Name: "a", Type: arrow.BinaryTypes.String},
		{Name: "b", Type: arrow.BinaryTypes.String},
	}, nil)
	rSchema := arrow.NewSchema([]arrow.Field{
		{Name: "x", Type: arrow.BinaryTypes.String},
		{Name: "y", Type: arrow.BinaryTypes.String},
	}, nil)

	// Left row 0: ("ab\x00", "c") — composite key must not collide with ("ab", "\x00c").
	// Right: only ("ab\x00", "c") matches row 0; ("ab", "\x00c") does NOT.
	leftBatches := []arrow.RecordBatch{
		makeBatch(lSchema, []arrow.Array{
			testStringCol([]string{"ab\x00", "ab"}),
			testStringCol([]string{"c", "\x00c"}),
		}),
	}
	rightBatches := []arrow.RecordBatch{
		makeBatch(rSchema, []arrow.Array{
			testStringCol([]string{"ab\x00"}),
			testStringCol([]string{"c"}),
		}),
	}

	op := NewHashJoinOperator(node, nil, nil)
	results := runJoin(op, leftBatches, rightBatches)

	// Only the exact match ("ab\x00","c") should join; ("ab","\x00c") must not.
	assert.Equal(t, int64(1), totalRows(results))
}
