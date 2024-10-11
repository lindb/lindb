package tree

import (
	"github.com/lindb/lindb/spi/types"
)

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

	// time function names
	DateAdd   FuncName = "date_add"
	Now       FuncName = "now"
	StrToDate FuncName = "str_to_date"
)

func GetDefaultFuncReturnType(name FuncName) types.DataType {
	return defaultFuncReturnTypes[name]
}

var defaultFuncReturnTypes = map[FuncName]types.DataType{
	DateAdd:   types.DTTimestamp,
	Now:       types.DTTimestamp,
	StrToDate: types.DTTimestamp,
}
