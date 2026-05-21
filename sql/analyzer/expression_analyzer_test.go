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

package analyzer

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lindb/lindb/sql/tree"
)

// newTestVisitor creates an ExpressionVisitor wired to a fresh AnalyzerContext
// and an empty scope, suitable for unit-testing individual visit* methods.
func newTestVisitor() (*ExpressionVisitor, *AnalyzerContext) {
	idAlloc := tree.NewNodeIDAllocator()
	// Use a minimal Statement placeholder so Analysis can be constructed.
	stmt := &tree.QuerySpecification{}
	ctx := NewAnalyzerContext("test_db", stmt, idAlloc, false)
	scope := createScope(nil)
	return NewExpressionVisitor(scope, NewExpressionAnalyzer(ctx)), ctx
}

// ----- FloatLiteral — the regression that caused the production panic -----

// TestVisitFloatLiteral verifies that FloatLiteral is handled without panic
// and is assigned the Float64 data type.
func TestVisitFloatLiteral(t *testing.T) {
	v, ctx := newTestVisitor()
	idAlloc := tree.NewNodeIDAllocator()
	node := tree.NewFloatLiteral(idAlloc.Next(), nil, "0.99")

	var dt arrow.DataType
	require.NotPanics(t, func() {
		dt = v.visitFloatLiteral(nil, node).(arrow.DataType)
	})
	assert.True(t, arrow.TypeEqual(dt, arrow.PrimitiveTypes.Float64),
		"FloatLiteral should resolve to Float64, got %v", dt)
	assert.True(t, arrow.TypeEqual(ctx.Analysis.GetType(node), arrow.PrimitiveTypes.Float64))
}

// TestVisitFloatLiteral_VariousValues verifies a range of phi-like values.
func TestVisitFloatLiteral_VariousValues(t *testing.T) {
	cases := []string{"0.0", "0.5", "0.95", "0.99", "1.0"}
	for _, raw := range cases {
		v, _ := newTestVisitor()
		idAlloc := tree.NewNodeIDAllocator()
		node := tree.NewFloatLiteral(idAlloc.Next(), nil, raw)
		require.NotPanics(t, func() { v.visitFloatLiteral(nil, node) },
			"FloatLiteral(%s) must not panic", raw)
	}
}

// TestVisit_FloatLiteral_ViaVisitDispatch verifies that the Visit switch correctly
// dispatches *tree.FloatLiteral to visitFloatLiteral (not the default panic branch).
// This is the exact code path that was broken in production.
func TestVisit_FloatLiteral_ViaVisitDispatch(t *testing.T) {
	v, ctx := newTestVisitor()
	idAlloc := tree.NewNodeIDAllocator()
	node := tree.NewFloatLiteral(idAlloc.Next(), nil, "0.99")

	stackCtx := tree.NewStackableVisitorContext(&Context{scope: createScope(nil)})
	var dt arrow.DataType
	require.NotPanics(t, func() {
		dt = v.Visit(stackCtx, node).(arrow.DataType)
	}, "Visit must not panic for *tree.FloatLiteral")
	assert.True(t, arrow.TypeEqual(dt, arrow.PrimitiveTypes.Float64))
	assert.True(t, arrow.TypeEqual(ctx.Analysis.GetType(node), arrow.PrimitiveTypes.Float64))
}

// TestResolveHistogramColumn_UnknownColumn verifies that a histogram column argument
// that is NOT in the schema is resolved to Sum type instead of panicking.
// This covers the case where "sent_duration" is a logical histogram name not yet in the schema.
func TestResolveHistogramColumn_UnknownColumn(t *testing.T) {
	v, ctx := newTestVisitor()
	idAlloc := tree.NewNodeIDAllocator()
	node := &tree.Identifier{Value: "sent_duration"}
	node.SetID(idAlloc.Next())

	stackCtx := tree.NewStackableVisitorContext(&Context{scope: createScope(nil)})
	var dt arrow.DataType
	require.NotPanics(t, func() {
		dt = v.resolveHistogramColumn(stackCtx, node)
	}, "resolveHistogramColumn must not panic for unknown histogram column name")
	assert.NotNil(t, dt, "should return Sum type placeholder for unresolved histogram name")
	_ = ctx
}

// TestVisitFunctionCall_HistogramQuantile_UnknownColumn verifies the full visitFunctionCall
// dispatch for histogram_quantile(0.99, sent_duration) where sent_duration is NOT in the scope.
// This is the exact production failure scenario.
func TestVisitFunctionCall_HistogramQuantile_UnknownColumn(t *testing.T) {
	v, _ := newTestVisitor()
	idAlloc := tree.NewNodeIDAllocator()

	phi := tree.NewFloatLiteral(idAlloc.Next(), nil, "0.99")
	col := &tree.Identifier{Value: "sent_duration"}
	col.SetID(idAlloc.Next())

	call := &tree.FunctionCall{Name: tree.HistogramQuantile, Arguments: []tree.Expression{phi, col}}
	call.SetID(idAlloc.Next())

	stackCtx := tree.NewStackableVisitorContext(&Context{scope: createScope(nil)})
	require.NotPanics(t, func() {
		v.Visit(stackCtx, call)
	}, "histogram_quantile(0.99, unknown_col) must not panic")
}

// ----- Regression guard: other literal types must still work -----

func TestVisit_LongLiteral(t *testing.T) {
	v, ctx := newTestVisitor()
	idAlloc := tree.NewNodeIDAllocator()
	node := tree.NewLongLiteral(idAlloc.Next(), nil, "42")
	stackCtx := tree.NewStackableVisitorContext(&Context{scope: createScope(nil)})

	require.NotPanics(t, func() { v.Visit(stackCtx, node) })
	assert.True(t, arrow.TypeEqual(ctx.Analysis.GetType(node), arrow.PrimitiveTypes.Int64))
}

func TestVisit_StringLiteral(t *testing.T) {
	v, ctx := newTestVisitor()
	idAlloc := tree.NewNodeIDAllocator()
	node := &tree.StringLiteral{Value: "hello"}
	node.SetID(idAlloc.Next())
	stackCtx := tree.NewStackableVisitorContext(&Context{scope: createScope(nil)})

	require.NotPanics(t, func() { v.Visit(stackCtx, node) })
	assert.True(t, arrow.TypeEqual(ctx.Analysis.GetType(node), arrow.BinaryTypes.String))
}
