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

	"github.com/lindb/common/models"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

type baseFunc struct {
	ctx  EvalContext
	args []Expression
}

func (*baseFunc) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalString(row types.Row) (val string, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalDuration(row types.Row) (val time.Duration, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalMap(row types.Row) (val map[string]string, isNull bool, err error) {
	panic("implement me")
}

func (*baseFunc) EvalExemplar(row types.Row) (val *models.Exemplar, isNull bool, err error) {
	panic("implement me")
}

type Func interface {
	EvalInt(row types.Row) (val int64, isNull bool, err error)
	EvalFloat(row types.Row) (val float64, isNull bool, err error)
	EvalString(row types.Row) (val string, isNull bool, err error)
	EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error)
	EvalDuration(row types.Row) (val time.Duration, isNull bool, err error)
	EvalTime(row types.Row) (val time.Time, isNull bool, err error)
	EvalMap(row types.Row) (val map[string]string, isNull bool, err error)
	EvalExemplar(row types.Row) (val *models.Exemplar, isNull bool, err error)
}

type NewFunc = func(ctx EvalContext, args []Expression) Func

// IsFuncSupported check if given function name is supported.
var funcs = map[tree.FuncName]NewFunc{
	tree.Plus:  newArithmeticPlusFunc,
	tree.Minus: newArithmeticMinusFunc,
	tree.Mul:   newArithmeticMulFunc,
	tree.Div:   newArithmeticDivFunc,
	tree.Mod:   newArithmeticModFunc,

	tree.Count:    newArithmeticPlusFunc,
	tree.Sampling: newSamplingFunc, // NOTE: just pass check function if exists

	// time functions
	// ref: https://dev.mysql.com/doc/refman/8.4/en/date-and-time-functions.html
	tree.DateAdd:   newAddSubDateFunc,
	tree.Now:       newNowFunc,
	tree.StrToDate: newStrToDateFunc,
	tree.TimeTrunc: newTimeTruncFunc,

	// map functions
	tree.MapValues: newMapValuesFunc,
}
