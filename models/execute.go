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

// ExecuteParam represents lin query language executor's param.
type ExecuteParam struct {
	Database string `form:"db" json:"db"`
	SQL      string `form:"sql" json:"sql" binding:"required"`
}

type Session struct {
	Databases string `header:"X-LinDB-Database"`
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
		return "Unknown"
	}
}
