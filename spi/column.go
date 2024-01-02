package spi

import (
	"github.com/lindb/lindb/spi/value"
)

type ColumnMetadata struct {
	Name          string              `json:"name"`
	ValueType     value.ValueType     `json:"valueType"`
	AggregateType value.AggregateType `json:"aggregateType,omitempty"`
}

func NewColumnInfo(name string, vt value.ValueType) ColumnMetadata {
	return ColumnMetadata{
		Name:      name,
		ValueType: vt,
	}
}

type Column struct {
	Blocks []value.Block `json:"block"`
	Length int           `json:"length"`
}

func NewColumn() *Column {
	return &Column{}
}

func (c *Column) AppendTimeSeries(val *value.TimeSeries) {
	c.Blocks = append(c.Blocks, val)
	c.Length++
}

func (c *Column) GetString(row int) *value.String {
	return nil
}

func (c *Column) GetTimeSeries(row int) *value.TimeSeries {
	//FIXME:
	return c.Blocks[row].(*value.TimeSeries)
}
