package operator

import "github.com/lindb/lindb/spi"

type HashAggregationOperatorFactory struct{}

func NewHashAggregationOperatorFactory() OperatorFactory {
	return &HashAggregationOperatorFactory{}
}

// CreateOperator implements OperatorFactory.
func (fct *HashAggregationOperatorFactory) CreateOperator() Operator {
	return NewHashAggregationOperator()
}

type HashAggregationOperator struct {
	page *spi.Page
}

func NewHashAggregationOperator() Operator {
	return &HashAggregationOperator{}
}

// AddInput implements Operator.
func (h *HashAggregationOperator) AddInput(page *spi.Page) {
	h.page = page
}

// Finish implements Operator.
func (h *HashAggregationOperator) Finish() {
}

// GetOutput implements Operator.
func (h *HashAggregationOperator) GetOutput() *spi.Page {
	// TODO: imple
	return h.page
}

// IsFinished implements Operator.
func (h *HashAggregationOperator) IsFinished() bool {
	return true
}
