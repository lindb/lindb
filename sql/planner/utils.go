package planner

import (
	"fmt"

	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/tree"
)

func coerceIfNecessary(analysis *analyzer.Analysis, origianl, rewritten tree.Expression) tree.Expression {
	coercion, ok := analysis.GetCoercion(origianl)
	fmt.Printf("check coercion%T=%v,=%v\n", origianl, coercion, ok)
	if !ok {
		return rewritten
	}
	fmt.Println("cast ....... rewrite")
	return &tree.Cast{
		BaseNode: tree.BaseNode{
			ID: origianl.GetID(),
		},
		Type:       coercion,
		Expression: rewritten,
	}
}
