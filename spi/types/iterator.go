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

// Iterator4Page represents a iterator that is used to iterator rows inside a page.
type Iterator4Page struct {
	page *Page

	cursor  int
	numRows int
}

// NewIterator4Page creates a iterator for Page.
func NewIterator4Page(page *Page) *Iterator4Page {
	return &Iterator4Page{
		page: page,
	}
}

// Begin resets the cursor of the iterator and returns the first Row.
func (it *Iterator4Page) Begin() Row {
	if it.page == nil {
		return it.End()
	}

	it.numRows = it.page.NumRows()
	if it.numRows == 0 {
		return it.End()
	}
	it.cursor = 1
	return it.page.GetRow(0)
}

// Next returns the next Row.
func (it *Iterator4Page) Next() Row {
	if it.cursor >= it.numRows {
		it.cursor = it.numRows + 1
		return it.End()
	}
	row := it.page.GetRow(it.cursor)
	it.cursor++
	return row
}

// End returns the invalid end Row.
func (it *Iterator4Page) End() Row {
	return Row{}
}
