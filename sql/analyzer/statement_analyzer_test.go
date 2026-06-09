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

// ----- hasAggregates regression test -----

// TestHasAggregates_ScalarFuncNotCounted is a regression test for the bug where
// hasAggregates() counted ALL FunctionCall nodes instead of only aggregation functions.
// rand() is a scalar function; it must NOT cause hasAggregates() to return true, which
// would otherwise trigger an incorrect Aggregate plan node and broken table scan columns.
func TestHasAggregates_ScalarFuncNotCounted(t *testing.T) {
	v := buildMinimalVisitor()
	idAlloc := tree.NewNodeIDAllocator()

	// Build SELECT rand() — scalar function, not an aggregation.
	randCall := &tree.FunctionCall{Name: tree.Rand}
	randCall.SetID(idAlloc.Next())
	query := &tree.QuerySpecification{
		Select: &tree.Select{
			SelectItems: []tree.SelectItem{
				&tree.SingleColumn{Expression: randCall},
			},
		},
	}

	assert.False(t, v.hasAggregates(query),
		"rand() is a scalar function — hasAggregates must return false to avoid a spurious Aggregate plan node")
}

// TestHasAggregates_AggFuncCounted verifies that a genuine aggregation function
// (e.g. sum()) still causes hasAggregates() to return true.
func TestHasAggregates_AggFuncCounted(t *testing.T) {
	v := buildMinimalVisitor()
	idAlloc := tree.NewNodeIDAllocator()

	sumCall := &tree.FunctionCall{Name: tree.Sum}
	sumCall.SetID(idAlloc.Next())
	query := &tree.QuerySpecification{
		Select: &tree.Select{
			SelectItems: []tree.SelectItem{
				&tree.SingleColumn{Expression: sumCall},
			},
		},
	}

	assert.True(t, v.hasAggregates(query),
		"sum() is an aggregation function — hasAggregates must return true")
}

// ----- rewriteAliasesInExpression tests -----

// TestRewriteAliases_IdentifierReplacedWithUnderlying verifies that an Identifier
// matching a SELECT alias is replaced by the underlying expression.
// This is the core regression for "HAVING c = 1169" where c is "count(1) AS c".
func TestRewriteAliases_IdentifierReplacedWithUnderlying(t *testing.T) {
	idAlloc := tree.NewNodeIDAllocator()

	// Build count(1) as the underlying expression for alias "c".
	arg1 := tree.NewLongLiteral(idAlloc.Next(), nil, "1")
	countCall := &tree.FunctionCall{Name: tree.Count, Arguments: []tree.Expression{arg1}}
	countCall.SetID(idAlloc.Next())

	aliasMap := map[string]tree.Expression{"c": countCall}

	// HAVING: c = 1169  →  count(1) = 1169
	aliasIdent := &tree.Identifier{Value: "c"}
	aliasIdent.SetID(idAlloc.Next())
	lit := tree.NewLongLiteral(idAlloc.Next(), nil, "1169")
	cmp := &tree.ComparisonExpression{Operator: tree.ComparisonEQ, Left: aliasIdent, Right: lit}
	cmp.SetID(idAlloc.Next())

	result := rewriteAliasesInExpression(cmp, aliasMap)

	// After rewriting, Left should be the count(1) FunctionCall, not the alias Identifier.
	rewritten, ok := result.(*tree.ComparisonExpression)
	require.True(t, ok, "result must still be a ComparisonExpression")
	assert.Equal(t, countCall, rewritten.Left,
		"alias Identifier 'c' must be replaced by the underlying count(1) FunctionCall")
	assert.Equal(t, lit, rewritten.Right, "right side (literal) must remain unchanged")
}

