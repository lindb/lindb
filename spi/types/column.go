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
	"fmt"
	"time"
)

type Column struct {
	Blocks    []Block `json:"block"`
	NumOfRows int     `json:"numOfRows"`
}

func NewColumn() *Column {
	return &Column{}
}

func (c *Column) AppendTimeSeries(val *TimeSeries) {
	c.Blocks = append(c.Blocks, val)
	c.NumOfRows++
}

func (c *Column) AppendString(val string) {
	v := String(val)
	c.Blocks = append(c.Blocks, &v)
	c.NumOfRows++
}

func (c *Column) AppendJSON(val json.RawMessage) {
	c.Blocks = append(c.Blocks, &val)
	c.NumOfRows++
}

func (c *Column) Append(val any) {
	c.Blocks = append(c.Blocks, val)
	c.NumOfRows++
}

func (c *Column) Reset(row int, val any) {
	c.Blocks[row] = val
}

func (c *Column) AppendInt(val int64) {
	v := Int(val)
	c.Blocks = append(c.Blocks, &v)
	c.NumOfRows++
}

func (c *Column) AppendFloat(val float64) {
	v := Float(val)
	c.Blocks = append(c.Blocks, &v)
	c.NumOfRows++
}

func (c *Column) AppendTimestamp(val time.Time) {
	c.Blocks = append(c.Blocks, &val)
	c.NumOfRows++
}

func (c *Column) AppendDuration(val time.Duration) {
	v := Int(val.Nanoseconds())
	c.Blocks = append(c.Blocks, &v)
	c.NumOfRows++
}

func (c *Column) GetString(row int) *String {
	if row >= len(c.Blocks) {
		return nil
	}
	val, ok := c.Blocks[row].(*String)
	if ok {
		return val
	}
	v := String(fmt.Sprintf("%v", val))
	return &v
}

func (c *Column) GetJSON(row int) *json.RawMessage {
	if row >= len(c.Blocks) {
		return nil
	}
	return c.Blocks[row].(*json.RawMessage)
}

func (c *Column) GetInt(row int) *Int {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row].(*Int)
}

func (c *Column) GetFloat(row int) *Float {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row].(*Float)
}

func (c *Column) GetTimestamp(row int) *time.Time {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row].(*time.Time)
}

func (c *Column) GetDuration(row int) *time.Duration {
	if row >= len(c.Blocks) {
		return nil
	}
	val := c.Blocks[row].(*Int)
	duration := time.Duration(int64(*val))
	return &duration
}

func (c *Column) GetTimeSeries(row int) *TimeSeries {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row].(*TimeSeries)
}

func (c *Column) Get(row int) any {
	if row >= len(c.Blocks) {
		return nil
	}
	// FIXME:
	return c.Blocks[row]
}
