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

package aggregation

import (
	"container/heap"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/sql/stmt"
)

// OrderByItem represents the order by expr item.
type OrderByItem struct {
	Expr     *stmt.OrderByExpr
	Name     string
	FuncType function.FuncType
	Desc     bool
}

// Row represents the series data for one group.
type Row interface {
	// GetValue returns the value based on given field name and function type.
	GetValue(fieldName string, funcType function.FuncType) float64
	// ResultSet returns the result set(tags/fields).
	ResultSet() (tags string, fields map[string]*collections.FloatArray)
}

// topNHeap represents topN heap.
type topNHeap struct {
	orderByItems []*OrderByItem
	rows         []Row

	limit int // max size of heap
	size  int // current size of heap
}

// newTopNHeap creates a topN heap.
func newTopNHeap(orderByItems []*OrderByItem, topN int) *topNHeap {
	return &topNHeap{
		limit:        topN,
		orderByItems: orderByItems,
	}
}

// Len is the number of elements in the collection.
func (h *topNHeap) Len() int {
	return h.size
}

// Swap swaps the elements with indexes i and j.
func (h *topNHeap) Swap(i, j int) {
	h.rows[i], h.rows[j] = h.rows[j], h.rows[i]
}

// Less compares the value of row based on order by items.
func (h *topNHeap) Less(i, j int) bool {
	for _, by := range h.orderByItems {
		ret := h.rows[i].GetValue(by.Name, by.FuncType) - h.rows[j].GetValue(by.Name, by.FuncType)
		if by.Desc {
			ret = -ret
		}
		if ret > 0 {
			return true
		} else if ret < 0 {
			return false
		}
		// if equals goto next order by item
	}
	return false
}

// Push pushes row into topN heap.
func (h *topNHeap) Push(row interface{}) {
	h.rows = append(h.rows, row.(Row))
	h.size++
}

func (h *topNHeap) Pop() interface{} {
	// never invoke here
	return nil
}

// Add adds row into topN heap.
func (h *topNHeap) Add(row Row) {
	if h.size >= h.limit {
		h.rows = append(h.rows, row)
		if h.Less(0, h.size) {
			h.Swap(0, h.size)
			heap.Fix(h, 0)
		}
		h.rows = h.rows[:h.size]
	} else {
		heap.Push(h, row)
	}
}

// ResultSet returns result set of topN.
func (h *topNHeap) ResultSet() []Row {
	return h.rows
}
