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

func newAddSubDateFunc(ctx EvalContext, args []Expression) Func {
	return &addSubDateFunc{
		baseFunc: baseFunc{
			ctx:  ctx,
			args: args,
		},
	}
}

func (n *addSubDateFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	// TODO: check error/now func
	tsStr, _, _ := n.args[0].EvalString(row)
	format := timeutil.DataTimeFormat2
	timestamp, err := timeutil.ParseTimestamp(tsStr, format)
	if err != nil {
		return time.Time{}, true, err
	}
	duration, _, _ := n.args[1].EvalDuration(row)
	return time.UnixMilli(timestamp).Add(duration), false, nil
}

type nowFunc struct {
	baseFunc
}

func newNowFunc(ctx EvalContext, args []Expression) Func {
	return &nowFunc{
		baseFunc: baseFunc{
			ctx:  ctx,
			args: args,
		},
	}
}

func (n *nowFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	return n.ctx.CurrentTime(), false, nil
}

type strToDateFunc struct {
	baseFunc
}

func newStrToDateFunc(ctx EvalContext, args []Expression) Func {
	return &strToDateFunc{
		baseFunc: baseFunc{
			ctx:  ctx,
			args: args,
		},
	}
}

func (n *strToDateFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	// TODO: check error
	tsStr, _, _ := n.args[0].EvalString(row)
	format, _, _ := n.args[1].EvalString(row)
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

	duration time.Duration
}

func newTimeTruncFunc(ctx EvalContext, args []Expression) Func {
	duration, _, _ := args[1].EvalDuration(types.EmptyRow)
	return &timeTruncFunc{
		duration: duration,
		baseFunc: baseFunc{
			ctx:  ctx,
			args: args,
		},
	}
}

func (n *timeTruncFunc) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	ts, _, _ := n.args[0].EvalTime(row)
	rs := ts.Truncate(n.duration)
	return rs, false, nil
}
