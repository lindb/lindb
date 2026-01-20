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

	"github.com/lindb/lindb/pkg/stream"
)

var (
	RowWithEmptyValue = NewRowWithEmptyValue()
	EmptyRow          = &PageRow{}
)

func MarshalPage(page *Page) ([]byte, error) {
	w := stream.NewBufferWriter(nil)
	page.Marshal(w)
	return w.Bytes()
}

func UnmarshalPage(data []byte) (*Page, error) {
	r := stream.NewReader(data)
	page := NewPage()
	page.Unmarshal(r)
	return page, r.Error()
}

type Page struct {
	Layout   []ColumnMetadata `json:"layout,omitempty"`
	Grouping []int            `json:"grouping,omitempty"` // grouping column indexes
	Columns  []*Column        `json:"columns,omitempty"`

	Error string `json:"error,omitempty"`

	numRows int // cache num of rows
}

func NewPage() *Page {
	return &Page{numRows: -1}
}

func NewRowWithEmptyValue() *Page {
	page := NewPage()
	column := NewColumn()
	page.AppendColumn(ColumnMetadata{DataType: DTString}, column)
	column.Append("") // mock empty value
	return page
}

func (p *Page) Marshal(w *stream.BufferWriter) {
	// write error
	if p.Error != "" {
		w.PutByte(1)
		w.PutString(p.Error)
		return
	}
	// no error
	w.PutByte(0)
	// write Layout
	w.PutUvarint32(uint32(len(p.Layout)))
	for _, c := range p.Layout {
		c.Marshal(w)
	}
	// write Grouping
	w.PutUvarint32(uint32(len(p.Grouping)))
	for _, idx := range p.Grouping {
		w.PutUvarint32(uint32(idx))
	}
	numOfRows := p.NumRows()
	// write NumOfRows
	w.PutUvarint32(uint32(numOfRows))
	if numOfRows == 0 {
		return
	}
	// write Columns
	for i, col := range p.Columns {
		if p.Layout[i].Hidden {
			continue
		}
		col.Marshal(p.Layout[i], w)
	}
}

func (p *Page) Unmarshal(r *stream.Reader) {
	// read error
	hasError := r.ReadByte()
	if hasError == 1 {
		p.Error = r.ReadString()
		return
	}
	// read layout
	size := r.ReadUvarint32()
	p.Layout = make([]ColumnMetadata, size)
	for i := range size {
		col := &ColumnMetadata{}
		col.Unmarshal(r)
		p.Layout[i] = *col
	}
	// read grouping
	groupingSize := r.ReadUvarint32()
	p.Grouping = make([]int, groupingSize)
	for i := range groupingSize {
		idx := int(r.ReadUvarint32())
		p.Grouping[i] = idx
	}
	numOfRows := r.ReadUvarint32()
	p.numRows = int(numOfRows)
	if numOfRows == 0 {
		return
	}
	// read columns
	p.Columns = make([]*Column, size)
	for i := range size {
		col := NewColumn()
		if !p.Layout[i].Hidden {
			col.Unmarshal(p.Layout[i], p.numRows, r)
		}
		p.Columns[i] = col
	}
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
	if idx >= p.NumRows() || idx < 0 {
		return EmptyRow
	}
	return &PageRow{p: p, idx: idx}
}

// NumRows returns the number of rows in the page.
func (p *Page) NumRows() int {
	if p.numRows != -1 {
		return p.numRows
	}
	if len(p.Columns) == 0 {
		return 0
	}
	// TODO: select rows/no column
	p.numRows = lo.MaxBy(p.Columns, func(a, b *Column) bool {
		return a.NumOfRows > b.NumOfRows
	}).NumOfRows
	return p.numRows
}

func (p *Page) Iterator() *Iterator4Page {
	return NewIterator4Page(p)
}
