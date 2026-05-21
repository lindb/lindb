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

	"github.com/apache/arrow-go/v18/arrow"
	larrow "github.com/lindb/arrow/pkg/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type Context struct {
	scope *Scope
}

type ExpressionAnalyzer struct {
	ctx *AnalyzerContext
}

func NewExpressionAnalyzer(ctx *AnalyzerContext) *ExpressionAnalyzer {
	return &ExpressionAnalyzer{
		ctx: ctx,
	}
}

func (a *ExpressionAnalyzer) Analyze(expression tree.Expression, scope *Scope) {
	visitor := NewExpressionVisitor(scope, a)
	expression.Accept(tree.NewStackableVisitorContext(&Context{
		scope: scope,
	}), visitor)
}

type ExpressionVisitor struct {
	tree.StackableAstVisitor[*Context]
	baseScope *Scope
	analyzer  *ExpressionAnalyzer
}

func NewExpressionVisitor(scope *Scope, analyzer *ExpressionAnalyzer) *ExpressionVisitor {
	return &ExpressionVisitor{
		baseScope: scope,
		analyzer:  analyzer,
	}
}

func (v *ExpressionVisitor) Visit(context any, n tree.Node) (r any) {
	// TODO:
	_ = n.Accept(context, &v.StackableAstVisitor)
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		return v.visitComparisonExpression(context, node)
	case *tree.NotExpression:
		return node.Value.Accept(context, v)
	case *tree.InPredicate:
		return v.visitInPredicate(context, node)
	case *tree.LikePredicate:
		return v.visitLikePredicate(context, node)
	case *tree.RegexPredicate:
		return v.visitRegexPredicate(context, node)
	case *tree.NullPredicate:
		return v.visitNullPredicate(context, node)
	case *tree.ArithmeticBinaryExpression:
		return v.visitArithemticBinary(context, node)
	case *tree.TimePredicate:
		return v.visitTimestampPredicate(context, node)
	case *tree.LogicalExpression:
		return v.visitLogicalExpression(context, node)
	case *tree.DereferenceExpression:
		return v.visitDereferenceExpression(context, node)
	case *tree.SubscriptExpression:
		return v.visitSubscriptExpression(context, node)
	case *tree.FunctionCall:
		return v.visitFunctionCall(context, node)
	case *tree.StringLiteral:
		return v.visitStringLiteral(context, node)
	case *tree.LongLiteral:
		return v.visitLongLiteral(context, node)
	case *tree.FloatLiteral:
		return v.visitFloatLiteral(context, node)
	case *tree.IntervalLiteral:
		return v.visitIntervalLiteral(context, node)
	case *tree.Identifier:
		return v.visitIdentifier(context, node)
	case *tree.FieldReference:
		return v.visitFieldReference(context, node)
	case *tree.Row:
		return v.visitRow(context, node)
	default:
		panic(fmt.Sprintf("expression analyzer unsupport node:%T", n))
	}
}

func (v *ExpressionVisitor) visitRow(context any, node *tree.Row) (r any) {
	for _, item := range node.Items {
		item.Accept(context, v)
	}
	// TODO: change data type
	return v.setExpressionType(node, arrow.BinaryTypes.String)
}

func (v *ExpressionVisitor) visitFieldReference(context any, node *tree.FieldReference) (r any) {
	ctx := context.(*tree.StackableVisitorContext[*Context])
	resolvedField := v.baseScope.getField(node.FieldIndex)
	return v.handleResolvedField(ctx, node, resolvedField)
}

func (v *ExpressionVisitor) visitComparisonExpression(context any, node *tree.ComparisonExpression) (r any) {
	var operatorType types.OperatorType
	switch node.Operator {
	case tree.ComparisonEQ:
		operatorType = types.Equal
	default:
		panic("not supported operator:" + node.Operator)
	}

	return v.getOperator(context.(*tree.StackableVisitorContext[*Context]), node, operatorType, node.Left, node.Right)
}

func (v *ExpressionVisitor) visitInPredicate(context any, node *tree.InPredicate) (r any) {
	node.Value.Accept(context, v)
	if inListExpression, ok := node.ValueList.(*tree.InListExpression); ok {
		for _, value := range inListExpression.Values {
			value.Accept(context, v)
		}
	}
	// TODO: check args types
	// TODO: check all
	return v.setExpressionType(node, arrow.PrimitiveTypes.Uint32)
}

func (v *ExpressionVisitor) visitLikePredicate(context any, node *tree.LikePredicate) (r any) {
	node.Value.Accept(context, v)
	node.Pattern.Accept(context, v)
	return v.setExpressionType(node, arrow.PrimitiveTypes.Uint32)
}

