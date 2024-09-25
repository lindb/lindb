package expression

import (
	"github.com/lindb/lindb/spi/types"
)

type Expression interface {
	EvalString(row types.Row) (val string, isNull bool, err error)
	EvalInt(row types.Row) (val int64, isNull bool, err error)
	EvalFloat(row types.Row) (val float64, isNull bool, err error)
	EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error)
	// Getype returns the data type of the expression returns.
	GetType() types.DataType
	// String returns the expression in string format.
	String() string
}
