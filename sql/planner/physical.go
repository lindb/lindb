package planner

import (
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/execution/pipeline/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

type PhysicalOperation struct {
	layout            map[string]int
	operatorFactories []operator.OperatorFactory
}

func NewPhysicalOperation(operatorFactory operator.OperatorFactory, layoutSymbols []*plan.Symbol, source *PhysicalOperation) *PhysicalOperation {
	op := &PhysicalOperation{
		layout: make(map[string]int),
	}
	lo.ForEach(layoutSymbols, func(symbol *plan.Symbol, index int) {
		op.layout[symbol.Name] = index
	})
	// Push-Based query engine
	if source != nil {
		op.operatorFactories = append(op.operatorFactories, source.operatorFactories...)
	}
	op.operatorFactories = append(op.operatorFactories, operatorFactory)
	return op
}

func (op *PhysicalOperation) GetLayout() map[string]int {
	return op.layout
}
