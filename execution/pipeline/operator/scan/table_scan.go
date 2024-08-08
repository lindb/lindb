package scan

import (
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/execution/pipeline/operator"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type TableScanOperatorFactory struct {
	table  spi.TableHandle
	filter tree.Expression

	sourceID plan.PlanNodeID
}

func NewTableScanOperatorFactory(sourceID plan.PlanNodeID, table spi.TableHandle, filter tree.Expression) operator.OperatorFactory {
	return &TableScanOperatorFactory{
		sourceID: sourceID,
		table:    table,
		filter:   filter,
	}
}

func (fct *TableScanOperatorFactory) CreateOperator() operator.Operator {
	provider := spi.GetPageSourceProvider(fct.table)
	return NewTableScanOperator(fct.sourceID, provider.CreatePageSource(fct.table))
}

type TableScanOperator struct {
	pageSource spi.PageSource

	sourceID plan.PlanNodeID
}

func NewTableScanOperator(sourceID plan.PlanNodeID, pageSource spi.PageSource) operator.SourceOperator {
	return &TableScanOperator{
		sourceID:   sourceID,
		pageSource: pageSource,
	}
}

func (op *TableScanOperator) GetSourceID() plan.PlanNodeID {
	return op.sourceID
}

func (op *TableScanOperator) NoMoreSplits() {
}

func (op *TableScanOperator) AddSplit(split spi.Split) {
	op.pageSource.AddSplit(split)
}

// AddInput implements operator.Operator
func (op *TableScanOperator) AddInput(page *spi.Page) {
	panic(fmt.Errorf("%w: table scan cannot take input", constants.ErrNotSupportOperation))
}

// Finish implements operator.Operator
func (op *TableScanOperator) Finish() {
}

// GetOutput implements operator.Operator
func (op *TableScanOperator) GetOutput() *spi.Page {
	return op.pageSource.GetNextPage()
}

// IsFinished implements operator.Operator
func (op *TableScanOperator) IsFinished() bool {
	return true
}
