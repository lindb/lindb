package memdb

import (
	"sync"

	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./block_store.go -destination=./block_store_mock_test.go -package memdb

// the longest length of basic-variable on x64 platform
const maxTimeWindow = 64

// blockStore represents a pool of block for reuse
type blockStore struct {
	timeWindow     uint16
	intBlockPool   sync.Pool
	floatBlockPool sync.Pool
}

// newBlockStore returns a pool of block with fixed time window
func newBlockStore(timeWindow uint16) *blockStore {
	tw := timeWindow
	if tw == 0 || tw > maxTimeWindow {
		tw = maxTimeWindow
	}
	return &blockStore{
		timeWindow: tw,
		intBlockPool: sync.Pool{
			New: func() interface{} {
				return newIntBlock(tw)
			},
		},
		floatBlockPool: sync.Pool{
			New: func() interface{} {
				return newFloatBlock(tw)
			},
		},
	}
}

// freeBlock resets block data and free it, puts it into pool for reusing
func (bs *blockStore) freeBlock(block block) {
	block.reset()
	switch b := block.(type) {
	case *intBlock:
		bs.intBlockPool.Put(b)
	case *floatBlock:
		bs.floatBlockPool.Put(b)
	}
}

func (bs *blockStore) allocBlock(valueType field.ValueType) block {
	switch valueType {
	case field.Integer:
		return bs.allocIntBlock()
	case field.Float:
		return bs.allocFloatBlock()
	default:
		return nil
	}
}

// allocIntBlock alloc int block from pool
func (bs *blockStore) allocIntBlock() *intBlock {
	block := bs.intBlockPool.Get()
	return block.(*intBlock)
}

// allocIntBlock alloc float block from pool
func (bs *blockStore) allocFloatBlock() *floatBlock {
	block := bs.floatBlockPool.Get()
	return block.(*floatBlock)
}

// block represents a fixed size time window of metric data.
// All block implementations need provide fast random access to data.
type block interface {
	// hasValue returns if has value with pos, if has value return true
	hasValue(pos uint16) bool
	// setIntValue sets int64 value with pos
	setIntValue(pos uint16, value int64)
	// getIntValue returns int64 value for pos
	getIntValue(pos uint16) int64
	// setFloatValue sets float64 value with pos
	setFloatValue(pos uint16, value float64)
	// getFloatValue returns float64 value for pos
	getFloatValue(pos uint16) float64
	// getSize returns the size of values
	getSize() uint16
	// reset cleans block data, just reset container mark
	reset()
	// memsize returns the memory size in bytes count
	memsize() int
}

const (
	emptyContainerSize = 8 + // container
		24 // empty byte
)

func (b *floatBlock) getIntValue(pos uint16) int64 {
	// do nothing
	return 0
}
func (b *floatBlock) setIntValue(pos uint16, value int64) {
	// do nothing
}

func (b *intBlock) getFloatValue(pos uint16) float64 {
	// do nothing
	return 0
}

func (b *intBlock) setFloatValue(pos uint16, value float64) {
	// do nothing
}
