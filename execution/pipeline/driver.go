package pipeline

import (
	"fmt"
	"reflect"

	"github.com/lindb/lindb/execution/pipeline/operator"
	"github.com/lindb/lindb/spi"
)

type DriverFactory struct {
	pielineID         int32
	operatorFactories []operator.OperatorFactory
}

func NewDriverFactory(pipelineID int32, operatorFactories []operator.OperatorFactory) *DriverFactory {
	return &DriverFactory{
		pielineID:         pipelineID,
		operatorFactories: operatorFactories,
	}
}

func (fct *DriverFactory) CreateDriver() *Driver {
	var operators []operator.Operator
	for _, operatorFct := range fct.operatorFactories {
		operators = append(operators, operatorFct.CreateOperator())
	}
	return NewDriver(operators)
}

type Driver struct {
	operators      []operator.Operator
	sourceOperator operator.SourceOperator
}

func NewDriver(operators []operator.Operator) *Driver {
	var sourceOperator operator.SourceOperator
	for _, op := range operators {
		if source, ok := op.(operator.SourceOperator); ok {
			sourceOperator = source
		}
	}
	return &Driver{
		operators:      operators,
		sourceOperator: sourceOperator,
	}
}

func (d *Driver) GetSourceOperator() operator.SourceOperator {
	return d.sourceOperator
}

func (d *Driver) AddSplit(split spi.Split) {
	d.sourceOperator.AddSplit(split)
}

func (d *Driver) Process() {
	for i := 0; i < len(d.operators)-1; i++ {
		current := d.operators[i]
		next := d.operators[i+1]
		page := current.GetOutput()
		if page != nil {
			next.AddInput(page)
		}
		fmt.Printf("%s->%s,page=%v\n", reflect.TypeOf(current), reflect.TypeOf(next), page)
	}
}
