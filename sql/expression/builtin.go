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
	"time"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type baseFunc struct {
	args []Expression
}

func (*baseFunc) EvalInt(ctx EvalContext, row types.Row) (val int64, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalFloat(ctx EvalContext, row types.Row) (val float64, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalString(ctx EvalContext, row types.Row) (val string, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalTimeSeries(ctx EvalContext, row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalDuration(ctx EvalContext, row types.Row) (val time.Duration, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error) {
	panic("implement me")
}

type Func interface {
	EvalInt(ctx EvalContext, row types.Row) (val int64, isNull bool, err error)
	EvalFloat(ctx EvalContext, row types.Row) (val float64, isNull bool, err error)
	EvalString(ctx EvalContext, row types.Row) (val string, isNull bool, err error)
	EvalTimeSeries(ctx EvalContext, row types.Row) (val *types.TimeSeries, isNull bool, err error)
	EvalDuration(ctx EvalContext, row types.Row) (val time.Duration, isNull bool, err error)
	EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error)
}

type FuncFactory interface {
	NewFunc(args []Expression) Func
}

// IsFuncSupported check if given function name is supported.
func IsFuncSupported(name tree.FuncName) bool {
	_, ok := funcs[name]
	return ok
}

var funcs = map[tree.FuncName]FuncFactory{
	tree.Plus:  &arithmeticPlusFuncFactory{},
	tree.Minus: &arithmeticMinusFuncFactory{},
	tree.Mul:   &arithmeticMulFuncFactory{},
	tree.Div:   &arithmeticDivFuncFactory{},
	tree.Mod:   &arithmeticModFuncFactory{},

	tree.Count: &arithmeticPlusFuncFactory{},

	// time functions
	// ref: https://dev.mysql.com/doc/refman/8.4/en/date-and-time-functions.html
	tree.DateAdd:   &addSubDateFuncFactory{},
	tree.Now:       &nowFuncFactory{},
	tree.StrToDate: &strToDateFuncFactory{},
}
