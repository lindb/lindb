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

package table

import (
	"container/heap"
)

//go:generate mockgen -source ./iterator.go -destination=./iterator_mock.go -package table

// Iterator iterates over a store's key/value pairs in key order.
type Iterator interface {
	// HasNext returns if the iteration has more element.
	// It returns false if the iterator is exhausted.
	HasNext() bool
	// Key returns the key of the current key/value pair
	Key() uint32
	// Value returns the value of the current key/value pair
	Value() []byte
}

/////////////
// The priorityQueue is used to keep iterators sorted.
// reference:(https://golang.org/src/container/heap/example_pq_test.go)
////////////

// mergedIterator iterates over some iterator in key order
type mergedIterator struct {
	its []Iterator
	pq  priorityQueue

	curKey   uint32
	curValue []byte
}

// NewMergedIterator create merged iterator for multi iterators
func NewMergedIterator(its []Iterator) Iterator {
	it := &mergedIterator{
		its: its,
	}
	it.initQueue()
	return it
}

// initQueue initializes the priority queue
func (m *mergedIterator) initQueue() {
	i := 0
	for _, it := range m.its {
		if it.HasNext() {
			m.pq = append(m.pq, &item{
				it:    it,
				key:   it.Key(),
				value: it.Value(),
				index: i,
			})
			i++
		}
	}
	if len(m.pq) > 0 {
		heap.Init(&m.pq)
	}
}

// HasNext returns if the iteration has more element.
// It returns false if the iterator is exhausted.
func (m *mergedIterator) HasNext() bool {
	result := len(m.pq) > 0
	if result {
		// pop item and get value
		val := heap.Pop(&m.pq)
		item := val.(*item)
		m.curKey = item.key
		m.curValue = item.value

		// if it has value, push back queue and adjust priority
		it := item.it
		if it.HasNext() {
			item.key = it.Key()
			item.value = it.Value()
			m.pq.Push(item)
			m.pq.update(item)
		}
	}
	return result
}

// Key returns the key of the current key/value pair
// NOTICE: the key maybe is same as previous
func (m *mergedIterator) Key() uint32 {
	return m.curKey
}

// Value returns the value of the current key/value pair
func (m *mergedIterator) Value() []byte {
	return m.curValue
}

// item represents an item under priority queue, using key as priority.
type item struct {
	it Iterator

	key   uint32
	value []byte

	index int
}

// priorityQueue implements heap.Interface and holds Items.
type priorityQueue []*item

// Len returns the number of elements in priority queue
func (pq priorityQueue) Len() int { return len(pq) }

// Less compares key of item
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].key < pq[j].key
}

// Swap swaps the elements with indexes i and j.
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = j
	pq[j].index = i
}

// Push pushes a item into priority queue
func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and returns element of length -1
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority by the key of item
func (pq *priorityQueue) update(item *item) {
	heap.Fix(pq, item.index)
}
