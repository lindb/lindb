package collections

const blockSize = 8

// FloatArray represents a float array, support mark pos if has value
type FloatArray struct {
	marks  []uint8
	values []float64

	capacity int
	size     int
}

// NewFloatArray creates a float array
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

// GetValue returns value with pos, if has not value return 0
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
		idx := pos % blockSize
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

// Iterator returns an iterator over the array
func (f *FloatArray) Iterator() *FloatArrayIterator {
	return newFloatArrayIterator(f)
}

// checkPos checks pos if out of bounds
func (f *FloatArray) checkPos(pos int) bool {
	if pos < 0 || pos >= f.capacity {
		return false
	}
	return true
}

// FloatArrayIterator represents a float array iterator
type FloatArrayIterator struct {
	fa  *FloatArray
	idx int

	count    int
	hasValue bool
	mark     uint8
}

// newFloatArrayIterator creates a float array iterator
func newFloatArrayIterator(fa *FloatArray) *FloatArrayIterator {
	return &FloatArrayIterator{
		fa:       fa,
		hasValue: true,
	}
}

// HasNext returns if this iterator has more values
func (it *FloatArrayIterator) HasNext() bool {
	for it.idx < it.fa.capacity && it.count < it.fa.Size() {
		blockIdx := it.idx / blockSize
		idx := it.idx % blockSize
		if idx == 0 {
			it.mark = it.fa.marks[blockIdx]
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
