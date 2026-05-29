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

// VectorFunc is the per-query instance of a batch-evaluation function.
//
// Its interface intentionally matches expression.Func (Eval + EvalScalar with no
// extra parameters) so that the expression layer can use a VectorFunc directly as
// a Func without any wrapping.
//
// Instances are created by VectorFuncFactory once per query invocation.  The factory
// bakes context and argument expressions in at construction time, allowing constant
// arguments to be evaluated once (probe via EvalScalar) and cached — rather than
// re-evaluated on every RecordBatch.
type VectorFunc interface {
	// Eval applies the function to all rows of record.
	// Constant args are cached from construction; non-constant args are resolved
	// via the baked-in Expr references.
	Eval(record arrow.RecordBatch) (arrow.Array, error)

	// EvalScalar returns the result as a compile-time constant.
	// Returns an error when the function depends on runtime row data.
	EvalScalar() (scalar.Scalar, error)
}

// VectorFuncFactory creates a per-query VectorFunc instance.
// It is called once when the expression tree is built for a query; the returned
// VectorFunc holds ctx and args in a closure and caches any constant arguments.
type VectorFuncFactory func(ctx EvalContext, args []Expr) VectorFunc
