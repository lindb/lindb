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

type Row interface {
	GetString(colIdx int) *String
	GetJSON(colIdx int) *json.RawMessage
	GetFloat(colIdx int) *Float
	GetInt(colIdx int) *Int
	GetTimeSeries(colIdx int) *TimeSeries
	GetTimestamp(colIdx int) *time.Time
	GetDuration(colIdx int) *time.Duration
	Get(colIdx int) any
}

type ArrayRow struct {
	columns []any
}

func NewArrayRow(columns []any) Row {
	return &ArrayRow{columns: columns}
}

// Get implements Row.
func (a *ArrayRow) Get(colIdx int) any {
	return a.columns[colIdx]
}

// GetDuration implements Row.
func (a *ArrayRow) GetDuration(colIdx int) *time.Duration {
	return a.columns[colIdx].(*time.Duration)
}

// GetFloat implements Row.
func (a *ArrayRow) GetFloat(colIdx int) *Float {
	return a.columns[colIdx].(*Float)
}

// GetInt implements Row.
func (a *ArrayRow) GetInt(colIdx int) *Int {
	return a.columns[colIdx].(*Int)
}

// GetJSON implements Row.
func (a *ArrayRow) GetJSON(colIdx int) *json.RawMessage {
	return a.columns[colIdx].(*json.RawMessage)
}

// GetString implements Row.
func (a *ArrayRow) GetString(colIdx int) *String {
	return a.columns[colIdx].(*String)
}

// GetTimeSeries implements Row.
func (a *ArrayRow) GetTimeSeries(colIdx int) *TimeSeries {
	return a.columns[colIdx].(*TimeSeries)
}

// GetTimestamp implements Row.
func (a *ArrayRow) GetTimestamp(colIdx int) *time.Time {
	return a.columns[colIdx].(*time.Time)
}

// PageRow represents a row in the page.
type PageRow struct {
	p   *Page
	idx int
}

// GetString returns the string value in the row with the column index.
func (r *PageRow) GetString(colIdx int) *String {
	if r.idx >= len(r.p.Columns) {
		return nil
	}
	return r.p.Columns[colIdx].GetString(r.idx)
}

// GetJSON returns the json value in the row with the column index.
func (r *PageRow) GetJSON(colIdx int) *json.RawMessage {
	if r.idx >= len(r.p.Columns) {
		return nil
	}
	return r.p.Columns[colIdx].GetJSON(r.idx)
}

// GetFloat returns the float value in the row with the column index.
func (r *PageRow) GetFloat(colIdx int) *Float {
	if r.idx >= len(r.p.Columns) {
		return nil
	}
	return r.p.Columns[colIdx].GetFloat(r.idx)
}

// GetInt returns the int value in the row with the column index.
func (r *PageRow) GetInt(colIdx int) *Int {
	if r.idx >= len(r.p.Columns) {
		return nil
	}
	return r.p.Columns[colIdx].GetInt(r.idx)
}

// GetTimeSeries returns the time series value in the row with the column index.
func (r *PageRow) GetTimeSeries(colIdx int) *TimeSeries {
	if r.idx >= len(r.p.Columns) {
		return nil
	}
	return r.p.Columns[colIdx].GetTimeSeries(r.idx)
}

func (r *PageRow) GetTimestamp(colIdx int) *time.Time {
	if r.idx >= len(r.p.Columns) {
		return nil
	}
	return r.p.Columns[colIdx].GetTimestamp(r.idx)
}

func (r *PageRow) GetDuration(colIdx int) *time.Duration {
	if r.idx >= len(r.p.Columns) {
		return nil
	}
	return r.p.Columns[colIdx].GetDuration(r.idx)
}

func (r *PageRow) Get(colIdx int) any {
	if r.idx >= len(r.p.Columns) {
		return nil
	}
	return r.p.Columns[colIdx].Get(r.idx)
}
