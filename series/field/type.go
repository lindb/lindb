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

package field

import (
	"math"

	"github.com/lindb/lindb/aggregation/function"
)

// AggType represents field's aggregator type.
type AggType uint8

// ID represents field id.
type ID uint8

// Name represents field name.
type Name string

func (n Name) String() string {
	return string(n)
}

// Defines all aggregator types for field
const (
	Sum AggType = iota + 1
	Count
	Min
	Max
	LastValue
)

// Aggregate aggregates two float64 values into one
func (t AggType) Aggregate(a, b float64) float64 {
	switch t {
	case Sum, Count:
		return a + b
	case LastValue:
		return b
	case Min:
		return math.Min(a, b)
	case Max:
		return math.Max(a, b)
	default:
		panic("unspecified AggType")
	}
}

// Type represents field type for LinDB support
type Type uint8

// Defines all field types for LinDB support(user write)
const (
	Unknown Type = iota
	SumField
	MinField
	MaxField
	GaugeField
	HistogramField // alias for sumField, only visible for tsdb
)

// String returns the field type's string value
func (t Type) String() string {
	switch t {
	case SumField:
		return "sum"
	case MinField:
		return "min"
	case MaxField:
		return "max"
	case GaugeField:
		return "gauge"
	case HistogramField:
		return "histogram"
	default:
		return "unknown"
	}
}

// AggType returns the aggregate function
func (t Type) AggType() AggType {
	switch t {
	case SumField, HistogramField:
		return Sum
	case MinField:
		return Min
	case MaxField:
		return Max
	default:
		//FIXME(stone1100)
		return Max
	}
}

func (t Type) DownSamplingFunc() function.FuncType {
	switch t {
	case SumField:
		return function.Sum
	case MinField:
		return function.Min
	case MaxField:
		return function.Max
	case GaugeField:
		return function.LastValue
	case HistogramField:
		return function.Sum
	default:
		return function.Unknown
	}
}

func (t Type) IsFuncSupported(funcType function.FuncType) bool {
	switch t {
	case SumField:
		switch funcType {
		case function.Sum, function.Min, function.Max:
			return true
		default:
			return false
		}
	case MinField:
		switch funcType {
		case function.Min:
			return true
		default:
			return false
		}
	case MaxField:
		switch funcType {
		case function.Max:
			return true
		default:
			return false
		}
	case GaugeField:
		switch funcType {
		case function.Sum, function.Min, function.Max, function.LastValue:
			return true
		default:
			return false
		}
	case HistogramField:
		switch funcType {
		case function.Sum:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// GetFuncFieldParams returns agg type for field aggregator by given function type.
func (t Type) GetFuncFieldParams(funcType function.FuncType) []AggType {
	switch t {
	case SumField:
		return getFieldParamsForSumField(funcType)
	case MinField:
		return getFieldParamsForMinField(funcType)
	case GaugeField:
		return getFieldParamsForGaugeField(funcType)
	case HistogramField:
		// Histogram field only supports sum
		return []AggType{Sum}
	}
	return nil
}

// GetDefaultFuncFieldParams returns default agg type for field aggregator.
func (t Type) GetDefaultFuncFieldParams() []AggType {
	switch t {
	case SumField:
		return []AggType{Sum}
	case MinField:
		return []AggType{Min}
	case GaugeField:
		return []AggType{LastValue}
	case HistogramField:
		return []AggType{Sum}
	}
	return nil
}

func getFieldParamsForSumField(funcType function.FuncType) []AggType {
	switch funcType {
	case function.Max:
		return []AggType{Max}
	default:
		return []AggType{Sum}
	}
}

func getFieldParamsForMinField(funcType function.FuncType) []AggType {
	switch funcType {
	case function.Max:
		return []AggType{Max}
	default:
		return []AggType{Min}
	}
}

func getFieldParamsForGaugeField(funcType function.FuncType) []AggType {
	switch funcType {
	case function.Max:
		return []AggType{Max}
	default:
		return []AggType{LastValue}
	}
}
