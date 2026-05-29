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
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"

	"github.com/lindb/lindb/sql/function"
)

// ── Sum ───────────────────────────────────────────────────────────────────────

type sumIncrementalAgg struct{}

func (s *sumIncrementalAgg) NewAccumulator(_ function.EvalContext, args []function.Expr) function.Accumulator {
	var arg function.Expr
	if len(args) > 0 {
		arg = args[0]
	}
	return &sumAccumulator{arg: arg}
}

type sumAccumulator struct {
	value float64
	arg   function.Expr // baked in at NewAccumulator time
	col   arrow.Array  // cached per-batch; released in Initialize and Reset
}

func (a *sumAccumulator) Initialize(record arrow.RecordBatch) {
	if a.col != nil {
		a.col.Release()
		a.col = nil
	}
	if a.arg != nil {
		col, err := a.arg.Eval(record)
		if err == nil {
			a.col = col
		}
	}
}

func (a *sumAccumulator) Update(rowIdx int) {
	if a.col != nil && !a.col.IsNull(rowIdx) {
		a.value += readScalarValue(a.col, rowIdx)
	}
}

func (a *sumAccumulator) Result(b array.Builder) {
	appendAggResult(b, a.value)
}

func (a *sumAccumulator) Reset() {
	a.value = 0
	if a.col != nil {
		a.col.Release()
		a.col = nil
	}
}
