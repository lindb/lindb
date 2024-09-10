package expression

import (
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type baseFunc struct {
	args []Expression
}

type Func interface {
	EvalInt(row spi.Row) (val int64, isNull bool, err error)
	EvalFloat(row spi.Row) (val float64, isNull bool, err error)
	EvalTimeSeries(row spi.Row) (val *types.TimeSeries, isNull bool, err error)
}

type FuncFactory interface {
	NewFunc(args []Expression) Func
}

var funcs = map[tree.FunctionName]FuncFactory{
	tree.Plus:  &arithmeticPlusFuncFactory{},
	tree.Minus: &arithmeticMinusFuncFactory{},
	tree.Mul:   &arithmeticMulFuncFactory{},
	tree.Div:   &arithmeticDivFuncFactory{},
	tree.Mod:   &arithmeticModFuncFactory{},
}
