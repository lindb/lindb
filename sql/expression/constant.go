package expression

import (
	"fmt"
	"time"

	"github.com/lindb/common/pkg/timeutil"

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
func (c *Constant) EvalString(_ EvalContext, _ types.Row) (val string, isNull bool, err error) {
	return c.value.(string), false, nil
}

func (c *Constant) EvalInt(_ EvalContext, _ types.Row) (val int64, isNull bool, err error) {
	return c.value.(int64), false, nil
}

func (c *Constant) EvalFloat(_ EvalContext, _ types.Row) (val float64, isNull bool, err error) {
	return
}

func (c *Constant) EvalTimeSeries(_ EvalContext, _ types.Row) (val *types.TimeSeries, isNull bool, err error) {
	return
}

func (c *Constant) EvalDuration(_ EvalContext, _ types.Row) (val time.Duration, isNull bool, err error) {
	val = c.value.(time.Duration)
	return
}

func (c *Constant) EvalTime(_ EvalContext, _ types.Row) (val time.Time, isNull bool, err error) {
	switch v := c.value.(type) {
	case string:
		timestamp, err := timeutil.ParseTimestamp(v, timeutil.DataTimeFormat2)
		if err != nil {
			return time.Time{}, true, err
		}
		return time.UnixMilli(timestamp), false, nil
	}
	return
}

func (c *Constant) GetType() types.DataType {
	return c.retType
}

// String returns the constant in string format.
func (c *Constant) String() string {
	return fmt.Sprintf("%v", c.value)
}
