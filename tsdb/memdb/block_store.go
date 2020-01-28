package memdb

import (
	"math/bits"
	"sync"
)

// the longest length of basic-variable on x64 platform
const maxTimeWindow = 64

const (
	emptyContainerSize = 8 + // container
		24 // empty byte
)

// blockStore represents a pool of block for reuse
type blockStore struct {
	timeWindow uint16
	pool       sync.Pool
}

// newBlockStore returns a pool of block with fixed time window
func newBlockStore(timeWindow uint16) *blockStore {
	tw := timeWindow
	if tw == 0 || tw > maxTimeWindow {
		tw = maxTimeWindow
	}
	return &blockStore{
		timeWindow: tw,
		pool: sync.Pool{
			New: func() interface{} {
				return newBlock(tw)
			},
		},
	}
}

// freeBlock resets block data and free it, puts it into pool for reusing
func (bs *blockStore) freeBlock(block *block) {
	block.reset()
	bs.pool.Put(block)
}

// allocBlock alloc block from pool
func (bs *blockStore) allocBlock() *block {
	b := bs.pool.Get()
	return b.(*block)
}

// block represents a fixed size time window of metric data.
// All block implementations need provide fast random access to data.
type block struct {
	// container(bit array) is a mapping from 64 value to uint64 in big-endian,
	// it is a temporary data-structure for compressing data.
	container uint64
	values    []float64
}

// newBlock returns block with fixed time window
func newBlock(size uint16) *block {
	return &block{
		values: make([]float64, size),
	}
}

// getSize returns the size of values
func (b *block) getSize() uint16 {
	if b.container == 0 {
		return 0
	}
	// get trailing zeros for container
	trailing := bits.TrailingZeros64(b.container)
	return uint16(maxTimeWindow - trailing - 1)
}

// hasValue returns whether value is absent or present at pos, if present return true
func (b *block) hasValue(pos uint16) bool {
	return b.container&(1<<uint64(maxTimeWindow-pos-1)) != 0
}

// setValue updates value with pos
func (b *block) setValue(pos uint16, value float64) {
	b.container |= 1 << uint64(maxTimeWindow-pos-1)
	b.values[pos] = value
}

// getValue returns value for pos
func (b *block) getValue(pos uint16) float64 {
	return b.values[pos]
}

// memsize returns the memory size in bytes count
func (b *block) memsize() int {
	return emptyContainerSize + cap(b.values)*8
}

// reset cleans block data, just reset container mark
func (b *block) reset() {
	b.container = 0
}

// isEmpty returns the block if empty
func (b *block) isEmpty() bool {
	return b.container == 0
}
