package value

import (
	"errors"

	"github.com/lindb/common/pkg/encoding"
)

type AggregateType byte
type ValueType byte

const (
	VTString ValueType = iota + 1
	VTTimeSeries
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

func (vt ValueType) String() string {
	switch vt {
	case VTString:
		return "string"
	case VTTimeSeries:
		return "timeseries"
	default:
		panic("invalid value type")
	}
}

func (vt ValueType) MarshalJSON() ([]byte, error) {
	return encoding.JSONMarshal(vt.String()), nil
}

func (tv *ValueType) UnmarshalJSON(data []byte) error {
	var str string
	if err := encoding.JSONUnmarshal(data, &str); err != nil {
		return err
	}
	switch str {
	case "string":
		*tv = VTString
	case "timeseries":
		*tv = VTTimeSeries
	default:
		return errors.New("invalid value type")
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

type Type interface {
}

type Block interface{}
