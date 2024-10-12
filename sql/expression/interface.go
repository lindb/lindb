package expression

import (
	"time"

	"github.com/lindb/lindb/spi/types"
)

type Expression interface {
	EvalInt(ctx EvalContext, row types.Row) (val int64, isNull bool, err error)
	EvalString(ctx EvalContext, row types.Row) (val string, isNull bool, err error)
	EvalFloat(ctx EvalContext, row types.Row) (val float64, isNull bool, err error)
	EvalDuration(ctx EvalContext, row types.Row) (val time.Duration, isNull bool, err error)
	EvalTimeSeries(ctx EvalContext, row types.Row) (val *types.TimeSeries, isNull bool, err error)
	EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error)
	// Getype returns the data type of the expression returns.
	GetType() types.DataType
	// String returns the expression in string format.
	String() string
}
