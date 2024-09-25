package buffer

import (
	"github.com/lindb/lindb/spi/types"
)

type QueryOutputBuffer struct {
	rsBuild *ResultSetBuild
}

func NewQueryOutputBuffer(rsBuild *ResultSetBuild) OutputBuffer {
	return &QueryOutputBuffer{
		rsBuild: rsBuild,
	}
}

func (buf *QueryOutputBuffer) AddPage(page *types.Page) {
	buf.rsBuild.AddPage(page)
}
