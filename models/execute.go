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

package models

import (
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/timeutil"
)

// ExecuteParam represents lin query language executor's param.
type ExecuteParam struct {
	Database  string             `form:"db" json:"db"`
	SQL       string             `form:"sql" json:"sql" binding:"required"`
	TimeRange timeutil.TimeRange `form:"timeRange" json:"timeRange"`
	// Cursor is the composite pagination cursor "shardID:ts:logID" from the previous page response.
	Cursor string `form:"cursor" json:"cursor,omitempty"`
	// Limit overrides the per-page row count; 0 means use the server default.
	Limit int64 `form:"limit" json:"limit,omitempty"`
}

// LogPageResult is the HTTP response body for log query requests that include pagination.
type LogPageResult struct {
	Columns    []string        `json:"columns"`
	Values     [][]interface{} `json:"values"`
	NextCursor string          `json:"next_cursor,omitempty"` // "shardID:ts:logID,…"; empty on last page
	HasMore    bool            `json:"has_more"`
}

type Session struct {
	Database string `header:"X-LinDB-Database"`
}

type StatementType int

const (
	Unknown StatementType = iota
	DataDefinition
	Describe
	Select
)

func (stmt StatementType) String() string {
	switch stmt {
	case DataDefinition:
		return "DataDefinition"
	case Describe:
		return "Describe"
	case Select:
		return "Select"
	default:
		return constants.Unknown
	}
}
