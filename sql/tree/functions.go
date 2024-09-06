package tree

type FunctionName string

const (
	// scalar function names
	Plus  FunctionName = "plus"
	Minus FunctionName = "minus"
	Div   FunctionName = "div"
	Mul   FunctionName = "mul"
	Mod   FunctionName = "mod"
)
