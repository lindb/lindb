package types

import (
	"errors"
	"fmt"

	"github.com/lindb/common/pkg/encoding"
)

type AggregateType byte

type DataType uint16

const (
	DataTypeString DataType = iota + 1
	DataTypeInt
	DataTypeFloat
	DataTypeTimeSeries
	DataTypeSum
	DataTypeMin
	DataTypeMax
	DataTypeLast
	DataTypeFirst
	DataTypeHistogram
)

const (
	ATUnknown AggregateType = iota
	ATSum
	ATMin
	ATMax
	ATLast
	ATFirst
	ATHistogram
)

func (dt DataType) CanAggregatin() bool {
	return dt == DataTypeSum || dt == DataTypeMin || dt == DataTypeMax ||
		dt == DataTypeLast || dt == DataTypeFirst || dt == DataTypeHistogram
}

func (dt DataType) String() string {
	switch dt {
	case DataTypeString:
		return "string"
	case DataTypeInt:
		return "int"
	case DataTypeFloat:
		return "float"
	case DataTypeTimeSeries:
		return "timeSeries"
	case DataTypeSum:
		return "sum"
	case DataTypeMin:
		return "min"
	case DataTypeMax:
		return "max"
	case DataTypeLast:
		return "last"
	case DataTypeFirst:
		return "first"
	case DataTypeHistogram:
		return "histogram"
	default:
		panic(fmt.Sprintf("invalid value type<%d>", dt))
	}
}

func (vt DataType) MarshalJSON() ([]byte, error) {
	return encoding.JSONMarshal(vt.String()), nil
}

func (tv *DataType) UnmarshalJSON(data []byte) error {
	var str string
	if err := encoding.JSONUnmarshal(data, &str); err != nil {
		return err
	}
	switch str {
	case "string":
		*tv = DataTypeString
	case "int":
		*tv = DataTypeInt
	case "float":
		*tv = DataTypeFloat
	case "timeSeries":
		*tv = DataTypeTimeSeries
	case "sum":
		*tv = DataTypeSum
	case "min":
		*tv = DataTypeMin
	case "max":
		*tv = DataTypeMax
	case "last":
		*tv = DataTypeLast
	case "first":
		*tv = DataTypeFirst
	case "histogram":
		*tv = DataTypeHistogram
	default:
		return fmt.Errorf("invalid value type<%s>", str)
	}
	return nil
}

func (at AggregateType) String() string {
	switch at {
	case ATSum:
		return "sum"
	case ATMin:
		return "min"
	case ATMax:
		return "max"
	case ATFirst:
		return "first"
	case ATLast:
		return "last"
	case ATHistogram:
		return "histogram"
	default:
		panic("invalid aggregate type")
	}
}

func (at AggregateType) MarshalJSON() ([]byte, error) {
	return encoding.JSONMarshal(at.String()), nil
}

func (at *AggregateType) UnmarshalJSON(data []byte) error {
	var str string
	if err := encoding.JSONUnmarshal(data, &str); err != nil {
		return err
	}
	switch str {
	case "sum":
		*at = ATSum
	case "min":
		*at = ATMin
	case "max":
		*at = ATMax
	case "last":
		*at = ATLast
	case "first":
		*at = ATFirst
	case "histogram":
		*at = ATHistogram
	default:
		return errors.New("invalid aggregate type")
	}
	return nil
}

type Type interface{}

type Block interface{}
