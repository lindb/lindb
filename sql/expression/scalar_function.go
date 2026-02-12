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
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/tree"
)

type ScalarFunc struct {
	ctx      EvalContext
	function Func
	funcName tree.FuncName
	args     []Expression

	rt ResultType
}

func NewScalarFunc(ctx EvalContext, funcName tree.FuncName, rt ResultType, args []Expression) Expression {
	newFn, ok := funcs[funcName]
	if !ok {
		panic(fmt.Sprintf("func not support, func name: %s", funcName))
	}
	fn := newFn(ctx, args)
	return &ScalarFunc{
		ctx:      ctx,
		rt:       rt,
		function: fn,
		funcName: funcName,
		args:     args,
	}
}

func (f *ScalarFunc) EvalScalar() (scalar.Scalar, error) {
	return f.function.EvalScalar()
}

func (f *ScalarFunc) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	return f.function.Eval(record)
}

func (f *ScalarFunc) ResultType() ResultType {
	return f.rt
}

// String returns the scalar function in string format.
func (f *ScalarFunc) String() string {
	return fmt.Sprintf("%s(%s)", f.funcName, strings.Join(lo.Map(f.args, func(item Expression, index int) string {
		return item.String()
	}), ","))
}
