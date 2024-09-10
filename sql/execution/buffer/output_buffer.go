package buffer

import "github.com/lindb/lindb/spi"

type OutputBuffer interface {
	AddPage(page *spi.Page)
}
