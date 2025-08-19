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
	"github.com/samber/lo"
)

var (
	RowWithEmptyValue = NewRowWithEmptyValue()
	EmptyRow          = Row{}
)

type Page struct {
	Layout   []ColumnMetadata `json:"layout,omitempty"`
	Grouping []int            `json:"grouping,omitempty"` // grouping column indexes
	Columns  []*Column        `json:"columns,omitempty"`

	Error string `json:"error,omitempty"`
}

func NewPage() *Page {
	return &Page{}
}

func NewRowWithEmptyValue() *Page {
	page := NewPage()
	column := NewColumn()
	page.AppendColumn(ColumnMetadata{DataType: DTString}, column)
	column.AppendString("") // mock empty value
	return page
}

func (p *Page) SetGrouping(columnIndexes []int) {
	p.Grouping = columnIndexes
}

func (p *Page) AppendColumn(info ColumnMetadata, column *Column) {
	p.Layout = append(p.Layout, info)
	p.Columns = append(p.Columns, column)
}

// GetRow gets the Row in the page with the row index.
func (p *Page) GetRow(idx int) Row {
	// TODO: select rows?
	return Row{p: p, idx: idx}
}

// NumRows returns the number of rows in the page.
func (p *Page) NumRows() int {
	if len(p.Columns) == 0 {
		return 0
	}
	// TODO: select rows/no column
	return lo.MaxBy(p.Columns, func(a, b *Column) bool {
		return a.NumOfRows > b.NumOfRows
	}).NumOfRows
}

func (p *Page) Iterator() *Iterator4Page {
	return NewIterator4Page(p)
}
