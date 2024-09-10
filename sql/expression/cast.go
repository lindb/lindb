package expression

import (
	"fmt"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

type Cast struct {
	function Func
	arg      Expression
	retType  types.DataType
}

func NewCast(retType types.DataType, arg Expression) Expression {
	return &Cast{
		retType: retType,
		arg:     arg,
		function: &castFunc{
			baseFunc: baseFunc{args: []Expression{arg}},
		},
	}
}

// EvalString implements Expression.
func (c *Cast) EvalString(row spi.Row) (val string, isNull bool, err error) {
	panic("unimplemented")
}

func (c *Cast) EvalInt(row spi.Row) (val int64, isNull bool, err error) {
	fmt.Printf("cast eval int=%v\n", c.retType)
	return c.function.EvalInt(row)
}

func (c *Cast) EvalFloat(row spi.Row) (val float64, isNull bool, err error) {
	return c.function.EvalFloat(row)
}

func (c *Cast) EvalTimeSeries(row spi.Row) (val *types.TimeSeries, isNull bool, err error) {
	return c.function.EvalTimeSeries(row)
}

// GetType implements Expression.
func (c *Cast) GetType() types.DataType {
	return c.retType
}

func (c *Cast) String() string {
	return fmt.Sprintf("CAST(%s as %s)", c.arg.String(), c.retType)
}

type castFuncFactory struct{}

func (fct *castFuncFactory) NewFunc(args []Expression) Func {
	return &castFunc{
		baseFunc: baseFunc{args: args},
	}
}

type castFunc struct {
	baseFunc
}

func (f *castFunc) EvalInt(row spi.Row) (val int64, isNull bool, err error) {
	lv, _, _ := f.args[0].EvalInt(row)
	fmt.Println("cast int..........")
	return lv, false, nil
}

// EvalFloat implements Func.
func (f *castFunc) EvalFloat(row spi.Row) (val float64, isNull bool, err error) {
	fmt.Println("cast float..........")
	return
}

// EvalTimeSeries evaluates the expression, cast result to types.TimeSeries type.
func (f *castFunc) EvalTimeSeries(row spi.Row) (val *types.TimeSeries, isNull bool, err error) {
	fmt.Printf("cast time series..........,type =%T,%s,%s\n", f.args[0], f.args[0].GetType(), f.args[0].String())
	switch f.args[0].GetType() {
	case types.DataTypeInt:
		val, isNull, err := f.args[0].EvalInt(row)
		if err != nil {
			return nil, false, err
		}
		if isNull {
			return nil, true, nil
		}
		return types.NewTimeSeriesWithSingleValue(float64(val)), false, nil
	case types.DataTypeFloat:
		val, isNull, err := f.args[0].EvalFloat(row)
		if err != nil {
			return nil, false, err
		}
		if isNull {
			return nil, true, nil
		}
		return types.NewTimeSeriesWithSingleValue(val), false, nil
	case types.DataTypeTimeSeries, types.DataTypeSum, types.DataTypeFirst, types.DataTypeLast, types.DataTypeMin, types.DataTypeMax, types.DataTypeHistogram:
		return f.args[0].EvalTimeSeries(row)
	}
	return
}
