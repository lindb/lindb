package aggregation

import "github.com/lindb/lindb/sql/tree"

type Aggregator interface {
	Process(event any)
	GetValue() any
}

type AggregatorFactory interface {
	NewAggregator(args []tree.Expression) Aggregator
}
