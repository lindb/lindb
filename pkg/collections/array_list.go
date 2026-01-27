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

package collections

const blockSize = 8

// Array encapsulates methods for using the array
// support mark pos if it has value
type Array[V any] struct {
	marks    []uint8
	values   []V
	capacity int
	size     int
}

// NewArray creates a float array with a certain capacity
func NewArray[V any](capacity int) *Array[V] {
	markLen := capacity / blockSize
	if capacity%blockSize > 0 {
		markLen++
	}
	return &Array[V]{
		capacity: capacity,
		values:   make([]V, capacity),
		marks:    make([]uint8, markLen),
	}
}

// Values returns the values of array.
func (f *Array[V]) Values() []V {
	return f.values
}

// HasValue returns if has value with pos
func (f *Array[V]) HasValue(pos int) bool {
	if !f.checkPos(pos) {
		return false
	}
	blockIdx := pos / blockSize
	idx := pos % blockSize
	mark := f.marks[blockIdx]
	return mark&(1<<uint64(idx)) != 0
}

// GetValue returns value with pos, if it has not value return 0
// NOTE: need check if has value before get value
func (f *Array[V]) GetValue(pos int) V {
	return f.values[pos]
}

// SetValue sets value with pos, if pos out of bounds, return it
func (f *Array[V]) SetValue(pos int, value V) {
	if !f.checkPos(pos) {
		return
	}
	f.values[pos] = value

	if !f.HasValue(pos) {
		blockIdx := pos / blockSize
		idx := pos - pos/blockSize*blockSize
		mark := f.marks[blockIdx]
		mark |= 1 << uint64(idx)
		f.marks[blockIdx] = mark

		f.size++
	}
}

// IsEmpty tests if array is empty
func (f *Array[V]) IsEmpty() bool {
	return f.size == 0
}

// Size returns size of array
func (f *Array[V]) Size() int {
	return f.size
}

// Capacity returns the capacity of array
func (f *Array[V]) Capacity() int {
	return f.capacity
}

// checkPos checks pos if out of bounds
func (f *Array[V]) checkPos(pos int) bool {
	if pos < 0 || pos >= f.capacity {
		return false
	}
	return true
}

// Reset resets all values and mark for reusing
func (f *Array[V]) Reset() {
	f.size = 0
	for i := range f.marks {
		f.marks[i] = 0
	}
}
