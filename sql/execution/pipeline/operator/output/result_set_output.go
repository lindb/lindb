package output

import (
	"context"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/buffer"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

type RSOutputOperatorFactory struct {
	output       buffer.OutputBuffer
	sourceLayout map[string]int
	layout       []*plan.Symbol
	rebuildPage  bool
}

func NewRSOutputOperatorFactory(output buffer.OutputBuffer, layout []*plan.Symbol, sourceLayout map[string]int) operator.OperatorFactory {
	rebuildPage := false
	for idx, symbol := range layout {
		sourceIdx, ok := sourceLayout[symbol.Name]
		if ok && idx != sourceIdx {
			rebuildPage = true
			break
		}
	}
	return &RSOutputOperatorFactory{
		output:       output,
		layout:       layout,
		sourceLayout: sourceLayout,
		rebuildPage:  rebuildPage,
	}
}

// CreateOperator implements operator.OperatorFactory
func (fct *RSOutputOperatorFactory) CreateOperator(ctx context.Context) operator.Operator {
	return &ResultSetOutputOperator{
		output:       fct.output,
		sourceLayout: fct.sourceLayout,
		layout:       fct.layout,
		rebuildPage:  fct.rebuildPage,
	}
}

type ResultSetOutputOperator struct {
	output       buffer.OutputBuffer
	sourceLayout map[string]int
	layout       []*plan.Symbol
	rebuildPage  bool
}

// AddInput implements operator.Operator
func (op *ResultSetOutputOperator) AddInput(page *types.Page) {
	if page == nil || page.NumRows() == 0 {
		return
	}
	if op.rebuildPage {
		targetPage := types.NewPage()
		for _, col := range op.layout {
			if idx, ok := op.sourceLayout[col.Name]; ok {
				targetPage.AppendColumn(page.Layout[idx], page.Columns[idx])
			}
		}
		op.output.AddPage(targetPage)
	} else {
		op.output.AddPage(page)
	}
}

// Finish implements operator.Operator
func (op *ResultSetOutputOperator) Finish() {
	panic("unimplemented")
}

// GetOutput implements operator.Operator
func (op *ResultSetOutputOperator) GetOutput() *types.Page {
	return nil
}

// IsFinished implements operator.Operator
func (op *ResultSetOutputOperator) IsFinished() bool {
	panic("unimplemented")
}
