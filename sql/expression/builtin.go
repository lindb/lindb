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

package expression

import (
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/function"
	_ "github.com/lindb/lindb/sql/function/builtin" // register built-in vector functions
	"github.com/lindb/lindb/sql/tree"
)

// Func is the expression-layer function interface.
// function.VectorFunc satisfies this interface (same Eval + EvalScalar signatures),
// so VectorFuncFactory output can be used directly as a Func — no wrapper needed.
type Func interface {
	EvalScalar() (scalar.Scalar, error)
	Eval(record arrow.RecordBatch) (arrow.Array, error)
}

type NewFunc = func(ctx EvalContext, args []Expression) Func

// newRegistryFunc returns a NewFunc that calls the VectorFuncFactory from the registry.
// The factory creates a per-query VectorFunc instance with ctx and args baked in.
// Panics at query-construction time (programming error) if the function is not registered.
func newRegistryFunc(name tree.FuncName) NewFunc {
	return func(ctx EvalContext, args []Expression) Func {
		factory, ok := function.DefaultRegistry.Vectorized(name)
		if !ok {
			panic(fmt.Sprintf("VectorFuncFactory for %q not registered in DefaultRegistry", name))
		}
		// expression.Expression satisfies function.Expr; function.VectorFunc satisfies Func.
		fArgs := make([]function.Expr, len(args))
		for i, a := range args {
			fArgs[i] = a
		}
		return factory(ctx, fArgs)
	}
}

// funcs maps function names to their NewFunc factories.
// All functions are now resolved from function.DefaultRegistry; the expression
// layer only provides the wrapper that adapts VectorFunc to the Func interface.
var funcs = map[tree.FuncName]NewFunc{
	// Arithmetic operators
	tree.Plus:  newRegistryFunc(tree.Plus),
	tree.Minus: newRegistryFunc(tree.Minus),
	tree.Mul:   newRegistryFunc(tree.Mul),
	tree.Div:   newRegistryFunc(tree.Div),
	tree.Mod:   newRegistryFunc(tree.Mod),

	// Aggregation functions
	tree.Sum:   newRegistryFunc(tree.Sum),
	tree.Count: newRegistryFunc(tree.Count),
	tree.Min:   newRegistryFunc(tree.Min),
	tree.Max:   newRegistryFunc(tree.Max),
	tree.First: newRegistryFunc(tree.First),
	tree.Last:  newRegistryFunc(tree.Last),

	tree.Sampling: newSamplingFunc, // NOTE: sampling stays in expression (tracing-specific)

	// Time functions
	// ref: https://dev.mysql.com/doc/refman/8.4/en/date-and-time-functions.html
	tree.DateAdd:   newRegistryFunc(tree.DateAdd),
	tree.Now:       newRegistryFunc(tree.Now),
	tree.StrToDate: newRegistryFunc(tree.StrToDate),
	tree.TimeTrunc: newRegistryFunc(tree.TimeTrunc),

	// Map functions
	tree.MapValues: newRegistryFunc(tree.MapValues),

	// Math functions
	// ref: https://dev.mysql.com/doc/refman/8.4/en/mathematical-functions.html
	tree.Rand: newRegistryFunc(tree.Rand),
}
