package operator

import (
	"context"

	"github.com/lindb/lindb/spi/types"
)

type HashAggregationOperatorFactory struct{}

func NewHashAggregationOperatorFactory() OperatorFactory {
	return &HashAggregationOperatorFactory{}
}

// CreateOperator implements OperatorFactory.
func (fct *HashAggregationOperatorFactory) CreateOperator(ctx context.Context) Operator {
	return NewHashAggregationOperator()
}

type HashAggregationOperator struct {
	page *types.Page
}

func NewHashAggregationOperator() Operator {
	return &HashAggregationOperator{}
}

// AddInput implements Operator.
func (h *HashAggregationOperator) AddInput(page *types.Page) {
	h.page = page
}

// Finish implements Operator.
func (h *HashAggregationOperator) Finish() {
}

// GetOutput implements Operator.
func (h *HashAggregationOperator) GetOutput() *types.Page {
	// TODO: imple
	return h.page
}

// IsFinished implements Operator.
func (h *HashAggregationOperator) IsFinished() bool {
	return true
}
