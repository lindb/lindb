package pipeline

import (
	"fmt"
	"reflect"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
)

type DriverFactory struct {
	operatorFactories []operator.OperatorFactory
	pielineID         int32
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
	sourceOperator operator.SourceOperator
	operators      []operator.Operator
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
		// if page != nil {
		// 	next.AddInput(page)
		// }
		next.AddInput(page)
		fmt.Printf("process operator:(driver=%p) %s->%s,page=%v\n", d, reflect.TypeOf(current), reflect.TypeOf(next), page)
	}
}
