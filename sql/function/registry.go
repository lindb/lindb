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
	"github.com/lindb/lindb/sql/tree"
)

// FunctionSignature describes the arity constraints of a function.
type FunctionSignature struct {
	MinArgs int // minimum number of required arguments
	MaxArgs int // maximum number of arguments; -1 means variadic
}

// FunctionDef is the single registration entry for a built-in or user-defined function.
// Each execution layer consults the registry to obtain the implementation it needs:
//
//   - Expression layer      → Vectorized  (VectorFuncFactory)
//   - Streaming / CEP layer → Incremental (IncrementalAgg)
//   - Broker agg layer      → MergeFunc   (partial-result combiner)
//
// Nil fields indicate that the function does not support that execution mode.
// A nil MergeFunc causes the broker layer to fall back to additive merging.
type FunctionDef struct {
	Name      tree.FuncName
	Signature FunctionSignature

	Vectorized  VectorFuncFactory          // expression-layer: factory creates a per-query instance
	Incremental IncrementalAgg             // streaming per-group stateful aggregation
	MergeFunc   func(a, b float64) float64 // broker-layer partial-result combiner
}

// FunctionRegistry is the central store for all registered functions.
// A single global instance (DefaultRegistry) is initialised by builtin.Register.
type FunctionRegistry struct {
	defs map[tree.FuncName]*FunctionDef
}

// NewFunctionRegistry creates an empty registry.
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{defs: make(map[tree.FuncName]*FunctionDef)}
}

// Register adds or replaces the definition for def.Name.
func (r *FunctionRegistry) Register(def *FunctionDef) {
	r.defs[def.Name] = def
}

// Vectorized returns the VectorFuncFactory for name, or (nil, false) if not registered.
func (r *FunctionRegistry) Vectorized(name tree.FuncName) (VectorFuncFactory, bool) {
	if def, ok := r.defs[name]; ok && def.Vectorized != nil {
		return def.Vectorized, true
	}
	return nil, false
}

// Incremental returns the IncrementalAgg for name, or (nil, false) if not registered.
func (r *FunctionRegistry) Incremental(name tree.FuncName) (IncrementalAgg, bool) {
	if def, ok := r.defs[name]; ok && def.Incremental != nil {
		return def.Incremental, true
	}
	return nil, false
}

// MergeFunc returns the partial-result combiner for name.
// When the function has no registered MergeFunc, the default additive combiner
// (a + b) is returned so callers never need to handle a nil case.
func (r *FunctionRegistry) MergeFunc(name tree.FuncName) func(a, b float64) float64 {
	if def, ok := r.defs[name]; ok && def.MergeFunc != nil {
		return def.MergeFunc
	}
	return func(a, b float64) float64 { return a + b }
}

// DefaultRegistry is the global registry populated by builtin.Register().
// Execution layers call DefaultRegistry.Vectorized / Incremental / MergeFunc
// to obtain function implementations.
var DefaultRegistry = NewFunctionRegistry()
