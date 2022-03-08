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
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/lindb/lindb/sql/stmt"
)

// Metadata represents metadata query result model
type Metadata struct {
	Type   string      `json:"type"`
	Values interface{} `json:"values"`
}

func (m *Metadata) ToTable() (int, string) {
	switch m.Type {
	case stmt.Database.String():
		writer := NewTableFormatter()
		writer.AppendHeader(table.Row{"Database"})
		dbs := m.Values.([]interface{})
		for i := range dbs {
			writer.AppendRow(table.Row{dbs[i]})
		}
		return len(dbs), writer.Render()
	default:
		return 0, ""
	}
}

// Field represents field metadata
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
