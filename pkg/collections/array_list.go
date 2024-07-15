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

// FloatArray encapsulates methods for using the float array
// support mark pos if it has value
type FloatArray struct {
	it       *FloatArrayIterator
	marks    []uint8
	values   []float64
	capacity int
	size     int
	isSingle bool
}

// NewFloatArray creates a float array with a certain capacity
func NewFloatArray(capacity int) *FloatArray {
	markLen := capacity / blockSize
	if capacity%blockSize > 0 {
		markLen++
	}
	return &FloatArray{
		capacity: capacity,
		values:   make([]float64, capacity),
		marks:    make([]uint8, markLen),
	}
}

// HasValue returns if has value with pos
func (f *FloatArray) HasValue(pos int) bool {
	if !f.checkPos(pos) {
		return false
	}
	blockIdx := pos / blockSize
	idx := pos % blockSize
	mark := f.marks[blockIdx]
	return mark&(1<<uint64(idx)) != 0
}

// GetValue returns value with pos, if it has not value return 0
func (f *FloatArray) GetValue(pos int) float64 {
	if !f.checkPos(pos) {
		return 0
	}
	return f.values[pos]
}

// SetValue sets value with pos, if pos out of bounds, return it
func (f *FloatArray) SetValue(pos int, value float64) {
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
func (f *FloatArray) IsEmpty() bool {
	return f.size == 0
}

// Size returns size of array
func (f *FloatArray) Size() int {
	return f.size
}

// NewIterator returns an iterator over the array
func (f *FloatArray) NewIterator() *FloatArrayIterator {
	if f.it == nil {
		f.it = newFloatArrayIterator(f)
	} else {
		f.it.reset()
	}
	return f.it
}

// Capacity returns the capacity of array
func (f *FloatArray) Capacity() int {
	return f.capacity
}

// Marks returns the marks of array
func (f *FloatArray) Marks() []uint8 {
	return f.marks
}

// checkPos checks pos if out of bounds
func (f *FloatArray) checkPos(pos int) bool {
	if pos < 0 || pos >= f.capacity {
		return false
	}
	return true
}

// Reset resets all values and mark for reusing
func (f *FloatArray) Reset() {
	f.size = 0
	f.isSingle = false
	for i := range f.marks {
		f.marks[i] = 0
	}
}

// SetSingle sets is array is single value, mean all values is same
func (f *FloatArray) SetSingle(single bool) {
	f.isSingle = single
}

// IsSingle return if is single value
func (f *FloatArray) IsSingle() bool {
	return f.isSingle
}

// FloatArrayIterator represents a float array iterator
type FloatArrayIterator struct {
	fa       *FloatArray
	marks    []uint8
	idx      int
	count    int
	hasValue bool
	mark     uint8
}

// newFloatArrayIterator creates a float array iterator
func newFloatArrayIterator(fa *FloatArray) *FloatArrayIterator {
	return &FloatArrayIterator{
		fa:       fa,
		hasValue: true,
		marks:    fa.Marks(),
	}
}

func (it *FloatArrayIterator) reset() {
	it.idx = 0
	it.count = 0
	it.marks = it.fa.Marks()
	it.hasValue = true
}

// HasNext returns if this iterator has more values
func (it *FloatArrayIterator) HasNext() bool {
	for it.idx < it.fa.Capacity() && it.count < it.fa.Size() {
		blockIdx := it.idx / blockSize
		idx := it.idx % blockSize
		if idx == 0 {
			it.mark = it.marks[blockIdx]
		}
		it.idx++
		if it.mark&(1<<uint64(idx)) != 0 {
			it.count++
			return true
		}
	}
	it.hasValue = false
	return false
}

// Next returns the next value and index
func (it *FloatArrayIterator) Next() (idx int, value float64) {
	if !it.hasValue {
		return -1, 0
	}
	idx = it.idx - 1
	value = it.fa.GetValue(idx)
	return idx, value
}
