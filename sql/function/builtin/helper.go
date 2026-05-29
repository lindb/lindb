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
	"errors"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	larray "github.com/lindb/arrow/pkg/arrow/array"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/function"
)

// errNoScalarEval is returned by aggregation VectorFunc.EvalScalar implementations
// to signal that the function cannot be folded to a constant at plan time.
var errNoScalarEval = errors.New("aggregation function cannot be evaluated as a scalar constant")

// identityInstance is a per-query VectorFunc that passes through the first argument
// column unchanged.  Used for aggregation functions (sum, count, min, max, first, last)
// where the broker already holds pre-aggregated values and only needs to relay them.
type identityInstance struct {
	arg function.Expr
}

func (f *identityInstance) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	if f.arg == nil {
		return nil, errors.New("aggregation VectorFunc requires at least one argument")
	}
	return f.arg.Eval(record)
}

func (f *identityInstance) EvalScalar() (scalar.Scalar, error) {
	return nil, errNoScalarEval
}

// IdentityFactory is the VectorFuncFactory for identity (pass-through) aggregation functions.
var IdentityFactory function.VectorFuncFactory = func(_ function.EvalContext, args []function.Expr) function.VectorFunc {
	var arg function.Expr
	if len(args) > 0 {
		arg = args[0]
	}
	return &identityInstance{arg: arg}
}

// readScalarValue reads a float64 from an Aggregation or plain Float64 column.
// Returns 0 for null or unsupported column types.
func readScalarValue(col arrow.Array, rowIdx int) float64 {
	if col == nil || col.IsNull(rowIdx) {
		return 0
	}
	switch a := col.(type) {
	case *larray.Aggregation:
		return a.Value(rowIdx)
	case *array.Float64:
		return a.Value(rowIdx)
	case *array.Int64:
		return float64(a.Value(rowIdx))
	}
	return 0
}

// appendAggResult writes val to b, which must be an ExtensionBuilder wrapping an
// Aggregation-type column.  AppendNull is called on type mismatch.
func appendAggResult(b array.Builder, val float64) {
	eb, ok := b.(*array.ExtensionBuilder)
	if !ok {
		b.AppendNull()
		return
	}
	larray.NewAggregationBuilder(eb).Append(val)
}
