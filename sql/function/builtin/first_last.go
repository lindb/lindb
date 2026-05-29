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

// ── First ─────────────────────────────────────────────────────────────────────

type firstIncrementalAgg struct{}

func (f *firstIncrementalAgg) NewAccumulator(_ function.EvalContext, args []function.Expr) function.Accumulator {
	var arg function.Expr
	if len(args) > 0 {
		arg = args[0]
	}
	return &firstAccumulator{arg: arg}
}

type firstAccumulator struct {
	value float64
	arg   function.Expr
	col   arrow.Array
	init  bool
}

func (a *firstAccumulator) Initialize(record arrow.RecordBatch) {
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

// Update records the first non-null value seen; subsequent values are ignored.
func (a *firstAccumulator) Update(rowIdx int) {
	if a.init || a.col == nil || a.col.IsNull(rowIdx) {
		return
	}
	a.value = readScalarValue(a.col, rowIdx)
	a.init = true
}

func (a *firstAccumulator) Result(b array.Builder) {
	if a.init {
		appendAggResult(b, a.value)
	} else {
		b.AppendNull()
	}
}

func (a *firstAccumulator) Reset() {
	a.value = 0
	a.init = false
	if a.col != nil {
		a.col.Release()
		a.col = nil
	}
}

// ── Last ──────────────────────────────────────────────────────────────────────

type lastIncrementalAgg struct{}

func (l *lastIncrementalAgg) NewAccumulator(_ function.EvalContext, args []function.Expr) function.Accumulator {
	var arg function.Expr
	if len(args) > 0 {
		arg = args[0]
	}
	return &lastAccumulator{arg: arg}
}

type lastAccumulator struct {
	value float64
	arg   function.Expr
	col   arrow.Array
	init  bool
}

func (a *lastAccumulator) Initialize(record arrow.RecordBatch) {
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

// Update always overwrites with the latest non-null value seen.
func (a *lastAccumulator) Update(rowIdx int) {
	if a.col == nil || a.col.IsNull(rowIdx) {
		return
	}
	a.value = readScalarValue(a.col, rowIdx)
	a.init = true
}

func (a *lastAccumulator) Result(b array.Builder) {
	if a.init {
		appendAggResult(b, a.value)
	} else {
		b.AppendNull()
	}
}

func (a *lastAccumulator) Reset() {
	a.value = 0
	a.init = false
	if a.col != nil {
		a.col.Release()
		a.col = nil
	}
}
