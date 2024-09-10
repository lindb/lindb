package expression

import (
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

type Expression interface {
	EvalString(row spi.Row) (val string, isNull bool, err error)
	EvalInt(row spi.Row) (val int64, isNull bool, err error)
	EvalFloat(row spi.Row) (val float64, isNull bool, err error)
	EvalTimeSeries(row spi.Row) (val *types.TimeSeries, isNull bool, err error)
	// Getype returns the data type of the expression returns.
	GetType() types.DataType
	// String returns the expression in string format.
	String() string
}