// TestRewriteAliases_UnknownIdentifierUnchanged verifies that an Identifier NOT in the
// aliasMap is left as-is (e.g. a real column name in HAVING must not be modified).
func TestRewriteAliases_UnknownIdentifierUnchanged(t *testing.T) {
	idAlloc := tree.NewNodeIDAllocator()
	aliasMap := map[string]tree.Expression{} // empty — no aliases

	col := &tree.Identifier{Value: "level"}
	col.SetID(idAlloc.Next())

	result := rewriteAliasesInExpression(col, aliasMap)
	assert.Equal(t, col, result, "unknown Identifier must be returned unchanged")
}

// TestRewriteAliases_LogicalExpressionRewritesAllTerms verifies that aliases nested
// inside a LogicalExpression (AND/OR) are all expanded.
func TestRewriteAliases_LogicalExpressionRewritesAllTerms(t *testing.T) {
	idAlloc := tree.NewNodeIDAllocator()

	// Build underlying expressions for two aliases.
	countCall := &tree.FunctionCall{Name: tree.Count}
	countCall.SetID(idAlloc.Next())
	sumCall := &tree.FunctionCall{Name: tree.Sum}
	sumCall.SetID(idAlloc.Next())

	aliasMap := map[string]tree.Expression{
		"c": countCall,
		"s": sumCall,
	}

	// HAVING: c > 1 AND s < 100
	cIdent := &tree.Identifier{Value: "c"}
	cIdent.SetID(idAlloc.Next())
	sIdent := &tree.Identifier{Value: "s"}
	sIdent.SetID(idAlloc.Next())

	cmpC := &tree.ComparisonExpression{Operator: tree.ComparisonGT, Left: cIdent, Right: tree.NewLongLiteral(idAlloc.Next(), nil, "1")}
	cmpC.SetID(idAlloc.Next())
	cmpS := &tree.ComparisonExpression{Operator: tree.ComparisonLT, Left: sIdent, Right: tree.NewLongLiteral(idAlloc.Next(), nil, "100")}
	cmpS.SetID(idAlloc.Next())

	logical := &tree.LogicalExpression{
		Operator: tree.LogicalAND,
		Terms:    []tree.Expression{cmpC, cmpS},
	}
	logical.SetID(idAlloc.Next())

	rewriteAliasesInExpression(logical, aliasMap)

	// Both terms should now reference the underlying FunctionCall, not the Identifier.
	assert.Equal(t, countCall, logical.Terms[0].(*tree.ComparisonExpression).Left,
		"first term's left side must be count FunctionCall after rewriting alias 'c'")
	assert.Equal(t, sumCall, logical.Terms[1].(*tree.ComparisonExpression).Left,
		"second term's left side must be sum FunctionCall after rewriting alias 's'")
}

// TestRewriteAliases_NilExpressionReturnedAsNil verifies that nil input is handled safely.
func TestRewriteAliases_NilExpressionReturnedAsNil(t *testing.T) {
	result := rewriteAliasesInExpression(nil, map[string]tree.Expression{})
	assert.Nil(t, result, "nil expression must return nil")
}

// ----- collectGroupByColumnNames tests -----

// TestCollectGroupByColumnNames_NilGroupBy verifies that a nil GroupBy returns an empty set.
func TestCollectGroupByColumnNames_NilGroupBy(t *testing.T) {
	cols := collectGroupByColumnNames(nil)
	assert.Empty(t, cols, "nil GroupBy must return empty set")
}

// TestCollectGroupByColumnNames_SingleColumn verifies that a single-column GROUP BY is extracted.
func TestCollectGroupByColumnNames_SingleColumn(t *testing.T) {
	idAlloc := tree.NewNodeIDAllocator()
	col := &tree.Identifier{Value: "level"}
	col.SetID(idAlloc.Next())
	gb := &tree.GroupBy{
		GroupingElements: []tree.GroupingElement{
			&tree.SimpleGroupBy{Columns: []tree.Expression{col}},
		},
	}
	cols := collectGroupByColumnNames(gb)
	assert.Contains(t, cols, "level", "GROUP BY level must be collected")
	assert.Len(t, cols, 1)
}

