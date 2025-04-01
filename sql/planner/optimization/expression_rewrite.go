package optimization

import (
	"fmt"

	"github.com/lindb/lindb/sql/tree"
)

type ExpressionRewrite struct {
	originSymbol, newSymbol string
}

func (r *ExpressionRewrite) Rewrite(originSymbol, newSymbol string, e tree.Expression) tree.Expression {
	// reset origin and new symbol
	r.originSymbol = originSymbol
	r.newSymbol = newSymbol
	return e.Accept(nil, r).(tree.Expression)
}

func (r *ExpressionRewrite) Visit(context any, e tree.Node) any {
	fmt.Printf("expression rewrite=%T\n", e)
	switch expr := e.(type) {
	case *tree.ComparisonExpression:
		return &tree.ComparisonExpression{
			BaseNode: tree.BaseNode{ID: expr.ID},
			Operator: expr.Operator,
			Left:     expr.Left.Accept(context, r).(tree.Expression),
			Right:    expr.Right.Accept(context, r).(tree.Expression),
		}
	case *tree.SymbolReference:
		if expr.Name == r.originSymbol {
			return &tree.SymbolReference{
				BaseNode: tree.BaseNode{ID: expr.ID},
				Name:     r.newSymbol,
			}
		}
		return expr
	default:
		return e
	}
}
