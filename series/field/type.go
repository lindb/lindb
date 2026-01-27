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

	"github.com/lindb/common/models"

	"github.com/lindb/lindb/spi/types"
)

// EmptyFieldID represents empty value for field id.
const EmptyFieldID = ID(0)

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
	Last
	First
	Exemplar
)

// Aggregate aggregates two float64 values into one
func (t AggType) Aggregate(a, b float64) float64 {
	switch t {
	case Sum, Count:
		return a + b
	case Last:
		return b
	case First:
		return a
	case Min:
		return math.Min(a, b)
	case Max:
		return math.Max(a, b)
	default:
		panic("unspecified AggType")
	}
}

func ExemplarAggregate(a, b *models.Exemplar) *models.Exemplar {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.Duration >= b.Duration {
		return a
	}
	return b
}

// Type represents field type for LinDB support
type Type uint8

// Defines all field types for LinDB support(user write)
const (
	Unknown Type = iota
	SumField
	MinField
	MaxField
	LastField
	HistogramField // alias for sumField, only visible for tsdb
	FirstField
	ExemplarField
)

func (t Type) IsExemplar() bool {
	return t == ExemplarField
}

func (t Type) AggregateType() types.AggregateType {
	switch t {
	case SumField:
		return types.ATSum
	case MinField:
		return types.ATMin
	case MaxField:
		return types.ATMax
	case LastField:
		return types.ATLast
	case HistogramField:
		return types.ATSum
	case FirstField:
		return types.ATFirst
	case ExemplarField:
		return types.ATExemplar
	default:
		panic("unknown aggregate type")
	}
}

// String returns the field type's string value
func (t Type) String() string {
	switch t {
	case SumField:
		return "sum"
	case MinField:
		return "min"
	case MaxField:
		return "max"
	case LastField:
		return "last"
	case HistogramField:
		return "histogram"
	case FirstField:
		return "first"
	case ExemplarField:
		return "exemplar"
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
	case LastField:
		return Last
	case FirstField:
		return First
	case ExemplarField:
		return Exemplar
	default:
		panic("need impl")
	}
}