// TestCollectGroupByColumnNames_MultipleColumns verifies that all GROUP BY columns are extracted.
func TestCollectGroupByColumnNames_MultipleColumns(t *testing.T) {
	idAlloc := tree.NewNodeIDAllocator()
	col1 := &tree.Identifier{Value: "region"}
	col1.SetID(idAlloc.Next())
	col2 := &tree.Identifier{Value: "host"}
	col2.SetID(idAlloc.Next())
	gb := &tree.GroupBy{
		GroupingElements: []tree.GroupingElement{
			&tree.SimpleGroupBy{Columns: []tree.Expression{col1, col2}},
		},
	}
	cols := collectGroupByColumnNames(gb)
	assert.Contains(t, cols, "region")
	assert.Contains(t, cols, "host")
	assert.Len(t, cols, 2)
}

// ----- validateHavingIdentifiers tests -----

// TestValidateHavingIdentifiers_SelectAlias verifies that a SELECT alias in HAVING does NOT panic.
func TestValidateHavingIdentifiers_SelectAlias(t *testing.T) {
	idAlloc := tree.NewNodeIDAllocator()

	// aliasMap: "c" → count(1)
	countCall := &tree.FunctionCall{Name: tree.Count}
	countCall.SetID(idAlloc.Next())
	aliasMap := map[string]tree.Expression{"c": countCall}
	groupByCols := map[string]struct{}{}

	// HAVING: c = 1169
	cIdent := &tree.Identifier{Value: "c"}
	cIdent.SetID(idAlloc.Next())
	lit := tree.NewLongLiteral(idAlloc.Next(), nil, "1169")
	cmp := &tree.ComparisonExpression{Operator: tree.ComparisonEQ, Left: cIdent, Right: lit}
	cmp.SetID(idAlloc.Next())

	require.NotPanics(t, func() {
		validateHavingIdentifiers(cmp, aliasMap, groupByCols)
	}, "SELECT alias in HAVING must not panic")
}

// TestValidateHavingIdentifiers_GroupByColumn verifies that a GROUP BY column in HAVING does NOT panic.
func TestValidateHavingIdentifiers_GroupByColumn(t *testing.T) {
	idAlloc := tree.NewNodeIDAllocator()

	aliasMap := map[string]tree.Expression{}
	groupByCols := map[string]struct{}{"level": {}}

	// HAVING: level = 'error'
	levelIdent := &tree.Identifier{Value: "level"}
	levelIdent.SetID(idAlloc.Next())
	lit := &tree.StringLiteral{Value: "error"}
	lit.SetID(idAlloc.Next())
	cmp := &tree.ComparisonExpression{Operator: tree.ComparisonEQ, Left: levelIdent, Right: lit}
	cmp.SetID(idAlloc.Next())

	require.NotPanics(t, func() {
		validateHavingIdentifiers(cmp, aliasMap, groupByCols)
	}, "GROUP BY column in HAVING must not panic")
}

// TestValidateHavingIdentifiers_UnknownColumn is the regression test for the production bug:
// HAVING references a bare identifier that is neither a SELECT alias nor a GROUP BY column.
// Expected: panics with "Unknown column 'ct' in 'having clause'" (MySQL-style).
// Before the fix: panicked with "'(ct <> 100)' must be an aggregate expression or appear in GROUP BY clause".
func TestValidateHavingIdentifiers_UnknownColumn(t *testing.T) {
	idAlloc := tree.NewNodeIDAllocator()

	aliasMap := map[string]tree.Expression{}    // no aliases
	groupByCols := map[string]struct{}{}         // no GROUP BY columns

	// HAVING: ct <> 100
	ctIdent := &tree.Identifier{Value: "ct"}
	ctIdent.SetID(idAlloc.Next())
	lit := tree.NewLongLiteral(idAlloc.Next(), nil, "100")
	cmp := &tree.ComparisonExpression{Operator: tree.ComparisonNEQ, Left: ctIdent, Right: lit}
	cmp.SetID(idAlloc.Next())

	require.Panics(t, func() {
		validateHavingIdentifiers(cmp, aliasMap, groupByCols)
	}, "unknown HAVING identifier must panic")

	// Verify the panic message is the MySQL-style error, not the aggregation error.
	defer func() {
		if r := recover(); r != nil {
			msg, ok := r.(string)
			require.True(t, ok, "panic value must be a string")
			assert.Contains(t, msg, "Unknown column 'ct' in 'having clause'",
				"error must be MySQL-style 'Unknown column' message")
		}
	}()
	validateHavingIdentifiers(cmp, aliasMap, groupByCols)
}

