package types

import "time"

// Row represents a row in the page.
type Row struct {
	p   *Page
	idx int
}

// GetString returns the string value in the row with the column index.
func (r *Row) GetString(colIdx int) *String {
	return r.p.Columns[colIdx].GetString(r.idx)
}

// GetFloat returns the float value in the row with the column index.
func (r *Row) GetFloat(colIdx int) *Float {
	return r.p.Columns[colIdx].GetFloat(r.idx)
}

// GetInt returns the int value in the row with the column index.
func (r *Row) GetInt(colIdx int) *Int {
	return r.p.Columns[colIdx].GetInt(r.idx)
}

// GetimeSeries returns the time series value in the row with the column index.
func (r *Row) GetTimeSeries(colIdx int) *TimeSeries {
	return r.p.Columns[colIdx].GetTimeSeries(r.idx)
}

func (r *Row) GetTimestamp(colIdx int) *time.Time {
	return r.p.Columns[colIdx].GetTimestamp(r.idx)
}

func (r *Row) GetDuration(colIdx int) *time.Duration {
	return r.p.Columns[colIdx].GetDuration(r.idx)
}
