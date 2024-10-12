package expression

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type RewriteContext struct {
	SourceLayout []*plan.Symbol
}

func Rewrite(ctx *RewriteContext, node tree.Expression) Expression {
	return (&rewriter{ctx: ctx}).rewrite(node)
}

type rewriter struct {
	ctx *RewriteContext
}

func (r *rewriter) rewrite(node tree.Expression) Expression {
	switch expr := node.(type) {
	case *tree.FunctionCall:
		return r.rewriteCall(expr)
	case *tree.Identifier:
		// TODO: right?
		return NewConstant(expr.Value, types.DTString)
	case *tree.StringLiteral:
		// TODO: right?
		return NewConstant(expr.Value, types.DTString)
	case *tree.Constant:
		return NewConstant(expr.Value, expr.Type)
	case *tree.SymbolReference:
		// TODO: add check
		_, index, _ := lo.FindIndexOf(r.ctx.SourceLayout, func(item *plan.Symbol) bool {
			return item.Name == expr.Name
		})
		return NewColumn(expr.Name, index, expr.DataType)
	case *tree.Cast:
		return NewCast(expr.Type, r.rewrite(expr.Expression))
	default:
		panic(fmt.Sprintf("expression rewrite unimplemented: %T", node))
	}
}

func (r *rewriter) rewriteCall(node *tree.FunctionCall) Expression {
	scalarFunc, err := NewScalarFunc(node.Name, node.RetType, lo.Map(node.Arguments,
		func(item tree.Expression, index int) Expression {
			return r.rewrite(item)
		},
	))
	if err != nil {
		panic(err)
	}
	return scalarFunc
}
