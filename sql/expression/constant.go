package expression

import (
	"fmt"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

type Constant struct {
	value   any
	retType types.DataType
}

func NewConstant(value any, retType types.DataType) Expression {
	return &Constant{
		retType: retType,
		value:   value,
	}
}

// EvalString implements Expression.
func (c *Constant) EvalString(row spi.Row) (val string, isNull bool, err error) {
	panic("unimplemented")
}

func (c *Constant) EvalInt(_ spi.Row) (val int64, isNull bool, err error) {
	return c.value.(int64), false, nil
}

func (c *Constant) EvalFloat(_ spi.Row) (val float64, isNull bool, err error) {
	return
}

func (c *Constant) EvalTimeSeries(_ spi.Row) (val *types.TimeSeries, isNull bool, err error) {
	return
}

func (c *Constant) GetType() types.DataType {
	return c.retType
}

// String returns the constant in string format.
func (c *Constant) String() string {
	return fmt.Sprintf("%v", c.value)
}
