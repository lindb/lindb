package tree

// FuncName represents function name.
type FuncName string

// AggFuncName represents aggregation function name.
type AggFuncName string

const (
	// scalar function names
	Plus  FuncName = "plus"
	Minus FuncName = "minus"
	Div   FuncName = "div"
	Mul   FuncName = "mul"
	Mod   FuncName = "mod"

	// aggregation function names
	Sum FuncName = "sum"
	Min FuncName = "min"
	Max FuncName = "max"
)

const (
// Sum
)
