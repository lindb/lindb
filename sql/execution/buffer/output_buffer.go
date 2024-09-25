package buffer

import (
	"github.com/lindb/lindb/spi/types"
)

type OutputBuffer interface {
	AddPage(page *types.Page)
}
