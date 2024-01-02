package operator

import (
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/spi"
)

type OutputFactory interface {
	CreateOutputOperator() OperatorFactory
}

type MetricResultSetOutputFactory struct {
}

func NewMetricResultSetOutputFactory() OutputFactory {
	return &MetricResultSetOutputFactory{}
}

func (fct *MetricResultSetOutputFactory) CreateOutputOperator() OperatorFactory {
	return &MetricRSOperatorFactory{}
}

type MetricRSOperatorFactory struct {
}

func (fct *MetricRSOperatorFactory) CreateOperator() Operator {
	return NewMetricResultSetOperator()
}

type MetricResultSetOperator struct {
	groupedSeriesList series.GroupedIterators
}

func NewMetricResultSetOperator() Operator {
	return &MetricResultSetOperator{}
}

// AddInput implements Operator
func (*MetricResultSetOperator) AddInput(page *spi.Page) {
}

// GetOutput implements Operator
func (*MetricResultSetOperator) GetOutput() *spi.Page {
	return nil
}

func (op *MetricResultSetOperator) Finish() {
}

// GetOutput implements Operator
func (op *MetricResultSetOperator) IsFinished() bool {
	return true
}