func (v *ExpressionVisitor) visitRegexPredicate(context any, node *tree.RegexPredicate) (r any) {
	node.Value.Accept(context, v)
	node.Pattern.Accept(context, v)
	return v.setExpressionType(node, arrow.PrimitiveTypes.Uint32)
}

func (v *ExpressionVisitor) visitNullPredicate(context any, node *tree.NullPredicate) (r any) {
	node.Value.Accept(context, v)
	return v.setExpressionType(node, arrow.PrimitiveTypes.Uint32)
}

func (v *ExpressionVisitor) visitSubscriptExpression(context any, node *tree.SubscriptExpression) (r any) {
	// Analyze the base (map column) and key expressions; result is the map's value type (String).
	node.Base.Accept(context, v)
	node.Key.Accept(context, v)
	return v.setExpressionType(node, arrow.BinaryTypes.String)
}

func (v *ExpressionVisitor) visitDereferenceExpression(context any, node *tree.DereferenceExpression) (r any) {
	ctx := context.(*tree.StackableVisitorContext[*Context])
	// FIXME: check all
	qualifiedName := node.ToQualifiedName()
	if qualifiedName != nil {
		resolvedField := ctx.GetContext().scope.tryResolveField(node, qualifiedName)
		if resolvedField != nil {
			return v.handleResolvedField(ctx, node, resolvedField)
		}
	}
	// rowType := &types.RowType{}
	// TODO: fixme
	return v.setExpressionType(node, arrow.BinaryTypes.String)
}

func (v *ExpressionVisitor) visitFunctionCall(context any, node *tree.FunctionCall) (r any) {
	var argumentTypes []arrow.DataType
	for _, arg := range node.Arguments {
		// For histogram functions, the column-name argument may be a logical histogram
		// field name (e.g. "sent_duration") that does not exist as a direct schema column.
		// Resolve it gracefully instead of panicking, so the planner can expand it later.
		if tree.IsHistogramFunc(node.Name) {
			if ident, ok := arg.(*tree.Identifier); ok {
				argumentTypes = append(argumentTypes, v.resolveHistogramColumn(context, ident))
				continue
			}
		}
		argumentTypes = append(argumentTypes, arg.Accept(context, v).(arrow.DataType))
	}
	expectedType := v.analyzer.ctx.GetFuncReturnType(node.Name)
	if expectedType == nil {
		if len(argumentTypes) > 0 {
			// TODO: check args types
			for i := range len(argumentTypes) {
				expectedType = types.GetAccurateType(expectedType, argumentTypes[i])
			}
		}
	}

	// TODO: coerce args types
	// for i, argumentType := range argumentTypes {
	// 	v.coerceType(node.Arguments[i], argumentType, expectedType)
	// }
	// FIXME:func call???
	// rowType := &types.RowType{}
	return v.setExpressionType(node, expectedType)
}

// resolveHistogramColumn resolves an Identifier argument of a histogram function.
// It first tries normal schema column resolution. If the column is not found as a
// direct field (e.g. "sent_duration" is a logical histogram name, not a physical column),
// it falls back to registering it as a Sum-type reference so the planner can expand it
// to the physical bucket/stat columns via expandHistogramColumns.
func (v *ExpressionVisitor) resolveHistogramColumn(context any, node *tree.Identifier) arrow.DataType {
	ctx := context.(*tree.StackableVisitorContext[*Context])
	resolvedField := ctx.GetContext().scope.resolveField(
		node, tree.NewQualifiedName([]*tree.Identifier{node}), true)
	if resolvedField != nil {
		return v.handleResolvedField(ctx, node, resolvedField)
	}
	// Logical histogram name not found as a direct column — treat as valid histogram reference.
	// The planner's expandHistogramColumns will match it against physical fields by prefix.
	return v.setExpressionType(node, larrow.ExtensionTypes.Histogram)
}

func (v *ExpressionVisitor) visitStringLiteral(_ any, node *tree.StringLiteral) (r any) {
	return v.setExpressionType(node, arrow.BinaryTypes.String)
}

func (v *ExpressionVisitor) visitLongLiteral(_ any, node *tree.LongLiteral) (r any) {
	return v.setExpressionType(node, arrow.PrimitiveTypes.Int64)
}

func (v *ExpressionVisitor) visitFloatLiteral(_ any, node *tree.FloatLiteral) (r any) {
	return v.setExpressionType(node, arrow.PrimitiveTypes.Float64)
}

func (v *ExpressionVisitor) visitIntervalLiteral(_ any, node *tree.IntervalLiteral) (r any) {
	return v.setExpressionType(node, arrow.FixedWidthTypes.Duration_ns)
}

