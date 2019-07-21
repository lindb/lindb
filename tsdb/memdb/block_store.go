package memdb

import (
	"math/bits"
	"sync"

	"github.com/eleme/lindb/pkg/bit"
	"github.com/eleme/lindb/pkg/encoding"
	"github.com/eleme/lindb/pkg/field"
)

// the longest length of basic-variable on x64 platform
const maxTimeWindow = 64

// blockStore represents a pool of block for reuse
type blockStore struct {
	timeWindow     int
	intBlockPool   sync.Pool
	floatBlockPool sync.Pool
}

// newBlockStore returns a pool of block with fixed time window
func newBlockStore(timeWindow int) *blockStore {
	return &blockStore{
		timeWindow: timeWindow,
		intBlockPool: sync.Pool{
			New: func() interface{} {
				return newIntBlock(timeWindow)
			},
		},
		floatBlockPool: sync.Pool{
			New: func() interface{} {
				return newFloatBlock(timeWindow)
			},
		},
	}
}

// freeIntBlock resets int block data and free it, puts it into pool for reusing
func (bs *blockStore) freeIntBlock(block *intBlock) {
	block.reset()
	bs.intBlockPool.Put(block)
}

// freeFloatBlock resets float block data and free it, puts it into pool for reusing
func (bs *blockStore) freeFloatBlock(block *floatBlock) {
	block.reset()
	bs.floatBlockPool.Put(block)
}

// allocIntBlock alloc int block from pool
func (bs *blockStore) allocIntBlock() *intBlock {
	block := bs.intBlockPool.Get()
	b, ok := block.(*intBlock)
	if ok {
		return b
	}
	return nil
}

// allocIntBlock alloc float block from pool
func (bs *blockStore) allocFloatBlock() *floatBlock {
	block := bs.floatBlockPool.Get()
	b, ok := block.(*floatBlock)
	if ok {
		return b
	}
	return nil
}

// block represents a fixed size time window of metric data.
// All block implementations need provide fast random access to data.
type block interface {
	// hasValue returns if has value with pos, if has value return true
	hasValue(pos int) bool
	// setValue marks pos has value
	setValue(pos int)
	// setStartTime sets start time slot
	setStartTime(startTime int)
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
}

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

// getEndTime returns end time slot
func (c *container) getEndTime() int {
	// get trailing zeros for container
	trailing := bits.TrailingZeros64(c.container)
	return c.startTime + (maxTimeWindow - trailing) - 1
}

// reset cleans block data, just reset container mark
func (c *container) reset() {
	c.container = 0
	c.compress = c.compress[:0]
}

// intBlock represents a int block for storing metric point in memory
type intBlock struct {
	container
	values []int64
}

// newIntBlock returns int block with fixed time window
func newIntBlock(size int) *intBlock {
	return &intBlock{
		values: make([]int64, size),
	}
}

// getValue returns value with pos
func (b *intBlock) getValue(pos int) int64 {
	return b.values[pos]
}

// updateValue updates value with pos
func (b *intBlock) updateValue(pos int, value int64) {
	b.values[pos] = value
}

// compact compress block data
func (b *intBlock) compact(aggFunc field.AggFunc) (startSlot, endSlot int, err error) {
	//TODO handle error
	merger := newMerger(b, b.values, b.compress, aggFunc)
	// do merge logic
	merger.merge()

	var buf []byte
	buf, err = merger.tsd.Bytes()
	if err != nil {
		return
	}

	b.compress = buf
	return merger.startTime, merger.endTime, nil
}

// bytes returns compress data for block data
func (b *intBlock) bytes() []byte {
	return b.compress
}

// floatBlock represents a float block for storing metric point in memory
type floatBlock struct {
	container
	values []float64
}

// newFloatBlock returns float block with fixed time window
func newFloatBlock(size int) *floatBlock {
	return &floatBlock{
		values: make([]float64, size),
	}
}

// getValue returns value with pos
func (b *floatBlock) getValue(pos int) float64 {
	return b.values[pos]
}

// updateValue updates value with pos
func (b *floatBlock) updateValue(pos int, value float64) {
	b.values[pos] = value
}

// compact compress block data
func (b *floatBlock) compact(aggFunc field.AggFunc) (startSlot, endSlot int, err error) {
	//TODO need implement
	return
}

