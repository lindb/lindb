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
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lindb/lindb/sql/tree"
)

// buildMinimalVisitor creates a StatementVisitor with a minimal AnalyzerContext,
// suitable for calling analyzeAggregations in isolation.
func buildMinimalVisitor() *StatementVisitor {
	idAlloc := tree.NewNodeIDAllocator()
	stmt := &tree.QuerySpecification{}
	ctx := NewAnalyzerContext("test_db", stmt, idAlloc, false)
	return NewStatementVisitor(nil, &StatementAnalyzer{ctx: ctx, metadataMgr: nil}, true)
}

// buildScopeWithField creates a scope that contains a single field with the given name and type.
func buildScopeWithField(name string, dataType arrow.DataType) *Scope {
	fields := []*tree.Field{
		{Index: 0, Name: name, DataType: dataType},
	}
	return createAndAssignScopeForTest(fields)
}

// createAndAssignScopeForTest creates a scope with the given fields (helper for tests).
func createAndAssignScopeForTest(fields []*tree.Field) *Scope {
	sb := NewScopeBuilder(nil)
	sb.withRelation(NewRelationID(nil), NewRelation(fields))
	return sb.build()
}

// makeHistogramQuantileExpr builds histogram_quantile(0.99, <colName>) as a tree.Expression.
func makeHistogramQuantileExpr(colName string) *tree.FunctionCall {
	idAlloc := tree.NewNodeIDAllocator()
	phi := tree.NewFloatLiteral(idAlloc.Next(), nil, "0.99")
	col := &tree.Identifier{Value: colName}
	col.SetID(idAlloc.Next())
	call := &tree.FunctionCall{
		Name:      tree.HistogramQuantile,
		Arguments: []tree.Expression{phi, col},
	}
	call.SetID(idAlloc.Next())
	return call
}

// ----- analyzeAggregations nil-pointer regression tests -----

// TestAnalyzeAggregations_NilResolvedField is the direct regression test for the production
// panic: when "sent_duration" is used as a histogram_quantile argument but is NOT present in
// the source scope, resolveField returns nil and the old code panicked at resolvedField.Field.
func TestAnalyzeAggregations_NilResolvedField(t *testing.T) {
	v := buildMinimalVisitor()
	// Empty scope — "sent_duration" is NOT present.
	emptyScope := createScope(nil)
	query := &tree.QuerySpecification{}

	call := makeHistogramQuantileExpr("sent_duration")
	require.NotPanics(t, func() {
		v.analyzeAggregations(query, emptyScope, nil, nil,
			[]tree.Expression{call}, nil)
	}, "analyzeAggregations must not panic when the histogram column is absent from the scope")
}

// TestAnalyzeAggregations_HistogramTypeSkipped verifies that an Identifier whose schema type
// is Histogram is NOT auto-injected as a built-in aggregation.  Injecting "histogram" as a
// function name would be invalid; histogram columns must be queried via histogram_quantile etc.
func TestAnalyzeAggregations_HistogramTypeSkipped(t *testing.T) {
	v := buildMinimalVisitor()
	// Scope with "sent_duration" as Histogram type (as returned by storage RPC TableSchema).
	scope := buildScopeWithField("sent_duration", larrow.ExtensionTypes.Histogram)
	query := &tree.QuerySpecification{}

	call := makeHistogramQuantileExpr("sent_duration")
	require.NotPanics(t, func() {
		v.analyzeAggregations(query, scope, nil, nil,
			[]tree.Expression{call}, nil)
	})

	// No auto-injected "histogram" function should have been added.
	aggs := v.analyzer.ctx.Analysis.GetAggregates(query)
	for _, agg := range aggs {
		assert.NotEqual(t, tree.FuncName("histogram"), agg.Name,
			"Histogram-type column must not produce an auto-injected 'histogram' function call")
	}
}

// TestAnalyzeAggregations_NormalFieldInjected verifies that a regular Sum-type field that is
// referenced directly in SELECT (not as a function arg) still gets a builtin agg injected.
// This ensures the fix doesn't break the existing auto-injection behaviour for plain fields.
func TestAnalyzeAggregations_NormalFieldInjected(t *testing.T) {
	v := buildMinimalVisitor()
	scope := buildScopeWithField("cpu", larrow.ExtensionTypes.Sum)
	query := &tree.QuerySpecification{}

	// Plain identifier "cpu" used directly in SELECT (not wrapped in a function).
	idAlloc := tree.NewNodeIDAllocator()
	col := &tree.Identifier{Value: "cpu"}
	col.SetID(idAlloc.Next())

	require.NotPanics(t, func() {
		v.analyzeAggregations(query, scope, nil, nil,
			[]tree.Expression{col}, nil)
	})

	// "sum" function should have been auto-injected for the "cpu" field.
	aggs := v.analyzer.ctx.Analysis.GetAggregates(query)
	found := false
	for _, agg := range aggs {
		if agg.Name == tree.Sum {
			found = true
		}
	}
	assert.True(t, found, "Sum-type plain field should have a builtin sum() injected")
}

// TestAnalyzeAggregations_HistogramAsDirectSelect verifies that using a histogram column
// directly in SELECT (e.g. SELECT sent_duration FROM ...) does NOT produce a "histogram"
// auto-injection even when the column IS in the scope.
func TestAnalyzeAggregations_HistogramAsDirectSelect(t *testing.T) {
	v := buildMinimalVisitor()
	scope := buildScopeWithField("sent_duration", larrow.ExtensionTypes.Histogram)
	query := &tree.QuerySpecification{}

	idAlloc := tree.NewNodeIDAllocator()
	col := &tree.Identifier{Value: "sent_duration"}
	col.SetID(idAlloc.Next())

	require.NotPanics(t, func() {
		v.analyzeAggregations(query, scope, nil, nil,
			[]tree.Expression{col}, nil)
	})

	aggs := v.analyzer.ctx.Analysis.GetAggregates(query)
	for _, agg := range aggs {
		assert.NotEqual(t, tree.FuncName("histogram"), agg.Name,
			"Histogram column used directly in SELECT must not inject an invalid 'histogram' func")
	}
}
