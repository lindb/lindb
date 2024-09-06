package expression

import (
	"fmt"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

type Constant struct {
	value any
}

func NewConstant(value any) Expression {
	return &Constant{
		value: value,
	}
}

func (c *Constant) EvalInt(_ spi.Row) (val int64, err error) {
	fmt.Printf("constant==%v\n", c.value)
	return c.value.(int64), nil
}

func (c *Constant) EvalFloat(_ spi.Row) (val float64, err error) {
	return 0, nil
}

func (c *Constant) EvalTimeSeries(_ spi.Row) (val *types.TimeSeries, err error) {
	return nil, nil
}
