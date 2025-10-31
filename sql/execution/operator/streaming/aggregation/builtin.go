package aggregation

import (
	"fmt"

	"github.com/lindb/lindb/sql/tree"
)

type NewAggregator func(args []tree.Expression) Aggregator

var funcs = map[tree.FuncName]NewAggregator{
	tree.Count: newCountAggregator,
}

func CreateAggregator(name tree.FuncName, args []tree.Expression) (Aggregator, error) {
	factory, ok := funcs[name]
	if !ok {
		return nil, fmt.Errorf("not support %s", name)
	}
	return factory(args), nil
}
