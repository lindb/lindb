package expression

import (
	"github.com/lindb/lindb/spi/types"
)

type Column struct {
	name    string
	retType types.DataType
	index   int
}

func NewColumn(name string, index int, retType types.DataType) Expression {
	return &Column{name: name, index: index, retType: retType}
}

func (c *Column) EvalString(row types.Row) (val string, isNull bool, err error) {
	return string(*row.GetString(c.index)), false, nil
}

func (c *Column) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	return 40, false, nil
}

func (c *Column) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
	return
}

func (c *Column) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	return row.GetTimeSeries(c.index), false, nil
}

// GetType returns the data type of the column returns.
func (c *Column) GetType() types.DataType {
	return c.retType
}

// String returns the column in string format.
func (c *Column) String() string {
	return c.name
}
