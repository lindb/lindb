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

// ToTable returns metadata list as table if it has value, else return empty string.
func (m *Metadata) ToTable() (rows int, tableStr string) {
	writer := NewTableFormatter()
	switch m.Type {
	case stmt.Namespace.String():
		return m.toTableForStringValues(table.Row{"Namespace"}, writer)
	case stmt.Metric.String():
		return m.toTableForStringValues(table.Row{"Metric"}, writer)
	case stmt.TagKey.String():
		return m.toTableForStringValues(table.Row{"Tag Key"}, writer)
	case stmt.TagValue.String():
		return m.toTableForStringValues(table.Row{"Tag Value"}, writer)
	case stmt.Field.String():
		return m.toTableForMapValues(table.Row{"Name", "Type"}, []string{"name", "type"}, writer)
	default:
		return 0, ""
	}
}

// toTableForStringValues returns table for string values.
func (m *Metadata) toTableForStringValues(header table.Row, writer table.Writer) (rows int, tableStr string) {
	writer.AppendHeader(header)
	values := m.Values.([]interface{})
	for i := range values {
		writer.AppendRow(table.Row{values[i]})
	}
	return len(values), writer.Render()
}

// toTableForMapValues returns table for map values.
func (m *Metadata) toTableForMapValues(header table.Row, cols []string, writer table.Writer) (rows int, tableStr string) {
	writer.AppendHeader(header)
	values := m.Values.([]interface{})
	for _, value := range values {
		mapValue := value.(map[string]interface{})
		var row table.Row
		for _, col := range cols {
			row = append(row, mapValue[col])
		}
		writer.AppendRow(row)
	}
	return len(values), writer.Render()
}

// Field represents field metadata
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
