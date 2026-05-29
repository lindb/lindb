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

package function

import (
	"github.com/apache/arrow-go/v18/arrow"

	"github.com/lindb/lindb/spi/scalar"
)

// Expr is the evaluation interface for function arguments.
//
// expression.Expression satisfies this interface via structural typing, so
// execution layers can pass argument expressions to VectorFunc and Accumulator
// implementations directly without any adapter or pre-evaluation in the framework.
//
// Functions receive the full RecordBatch and decide internally how to extract
// column data — consistent with the existing expression.Func.Eval(record) pattern.
type Expr interface {
	// Eval evaluates this expression against the given batch and returns the
	// resulting column array.  The caller owns the returned array and must
	// release it when no longer needed.
	Eval(record arrow.RecordBatch) (arrow.Array, error)

	// EvalScalar evaluates this expression as a compile-time constant.
	// Returns an error when the expression depends on runtime row data.
	// Used to detect constant arguments and enable scalar-optimised code paths.
	EvalScalar() (scalar.Scalar, error)
}
