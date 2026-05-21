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

package rule

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

// makeAgg is a helper to build a minimal AggregationAssignment.
func makeAgg(funcName tree.FuncName, args ...tree.Expression) *plan.AggregationAssignment {
	return &plan.AggregationAssignment{
		Aggregation: &plan.Aggregation{
			Function:  funcName,
			Arguments: args,
		},
	}
}

// symRef builds a SymbolReference with the given name (no DataType).
func symRef(name string) *tree.SymbolReference {
	return &tree.SymbolReference{Name: name}
}

// symRefWithType builds a SymbolReference with an explicit DataType,
// simulating what planAggregation produces after converting arguments to SymbolReferences.
func symRefWithType(name string, dt arrow.DataType) *tree.SymbolReference {
	return &tree.SymbolReference{Name: name, DataType: dt}
}

// floatLit builds a FloatLiteral with the given value.
func floatLit(v string) *tree.FloatLiteral {
	return tree.NewFloatLiteral(0, nil, v)
}

// longLit builds a LongLiteral with the given value.
func longLit(v string) *tree.LongLiteral {
	return tree.NewLongLiteral(0, nil, v)
}

// ----- buildColumnAggregations -----

// TestBuildColumnAggregations_SimpleSum verifies a plain sum(col) call.
func TestBuildColumnAggregations_SimpleSum(t *testing.T) {
	aggs := []*plan.AggregationAssignment{
		makeAgg(tree.Sum, symRef("cpu")),
	}
	result := buildColumnAggregations(aggs)

	require.Len(t, result, 1)
	assert.Equal(t, "cpu", result[0].Column)
	assert.Equal(t, tree.Sum, result[0].AggFuncName)
}

// TestBuildColumnAggregations_HistogramQuantile_FloatPhi verifies that
// histogram_quantile(0.99, col) correctly extracts "col" as the column name
// while ignoring the phi literal.
func TestBuildColumnAggregations_HistogramQuantile_FloatPhi(t *testing.T) {
	phi := floatLit("0.99")
	aggs := []*plan.AggregationAssignment{
		makeAgg(tree.HistogramQuantile, phi, symRef("duration")),
	}
	result := buildColumnAggregations(aggs)

	require.Len(t, result, 1)
	assert.Equal(t, "duration", result[0].Column)
	assert.Equal(t, tree.HistogramQuantile, result[0].AggFuncName)
}

// TestBuildColumnAggregations_HistogramQuantile_LongPhi verifies that
// an integer phi (e.g. histogram_quantile(1, col)) is handled correctly.
func TestBuildColumnAggregations_HistogramQuantile_LongPhi(t *testing.T) {
	phi := longLit("1")
	aggs := []*plan.AggregationAssignment{
		makeAgg(tree.HistogramQuantile, phi, symRef("latency")),
	}
	result := buildColumnAggregations(aggs)

	require.Len(t, result, 1)
	assert.Equal(t, "latency", result[0].Column)
	assert.Equal(t, tree.HistogramQuantile, result[0].AggFuncName)
}

// TestBuildColumnAggregations_MultiplePhiArgs verifies that multiple literal
// arguments still yield the correct column name.
func TestBuildColumnAggregations_MultiplePhiArgs(t *testing.T) {
	aggs := []*plan.AggregationAssignment{
		makeAgg(tree.HistogramQuantile, floatLit("0.5"), floatLit("0.99"), symRef("rtt")),
	}
	result := buildColumnAggregations(aggs)

	require.Len(t, result, 1)
	assert.Equal(t, "rtt", result[0].Column)
}

// TestBuildColumnAggregations_NoColumnName verifies that an aggregation with no
// SymbolReference or Constant argument produces no entry.
func TestBuildColumnAggregations_NoColumnName(t *testing.T) {
	aggs := []*plan.AggregationAssignment{
		makeAgg(tree.HistogramQuantile, floatLit("0.99")), // phi only, no column
	}
	result := buildColumnAggregations(aggs)
	assert.Empty(t, result, "no column name → entry must be skipped")
}

