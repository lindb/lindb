package memdb

import (
	"math"
	"math/bits"
	"sync"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./block_store.go -destination=./block_store_mock_test.go -package memdb

type mergeType uint8

// Defines all value type of primitive field
const (
	appendEmpty mergeType = iota + 1
	appendOld
	appendNew
	merge
)

// the longest length of basic-variable on x64 platform
const maxTimeWindow = 64

// define mergeFunc func for merging block store value and compress value
type mergeFunc func(mergeType mergeType, idx int, oldValue uint64)

// blockStore represents a pool of block for reuse
type blockStore struct {
	timeWindow     int
	intBlockPool   sync.Pool
	floatBlockPool sync.Pool
}

// newBlockStore returns a pool of block with fixed time window
func newBlockStore(timeWindow int) *blockStore {
	tw := timeWindow
	if tw <= 0 {
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
	hasValue(pos int) bool
	// setValue marks pos has value
	setValue(pos int)
	// setIntValue sets int64 value with pos
	setIntValue(pos int, value int64)
	// getIntValue returns int64 value for pos
	getIntValue(pos int) int64
	// setFloatValue sets float64 value with pos
	setFloatValue(pos int, value float64)
	// getFloatValue returns float64 value for pos
	getFloatValue(pos int) float64
	// setStartTime sets start time slot
	setStartTime(startTime int)
	// slotRange returns the block store time slot range
	slotRange() (start, end int)
	// getStartTime returns start time slot
	getStartTime() int
	// getEndTime returns end time slot
	getEndTime() int
	// compact compress block data with agg func for rollup operation
	compact(aggFunc field.AggFunc) (startSlot, endSlot int, err error)
	// reset cleans block data, just reset container mark
	reset()
	// bytes returns compress data for block data
	bytes() []byte
	// memsize returns the memory size in bytes count
	memsize() int
	// scan scans block data, then aggregates the data
	scan(aggFunc field.AggFunc, agg []aggregation.PrimitiveAggregator, memScanCtx *memScanContext)
}

const (
	emptyContainerSize = 8 + // container
		4 + // int
		24 // empty byte
)

// container(bit array) is a mapping from 64 value to uint64 in big-endian,
// it is a temporary data-structure for compressing data.
type container struct {
	container uint64
	startTime int

	compress []byte
}

// hasValue returns whether value is absent or present at pos, if present return true
func (c *container) hasValue(pos int) bool {
	return c.container&(1<<uint64(maxTimeWindow-pos-1)) != 0
}

// setValue marks pos is present
func (c *container) setValue(pos int) {
	c.container |= 1 << uint64(maxTimeWindow-pos-1)
}

// setStartTime sets start time slot
func (c *container) setStartTime(startTime int) {
	c.startTime = startTime
	c.container = 0
}

// getStartTime returns start time slot
func (c *container) getStartTime() int {
	return c.startTime
}

func (c *container) slotRange() (start, end int) {
	start = math.MaxInt32
	if len(c.compress) > 0 {
		oldStart, oldEnd := encoding.DecodeTSDTime(c.compress)
		if oldStart < start {
			start = oldStart
		}
		if oldEnd > end {
			end = oldEnd
		}
	}
	curStart := c.getStartTime()
	curEnd := c.getEndTime()
	if curStart < start {
		start = curStart
	}
	if curEnd > end {
		end = curEnd
	}
	return
}

// getEndTime returns end time slot
func (c *container) getEndTime() int {
	// get trailing zeros for container
	trailing := bits.TrailingZeros64(c.container)
	return c.startTime + (maxTimeWindow - trailing) - 1
}

func (c *container) getIntValue(pos int) int64 {
	// do nothing
	return 0
}
func (c *container) setIntValue(pos int, value int64) {
	// do nothing
}

func (c *container) getFloatValue(pos int) float64 {
	// do nothing
	return 0
}

func (c *container) setFloatValue(pos int, value float64) {
	// do nothing
}

// reset cleans block data, just reset container mark
func (c *container) reset() {
	c.container = 0
	c.compress = c.compress[:0]
}

func (c *container) isEmpty() bool {
	return c.container == 0
}

// bytes returns compress data for block data
func (c *container) bytes() []byte {
	return c.compress
}

// memsize returns the memory size in bytes count
func (c *container) memsize() int {
	return emptyContainerSize + cap(c.compress)
}

// isInRange return slot if in range, yes return true
func isInRange(slot, start, end int) bool {
	return slot >= start && slot <= end
}
