package exchange

import (
	"context"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
)

type LocalExchangeOperatorFactory struct{}

func NewLocalExchangeOperatorFactory() operator.OperatorFactory {
	return &LocalExchangeOperatorFactory{}
}

// CreateOperator implements operator.OperatorFactory.
func (l *LocalExchangeOperatorFactory) CreateOperator(ctx context.Context) operator.Operator {
	return NewLocalExchangeOperator()
}

type LocalExchangeOperator struct {
	page *types.Page
}

func NewLocalExchangeOperator() operator.Operator {
	return &LocalExchangeOperator{}
}

// AddInput implements operator.Operator.
func (l *LocalExchangeOperator) AddInput(page *types.Page) {
	l.page = page
}

// Finish implements operator.Operator.
func (l *LocalExchangeOperator) Finish() {
}

// GetOutput implements operator.Operator.
func (l *LocalExchangeOperator) GetOutput() *types.Page {
	return l.page
}

// IsFinished implements operator.Operator.
func (l *LocalExchangeOperator) IsFinished() bool {
	return true
}
