package tree

import (
	"fmt"
)

type ExpressionRewriter interface {
	RewriteExpression(context any, node Expression) Expression
}

type ExpressionTreeRewriter struct {
	rewriter ExpressionRewriter
	visitor  Visitor
}

func NewExpressionTreeRewriter(rewriter ExpressionRewriter) *ExpressionTreeRewriter {
	return &ExpressionTreeRewriter{
		rewriter: rewriter,
		visitor:  NewExpressionRewriteVisitor(rewriter),
	}
}

func (etr *ExpressionTreeRewriter) rewrite(context any, node Expression) Expression {
	return etr.visitor.Visit(context, node).(Expression)
}

func RewriteExpression(context any, rewriter ExpressionRewriter, node Expression) Expression {
	return NewExpressionTreeRewriter(rewriter).rewrite(context, node)
}

type ExpressionRewriteVisitor struct {
	rewriter ExpressionRewriter
}

func NewExpressionRewriteVisitor(rewriter ExpressionRewriter) *ExpressionRewriteVisitor {
	return &ExpressionRewriteVisitor{
		rewriter: rewriter,
	}
}

func (v *ExpressionRewriteVisitor) Visit(context any, n Node) any {
	// TODO: default????
	if expr, ok := n.(Expression); ok {
		result := v.rewriter.RewriteExpression(context, expr)
		if result != nil {
			return result
		}
	}
	panic(fmt.Sprintf("expression rewrite not support: %T", n))
}
