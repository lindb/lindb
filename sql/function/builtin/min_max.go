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
	"math"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"

	"github.com/lindb/lindb/sql/function"
)

// ── Min ───────────────────────────────────────────────────────────────────────

type minIncrementalAgg struct{}

func (m *minIncrementalAgg) NewAccumulator(_ function.EvalContext, args []function.Expr) function.Accumulator {
	var arg function.Expr
	if len(args) > 0 {
		arg = args[0]
	}
	return &minAccumulator{arg: arg, value: math.MaxFloat64}
}

type minAccumulator struct {
	value float64
	arg   function.Expr
	col   arrow.Array
	init  bool
}

func (a *minAccumulator) Initialize(record arrow.RecordBatch) {
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

func (a *minAccumulator) Update(rowIdx int) {
	if a.col == nil || a.col.IsNull(rowIdx) {
		return // skip NULL rows — they don't participate in MIN
	}
	v := readScalarValue(a.col, rowIdx)
	if !a.init || v < a.value {
		a.value = v
		a.init = true
	}
}

func (a *minAccumulator) Result(b array.Builder) {
	if a.init {
		appendAggResult(b, a.value)
	} else {
		b.AppendNull()
	}
}

func (a *minAccumulator) Reset() {
	a.value = math.MaxFloat64
	a.init = false
	if a.col != nil {
		a.col.Release()
		a.col = nil
	}
}

// ── Max ───────────────────────────────────────────────────────────────────────

type maxIncrementalAgg struct{}

func (m *maxIncrementalAgg) NewAccumulator(_ function.EvalContext, args []function.Expr) function.Accumulator {
	var arg function.Expr
	if len(args) > 0 {
		arg = args[0]
	}
	return &maxAccumulator{arg: arg, value: -math.MaxFloat64}
}

type maxAccumulator struct {
	value float64
	arg   function.Expr
	col   arrow.Array
	init  bool
}

func (a *maxAccumulator) Initialize(record arrow.RecordBatch) {
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

func (a *maxAccumulator) Update(rowIdx int) {
	if a.col == nil || a.col.IsNull(rowIdx) {
		return // skip NULL rows — they don't participate in MAX
	}
	v := readScalarValue(a.col, rowIdx)
	if !a.init || v > a.value {
		a.value = v
		a.init = true
	}
}

func (a *maxAccumulator) Result(b array.Builder) {
	if a.init {
		appendAggResult(b, a.value)
	} else {
		b.AppendNull()
	}
}

func (a *maxAccumulator) Reset() {
	a.value = -math.MaxFloat64
	a.init = false
	if a.col != nil {
		a.col.Release()
		a.col = nil
	}
}
