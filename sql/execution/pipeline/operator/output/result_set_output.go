package output

import (
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/buffer"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
)

type RSOutputOperatorFactory struct {
	output buffer.OutputBuffer
}

func NewRSOutputOperatorFactory(output buffer.OutputBuffer) operator.OperatorFactory {
	return &RSOutputOperatorFactory{
		output: output,
	}
}

// CreateOperator implements operator.OperatorFactory
func (fct *RSOutputOperatorFactory) CreateOperator() operator.Operator {
	return NewResultSetOutputOperator(fct.output)
}

type ResultSetOutputOperator struct {
	output buffer.OutputBuffer
}

func NewResultSetOutputOperator(output buffer.OutputBuffer) operator.Operator {
	return &ResultSetOutputOperator{
		output: output,
	}
}

// AddInput implements operator.Operator
func (op *ResultSetOutputOperator) AddInput(page *types.Page) {
	op.output.AddPage(page)
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
