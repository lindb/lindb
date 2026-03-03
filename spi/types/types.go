// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package types

import (
	"math"

	"github.com/apache/arrow-go/v18/arrow"
	larray "github.com/lindb/arrow/pkg/arrow/array"
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
	// DTDuration represents duration data type.
	DTDuration
	// DTTimestamp represents timestamp data type.
	DTTimestamp
	// DTTimeSeries represents time series data type.
	DTTimeSeries
	// DTJSON represents json data type.
	DTJSON
	// DTDynamic represents dynamic data type.
	DTDynamic
	// DTMap represents map data type.
	DTMap
	// DTExemplar represents exemplar data type.
	DTExemplar
	DTBinary
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
	// ATExemplar represents exemplar aggregation type.
	ATExemplar
)

func (dt DataType) ToArrowDataType() arrow.DataType {
	switch dt {
	case DTString:
		return arrow.BinaryTypes.String
	case DTInt:
		return arrow.PrimitiveTypes.Int64
	case DTFloat:
		return arrow.PrimitiveTypes.Float64
	case DTDuration:
		return arrow.FixedWidthTypes.Duration_ns
	case DTTimestamp:
		return arrow.FixedWidthTypes.Timestamp_ms
	case DTTimeSeries:
		return arrow.StructOf(
			arrow.Field{Name: "timestamp", Type: arrow.FixedWidthTypes.Timestamp_ms},
			arrow.Field{Name: "value", Type: arrow.PrimitiveTypes.Float64},
		)
	case DTJSON:
		return arrow.BinaryTypes.String
	case DTDynamic:
		return arrow.BinaryTypes.String
	case DTMap:
		return arrow.MapOf(arrow.BinaryTypes.String, arrow.BinaryTypes.String)
	case DTExemplar:
		return larray.NewExemplarType()
	case DTBinary:
		return arrow.BinaryTypes.Binary
	default:
		panic("unsupported data type")
	}
}

func (dt DataType) String() string {
	switch dt {
	case DTString:
		return "string"
	case DTInt:
		return "int"
	case DTFloat:
		return "float"
	case DTDuration:
		return "duration"
	case DTTimestamp:
		return "timestamp"
	case DTTimeSeries:
		return "time_series"
	case DTJSON:
		return "json"
	case DTDynamic:
		return "dynamic"
	case DTMap:
		return "map"
	case DTExemplar:
		return "exemplar"
	default:
		return "unknown"
	}
}

func (dt DataType) MarshalJSON() ([]byte, error) {
	return encoding.JSONMarshal(dt.String()), nil
}

func (dt *DataType) UnmarshalJSON(data []byte) error {
	var str string
	if err := encoding.JSONUnmarshal(data, &str); err != nil {
		return err
	}
	switch str {
	case "string":
		*dt = DTString
	case "int":
		*dt = DTInt
	case "float":
		*dt = DTFloat
	case "duration":
		*dt = DTDuration
	case "timestamp":
		*dt = DTTimestamp
	case "time_series":
		*dt = DTTimeSeries
	case "json":
		*dt = DTJSON
	case "dynamic":
		*dt = DTDynamic
	case "map":
		*dt = DTMap
	case "exemplar":
		*dt = DTExemplar
	default:
		*dt = DTUnknown
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
	case ATExemplar:
		return "exemplar"
	default:
		return ""
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
	case "exemplar":
		*at = ATExemplar
	default:
		*at = ATUnknown
	}
	return nil
}

// Aggregate aggregates two float64 values into one
func (at AggregateType) Aggregate(a, b float64) float64 {
	switch at {
	case ATSum:
		return a + b
	case ATLast:
		return b
	case ATFirst:
		return a
	case ATMin:
		return math.Min(a, b)
	case ATMax:
		return math.Max(a, b)
	default:
		panic("unspecified AggregateType")
	}
}

type Type any

type Value any
