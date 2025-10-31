package aggregation

import (
	"fmt"

	"github.com/lindb/lindb/sql/tree"
)

var funcs = map[tree.FuncName]AggregatorFactory{
	tree.Count: &countAggFactory{},
}

func NewAggregator(name tree.FuncName, args []tree.Expression) (Aggregator, error) {
	factory, ok := funcs[name]
	if !ok {
		return nil, fmt.Errorf("not support %s", name)
	}
	return factory.NewAggregator(args), nil
}