// TestBuildColumnAggregations_EmptyAggregations verifies that an empty input
// produces an empty output without panic.
func TestBuildColumnAggregations_EmptyAggregations(t *testing.T) {
	result := buildColumnAggregations(nil)
	assert.Empty(t, result)

	result = buildColumnAggregations([]*plan.AggregationAssignment{})
	assert.Empty(t, result)
}

// TestBuildColumnAggregations_MultipleAggregations verifies that multiple
// aggregations are all processed independently.
func TestBuildColumnAggregations_MultipleAggregations(t *testing.T) {
	aggs := []*plan.AggregationAssignment{
		makeAgg(tree.Sum, symRef("cpu")),
		makeAgg(tree.HistogramQuantile, floatLit("0.95"), symRef("latency")),
		makeAgg(tree.HistogramAvg, symRef("rtt")),
	}
	result := buildColumnAggregations(aggs)

	require.Len(t, result, 3)
	assert.Equal(t, "cpu", result[0].Column)
	assert.Equal(t, tree.Sum, result[0].AggFuncName)
	assert.Equal(t, "latency", result[1].Column)
	assert.Equal(t, tree.HistogramQuantile, result[1].AggFuncName)
	assert.Equal(t, "rtt", result[2].Column)
	assert.Equal(t, tree.HistogramAvg, result[2].AggFuncName)
}

// TestBuildColumnAggregations_FirstSymbolRefWins verifies that only the first
// SymbolReference is used as the column name.
func TestBuildColumnAggregations_FirstSymbolRefWins(t *testing.T) {
	aggs := []*plan.AggregationAssignment{
		makeAgg(tree.Sum, symRef("first_col"), symRef("second_col")),
	}
	result := buildColumnAggregations(aggs)

	require.Len(t, result, 1)
	assert.Equal(t, "first_col", result[0].Column, "first SymbolReference should win")
}

// TestBuildColumnAggregations_HistogramQuantile_SymbolRefArgs is a regression test for
// the production bug where planAggregation converted FloatLiteral(0.99) to
// SymbolRef{Name:"expr", DataType:Float64}.  The OLD code picked "expr" as column name.
// The NEW code uses isNumericSymbolRef to skip numeric SymbolReferences.
func TestBuildColumnAggregations_HistogramQuantile_SymbolRefArgs(t *testing.T) {
	phiRef := symRefWithType("expr", arrow.PrimitiveTypes.Float64)
	colRef := symRefWithType("sent_duration", larrow.ExtensionTypes.Histogram)

	aggs := []*plan.AggregationAssignment{
		makeAgg(tree.HistogramQuantile, phiRef, colRef),
	}
	result := buildColumnAggregations(aggs)

	require.Len(t, result, 1)
	assert.Equal(t, "sent_duration", result[0].Column,
		"numeric SymbolRef (Float64) must not be used as column name")
	assert.Equal(t, tree.HistogramQuantile, result[0].AggFuncName)
}

// TestIsNumericSymbolRef verifies the helper that distinguishes phi SymbolRefs from column SymbolRefs.
func TestIsNumericSymbolRef(t *testing.T) {
	assert.True(t, isNumericSymbolRef(symRefWithType("expr", arrow.PrimitiveTypes.Float64)))
	assert.True(t, isNumericSymbolRef(symRefWithType("x", arrow.PrimitiveTypes.Int64)))
	assert.False(t, isNumericSymbolRef(symRefWithType("sent_duration", larrow.ExtensionTypes.Histogram)))
	assert.False(t, isNumericSymbolRef(symRefWithType("cpu", larrow.ExtensionTypes.Sum)))
	assert.False(t, isNumericSymbolRef(symRef("no_type"))) // nil DataType → false
}
