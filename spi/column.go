package spi

import (
	"github.com/lindb/lindb/spi/types"
)

type ColumnMetadata struct {
	Name     string         `json:"name"`
	DataType types.DataType `json:"type"`
}

func NewColumnInfo(name string, vt types.DataType) ColumnMetadata {
	return ColumnMetadata{
		Name:     name,
		DataType: vt,
	}
}

type Column struct {
	Blocks    []types.Block `json:"block"`
	NumOfRows int           `json:"numOfRows"`
}

func NewColumn() *Column {
	return &Column{}
}

func (c *Column) AppendTimeSeries(val *types.TimeSeries) {
	c.Blocks = append(c.Blocks, val)
	c.NumOfRows++
}

func (c *Column) AppendString(val string) {
	v := types.String(val)
	c.Blocks = append(c.Blocks, &v)
	c.NumOfRows++
}

func (c *Column) AppendInt(val int64) {
	// TODO:
	// v := types.Int(val)
	// c.Blocks = append(c.Blocks, &v)
	// c.NumOfRows++
}

func (c *Column) AppendFloat(val float64) {
	// TODO:
	// v := types.Float(val)
	// c.Blocks = append(c.Blocks, &v)
	// c.NumOfRows++
}

func (c *Column) GetString(row int) *types.String {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row].(*types.String)
}

func (c *Column) GetTimeSeries(row int) *types.TimeSeries {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row].(*types.TimeSeries)
}
