package aggregation

import "github.com/lindb/lindb/spi/types"

type Aggregator interface {
	Enter(row types.Row)
	Flush(column *types.Column)
}
