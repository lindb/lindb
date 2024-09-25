package expression

import (
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type baseFunc struct {
	args []Expression
}

type Func interface {
	EvalInt(row types.Row) (val int64, isNull bool, err error)
	EvalFloat(row types.Row) (val float64, isNull bool, err error)
	EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error)
}

type FuncFactory interface {
	NewFunc(args []Expression) Func
}

var funcs = map[tree.FuncName]FuncFactory{
	tree.Plus:  &arithmeticPlusFuncFactory{},
	tree.Minus: &arithmeticMinusFuncFactory{},
	tree.Mul:   &arithmeticMulFuncFactory{},
	tree.Div:   &arithmeticDivFuncFactory{},
	tree.Mod:   &arithmeticModFuncFactory{},
}
