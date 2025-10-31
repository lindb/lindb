package aggregation

import (
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

func newCountAggregator(args []tree.Expression) Aggregator {
	return &countAgg{}
}

type countAgg struct {
	value int64
}

func NewCountAgg() Aggregator {
	return &countAgg{}
}

func (c *countAgg) Enter(row types.Row) {
	c.value++
	// fmt.Println("count....")
}

func (c *countAgg) Flush(column *types.Column) {
	column.AppendInt(c.value)
	// fmt.Println("flush count....")
}
