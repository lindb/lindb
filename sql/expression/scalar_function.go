package expression

import (
	"errors"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type ScalarFunc struct {
	function Func
}

func NewScalarFunc(funcName tree.FunctionName, args []Expression) (Expression, error) {
	fct, ok := funcs[funcName]
	if !ok {
		return nil, errors.New("func not support")
	}
	fn := fct.NewFunc(args)
	return &ScalarFunc{
		function: fn,
	}, nil
}

func (f *ScalarFunc) EvalInt(row spi.Row) (val int64, err error) {
	return f.function.EvalInt(row)
}

func (f *ScalarFunc) EvalFloat(row spi.Row) (val float64, err error) {
	return f.function.EvalFloat(row)
}

func (f *ScalarFunc) EvalTimeSeries(row spi.Row) (val *types.TimeSeries, err error) {
	return f.function.EvalTimeSeries(row)
}
