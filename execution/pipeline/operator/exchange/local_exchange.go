package exchange

import (
	"github.com/lindb/lindb/execution/pipeline/operator"
	"github.com/lindb/lindb/spi"
)

type LocalExchangeOperatorFactory struct{}

func NewLocalExchangeOperatorFactory() operator.OperatorFactory {
	return &LocalExchangeOperatorFactory{}
}

// CreateOperator implements operator.OperatorFactory.
func (l *LocalExchangeOperatorFactory) CreateOperator() operator.Operator {
	return NewLocalExchangeOperator()
}

type LocalExchangeOperator struct {
	page *spi.Page
}

func NewLocalExchangeOperator() operator.Operator {
	return &LocalExchangeOperator{}
}

// AddInput implements operator.Operator.
func (l *LocalExchangeOperator) AddInput(page *spi.Page) {
	l.page = page
}

// Finish implements operator.Operator.
func (l *LocalExchangeOperator) Finish() {
}

// GetOutput implements operator.Operator.
func (l *LocalExchangeOperator) GetOutput() *spi.Page {
	return l.page
}

// IsFinished implements operator.Operator.
func (l *LocalExchangeOperator) IsFinished() bool {
	return true
}
