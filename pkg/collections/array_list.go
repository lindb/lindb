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

// FloatArray represents a float array
type FloatArray interface {
	// Iterator returns an iterator over the array
	Iterator() FloatArrayIterator
	// GetValue returns value with pos, if has not value return 0
	GetValue(pos int) float64
	// HasValue returns if has value with pos
	HasValue(pos int) bool
	// SetValue sets value with pos, if pos out of bounds, return it
	SetValue(pos int, value float64)
	// IsEmpty tests if array is empty
	IsEmpty() bool
	// Size returns size of array
	Size() int
	// Capacity returns the capacity of array
	Capacity() int
	// Marks returns the marks of array
	Marks() []uint8
	// Reset resets all values and mark for reusing
	Reset()
	// SetSingle sets is array is single value, mean all values is same
	SetSingle(single bool)
	// IsSingle return if is single value
	IsSingle() bool
}

// floatArray represents a float array, support mark pos if has value
type floatArray struct {
	marks    []uint8
	values   []float64
	isSingle bool

	capacity int
	size     int

	it *floatArrayIterator
}

// NewFloatArray creates a float array
func NewFloatArray(capacity int) FloatArray {
	markLen := capacity / blockSize
	if capacity%blockSize > 0 {
		markLen++
	}
	return &floatArray{
		capacity: capacity,
		values:   make([]float64, capacity),
		marks:    make([]uint8, markLen),
	}
}

// HasValue returns if has value with pos
func (f *floatArray) HasValue(pos int) bool {
	if !f.checkPos(pos) {
		return false
	}
	blockIdx := pos / blockSize
	idx := pos % blockSize
	mark := f.marks[blockIdx]
	return mark&(1<<uint64(idx)) != 0
}

// GetValue returns value with pos, if has not value return 0
func (f *floatArray) GetValue(pos int) float64 {
	if !f.checkPos(pos) {
		return 0
	}
	return f.values[pos]
}

// SetValue sets value with pos, if pos out of bounds, return it
func (f *floatArray) SetValue(pos int, value float64) {
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
func (f *floatArray) IsEmpty() bool {
	return f.size == 0
}

// Size returns size of array
func (f *floatArray) Size() int {
	return f.size
}

// Iterator returns an iterator over the array
func (f *floatArray) Iterator() FloatArrayIterator {
	if f.it == nil {
		f.it = newFloatArrayIterator(f)
	} else {
		f.it.reset()
	}
	return f.it
}

// Capacity returns the capacity of array
func (f *floatArray) Capacity() int {
	return f.capacity
}

// Marks returns the marks of array
func (f *floatArray) Marks() []uint8 {
	return f.marks
}

// checkPos checks pos if out of bounds
func (f *floatArray) checkPos(pos int) bool {
	if pos < 0 || pos >= f.capacity {
		return false
	}
	return true
}

// Reset resets all values and mark for reusing
func (f *floatArray) Reset() {
	f.size = 0
	f.isSingle = false
	for i := range f.marks {
		f.marks[i] = 0
	}
}

// SetSingle sets is array is single value, mean all values is same
func (f *floatArray) SetSingle(single bool) {
	f.isSingle = single
}

// IsSingle return if is single value
func (f *floatArray) IsSingle() bool {
	return f.isSingle
}

// FloatArrayIterator represents a float array iterator
type FloatArrayIterator interface {
	// HasNext returns if this iterator has more values
	HasNext() bool
	// Next returns the next value and index
	Next() (idx int, value float64)
}

// floatArrayIterator represents a float array iterator
type floatArrayIterator struct {
	fa  FloatArray
	idx int

	count    int
	hasValue bool
	marks    []uint8
	mark     uint8
}

// newFloatArrayIterator creates a float array iterator
func newFloatArrayIterator(fa FloatArray) *floatArrayIterator {
	return &floatArrayIterator{
		fa:       fa,
		hasValue: true,
		marks:    fa.Marks(),
	}
}

func (it *floatArrayIterator) reset() {
	it.idx = 0
	it.count = 0
	it.marks = it.fa.Marks()
	it.hasValue = true
}

// HasNext returns if this iterator has more values
func (it *floatArrayIterator) HasNext() bool {
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
func (it *floatArrayIterator) Next() (idx int, value float64) {
	if !it.hasValue {
		return -1, 0
	}
	idx = it.idx - 1
	value = it.fa.GetValue(idx)
	return idx, value
}
