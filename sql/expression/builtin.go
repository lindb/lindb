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
	"github.com/apache/arrow-go/v18/arrow"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/tree"
)

type Func interface {
	EvalScalar() (scalar.Scalar, error)
	Eval(record arrow.RecordBatch) (arrow.Array, error)
}

type NewFunc = func(ctx EvalContext, args []Expression) Func

// IsFuncSupported check if given function name is supported.
var funcs = map[tree.FuncName]NewFunc{
	tree.Plus: func(ctx EvalContext, args []Expression) Func {
		return newArithmeticFunc(ctx, "plus", args)
	},
	tree.Minus: func(ctx EvalContext, args []Expression) Func {
		return newArithmeticFunc(ctx, "minus", args)
	},
	tree.Mul: func(ctx EvalContext, args []Expression) Func {
		return newArithmeticFunc(ctx, "mul", args)
	},
	tree.Div: func(ctx EvalContext, args []Expression) Func {
		return newArithmeticFunc(ctx, "div", args)
	},
	tree.Mod: func(ctx EvalContext, args []Expression) Func {
		return newArithmeticFunc(ctx, "mod", args)
	},
	tree.Sum: func(ctx EvalContext, args []Expression) Func {
		return newArithmeticFunc(ctx, "plus", args)
	},
	tree.Count: func(ctx EvalContext, args []Expression) Func {
		return newArithmeticFunc(ctx, "plus", args)
	},

	tree.Sampling: newSamplingFunc, // NOTE: just pass check function if exists

	// time functions
	// ref: https://dev.mysql.com/doc/refman/8.4/en/date-and-time-functions.html
	tree.DateAdd:   newAddSubDateFunc,
	tree.Now:       newNowFunc,
	tree.StrToDate: newStrToDateFunc,
	tree.TimeTrunc: newTimeTruncFunc,

	// map functions
	tree.MapValues: newMapValuesFunc,

	// math functions
	// ref: https://dev.mysql.com/doc/refman/8.4/en/mathematical-functions.html
	tree.Rand: newRandFunc,
}
