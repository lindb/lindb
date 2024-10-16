package analyzer

import (
	"fmt"

	"github.com/lindb/lindb/spi/function"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type Context struct {
	scope *Scope
}

type ExpressionAnalyzer struct {
	ctx             *AnalyzerContext
	funcionResolver *function.FunctionResolver
}

func NewExpressionAnalyzer(ctx *AnalyzerContext) *ExpressionAnalyzer {
	return &ExpressionAnalyzer{
		ctx:             ctx,
		funcionResolver: function.NewFunctionResolver(), // FIXME:???
	}
}

func (a *ExpressionAnalyzer) Analyze(expression tree.Expression, scope *Scope) {
	fmt.Println("expression analyze")
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
	fmt.Printf("visit expression=%T,id=%v\n", n, n.GetID())
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
	case *tree.ArithmeticBinaryExpression:
		return v.visitArithemticBinary(context, node)
	case *tree.TimePredicate:
		return v.visitTimestampPredicate(context, node)
	case *tree.LogicalExpression:
		return v.visitLogicalExpression(context, node)
	case *tree.DereferenceExpression:
		return v.visitDereferenceExpression(context, node)
	case *tree.FunctionCall:
		return v.visitFunctionCall(context, node)
	case *tree.StringLiteral:
		return v.visitStringLiteral(context, node)
	case *tree.LongLiteral:
		return v.visitLongLiteral(context, node)
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
	return v.setExpressionType(node, types.DTString)
}

func (v *ExpressionVisitor) visitFieldReference(context any, node *tree.FieldReference) (r any) {
	ctx := context.(*tree.StackableVisitorContext[*Context])
	resolvedField := v.baseScope.getField(node.FieldIndex)
	return v.handleResolvedField(node, resolvedField, ctx)
}

func (v *ExpressionVisitor) visitComparisonExpression(context any, node *tree.ComparisonExpression) (r any) {
	var operatorType types.OperatorType
	switch node.Operator {
	case tree.ComparisonEQ:
		operatorType = types.Equal
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
	return v.setExpressionType(node, types.DTInt)
}

func (v *ExpressionVisitor) visitLikePredicate(context any, node *tree.LikePredicate) (r any) {
	node.Value.Accept(context, v)
	node.Pattern.Accept(context, v)
	return v.setExpressionType(node, types.DTInt)
}

func (v *ExpressionVisitor) visitRegexPredicate(context any, node *tree.RegexPredicate) (r any) {
	node.Value.Accept(context, v)
	node.Pattern.Accept(context, v)
	return v.setExpressionType(node, types.DTInt)
}

func (v *ExpressionVisitor) visitDereferenceExpression(context any, node *tree.DereferenceExpression) (r any) {
	ctx := context.(*tree.StackableVisitorContext[*Context])
	// FIXME: check all
	qualifiedName := node.ToQualifiedName()
	if qualifiedName != nil {
		resolvedField := ctx.GetContext().scope.tryResolveField(node, qualifiedName)
		fmt.Printf("visit de expre =%v\n", resolvedField)
		if resolvedField != nil {
			return v.handleResolvedField(node, resolvedField, ctx)
		}
	}
	// rowType := &types.RowType{}
	// TODO: fixme
	return v.setExpressionType(node, types.DTString)
}

func (v *ExpressionVisitor) visitFunctionCall(context any, node *tree.FunctionCall) (r any) {
	var argumentTypes []types.DataType
	for _, arg := range node.Arguments {
		argumentTypes = append(argumentTypes, arg.Accept(context, v).(types.DataType))
	}
	expectedType := tree.GetDefaultFuncReturnType(node.Name)
	if len(argumentTypes) > 0 {
		// TODO: check args types
		for i := 0; i < len(argumentTypes); i++ {
			expectedType = types.GetAccurateType(expectedType, argumentTypes[i])
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

func (v *ExpressionVisitor) visitStringLiteral(context any, node *tree.StringLiteral) (r any) {
	return v.setExpressionType(node, types.DTString)
}

func (v *ExpressionVisitor) visitLongLiteral(context any, node *tree.LongLiteral) (r any) {
	return v.setExpressionType(node, types.DTInt)
}

func (v *ExpressionVisitor) visitIntervalLiteral(context any, node *tree.IntervalLiteral) (r any) {
	return v.setExpressionType(node, types.DTDuration)
}

func (v *ExpressionVisitor) visitIdentifier(context any, node *tree.Identifier) (r any) {
	ctx := context.(*tree.StackableVisitorContext[*Context])
	fmt.Printf("expr visitor %V\n", node.Value)
	// FIXME:???
	resolvedField := ctx.GetContext().scope.resolveField(node, tree.NewQualifiedName([]*tree.Identifier{node}), true)

	if resolvedField == nil {
		panic(fmt.Sprintf("unknown column: '%v'", node.Value))
	}
	return v.handleResolvedField(node, resolvedField, ctx)
}

func (v *ExpressionVisitor) visitArithemticBinary(context any, node *tree.ArithmeticBinaryExpression) (r any) {
	// TODO: remove op
	return v.getOperator(context.(*tree.StackableVisitorContext[*Context]), node, types.Subtract, node.Left, node.Right)
}

func (v *ExpressionVisitor) visitTimestampPredicate(context any, node *tree.TimePredicate) (r any) {
	return v.setExpressionType(node, types.DTTimestamp)
}

func (v *ExpressionVisitor) visitLogicalExpression(context any, node *tree.LogicalExpression) (r any) {
	for _, term := range node.Terms {
		// TODO: add coerce type?
		_ = term.Accept(context, v).(types.DataType)
		// v.coerceType(term, activeType, types.DTInt)
	}
	// TODO: set bool
	return v.setExpressionType(node, types.DTInt)
}

func (v *ExpressionVisitor) getOperator(context *tree.StackableVisitorContext[*Context],
	node tree.Expression, operatorType types.OperatorType, arguments ...tree.Expression,
) types.DataType {
	var argumentTypes []types.DataType
	for i := range arguments {
		expression := arguments[i]
		argumentTypes = append(argumentTypes, expression.Accept(context, v).(types.DataType))
	}

	// operatorSignature := v.analyzer.funcionResolver.ResolveOperator(operatorType, nil).Signature

	// TODO: check args types
	expectedType := types.GetAccurateType(argumentTypes[0], argumentTypes[1])
	fmt.Printf("visit arithmetic id=%v left=%v,right=%v,result=%v\n", node.GetID(), argumentTypes[0], argumentTypes[1], expectedType)
	if expectedType == types.DTTimeSeries {
		for i, argumentType := range argumentTypes {
			v.coerceType(arguments[i], argumentType, expectedType)
		}
	}

	return v.setExpressionType(node, expectedType)
}

func (v *ExpressionVisitor) coerceType(expression tree.Expression, actualType, expectedType types.DataType) {
	// TODO: add check
	if actualType != expectedType {
		fmt.Printf("add coercion %v=>%v\n", expression, expectedType)
		v.analyzer.ctx.Analysis.AddCoercion(expression, expectedType)
	}
}

func (v *ExpressionVisitor) handleResolvedField(node tree.Expression, resolvedField *ResolvedField, context *tree.StackableVisitorContext[*Context]) types.DataType {
	v.analyzer.ctx.Analysis.AddColumnReference(node, resolvedField)
	v.analyzer.ctx.Analysis.AddType(node, resolvedField.Field.DataType)
	return resolvedField.Field.DataType
}

func (v *ExpressionVisitor) setExpressionType(expression tree.Expression, expressionType types.DataType) types.DataType {
	v.analyzer.ctx.Analysis.AddType(expression, expressionType)
	return expressionType
}
