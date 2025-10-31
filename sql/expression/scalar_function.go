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
	"time"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type ScalarFunc struct {
	function Func
	funcName tree.FuncName
	args     []Expression
	retType  types.DataType
}

func NewScalarFunc(funcName tree.FuncName, retType types.DataType, args []Expression) (Expression, error) {
	newFn, ok := funcs[funcName]
	if !ok {
		return nil, fmt.Errorf("func not support, func name: %s", funcName)
	}
	fn := newFn(args)
	return &ScalarFunc{
		retType:  retType,
		function: fn,
		funcName: funcName,
		args:     args,
	}, nil
}

// EvalString implements Expression.
func (f *ScalarFunc) EvalString(ctx EvalContext, row types.Row) (val string, isNull bool, err error) {
	return f.function.EvalString(ctx, row)
}

func (f *ScalarFunc) EvalInt(ctx EvalContext, row types.Row) (val int64, isNull bool, err error) {
	return f.function.EvalInt(ctx, row)
}

func (f *ScalarFunc) EvalFloat(ctx EvalContext, row types.Row) (val float64, isNull bool, err error) {
	return f.function.EvalFloat(ctx, row)
}

func (f *ScalarFunc) EvalTimeSeries(ctx EvalContext, row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	return f.function.EvalTimeSeries(ctx, row)
}

func (f *ScalarFunc) EvalDuration(ctx EvalContext, row types.Row) (val time.Duration, isNull bool, err error) {
	return f.function.EvalDuration(ctx, row)
}

func (f *ScalarFunc) EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error) {
	return f.function.EvalTime(ctx, row)
}

func (f *ScalarFunc) EvalMap(ctx EvalContext, row types.Row) (val map[string]string, isNull bool, err error) {
	return f.function.EvalMap(ctx, row)
}

// GetType implements Expression.
func (f *ScalarFunc) GetType() types.DataType {
	return f.retType
}

// String returns the scalar function in string format.
func (f *ScalarFunc) String() string {
	return fmt.Sprintf("%s(%s)", f.funcName, strings.Join(lo.Map(f.args, func(item Expression, index int) string {
		return item.String()
	}), ","))
}
