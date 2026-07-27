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
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi"
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

// TestBuildColumnAggregations_NoColumnName verifies that an aggregation whose
// arguments are literal-only (e.g. histogram_quantile(0.99) missing the column ref)
// is an incomplete call and must be skipped.
func TestBuildColumnAggregations_NoColumnName(t *testing.T) {
	aggs := []*plan.AggregationAssignment{
		makeAgg(tree.HistogramQuantile, floatLit("0.99")), // phi only, no column
	}
	result := buildColumnAggregations(aggs)
	assert.Empty(t, result, "literal-only args with no column name → entry must be skipped")
}

// TestBuildColumnAggregations_Count verifies that COUNT(*) (zero arguments) is included
// using the function name as a synthetic column placeholder so push-down can proceed.
// The storage aggregator identifies aggregation via handle.Aggregation != "", not by column name.
func TestBuildColumnAggregations_Count(t *testing.T) {
	aggs := []*plan.AggregationAssignment{
		makeAgg(tree.Count), // COUNT(*) — no arguments
	}
	result := buildColumnAggregations(aggs)
	require.Len(t, result, 1, "COUNT(*) must be included for push-down")
	assert.Equal(t, string(tree.Count), result[0].Column, "column should be function name placeholder")
	assert.Equal(t, tree.Count, result[0].AggFuncName)
}

