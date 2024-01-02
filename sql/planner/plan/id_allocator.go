package plan

import (
	"fmt"

	"github.com/lindb/lindb/sql/tree"
)

type PlanNodeIDAllocator struct {
	next PlanNodeID
}

func NewPlanNodeIDAllocator() *PlanNodeIDAllocator {
	return &PlanNodeIDAllocator{}
}

func (a *PlanNodeIDAllocator) Next() PlanNodeID {
	a.next++
	return a.next
}

type SymbolAllocator struct{}

func NewSymbolAllocator() *SymbolAllocator {
	return &SymbolAllocator{}
}

func (a *SymbolAllocator) NewSymbol(expression tree.Expression, suffix string) *Symbol {
	fmt.Printf("new symbol=%T\n", expression)
	nameHint := "expr"
	switch expr := expression.(type) {
	case *tree.Identifier:
		nameHint = expr.Value
	case *tree.SymbolReference:
		nameHint = expr.Name
	}
	// FIXME: ????
	return a.newSymbol(nameHint, suffix)
}

func (a *SymbolAllocator) FromSymbol(symbolHint *Symbol, suffix string) *Symbol {
	return a.newSymbol(symbolHint.Name, suffix)
}

func (a *SymbolAllocator) newSymbol(nameHint, suffix string) *Symbol {
	unique := nameHint

	if suffix != "" {
		unique += "$" + suffix
	}
	// TODO: fixme cache symbols
	return &Symbol{
		Name: unique,
	}
}
