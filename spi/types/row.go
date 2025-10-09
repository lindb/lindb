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

// Row represents a row in the page.
type Row struct {
	p   *Page
	idx int
}

// GetString returns the string value in the row with the column index.
func (r *Row) GetString(colIdx int) *String {
	return r.p.Columns[colIdx].GetString(r.idx)
}

// GetJSON returns the json value in the row with the column index.
func (r *Row) GetJSON(colIdx int) *json.RawMessage {
	return r.p.Columns[colIdx].GetJSON(r.idx)
}

// GetFloat returns the float value in the row with the column index.
func (r *Row) GetFloat(colIdx int) *Float {
	return r.p.Columns[colIdx].GetFloat(r.idx)
}

// GetInt returns the int value in the row with the column index.
func (r *Row) GetInt(colIdx int) *Int {
	return r.p.Columns[colIdx].GetInt(r.idx)
}

// GetTimeSeries returns the time series value in the row with the column index.
func (r *Row) GetTimeSeries(colIdx int) *TimeSeries {
	return r.p.Columns[colIdx].GetTimeSeries(r.idx)
}

func (r *Row) GetTimestamp(colIdx int) *time.Time {
	return r.p.Columns[colIdx].GetTimestamp(r.idx)
}

func (r *Row) GetDuration(colIdx int) *time.Duration {
	return r.p.Columns[colIdx].GetDuration(r.idx)
}

func (r *Row) Get(colIdx int) any {
	return r.p.Columns[colIdx].Get(r.idx)
}
