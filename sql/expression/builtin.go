package expression

import (
	"time"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type baseFunc struct {
	args []Expression
}

func (*baseFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalString(row types.Row) (val string, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalDuration(row types.Row) (val time.Duration, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	panic("implement me")
}

type Func interface {
	EvalInt(row types.Row) (val int64, isNull bool, err error)
	EvalFloat(row types.Row) (val float64, isNull bool, err error)
	EvalString(row types.Row) (val string, isNull bool, err error)
	EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error)
	EvalDuration(row types.Row) (val time.Duration, isNull bool, err error)
	EvalTime(row types.Row) (val time.Time, isNull bool, err error)
}

type FuncFactory interface {
	NewFunc(args []Expression) Func
}

// IsFuncSupported check if given function name is supported.
func IsFuncSupported(name tree.FuncName) bool {
	_, ok := funcs[name]
	return ok
}

var funcs = map[tree.FuncName]FuncFactory{
	tree.Plus:  &arithmeticPlusFuncFactory{},
	tree.Minus: &arithmeticMinusFuncFactory{},
	tree.Mul:   &arithmeticMulFuncFactory{},
	tree.Div:   &arithmeticDivFuncFactory{},
	tree.Mod:   &arithmeticModFuncFactory{},

	// time functions
	// ref: https://dev.mysql.com/doc/refman/8.4/en/date-and-time-functions.html
	tree.DateAdd:   &addSubDateFuncFactory{},
	tree.Now:       &nowFuncFactory{},
	tree.StrToDate: &strToDateFuncFactory{},
}
