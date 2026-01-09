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

import "github.com/lindb/lindb/models"

type TableMetadata struct {
	Schema     *TableSchema
	Partitions map[models.InternalNode][]int

	SupportDynamicField bool
}

type TableSchema struct {
	Columns []ColumnMetadata `json:"columns,omitempty"`
}

func NewTableSchema() *TableSchema {
	return &TableSchema{}
}

func (s *TableSchema) AddColumn(column ColumnMetadata) {
	s.Columns = append(s.Columns, column)
}

func (s *TableSchema) AddColumns(columns []ColumnMetadata) {
	s.Columns = append(s.Columns, columns...)
}

type ColumnMetadata struct {
	Name     string        `json:"name"`
	DataType DataType      `json:"type"`
	AggType  AggregateType `json:"aggType,omitempty"`
	Hidden   bool          `json:"hidden"`
	Ref      int           `json:"-"`
}

func NewColumnInfo(name string, vt DataType, hidden bool, aggType AggregateType) ColumnMetadata {
	return ColumnMetadata{
		Name:     name,
		DataType: vt,
		Hidden:   hidden,
		AggType:  aggType,
	}
}
