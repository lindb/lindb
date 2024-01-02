package spi

import "github.com/lindb/lindb/spi/value"

type Row struct {
	p   *Page
	idx int
}

func (r *Row) GetString(colIdx int) *value.String {
	return r.p.Columns[colIdx].GetString(r.idx)
}

func (r *Row) GetTimeSeries(colIdx int) *value.TimeSeries {
	return r.p.Columns[colIdx].GetTimeSeries(r.idx)
}
