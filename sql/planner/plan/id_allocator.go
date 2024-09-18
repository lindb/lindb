package plan

import (
	"fmt"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/analyzer"
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

type SymbolAllocator struct {
	analyzerContext *analyzer.AnalyzerContext
}

func NewSymbolAllocator(analyzerContext *analyzer.AnalyzerContext) *SymbolAllocator {
	return &SymbolAllocator{
		analyzerContext: analyzerContext,
	}
}

func (a *SymbolAllocator) NewSymbol(expression tree.Expression, suffix string, dataType types.DataType) *Symbol {
	fmt.Printf("new symbol=%T\n", expression)
	nameHint := "expr"
	switch expr := expression.(type) {
	case *tree.Identifier:
		nameHint = expr.Value
	case *tree.SymbolReference:
		nameHint = expr.Name
	case *tree.FunctionCall:
		if expr.RefField != nil {
			nameHint = expr.RefField.Name
			dataType = expr.RefField.DataType
		} else {
			nameHint = string(expr.Name)
		}
		// FIXME: func call
	}
	// FIXME: ????
	return a.newSymbol(nameHint, suffix, dataType)
}

func (a *SymbolAllocator) FromSymbol(symbolHint *Symbol, suffix string, dataType types.DataType) *Symbol {
	return a.newSymbol(symbolHint.Name, suffix, dataType)
}

func (a *SymbolAllocator) newSymbol(nameHint, suffix string, dataType types.DataType) *Symbol {
	unique := nameHint

	// if suffix != "" {
	// 	unique += "$" + suffix
	// }
	// TODO: fixme cache symbols
	return &Symbol{
		Name:     unique,
		Suffix:   suffix,
		DataType: dataType,
	}
}
