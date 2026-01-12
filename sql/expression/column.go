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
)

type Column struct {
	name    string
	retType types.DataType
	index   int
}

func NewColumn(ctx EvalContext, name string, index int, retType types.DataType) Expression {
	return &Column{name: name, index: index, retType: retType}
}

func (c *Column) EvalString(row types.Row) (val string, isNull bool, err error) {
	v := row.Get(c.index)
	switch sv := v.(type) {
	case string:
		return sv, false, nil
	default:
		val := row.GetString(c.index)
		if val == nil {
			return "", true, nil
		}
		return string(*val), false, nil
	}
}

func (c *Column) EvalMap(row types.Row) (val map[string]string, isNull bool, err error) {
	v := row.Get(c.index)
	if v == nil {
		return nil, true, nil
	}
	return v.(map[string]string), false, nil
}

func (c *Column) EvalInt(row types.Row) (val int64, isNull bool, err error) {
	return int64(*row.GetInt(c.index)), false, nil
}

func (c *Column) EvalFloat(row types.Row) (val float64, isNull bool, err error) {
	return float64(*row.GetFloat(c.index)), false, nil
}

func (c *Column) EvalTimeSeries(row types.Row) (val *types.TimeSeries, isNull bool, err error) {
	return row.GetTimeSeries(c.index), false, nil
}

func (c *Column) EvalDuration(row types.Row) (val time.Duration, isNull bool, err error) {
	return *row.GetDuration(c.index), false, nil
}

func (c *Column) EvalTime(row types.Row) (val time.Time, isNull bool, err error) {
	v := row.GetTimestamp(c.index)
	if v == nil {
		return time.Time{}, true, nil
	}
	return *v, false, nil
}

// GetType returns the data type of the column returns.
func (c *Column) GetType() types.DataType {
	return c.retType
}

// String returns the column in string format.
func (c *Column) String() string {
	return c.name
}
