package types

type Row struct {
	p   *Page
	idx int
}

// GetString returns the string value in the column with the column index.
func (r *Row) GetString(colIdx int) *String {
	return r.p.Columns[colIdx].GetString(r.idx)
}

// GetimeSeries returns the TimeSeries value in the column with the column index.
func (r *Row) GetTimeSeries(colIdx int) *TimeSeries {
	return r.p.Columns[colIdx].GetTimeSeries(r.idx)
}
