package buffer

import (
	"github.com/lindb/lindb/spi"
)

type QueryOutputBuffer struct {
	rsBuild *ResultSetBuild
}

func NewQueryOutoutBuffer(rsBuild *ResultSetBuild) OutputBuffer {
	return &QueryOutputBuffer{
		rsBuild: rsBuild,
	}
}

func (buf *QueryOutputBuffer) AddPage(page *spi.Page) {
	buf.rsBuild.AddPage(page)
}
