package executor

import "github.com/lindb/lindb/spi/types"

type Executor interface {
	Enter(page *types.Page)
	Leave(output chan<- *types.Page)
}