// merger is merge operation which provides compress block data.
// 1) compress data not exist, just compress current block values
// 2) compress data exist, merge compress data and block values
type merger struct {
	block        block
	values       []int64
	compressData []byte
	tsd          *encoding.TSDEncoder

	oldData *encoding.TSDDecoder

	startTime int
	endTime   int

	aggFunc field.AggFunc
}

// newMerger creates merge operation with given agg func based on block data and exist compress data
func newMerger(block block, values []int64, compressData []byte, aggFunc field.AggFunc) *merger {
	m := &merger{
		block:        block,
		values:       values,
		compressData: compressData,
		aggFunc:      aggFunc,
	}
	m.init()
	return m
}

// init initializes merge context, such time range, tsd decoder if has compress data
func (m *merger) init() {
	curStartTime := m.block.getStartTime()
	curEndTime := m.block.getEndTime()
	if len(m.compressData) == 0 {
		m.startTime = curStartTime
		m.endTime = curEndTime
	} else {
		m.oldData = encoding.NewTSDDecoder(m.compressData)
		// calc compress time window range
		oldStartTime := m.oldData.StartTime()
		oldEndTime := m.oldData.EndTime()

		m.startTime = curStartTime
		if m.startTime > oldStartTime {
			m.startTime = oldStartTime
		}
		m.endTime = curEndTime
		if m.endTime < oldEndTime {
			m.endTime = oldEndTime
		}
	}

	// build tsd encoder
	m.tsd = encoding.NewTSDEncoder(m.startTime)
}

// merge does merge logic
func (m *merger) merge() {
	if m.oldData == nil {
		// compress data not exist, just compress block data
		m.compress()
	} else {
		// has old compress data, need merge block data
		curStartTime := m.block.getStartTime()
		curEndTime := m.block.getEndTime()
		oldStartTime := m.oldData.StartTime()
		oldEndTime := m.oldData.EndTime()

		// do merge and compress data
		for i := m.startTime; i <= m.endTime; i++ {
			newPos := i - curStartTime
			oldPos := i - oldStartTime

			inCurrentRange := m.isInRange(i, curStartTime, curEndTime)
			inOldRange := m.isInRange(i, oldStartTime, oldEndTime)
			switch {
			case inCurrentRange && inOldRange:
				// merge current block value and value in compress data with pos
				m.mergeData(newPos, oldPos)
			case inCurrentRange:
				// just compress current block value with pos
				m.appendNewData(newPos)
			case inOldRange:
				// read compress data and compress it again with new pos
				m.appendOldData(oldPos)
			default:
				m.appendEmptyValue()
			}
		}
	}
}

// isInRange return slot if in range, yes return true
func (m *merger) isInRange(slot, start, end int) bool {
	return slot >= start && slot <= end
}

// mergeData merges current block values and compress data
func (m *merger) mergeData(newPos, oldPos int) {
	b := m.block
	hasValue := b.hasValue(newPos)
	hasOldValue := m.oldData.HasValueWithSlot(oldPos)
	switch {
	case hasValue && hasOldValue:
		// has value both in current and old, do rollup operation with agg func
		val := m.aggFunc.AggregateInt(m.values[newPos], encoding.ZigZagDecode(m.oldData.Value()))
		m.appendValue(encoding.ZigZagEncode(val))
	case hasValue:
		// append current block block
		m.appendValue(encoding.ZigZagEncode(m.values[newPos]))
	case hasOldValue:
		// read old compress value then append value with new pos
		m.appendValue(m.oldData.Value())
	default:
		// just append empty value with pos
		m.appendEmptyValue()
	}
}

// compress compress current block values
func (m *merger) compress() {
	for i := m.startTime; i <= m.endTime; i++ {
		m.appendNewData(i - m.startTime)
	}
}

// appendNewData appends current block value with pos
func (m *merger) appendNewData(pos int) {
	if m.block.hasValue(pos) {
		m.appendValue(encoding.ZigZagEncode(m.values[pos]))
	} else {
		m.appendEmptyValue()
	}
}

// appendOldData reads compress data then appends it with new pos
func (m *merger) appendOldData(pos int) {
	if m.oldData.HasValueWithSlot(pos) {
		m.appendValue(m.oldData.Value())
	} else {
		m.appendEmptyValue()
	}
}

// appendValue appends value with new pos
func (m *merger) appendValue(val uint64) {
	m.tsd.AppendTime(bit.One)
	m.tsd.AppendValue(val)
}

// appendEmptyValue appends time slot only
func (m *merger) appendEmptyValue() {
	m.tsd.AppendTime(bit.Zero)
}