// TestBuildColumnAggregations_CountLiteral verifies that COUNT(1) (literal argument,
// no column reference) is also included, because the literal 1 is not a column name
// but COUNT still means "count all matching rows".
func TestBuildColumnAggregations_CountLiteral(t *testing.T) {
	aggs := []*plan.AggregationAssignment{
		makeAgg(tree.Count, longLit("1")), // COUNT(1) — literal argument
	}
	result := buildColumnAggregations(aggs)
	require.Len(t, result, 1, "COUNT(1) must be included for push-down")
	assert.Equal(t, string(tree.Count), result[0].Column, "column should be function name placeholder")
	assert.Equal(t, tree.Count, result[0].AggFuncName)
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

// ----- hasTimestampGroupingKey -----

// TestHasTimestampGroupingKey_NilGroupingSets verifies nil GroupingSets → false.
func TestHasTimestampGroupingKey_NilGroupingSets(t *testing.T) {
	node := &plan.AggregationNode{}
	assert.False(t, hasTimestampGroupingKey(node))
}

// TestHasTimestampGroupingKey_EmptyKeys verifies empty GroupingKeys → false.
func TestHasTimestampGroupingKey_EmptyKeys(t *testing.T) {
	node := &plan.AggregationNode{
		GroupingSets: &plan.GroupingSetDescriptor{GroupingKeys: nil},
	}
	assert.False(t, hasTimestampGroupingKey(node))
}

// TestHasTimestampGroupingKey_HasTimestamp verifies that a timestamp key is detected.
func TestHasTimestampGroupingKey_HasTimestamp(t *testing.T) {
	tsSym := &plan.Symbol{Name: constants.TimestampColumnName}
	node := &plan.AggregationNode{
		GroupingSets: &plan.GroupingSetDescriptor{GroupingKeys: []*plan.Symbol{tsSym}},
	}
	assert.True(t, hasTimestampGroupingKey(node))
}

// TestHasTimestampGroupingKey_NoTimestamp verifies that non-timestamp keys return false.
func TestHasTimestampGroupingKey_NoTimestamp(t *testing.T) {
	node := &plan.AggregationNode{
		GroupingSets: &plan.GroupingSetDescriptor{
			GroupingKeys: []*plan.Symbol{{Name: "level"}, {Name: "service"}},
		},
	}
	assert.False(t, hasTimestampGroupingKey(node))
}

// ----- buildGroupingKeysWithSumAssignments -----

// TestBuildGroupingKeysWithSumAssignments_NoGrouping verifies that with no grouping
// keys, only assignment symbols (Sum type) are returned.
func TestBuildGroupingKeysWithSumAssignments_NoGrouping(t *testing.T) {
	node := &plan.AggregationNode{
		GroupingSets: &plan.GroupingSetDescriptor{},
	}
	assignments := []*spi.ColumnAssignment{
		{Column: "count"},
	}
	out := buildGroupingKeysWithSumAssignments(node, assignments)
	require.Len(t, out, 1)
	assert.Equal(t, "count", out[0].Name)
	assert.True(t, arrow.TypeEqual(out[0].DataType, larray.NewAggregationType(larray.Sum)),
		"expected Sum type, got %v", out[0].DataType)
}

// TestBuildGroupingKeysWithSumAssignments_WithGrouping verifies that grouping keys
// are prepended before assignment symbols.
func TestBuildGroupingKeysWithSumAssignments_WithGrouping(t *testing.T) {
	levelSym := &plan.Symbol{Name: "level", DataType: arrow.BinaryTypes.String}
	node := &plan.AggregationNode{
		GroupingSets: &plan.GroupingSetDescriptor{GroupingKeys: []*plan.Symbol{levelSym}},
	}
	assignments := []*spi.ColumnAssignment{
		{Column: "count"},
	}
	out := buildGroupingKeysWithSumAssignments(node, assignments)
	require.Len(t, out, 2)
	assert.Equal(t, levelSym, out[0], "first symbol should be the grouping key")
	assert.Equal(t, "count", out[1].Name)
	assert.True(t, arrow.TypeEqual(out[1].DataType, larray.NewAggregationType(larray.Sum)))
}

// ----- buildTableScanOutputSymbols: log table without timestamp grouping -----

// mockLogTable is a minimal spi.TableHandle for log table tests.
type mockLogTable struct{}

func (m *mockLogTable) SetTimeRange(_ timeutil.TimeRange) {}
func (m *mockLogTable) GetTimeRange() timeutil.TimeRange  { return timeutil.TimeRange{} }
func (m *mockLogTable) SetInterval(_ timeutil.Interval)   {}
func (m *mockLogTable) GetInterval() timeutil.Interval    { return 0 }
func (m *mockLogTable) Kind() spi.DatasourceKind          { return spi.Log }
func (m *mockLogTable) String() string                    { return "log" }

// mockMetricTable is a minimal spi.TableHandle for metric table tests.
type mockMetricTable struct{}

func (m *mockMetricTable) SetTimeRange(_ timeutil.TimeRange) {}
func (m *mockMetricTable) GetTimeRange() timeutil.TimeRange  { return timeutil.TimeRange{} }
func (m *mockMetricTable) SetInterval(_ timeutil.Interval)   {}
func (m *mockMetricTable) GetInterval() timeutil.Interval    { return 0 }
func (m *mockMetricTable) Kind() spi.DatasourceKind          { return spi.Metric }
func (m *mockMetricTable) String() string                    { return "metric" }

// TestBuildTableScanOutputSymbols_LogNoTimestamp verifies that for a log table
// without a timestamp grouping key, the output symbols use Sum type (not TimeSeries).
// This matches aggregatorByField which emits Sum-typed values, not TimeSeries.
func TestBuildTableScanOutputSymbols_LogNoTimestamp(t *testing.T) {
	node := &plan.AggregationNode{
		Aggregations: []*plan.AggregationAssignment{
			makeAgg(tree.Count),
		},
		GroupingSets: &plan.GroupingSetDescriptor{},
	}
	assignments := []*spi.ColumnAssignment{
		{Column: "count"},
	}
	out := buildTableScanOutputSymbols(node, assignments, &mockLogTable{})
	require.Len(t, out, 1)
	assert.Equal(t, "count", out[0].Name)
	assert.True(t, arrow.TypeEqual(out[0].DataType, larray.NewAggregationType(larray.Sum)),
		"log table without timestamp → Sum type, got %v", out[0].DataType)
}

// TestBuildTableScanOutputSymbols_LogWithTimestamp verifies that for a log table
// with a timestamp grouping key, the original output symbols are kept unchanged.
// aggregatorByTime emits TimeSeries values, so we keep the SQL-level symbols.
func TestBuildTableScanOutputSymbols_LogWithTimestamp(t *testing.T) {
	tsSym := &plan.Symbol{Name: constants.TimestampColumnName, DataType: arrow.FixedWidthTypes.Timestamp_ns}
	countSym := &plan.Symbol{Name: "count", DataType: larrow.ExtensionTypes.TimeSeries}
	node := &plan.AggregationNode{
		Aggregations: []*plan.AggregationAssignment{
			makeAgg(tree.Count),
		},
		GroupingSets: &plan.GroupingSetDescriptor{
			GroupingKeys: []*plan.Symbol{tsSym},
		},
		Outputs: []*plan.Symbol{tsSym, countSym},
	}
	assignments := []*spi.ColumnAssignment{
		{Column: "count"},
	}
	out := buildTableScanOutputSymbols(node, assignments, &mockLogTable{})
	// With timestamp grouping, aggregatorByTime is used → keep original symbols (TimeSeries).
	require.Len(t, out, 2)
	assert.Equal(t, tsSym, out[0])
	assert.Equal(t, countSym, out[1])
}

// TestBuildTableScanOutputSymbols_MetricTable verifies that for a metric table
// (non-histogram, no log-specific logic), the original output symbols are returned unchanged.
func TestBuildTableScanOutputSymbols_MetricTable(t *testing.T) {
	countSym := &plan.Symbol{Name: "count", DataType: larrow.ExtensionTypes.TimeSeries}
	node := &plan.AggregationNode{
		Aggregations: []*plan.AggregationAssignment{
			makeAgg(tree.Count),
		},
		GroupingSets: &plan.GroupingSetDescriptor{},
		Outputs:      []*plan.Symbol{countSym},
	}
	assignments := []*spi.ColumnAssignment{
		{Column: "count"},
	}
	out := buildTableScanOutputSymbols(node, assignments, &mockMetricTable{})
	// Metric table: no log-specific override → original symbols returned.
	require.Len(t, out, 1)
	assert.Equal(t, countSym, out[0])
}

// TestBuildTableScanOutputSymbols_HistogramLog verifies that histogram aggregations
// on a log table use Sum-typed symbols (broker layer computes the final function).
func TestBuildTableScanOutputSymbols_HistogramLog(t *testing.T) {
	node := &plan.AggregationNode{
		Aggregations: []*plan.AggregationAssignment{
			makeAgg(tree.HistogramQuantile, floatLit("0.99"), symRef("duration")),
		},
		GroupingSets: &plan.GroupingSetDescriptor{},
	}
	assignments := []*spi.ColumnAssignment{
		{Column: "duration_bucket"},
		{Column: "duration_count"},
		{Column: "duration_sum"},
	}
	out := buildTableScanOutputSymbols(node, assignments, &mockLogTable{})
	require.Len(t, out, 3)
	for _, sym := range out {
		assert.True(t, arrow.TypeEqual(sym.DataType, larray.NewAggregationType(larray.Sum)),
			"histogram push-down → all symbols must be Sum type, got %v for %s", sym.DataType, sym.Name)
	}
}

// TestBuildTableScanOutputSymbols_HistogramWithGrouping verifies that histogram
// aggregations with grouping keys keep grouping keys as-is and use Sum for assignments.
func TestBuildTableScanOutputSymbols_HistogramWithGrouping(t *testing.T) {
	levelSym := &plan.Symbol{Name: "level", DataType: arrow.BinaryTypes.String}
	node := &plan.AggregationNode{
		Aggregations: []*plan.AggregationAssignment{
			makeAgg(tree.HistogramQuantile, floatLit("0.99"), symRef("duration")),
		},
		GroupingSets: &plan.GroupingSetDescriptor{
			GroupingKeys: []*plan.Symbol{levelSym},
		},
	}
	assignments := []*spi.ColumnAssignment{
		{Column: "duration_bucket"},
	}
	out := buildTableScanOutputSymbols(node, assignments, &mockMetricTable{})
	require.Len(t, out, 2)
	assert.Equal(t, levelSym, out[0], "grouping key should be kept as-is")
	assert.Equal(t, "duration_bucket", out[1].Name)
	assert.True(t, arrow.TypeEqual(out[1].DataType, larray.NewAggregationType(larray.Sum)))
}
