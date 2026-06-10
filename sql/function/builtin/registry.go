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

// Package builtin registers all built-in functions into function.DefaultRegistry.
// Import this package for its side effects:
//
//	import _ "github.com/lindb/lindb/sql/function/builtin"
package builtin

import (
	"math"

	"github.com/lindb/lindb/sql/function"
	"github.com/lindb/lindb/sql/tree"
)

func init() {
	Register(function.DefaultRegistry)
}

// Register adds all built-in function definitions to r.
// Execution layers that need a specific registry (e.g. tests) can call this
// directly instead of relying on the package-level init side-effect.
func Register(r *function.FunctionRegistry) {
	// ── Aggregation functions ─────────────────────────────────────────────────

	r.Register(&function.FunctionDef{
		Name:        tree.Sum,
		Signature:   function.FunctionSignature{MinArgs: 1, MaxArgs: 1},
		Vectorized:  IdentityFactory,
		Incremental: &sumIncrementalAgg{},
		MergeFunc:   func(a, b float64) float64 { return a + b },
	})
	r.Register(&function.FunctionDef{
		Name:        tree.Count,
		Signature:   function.FunctionSignature{MinArgs: 0, MaxArgs: 1},
		Vectorized:  IdentityFactory,
		Incremental: &countIncrementalAgg{},
		MergeFunc:   func(a, b float64) float64 { return a + b },
	})
	r.Register(&function.FunctionDef{
		Name:        tree.Min,
		Signature:   function.FunctionSignature{MinArgs: 1, MaxArgs: 1},
		Vectorized:  IdentityFactory,
		Incremental: &minIncrementalAgg{},
		MergeFunc:   math.Min,
	})
	r.Register(&function.FunctionDef{
		Name:        tree.Max,
		Signature:   function.FunctionSignature{MinArgs: 1, MaxArgs: 1},
		Vectorized:  IdentityFactory,
		Incremental: &maxIncrementalAgg{},
		MergeFunc:   math.Max,
	})
	r.Register(&function.FunctionDef{
		Name:        tree.First,
		Signature:   function.FunctionSignature{MinArgs: 1, MaxArgs: 1},
		Vectorized:  IdentityFactory,
		Incremental: &firstIncrementalAgg{},
		MergeFunc:   func(a, _ float64) float64 { return a },
	})
	r.Register(&function.FunctionDef{
		Name:        tree.Last,
		Signature:   function.FunctionSignature{MinArgs: 1, MaxArgs: 1},
		Vectorized:  IdentityFactory,
		Incremental: &lastIncrementalAgg{},
		MergeFunc:   func(_, b float64) float64 { return b },
	})

	// ── Arithmetic operators ──────────────────────────────────────────────────

	for _, entry := range []struct {
		name  tree.FuncName
		opStr string
	}{
		{tree.Plus, "plus"},
		{tree.Minus, "minus"},
		{tree.Mul, "mul"},
		{tree.Div, "div"},
		{tree.Mod, "mod"},
	} {
		r.Register(&function.FunctionDef{
			Name:       entry.name,
			Signature:  function.FunctionSignature{MinArgs: 2, MaxArgs: 2},
			Vectorized: newArithmeticFactory(entry.opStr),
		})
	}

	// ── Time functions ────────────────────────────────────────────────────────

	r.Register(&function.FunctionDef{
		Name:       tree.DateAdd,
		Signature:  function.FunctionSignature{MinArgs: 2, MaxArgs: 2},
		Vectorized: AddSubDateFactory,
	})
	r.Register(&function.FunctionDef{
		Name:       tree.Now,
		Signature:  function.FunctionSignature{MinArgs: 0, MaxArgs: 0},
		Vectorized: NowFactory,
	})
	r.Register(&function.FunctionDef{
		Name:       tree.StrToDate,
		Signature:  function.FunctionSignature{MinArgs: 2, MaxArgs: 2},
		Vectorized: StrToDateFactory,
	})
	r.Register(&function.FunctionDef{
		Name:       tree.TimeTrunc,
		Signature:  function.FunctionSignature{MinArgs: 2, MaxArgs: 2},
		Vectorized: TimeTruncFactory,
	})

	// ── Math functions ────────────────────────────────────────────────────────

	r.Register(&function.FunctionDef{
		Name:       tree.Rand,
		Signature:  function.FunctionSignature{MinArgs: 0, MaxArgs: 1},
		Vectorized: RandFactory,
	})

	// ── Map functions ─────────────────────────────────────────────────────────

	r.Register(&function.FunctionDef{
		Name:       tree.MapValues,
		Signature:  function.FunctionSignature{MinArgs: 2, MaxArgs: -1},
		Vectorized: MapValuesFactory,
	})

	// ── String functions ──────────────────────────────────────────────────────

	r.Register(&function.FunctionDef{
		Name:       tree.Concat,
		Signature:  function.FunctionSignature{MinArgs: 1, MaxArgs: -1},
		Vectorized: ConcatFactory,
	})
}
