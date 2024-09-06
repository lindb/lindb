package expression

import (
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

type Expression interface {
	EvalInt(row spi.Row) (val int64, err error)
	EvalFloat(row spi.Row) (val float64, err error)
	EvalTimeSeries(row spi.Row) (val *types.TimeSeries, err error)
}