func (v *ExpressionVisitor) visitIdentifier(context any, node *tree.Identifier) (r any) {
	ctx := context.(*tree.StackableVisitorContext[*Context])
	// FIXME:???
	resolvedField := ctx.GetContext().scope.resolveField(node, tree.NewQualifiedName([]*tree.Identifier{node}), true)

	if resolvedField == nil {
		panic(fmt.Sprintf("unknown column: '%v'", node.Value))
	}
	return v.handleResolvedField(ctx, node, resolvedField)
}

func (v *ExpressionVisitor) visitArithemticBinary(context any, node *tree.ArithmeticBinaryExpression) (r any) {
	// TODO: remove op
	return v.getOperator(context.(*tree.StackableVisitorContext[*Context]), node, types.Subtract, node.Left, node.Right)
}

func (v *ExpressionVisitor) visitTimestampPredicate(_ any, node *tree.TimePredicate) (r any) {
	return v.setExpressionType(node, arrow.FixedWidthTypes.Timestamp_ns)
}

func (v *ExpressionVisitor) visitLogicalExpression(context any, node *tree.LogicalExpression) (r any) {
	for _, term := range node.Terms {
		// TODO: add coerce type?
		_ = term.Accept(context, v).(arrow.DataType)
		// TODO: v.coerceType(term, activeType, types.DTInt)
	}
	// TODO: set bool
	return v.setExpressionType(node, arrow.PrimitiveTypes.Uint32)
}

func (v *ExpressionVisitor) getOperator(context *tree.StackableVisitorContext[*Context],
	node tree.Expression, _ types.OperatorType, arguments ...tree.Expression,
) arrow.DataType {
	var argumentTypes []arrow.DataType
	for i := range arguments {
		expression := arguments[i]
		argumentTypes = append(argumentTypes, expression.Accept(context, v).(arrow.DataType))
	}

	// TODO: operatorSignature := v.analyzer.funcionResolver.ResolveOperator(operatorType, nil).Signature

	// Handle AggregationType OP numeric (e.g. idle*100): result keeps the AggregationType,
	// and the numeric side is used as a scalar factor directly — no Cast needed.
	if aggType, _, ok := resolveAggNumeric(argumentTypes[0], argumentTypes[1]); ok {
		return v.setExpressionType(node, aggType)
	}

	// TODO: check args types
	expectedType := types.GetAccurateType(argumentTypes[0], argumentTypes[1])
	if expectedType == larrow.ExtensionTypes.TimeSeries {
		for i, argumentType := range argumentTypes {
			v.coerceType(arguments[i], argumentType, expectedType)
		}
	}

	return v.setExpressionType(node, expectedType)
}

// resolveAggNumeric checks if exactly one side is an AggregationType and the other is a
// numeric scalar (Int64 or Float64). Returns the AggregationType, the index of the numeric
// argument, and true when the pattern matches; otherwise returns (nil, -1, false).
func resolveAggNumeric(lhs, rhs arrow.DataType) (aggType arrow.DataType, numericIdx int, ok bool) {
	lhsIsAgg := isAggregationType(lhs)
	rhsIsAgg := isAggregationType(rhs)
	lhsIsNum := isNumericType(lhs)
	rhsIsNum := isNumericType(rhs)

	switch {
	case lhsIsAgg && rhsIsNum:
		return lhs, 1, true
	case rhsIsAgg && lhsIsNum:
		return rhs, 0, true
	default:
		return nil, -1, false
	}
}

// isAggregationType reports whether t is a lindb Arrow extension AggregationType.
func isAggregationType(t arrow.DataType) bool {
	_, ok := t.(*larray.AggregationType)
	return ok
}

// isNumericType reports whether t is a plain numeric scalar (Int64 or Float64).
func isNumericType(t arrow.DataType) bool {
	return arrow.TypeEqual(t, arrow.PrimitiveTypes.Int64) ||
		arrow.TypeEqual(t, arrow.PrimitiveTypes.Float64)
}

func (v *ExpressionVisitor) coerceType(expression tree.Expression, actualType, expectedType arrow.DataType) {
	// TODO: add check
	if actualType != expectedType {
		v.analyzer.ctx.Analysis.AddCoercion(expression, expectedType)
	}
}

func (v *ExpressionVisitor) handleResolvedField(_ *tree.StackableVisitorContext[*Context],
	node tree.Expression, resolvedField *ResolvedField,
) arrow.DataType {
	v.analyzer.ctx.Analysis.AddColumnReference(node, resolvedField)
	v.analyzer.ctx.Analysis.AddType(node, resolvedField.Field.DataType)
	return resolvedField.Field.DataType
}

func (v *ExpressionVisitor) setExpressionType(expression tree.Expression, expressionType arrow.DataType) arrow.DataType {
	v.analyzer.ctx.Analysis.AddType(expression, expressionType)
	return expressionType
}
