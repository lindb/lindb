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
	"github.com/apache/arrow-go/v18/arrow/array"
)

// IncrementalAgg is a factory for per-group stateful aggregation.
// It is used by the streaming (CEP) execution layer where rows are processed
// one at a time within each group window.
type IncrementalAgg interface {
	// NewAccumulator creates a fresh Accumulator with the argument expressions
	// baked in.  args are fixed for the lifetime of the query; the accumulator
	// calls args[i].Eval(record) in Initialize to resolve the actual column.
	NewAccumulator(ctx EvalContext, args []Expr) Accumulator
}

// Accumulator holds the intermediate state for one group during streaming aggregation.
//
// The function receives the full RecordBatch in Initialize and resolves its argument
// columns internally — consistent with the VectorFunc.Eval(record, args) pattern.
// This keeps the framework free of pre-evaluation logic and memory management.
type Accumulator interface {
	// Initialize resolves the baked-in argument expressions against the new batch
	// and caches the resulting columns.  Called once per incoming RecordBatch.
	Initialize(record arrow.RecordBatch)

	// Update incorporates one input row into the group state using the columns
	// cached by the most recent Initialize call.
	Update(rowIdx int)

	// Result writes the final aggregated value for this group to b and must
	// leave b in a valid state regardless of whether any rows were processed.
	Result(b array.Builder)

	// Reset clears the state so the Accumulator can be reused for the next group
	// without allocating a new instance.
	Reset()
}
