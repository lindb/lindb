package spi

import "github.com/lindb/lindb/spi/types"

type Row struct {
	p   *Page
	idx int
}

func (r *Row) GetString(colIdx int) *types.String {
	return r.p.Columns[colIdx].GetString(r.idx)
}

func (r *Row) GetTimeSeries(colIdx int) *types.TimeSeries {
	return r.p.Columns[colIdx].GetTimeSeries(r.idx)
}
