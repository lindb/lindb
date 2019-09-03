package memdb

import (
	"math"
	"math/bits"
	"sync"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./block_store.go -destination=./block_store_mock_test.go -package memdb

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
	// getStartTime returns start time slot
	getStartTime() int
	// getEndTime returns end time slot
	getEndTime() int
	// compact compress block data with agg func for rollup operation
	compact(aggFunc field.AggFunc, needSlotRange bool) (startSlot, endSlot int, err error)
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

// DecodeTSDTime returns the start/end under compress tsd data
func (c *container) DecodeTSDTime(needSlotRange bool) (startSlot, endSlot int, needCompact bool) {
	if c.container == 0 {
		if needSlotRange && len(c.compress) > 0 {
			startSlot, endSlot = encoding.DecodeTSDTime(c.compress)
		}
		return
	}
	// block has value, need compact value
	needCompact = true
	return
}

// merge merges values and compress data of container based on value type nad agg func
func (c *container) merge(valueType field.ValueType,
	values []uint64, aggFunc field.AggFunc) (start, end int, err error) {
	merger := newMerger(c, valueType, values, c.compress, aggFunc)

	buf, err := merger.merge()
	if err != nil {
		return merger.startTime, merger.endTime, err
	}

	c.compress = buf
	c.container = 0
	return merger.startTime, merger.endTime, nil
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

// setIntValue updates int64 value with pos
func (b *intBlock) setIntValue(pos int, value int64) {
	b.setValue(pos)
	b.values[pos] = value
}

// getIntValue return int64 value for pos
func (b *intBlock) getIntValue(pos int) int64 {
	return b.values[pos]
}

// compact compress block data
func (b *intBlock) compact(aggFunc field.AggFunc, needSlotRange bool) (startSlot, endSlot int, err error) {
	needCompact := false
	startSlot, endSlot, needCompact = b.DecodeTSDTime(needSlotRange)
	if !needCompact {
		return
	}

	length := len(b.values)
	values := make([]uint64, length)
	for i := 0; i < length; i++ {
		values[i] = encoding.ZigZagEncode(b.values[i])
	}
	startSlot, endSlot, err = b.merge(field.Integer, values, aggFunc)
	return
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

// setFloatValue updates float64 value with pos
func (b *floatBlock) setFloatValue(pos int, value float64) {
	b.setValue(pos)
	b.values[pos] = value
}

// getFloatValue returns float64 value for pos
func (b *floatBlock) getFloatValue(pos int) float64 {
	return b.values[pos]
}

// compact compress block data
func (b *floatBlock) compact(aggFunc field.AggFunc, needSlotRange bool) (startSlot, endSlot int, err error) {
	needCompact := false

	startSlot, endSlot, needCompact = b.DecodeTSDTime(needSlotRange)
	if !needCompact {
		return
	}
	length := len(b.values)
	values := make([]uint64, length)
	for i := 0; i < length; i++ {
		values[i] = math.Float64bits(b.values[i])
	}
	startSlot, endSlot, err = b.merge(field.Float, values, aggFunc)
	return
}

// merger is merge operation which provides compress block data.
// 1) compress data not exist, just compress current block values
// 2) compress data exist, merge compress data and block values
type merger struct {
	valueType    field.ValueType
	block        *container
	values       []uint64
	compressData []byte
	tsd          *encoding.TSDEncoder

	oldData *encoding.TSDDecoder

	startTime int
	endTime   int

	aggFunc field.AggFunc
}

// newMerger creates merge operation with given agg func based on block data and exist compress data
func newMerger(block *container, valueType field.ValueType,
	values []uint64, compressData []byte,
	aggFunc field.AggFunc) *merger {
	m := &merger{
		valueType:    valueType,
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

// merge merges block's values and compress data based on agg func and value type
func (m *merger) merge() ([]byte, error) {
	if m.oldData == nil {
		// compress data not exist, just compress block data
		m.compress()
	} else {
		// has old compress data, need merge block data
		curStartTime := m.block.getStartTime()
		curEndTime := m.block.getEndTime()
		oldStartTime := m.oldData.StartTime()
		oldEndTime := m.oldData.EndTime()
		//TODO add check start/end range????

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

	buf, err := m.tsd.Bytes()
	if err != nil {
		return nil, err
	}
	return buf, nil
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
		switch m.valueType {
		case field.Integer:
			val := m.aggFunc.AggregateInt(encoding.ZigZagDecode(m.values[newPos]), encoding.ZigZagDecode(m.oldData.Value()))
			m.appendValue(encoding.ZigZagEncode(val))
		case field.Float:
			val := m.aggFunc.AggregateFloat(math.Float64frombits(m.values[newPos]), math.Float64frombits(m.oldData.Value()))
			m.appendValue(math.Float64bits(val))
		}
	case hasValue:
		// append current block block
		m.appendValue(m.values[newPos])
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
		m.appendValue(m.values[pos])
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
