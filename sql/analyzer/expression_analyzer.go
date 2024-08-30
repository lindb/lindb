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
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		return v.visitComparisonExpression(context, node)
	case *tree.DereferenceExpression:
		return v.visitDereferenceExpression(context, node)
	case *tree.FunctionCall:
		return v.visitFunctionCall(context, node)
	case *tree.StringLiteral:
		return v.visitStringLiteral(context, node)
	case *tree.LongLiteral:
		return v.visitLongLiteral(context, node)
	case *tree.Identifier:
		return v.visitIdentifier(context, node)
	case *tree.FieldReference:
		return v.visitFieldReference(context, node)
	case *tree.ArithmeticBinaryExpression:
		return v.visitArithemticBinary(context, node)
	default:
		panic(fmt.Sprintf("expression analyzer unsupport node:%T", n))
	}
}

func (v *ExpressionVisitor) visitFieldReference(context any, node *tree.FieldReference) (r any) {
	ctx := context.(*tree.StackableVisitorContext[*Context])
	resolvedField := v.baseScope.getField(node.FieldIndex)
	return v.handleResolvedField(node, resolvedField, ctx)
}

func (v *ExpressionVisitor) visitComparisonExpression(context any, node *tree.ComparisonExpression) (r any) {
	var operatorType function.OperatorType
	switch node.Operator {
	case tree.ComparisonEqual:
		operatorType = function.Equal
	}

	return v.getOperator(context.(*tree.StackableVisitorContext[*Context]), node, operatorType, node.Left, node.Right)
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
	return v.setExpressionType(node, types.DataTypeString)
}

func (v *ExpressionVisitor) visitFunctionCall(context any, node *tree.FunctionCall) (r any) {
	// FIXME:func call???
	// rowType := &types.RowType{}
	return v.setExpressionType(node, types.DataTypeFirst)
}

func (v *ExpressionVisitor) visitStringLiteral(context any, node *tree.StringLiteral) (r any) {
	// FIXME:???
	// rowType := &types.RowType{}
	return v.setExpressionType(node, types.DataTypeString)
}

func (v *ExpressionVisitor) visitLongLiteral(context any, node *tree.LongLiteral) (r any) {
	return v.setExpressionType(node, types.DataTypeFloat)
}

func (v *ExpressionVisitor) visitIdentifier(context any, node *tree.Identifier) (r any) {
	ctx := context.(*tree.StackableVisitorContext[*Context])
	fmt.Printf("expr visitor %V\n", node.Value)
	// FIXME:???
	resolvedField := ctx.GetContext().scope.resolveField(node, tree.NewQualifiedName([]*tree.Identifier{node}), true)
	return v.handleResolvedField(node, resolvedField, ctx)
}

func (v *ExpressionVisitor) visitArithemticBinary(context any, node *tree.ArithmeticBinaryExpression) (r any) {
	return v.getOperator(context.(*tree.StackableVisitorContext[*Context]), node, function.Add, node.Left, node.Right)
}

func (v *ExpressionVisitor) getOperator(context *tree.StackableVisitorContext[*Context], node tree.Expression, operatorType function.OperatorType, arguments ...tree.Expression) types.Type {
	var argumentTypes []types.Type
	for i := range arguments {
		expression := arguments[i]
		argumentTypes = append(argumentTypes, expression.Accept(context, v).(types.Type))
	}

	operatorSignature := v.analyzer.funcionResolver.ResolveOperator(operatorType, argumentTypes).Signature

	fmt.Println(operatorSignature)
	return v.setExpressionType(node, operatorSignature.ReturnType)
}

func (v *ExpressionVisitor) handleResolvedField(node tree.Expression, resolvedField *ResolvedField, context *tree.StackableVisitorContext[*Context]) types.Type {
	v.analyzer.ctx.Analysis.AddColumnReference(node, resolvedField)
	v.analyzer.ctx.Analysis.AddType(node, resolvedField.Field.DataType)
	return &types.RowType{}
}

func (v *ExpressionVisitor) setExpressionType(expression tree.Expression, expressionType types.DataType) types.DataType {
	v.analyzer.ctx.Analysis.AddType(expression, expressionType)
	return expressionType
}
