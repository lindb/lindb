package planner

import (
	"fmt"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

func ExtractSymbolsFromExpressions(expressions []tree.Expression) (symbols []*plan.Symbol) {
	visitor := &tree.DefaultTraversalVisitor{
		PostProcess: func(n tree.Node) {
			if ref, ok := n.(*tree.SymbolReference); ok {
				symbols = append(symbols, plan.SymbolFrom(ref))
			}
		},
	}
	for _, node := range expressions {
		visitor.Visit(nil, node)
	}
	return
}

func ExtractSymbolsFromAggreation(aggregation *plan.Aggregation) (symbols []*plan.Symbol) {
	visitor := &tree.DefaultTraversalVisitor{
		PostProcess: func(n tree.Node) {
			if ref, ok := n.(*tree.SymbolReference); ok {
				symbols = append(symbols, plan.SymbolFrom(ref))
			}
		},
	}
	for _, node := range aggregation.Arguments {
		fmt.Printf("extract symbols agg args.......=%T\n", node)
		visitor.Visit(nil, node)
	}
	// TODO: add agg other fields
	return
}

func ExtractSymbolsFromExpression(expression tree.Expression) []*plan.Symbol {
	panic("unimplemented implements extract symbols from expression")
}
