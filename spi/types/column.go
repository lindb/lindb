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

package types

import (
	"encoding/json"
	"time"
)

type Column struct {
	Values    []Value
	NumOfRows int
}

func NewColumn() *Column {
	return &Column{}
}

func (c *Column) AppendTimeSeries(val *TimeSeries) {
	c.Values = append(c.Values, val)
	c.NumOfRows++
}

func (c *Column) AppendString(val string) {
	c.Values = append(c.Values, val)
	c.NumOfRows++
}

func (c *Column) AppendJSON(val json.RawMessage) {
	c.Values = append(c.Values, &val)
	c.NumOfRows++
}

func (c *Column) ApppendMap(val map[string]string) {
	c.Values = append(c.Values, val)
	c.NumOfRows++
}

func (c *Column) Append(val any) {
	c.Values = append(c.Values, val)
	c.NumOfRows++
}

func (c *Column) Reset(row int, val any) {
	c.Values[row] = val
}

func (c *Column) AppendInt(val int64) {
	c.Values = append(c.Values, val)
	c.NumOfRows++
}

func (c *Column) AppendFloat(val float64) {
	c.Values = append(c.Values, val)
	c.NumOfRows++
}

func (c *Column) AppendTimestamp(val time.Time) {
	c.Values = append(c.Values, val)
	c.NumOfRows++
}

func (c *Column) AppendDuration(val time.Duration) {
	c.Values = append(c.Values, val)
	c.NumOfRows++
}

func (c *Column) GetString(row int) string {
	if row >= len(c.Values) {
		return ""
	}
	return c.Values[row].(string)
}

func (c *Column) GetJSON(row int) json.RawMessage {
	if row >= len(c.Values) {
		return nil
	}
	return c.Values[row].(json.RawMessage)
}

func (c *Column) GetInt(row int) int64 {
	if row >= len(c.Values) {
		return 0
	}
	// FIXME:
	return c.Values[row].(int64)
}

func (c *Column) GetFloat(row int) float64 {
	if row >= len(c.Values) {
		return 0
	}
	// FIXME:
	return c.Values[row].(float64)
}

func (c *Column) GetTimestamp(row int) time.Time {
	if row >= len(c.Values) {
		return time.Time{}
	}
	return c.Values[row].(time.Time)
}

func (c *Column) GetDuration(row int) time.Duration {
	if row >= len(c.Values) {
		return time.Duration(0)
	}
	return c.Values[row].(time.Duration)
}

func (c *Column) GetMap(row int) map[string]string {
	if row >= len(c.Values) {
		return nil
	}
	return c.Values[row].(map[string]string)
}

func (c *Column) GetTimeSeries(row int) *TimeSeries {
	if row >= len(c.Values) {
		return nil
	}
	// FIXME:
	return c.Values[row].(*TimeSeries)
}

func (c *Column) Get(row int) any {
	if row >= len(c.Values) {
		return nil
	}
	// FIXME:
	return c.Values[row]
}
