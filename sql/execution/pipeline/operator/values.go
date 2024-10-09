package operator

import (
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/planner/plan"
)

type ValuesOperatorFactory struct {
	pages    []*types.Page
	sourceID plan.PlanNodeID
}

func NewValuesOperatorFactory(sourceID plan.PlanNodeID, pages []*types.Page) SourceOperatorFactory {
	return &ValuesOperatorFactory{
		pages:    pages,
		sourceID: sourceID,
	}
}

// CreateOperator implements OperatorFactory.
func (fct *ValuesOperatorFactory) CreateOperator() Operator {
	return &ValuesOperator{
		sourceID: fct.sourceID,
		pages:    fct.pages,
	}
}

type ValuesOperator struct {
	pages    []*types.Page
	sourceID plan.PlanNodeID
	idx      int
}

func NewValuesOperator(sourceID plan.PlanNodeID, pages []*types.Page) SourceOperator {
	return &ValuesOperator{
		sourceID: sourceID,
		pages:    pages,
	}
}

// AddSplit implements SourceOperator.
func (op *ValuesOperator) AddSplit(split spi.Split) {}

// GetSourceID implements SourceOperator.
func (op *ValuesOperator) GetSourceID() plan.PlanNodeID {
	return op.sourceID
}

// NoMoreSplits implements SourceOperator.
func (op *ValuesOperator) NoMoreSplits() {}

// AddInput implements Operator.
func (op *ValuesOperator) AddInput(page *types.Page) {
	panic(fmt.Errorf("%w: values cannot take input", constants.ErrNotSupportOperation))
}

// Finish implements Operator.
func (op *ValuesOperator) Finish() {
	panic("unimplemented")
}

// GetOutput implements Operator.
func (op *ValuesOperator) GetOutput() (page *types.Page) {
	if op.IsFinished() {
		return
	}
	page = op.pages[op.idx]
	op.idx++
	return
}

// IsFinished implements Operator.
func (op *ValuesOperator) IsFinished() bool {
	return op.idx == len(op.pages)
}
