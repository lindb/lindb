package planner

import "github.com/lindb/lindb/sql/execution/pipeline/operator"

type PhysicalOperation struct {
	operatorFactories []operator.OperatorFactory
}

func NewPhysicalOperation(operatorFactory operator.OperatorFactory, source *PhysicalOperation) *PhysicalOperation {
	op := &PhysicalOperation{}
	// Push-Based query engine
	if source != nil {
		op.operatorFactories = append(op.operatorFactories, source.operatorFactories...)
	}
	op.operatorFactories = append(op.operatorFactories, operatorFactory)
	return op
}
