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
	"fmt"

	"github.com/lindb/lindb/sql/tree"
)

// verifySourceAggregations checks that every output expression is either an aggregate
// function call or a column that appears in the grouping set.
// groupingFields is the flattened list of FieldIDs from the GROUP BY / implicit-timestamp sets.
func verifySourceAggregations(analysis *Analysis, groupingFields []*FieldID, expressions []tree.Expression) {
	analyzer := NewAggregationAnalyzer(analysis, groupingFields)
	for _, expression := range expressions {
		analyzer.analyze(expression)
	}
}

// verifyOrderByAggregations applies the same grouping-key check to ORDER BY expressions.
func verifyOrderByAggregations(analysis *Analysis, groupingFields []*FieldID, expressions []tree.Expression) {
	analyzer := NewAggregationAnalyzer(analysis, groupingFields)
	for _, expression := range expressions {
		analyzer.analyze(expression)
	}
}

type AggregationAnanlyzer struct {
	analysis       *Analysis
	groupingFields []*FieldID // FieldIDs that are valid non-aggregate references
}

func NewAggregationAnalyzer(analysis *Analysis, groupingFields []*FieldID) *AggregationAnanlyzer {
	return &AggregationAnanlyzer{
		analysis:       analysis,
		groupingFields: groupingFields,
	}
}

func (aa *AggregationAnanlyzer) analyze(expression tree.Expression) {
	visitor := NewAggregationAnalyzeVisitor(aa.analysis, aa.groupingFields)
	if r, ok := expression.Accept(nil, visitor).(bool); ok {
		if !r {
			panic(fmt.Sprintf("'%s' must be an aggregate expression or appear in GROUP BY clause",
				tree.FormatExpression(expression)))
		}
	}
}

type aggregationAnalyzeVisitor struct {
	analysis       *Analysis
	groupingFields []*FieldID
}

func NewAggregationAnalyzeVisitor(analysis *Analysis, groupingFields []*FieldID) tree.Visitor {
	return &aggregationAnalyzeVisitor{
		analysis:       analysis,
		groupingFields: groupingFields,
	}
}

func (v *aggregationAnalyzeVisitor) Visit(context any, n tree.Node) (r any) {
	switch node := n.(type) {
	case *tree.DereferenceExpression:
		return v.visitDereferenceExpression(node)
	case *tree.Identifier:
		return v.visitIdentifier(node)
	case *tree.ArithmeticBinaryExpression:
		// TODO: modify
		return true
	case *tree.FunctionCall:
		// FIXME: add logic
		return true
	case *tree.StringLiteral, *tree.LongLiteral, *tree.FloatLiteral, *tree.BooleanLiteral, *tree.Constant:
		// Literals are never column references, so they are always valid in any aggregate query.
		return true
	case *tree.ComparisonExpression:
		// Both sides must independently be valid aggregate or grouping expressions.
		leftOk, _ := node.Left.Accept(context, v).(bool)
		rightOk, _ := node.Right.Accept(context, v).(bool)
		return leftOk && rightOk
	case *tree.LogicalExpression:
		// All terms must be valid for AND; any term valid is sufficient for OR.
		// In practice we require all terms to be valid (conservative).
		for _, term := range node.Terms {
			if ok, _ := term.Accept(context, v).(bool); !ok {
				return false
			}
		}
		return true
	case *tree.NotExpression:
		return node.Value.Accept(context, v)
	default:
		panic(fmt.Sprintf("unsupported node<%T> when aggregation ananlyzer", n))
	}
}

// visitIdentifier checks whether an identifier refers to a column that is part of
// the grouping set.  A column reference that is not in the grouping set is invalid
// in an aggregating query (SQL: "must appear in GROUP BY clause").
func (v *aggregationAnalyzeVisitor) visitIdentifier(node *tree.Identifier) (r any) {
	resolvedField := v.analysis.GetColumnReferenceField(node)
	if resolvedField == nil {
		// Not a plain column reference (e.g. a literal alias) — allow it.
		return true
	}

	// Check whether this column is one of the grouping keys by comparing FieldID
	// (RelationID pointer identity + FieldIndex value).
	fieldID := resolvedField.FieldID()
	for _, gf := range v.groupingFields {
		if gf.RelationID == fieldID.RelationID && gf.FieldIndex == fieldID.FieldIndex {
			return true
		}
	}
	return false
}

func (v *aggregationAnalyzeVisitor) visitDereferenceExpression(node *tree.DereferenceExpression) (r any) {
	// 		            ExpressionAnalyzer.LabelPrefixedReference labelDereference = analysis.getLabelDereference(node);
	// if (labelDereference != null) {
	//     return labelDereference.getColumn().map(this::process).orElse(true);
	// }
	//
	// if (!hasReferencesToScope(node, analysis, sourceScope)) {
	//     // reference to outer scope is group-invariant
	//     return true;
	// }
	//
	// if (columnReferences.containsKey(NodeRef.<Expression>of(node))) {
	//     return isGroupingKey(node);
	// }
	//

	// Allow SELECT col1.f1 FROM table1 GROUP BY col1
	return node.Base.Accept(nil, v)
}
