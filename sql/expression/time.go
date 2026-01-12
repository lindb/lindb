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
	"time"

	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/spi/types"
)

type addSubDateFunc struct {
	baseFunc
}

func newAddSubDateFunc(args []Expression) Func {
	return &addSubDateFunc{
		baseFunc: baseFunc{
			args: args,
		},
	}
}

func (n *addSubDateFunc) EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error) {
	// TODO: check error/now func
	tsStr, _, _ := n.args[0].EvalString(ctx, row)
	format := timeutil.DataTimeFormat2
	timestamp, err := timeutil.ParseTimestamp(tsStr, format)
	if err != nil {
		return time.Time{}, true, err
	}
	duration, _, _ := n.args[1].EvalDuration(ctx, row)
	return time.UnixMilli(timestamp).Add(duration), false, nil
}

type nowFunc struct {
	baseFunc
}

func newNowFunc(args []Expression) Func {
	return &nowFunc{}
}

func (n *nowFunc) EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error) {
	fmt.Println(ctx)
	return ctx.CurrentTime(), false, nil
}

type strToDateFunc struct {
	baseFunc
}

func newStrToDateFunc(args []Expression) Func {
	return &strToDateFunc{
		baseFunc: baseFunc{
			args: args,
		},
	}
}

func (n *strToDateFunc) EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error) {
	// TODO: check error
	tsStr, _, _ := n.args[0].EvalString(ctx, row)
	format, _, _ := n.args[1].EvalString(ctx, row)
	switch format {
	case "YYYYMMDD HH:mm:ss":
		format = timeutil.DataTimeFormat1
	case "YYYY-MM-DD HH:mm:ss":
		format = timeutil.DataTimeFormat2
	case "YYYY/MM/DD HH:mm:ss":
		format = timeutil.DataTimeFormat3
	case "YYYYMMDDHHmmss":
		format = timeutil.DataTimeFormat4
	}
	timestamp, err := timeutil.ParseTimestamp(tsStr, format)
	if err != nil {
		return time.Time{}, true, err
	}
	return time.UnixMilli(timestamp), false, nil
}

type timeTruncFunc struct {
	baseFunc
}

func newTimeTruncFunc(args []Expression) Func {
	return &timeTruncFunc{
		baseFunc: baseFunc{
			args: args,
		},
	}
}

func (n *timeTruncFunc) EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error) {
	ts, _, _ := n.args[0].EvalTime(ctx, row)
	duration, _, _ := n.args[1].EvalDuration(ctx, row)
	rs := ts.Truncate(duration)
	fmt.Printf("time trunc=%v,%v,%v\n", duration, ts, rs)
	return rs, false, nil
}
