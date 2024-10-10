package types

import (
	"github.com/lindb/common/pkg/encoding"
)

// AggregateType represents aggregation type of column value.
type AggregateType byte

// DataType represents data type of column value.
type DataType uint16

const (
	// DTUnknown represents unknown data type.
	DTUnknown DataType = iota
	// DTString represents string data type.
	DTString
	// DTInt represents int data type.
	DTInt
	// DTFloat represents float data type.
	DTFloat
	// DTTimestamp represents timestamp data type.
	DTTimestamp
	// DTTimeSeries represents time series data type.
	DTTimeSeries
)

const (
	// ATUnknown represents unknown aggregation type.
	ATUnknown AggregateType = iota
	// ATSum represents sum aggregation type.
	ATSum
	// ATMin represents min aggregation type.
	ATMin
	// ATMax represents max aggregation type.
	ATMax
	// ATLast represents last aggregation type.
	ATLast
	// ATFirst represents first aggregation type.
	ATFirst
	// ATHistogram represents histogram aggregation type.
	ATHistogram
)

func (dt DataType) String() string {
	switch dt {
	case DTString:
		return "string"
	case DTInt:
		return "int"
	case DTFloat:
		return "float"
	case DTTimestamp:
		return "timestamp"
	case DTTimeSeries:
		return "time_series"
	default:
		return "unknown"
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
		*tv = DTString
	case "int":
		*tv = DTInt
	case "float":
		*tv = DTFloat
	case "timestamp":
		*tv = DTTimestamp
	case "time_series":
		*tv = DTTimeSeries
	default:
		*tv = DTUnknown
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
		return "unknown"
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
		*at = ATUnknown
	}
	return nil
}

type Type any

type Block any