// TestValidateHavingIdentifiers_FunctionCallNotRecursed verifies that a FunctionCall in HAVING
// (e.g. count(1) > 5) does NOT trigger the unknown-column check for its arguments,
// since aggregate args are validated against the source schema — not HAVING column rules.
func TestValidateHavingIdentifiers_FunctionCallNotRecursed(t *testing.T) {
	idAlloc := tree.NewNodeIDAllocator()

	aliasMap := map[string]tree.Expression{}
	groupByCols := map[string]struct{}{}

	// HAVING: count(1) > 5  — "1" is a literal arg, not an Identifier, so no panic anyway.
	arg := tree.NewLongLiteral(idAlloc.Next(), nil, "1")
	countCall := &tree.FunctionCall{Name: tree.Count, Arguments: []tree.Expression{arg}}
	countCall.SetID(idAlloc.Next())
	lit5 := tree.NewLongLiteral(idAlloc.Next(), nil, "5")
	cmp := &tree.ComparisonExpression{Operator: tree.ComparisonGT, Left: countCall, Right: lit5}
	cmp.SetID(idAlloc.Next())

	require.NotPanics(t, func() {
		validateHavingIdentifiers(cmp, aliasMap, groupByCols)
	}, "count(1) > 5 in HAVING must not panic — FunctionCall is always allowed")
}

// TestValidateHavingIdentifiers_LogicalAnd verifies that both arms of an AND expression
// are checked — an unknown identifier in either arm must trigger the error.
func TestValidateHavingIdentifiers_LogicalAnd_UnknownInSecondArm(t *testing.T) {
	idAlloc := tree.NewNodeIDAllocator()

	// aliasMap: "c" is valid; "bad" is not.
	countCall := &tree.FunctionCall{Name: tree.Count}
	countCall.SetID(idAlloc.Next())
	aliasMap := map[string]tree.Expression{"c": countCall}
	groupByCols := map[string]struct{}{}

	// HAVING: c > 1 AND bad = 2
	cIdent := &tree.Identifier{Value: "c"}
	cIdent.SetID(idAlloc.Next())
	badIdent := &tree.Identifier{Value: "bad"}
	badIdent.SetID(idAlloc.Next())

	cmpC := &tree.ComparisonExpression{
		Operator: tree.ComparisonGT,
		Left:     cIdent,
		Right:    tree.NewLongLiteral(idAlloc.Next(), nil, "1"),
	}
	cmpC.SetID(idAlloc.Next())
	cmpBad := &tree.ComparisonExpression{
		Operator: tree.ComparisonEQ,
		Left:     badIdent,
		Right:    tree.NewLongLiteral(idAlloc.Next(), nil, "2"),
	}
	cmpBad.SetID(idAlloc.Next())

	logical := &tree.LogicalExpression{
		Operator: tree.LogicalAND,
		Terms:    []tree.Expression{cmpC, cmpBad},
	}
	logical.SetID(idAlloc.Next())

	require.Panics(t, func() {
		validateHavingIdentifiers(logical, aliasMap, groupByCols)
	}, "unknown identifier in second arm of AND must panic")
}

// TestValidateHavingIdentifiers_NilExpression verifies nil is handled without panic.
func TestValidateHavingIdentifiers_NilExpression(t *testing.T) {
	require.NotPanics(t, func() {
		validateHavingIdentifiers(nil, map[string]tree.Expression{}, map[string]struct{}{})
	})
}
