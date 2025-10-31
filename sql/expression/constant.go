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

type Constant struct {
	value   any
	retType types.DataType
}

func NewConstant(value any, retType types.DataType) Expression {
	return &Constant{
		retType: retType,
		value:   value,
	}
}

// EvalString implements Expression.
func (c *Constant) EvalString(_ EvalContext, _ types.Row) (val string, isNull bool, err error) {
	return c.value.(string), false, nil
}

func (c *Constant) EvalInt(_ EvalContext, _ types.Row) (val int64, isNull bool, err error) {
	return c.value.(int64), false, nil
}

func (c *Constant) EvalFloat(_ EvalContext, _ types.Row) (val float64, isNull bool, err error) {
	return
}

func (c *Constant) EvalTimeSeries(_ EvalContext, _ types.Row) (val *types.TimeSeries, isNull bool, err error) {
	return
}

func (c *Constant) EvalDuration(_ EvalContext, _ types.Row) (val time.Duration, isNull bool, err error) {
	val = c.value.(time.Duration)
	return
}

func (c *Constant) EvalTime(_ EvalContext, _ types.Row) (val time.Time, isNull bool, err error) {
	switch v := c.value.(type) {
	case string:
		timestamp, err := timeutil.ParseTimestamp(v, timeutil.DataTimeFormat2)
		if err != nil {
			return time.Time{}, true, err
		}
		return time.UnixMilli(timestamp), false, nil
	default:
		return
	}
}

func (c *Constant) EvalMap(_ EvalContext, _ types.Row) (val map[string]string, isNull bool, err error) {
	return
}

func (c *Constant) GetType() types.DataType {
	return c.retType
}

// String returns the constant in string format.
func (c *Constant) String() string {
	return fmt.Sprintf("%v", c.value)
}
