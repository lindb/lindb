package expression

import (
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type ScalarFunc struct {
	function Func
	funcName tree.FuncName
	args     []Expression
	retType  types.DataType
}

func NewScalarFunc(funcName tree.FuncName, retType types.DataType, args []Expression) (Expression, error) {
	fct, ok := funcs[funcName]
	if !ok {
		return nil, fmt.Errorf("func not support, func name: %s", funcName)
	}
	fn := fct.NewFunc(args)
	return &ScalarFunc{
		retType:  retType,
		function: fn,
		funcName: funcName,
		args:     args,
	}, nil
}

// EvalString implements Expression.
func (f *ScalarFunc) EvalString(ctx EvalContext, row types.Row) (val string, isNull bool, err error) {
	return f.function.EvalString(ctx, row)
}

func (f *ScalarFunc) EvalInt(ctx EvalContext, row types.Row) (val int64, isNull bool, err error) {
	return f.function.EvalInt(ctx, row)
}

func (f *ScalarFunc) EvalFloat(ctx EvalContext, row types.Row) (val float64, isNull bool, err error) {
	return f.function.EvalFloat(ctx, row)
}

func (f *ScalarFunc) EvalTimeSeries(ctx EvalContext, row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	return f.function.EvalTimeSeries(ctx, row)
}

func (f *ScalarFunc) EvalDuration(ctx EvalContext, row types.Row) (val time.Duration, isNull bool, err error) {
	return f.function.EvalDuration(ctx, row)
}

func (f *ScalarFunc) EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error) {
	return f.function.EvalTime(ctx, row)
}

// GetType implements Expression.
func (f *ScalarFunc) GetType() types.DataType {
	return f.retType
}

// String returns the scalar function in string format.
func (f *ScalarFunc) String() string {
	return fmt.Sprintf("%s(%s)", f.funcName, strings.Join(lo.Map(f.args, func(item Expression, index int) string {
		return item.String()
	}), ","))
}
