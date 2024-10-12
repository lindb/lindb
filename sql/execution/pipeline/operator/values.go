package operator

import (
	"context"
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/planner/plan"
)

type ValuesOperatorFactory struct {
	page     *types.Page
	sourceID plan.PlanNodeID
}

func NewValuesOperatorFactory(sourceID plan.PlanNodeID, page *types.Page) SourceOperatorFactory {
	return &ValuesOperatorFactory{
		page:     page,
		sourceID: sourceID,
	}
}

// CreateOperator implements OperatorFactory.
func (fct *ValuesOperatorFactory) CreateOperator(ctx context.Context) Operator {
	return &ValuesOperator{
		sourceID: fct.sourceID,
		page:     fct.page,
	}
}

type ValuesOperator struct {
	page     *types.Page
	sourceID plan.PlanNodeID
}

func NewValuesOperator(sourceID plan.PlanNodeID, page *types.Page) SourceOperator {
	return &ValuesOperator{
		sourceID: sourceID,
		page:     page,
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
	return op.page
}

// IsFinished implements Operator.
func (op *ValuesOperator) IsFinished() bool {
	return op.page == nil
}
