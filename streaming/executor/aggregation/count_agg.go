package aggregation

import (
	"github.com/lindb/lindb/sql/tree"
)

type countAggFactory struct{}

func (f *countAggFactory) NewAggregator(args []tree.Expression) Aggregator {
	return &countAgg{}
}

type countAgg struct {
	value int64
}

func NewCountAgg() Aggregator {
	return &countAgg{}
}

// Process implements Aggregator.
func (c *countAgg) Process(event any) {
	c.value++
	// fmt.Println("count....")
}

func (c *countAgg) GetValue() any {
	return c.value
}
