package operator

import (
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/planner/plan"
)

type OperatorFactory interface {
	CreateOperator() Operator
}

type SourceOperatorFactory interface {
	OperatorFactory
}

type Operator interface {
	GetOutput() *types.Page

	AddInput(page *types.Page)
	// Finish notifies the operator that no more pages will be added
	// and the operator should finish processing and flush results.
	Finish()
	// IsFinished if this operator finished processing and no more output page will be produced.
	IsFinished() bool
}

type SourceOperator interface {
	Operator

	GetSourceID() plan.PlanNodeID

	AddSplit(split spi.Split)
	NoMoreSplits()
}
