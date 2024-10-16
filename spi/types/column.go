package types

import "time"

type Column struct {
	Blocks    []Block `json:"block"`
	NumOfRows int     `json:"numOfRows"`
}

func NewColumn() *Column {
	return &Column{}
}

func (c *Column) AppendTimeSeries(val *TimeSeries) {
	c.Blocks = append(c.Blocks, val)
	c.NumOfRows++
}

func (c *Column) AppendString(val string) {
	v := String(val)
	c.Blocks = append(c.Blocks, &v)
	c.NumOfRows++
}

func (c *Column) Append(val any) {
	c.Blocks = append(c.Blocks, val)
	c.NumOfRows++
}

func (c *Column) AppendInt(val int64) {
	v := Int(val)
	c.Blocks = append(c.Blocks, &v)
	c.NumOfRows++
}

func (c *Column) AppendFloat(val float64) {
	v := Float(val)
	c.Blocks = append(c.Blocks, &v)
	c.NumOfRows++
}

func (c *Column) AppendTimestamp(val time.Time) {
	c.Blocks = append(c.Blocks, &val)
	c.NumOfRows++
}

func (c *Column) AppendDuration(val time.Duration) {
	v := Int(val.Nanoseconds())
	c.Blocks = append(c.Blocks, &v)
	c.NumOfRows++
}

func (c *Column) GetString(row int) *String {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row].(*String)
}

func (c *Column) GetInt(row int) *Int {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row].(*Int)
}

func (c *Column) GetFloat(row int) *Float {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row].(*Float)
}

func (c *Column) GetTimestamp(row int) *time.Time {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row].(*time.Time)
}

func (c *Column) GetDuration(row int) *time.Duration {
	if row >= len(c.Blocks) {
		return nil
	}
	val := c.Blocks[row].(*Int)
	duration := time.Duration(int64(*val))
	return &duration
}

func (c *Column) GetTimeSeries(row int) *TimeSeries {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row].(*TimeSeries)
}
