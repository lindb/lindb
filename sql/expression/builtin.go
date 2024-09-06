package expression

import (
	"github.com/lindb/lindb/sql/tree"
)

type baseFunc struct {
	args []Expression
}

type Func interface {
	Expression
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
