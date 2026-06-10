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

package builtin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	larray "github.com/lindb/arrow/pkg/arrow/array"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/function"
)

// ConcatFactory creates a VectorFunc that concatenates all arguments as strings.
// Arguments can be any supported Arrow type: String, Int64, Float64, Aggregation, TimeSeries.
// Null values in any argument contribute an empty string (null is not propagated).
var ConcatFactory function.VectorFuncFactory = func(_ function.EvalContext, args []function.Expr) function.VectorFunc {
	return &concatFunc{args: args}
}

type concatFunc struct {
	args []function.Expr
}

func (f *concatFunc) EvalScalar() (scalar.Scalar, error) {
	return nil, fmt.Errorf("concat: scalar evaluation not supported")
}

// Eval concatenates each argument column value into a single string per row.
// For each row i, all argument arrays are converted to their string representation
// and joined without a separator — matching standard SQL CONCAT semantics.
func (f *concatFunc) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	// Evaluate all argument expressions upfront.
	evaluated := make([]arrow.Array, len(f.args))
	for i, arg := range f.args {
		arr, err := arg.Eval(record)
		if err != nil {
			return nil, fmt.Errorf("concat: arg %d eval: %w", i, err)
		}
		evaluated[i] = arr
	}

	n := int(record.NumRows())
	b := array.NewStringBuilder(memory.DefaultAllocator)
	defer b.Release()
	b.Reserve(n)

	var sb strings.Builder
	for i := range n {
		sb.Reset()
		for _, arr := range evaluated {
			// Null or missing array contribution is treated as empty string.
			if arr == nil || arr.IsNull(i) {
				continue
			}
			sb.WriteString(arrayValueToString(arr, i))
		}
		b.Append(sb.String())
	}
	return b.NewArray(), nil
}

// arrayValueToString converts the value at row i of an Arrow array to its string
// representation. Supported types: String, Int64, Float64, Aggregation, TimeSeries.
// For TimeSeries (e.g. COUNT result), all per-slot values are summed to produce
// the total count for the time window — this matches the user-visible "current count"
// semantic expected in alert message templates.
func arrayValueToString(arr arrow.Array, i int) string {
	switch a := arr.(type) {
	case *array.String:
		return a.Value(i)
	case *array.Int64:
		return strconv.FormatInt(a.Value(i), 10)
	case *array.Float64:
		return strconv.FormatFloat(a.Value(i), 'f', -1, 64)
	case *larray.Aggregation:
		// Aggregation stores a single float64 scalar per row (Sum, Min, Max, etc.).
		return strconv.FormatFloat(a.Value(i), 'f', -1, 64)
	case *larray.TimeSeries:
		// TimeSeries stores per-time-slot counts/values; sum all slots to get the total.
		var total float64
		for _, v := range a.Values(i) {
			total += v
		}
		return strconv.FormatFloat(total, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", arr)
	}
}
