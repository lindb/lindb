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

package tree

import (
	"github.com/apache/arrow-go/v18/arrow"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/samber/lo"
)

// FuncName represents function name.
type FuncName string

// AggFuncName represents aggregation function name.
type AggFuncName string

const (
	// scalar function names
	Plus  FuncName = "plus"
	Minus FuncName = "minus"
	Div   FuncName = "div"
	Mul   FuncName = "mul"
	Mod   FuncName = "mod"

	// aggregation function names
	Sum      FuncName = "sum"
	Min      FuncName = "min"
	Max      FuncName = "max"
	First    FuncName = "first"
	Last     FuncName = "last"
	Count    FuncName = "count"
	Sampling FuncName = "sampling"

	// histogram aggregation function names (Prometheus-style)
	HistogramQuantile FuncName = "histogram_quantile"
	HistogramAvg      FuncName = "histogram_avg"
	HistogramSum      FuncName = "histogram_sum"
	HistogramCount    FuncName = "histogram_count"

	// time function names
	DateAdd   FuncName = "date_add"
	Now       FuncName = "now"
	StrToDate FuncName = "str_to_date"
	TimeTrunc FuncName = "time_trunc"

	// map function names
	MapValues FuncName = "map_values"

	// string function names
	Concat FuncName = "concat"

	// math function names
	// ref: https://dev.mysql.com/doc/refman/8.4/en/mathematical-functions.html
	Rand FuncName = "rand"
)

type GetFuncReturnType func(name FuncName) arrow.DataType

func GetDefaultFuncReturnType(name FuncName) arrow.DataType {
	return defaultFuncReturnTypes[name]
}

func GetStreamingFuncReturnType(name FuncName) arrow.DataType {
	return streamingFuncReturnTypes[name]
}

var defaultFuncReturnTypes = map[FuncName]arrow.DataType{
	DateAdd:   arrow.FixedWidthTypes.Timestamp_ns,
	Now:       arrow.FixedWidthTypes.Timestamp_ns,
	StrToDate: arrow.FixedWidthTypes.Timestamp_ns,
	TimeTrunc: arrow.FixedWidthTypes.Timestamp_ns,

	Count:    larrow.ExtensionTypes.TimeSeries,
	Sampling: larrow.ExtensionTypes.Exemplar,

	// histogram functions all produce a single scalar aggregation result
	HistogramQuantile: larrow.ExtensionTypes.Sum,
	HistogramAvg:      larrow.ExtensionTypes.Sum,
	HistogramSum:      larrow.ExtensionTypes.Sum,
	HistogramCount:    larrow.ExtensionTypes.Sum,

	MapValues: arrow.MapOf(arrow.BinaryTypes.String, arrow.BinaryTypes.String),

	// string functions
	Concat: arrow.BinaryTypes.String,

	// math functions
	Rand: arrow.PrimitiveTypes.Float64,
}

var streamingFuncReturnTypes = map[FuncName]arrow.DataType{
	Count: larrow.ExtensionTypes.Sum,
}

var defaultFuncAggTypes = map[FuncName]arrow.DataType{
	Count:    larrow.ExtensionTypes.Sum,
	Sum:      larrow.ExtensionTypes.Sum,
	Min:      larrow.ExtensionTypes.Min,
	Max:      larrow.ExtensionTypes.Max,
	First:    larrow.ExtensionTypes.First,
	Last:     larrow.ExtensionTypes.Last,
	Sampling: larrow.ExtensionTypes.Exemplar,

	// histogram functions: accept a Histogram-type column, return a scalar Sum result
	HistogramQuantile: larrow.ExtensionTypes.Histogram,
	HistogramAvg:      larrow.ExtensionTypes.Histogram,
	HistogramSum:      larrow.ExtensionTypes.Histogram,
	HistogramCount:    larrow.ExtensionTypes.Histogram,
}

func init() {
	streamingFuncReturnTypes = lo.Assign(defaultFuncReturnTypes, streamingFuncReturnTypes)
}

// IsAggFunc returns if given function name is an aggregation function.
func IsAggFunc(name FuncName) bool {
	_, ok := defaultFuncAggTypes[name]
	return ok
}

// IsHistogramFunc reports whether the function is a histogram aggregation.
func IsHistogramFunc(name FuncName) bool {
	switch name {
	case HistogramQuantile, HistogramAvg, HistogramSum, HistogramCount:
		return true
	}
	return false
}

// IsFuncSupported checks if given function name is supported.
func IsFuncSupported(name FuncName) bool {
	_, ok := funcs[name]
	return ok
}

var funcs = map[FuncName]struct{}{
	Plus:  {},
	Minus: {},
	Mul:   {},
	Div:   {},
	Mod:   {},

	Count:    {},
	Sampling: {},

	// histogram functions
	HistogramQuantile: {},
	HistogramAvg:      {},
	HistogramSum:      {},
	HistogramCount:    {},

	// time functions
	// ref: https://dev.mysql.com/doc/refman/8.4/en/date-and-time-functions.html
	DateAdd:   {},
	Now:       {},
	StrToDate: {},
	TimeTrunc: {},

	MapValues: {},

	// string functions
	Concat: {},

	// math functions
	// ref: https://dev.mysql.com/doc/refman/8.4/en/mathematical-functions.html
	Rand: {},
}
