package expression

import (
	"fmt"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

type Column struct {
	name  string
	index int
}

func NewColumn(name string, index int) Expression {
	return &Column{name: name, index: index}
}

func (c *Column) EvalInt(row spi.Row) (val int64, err error) {
	fmt.Printf("column===%s,%d,%v\n", c.name, c.index, row.GetTimeSeries(1))
	return 40, nil
}

func (c *Column) EvalFloat(row spi.Row) (val float64, err error) {
	return 0, nil
}

func (c *Column) EvalTimeSeries(row spi.Row) (val *types.TimeSeries, err error) {
	return nil, nil
}
