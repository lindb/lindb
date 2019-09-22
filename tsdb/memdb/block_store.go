package memdb

import (
	"math"
	"math/bits"
	"sync"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/bit"
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

// define aggregator func for scanning block store data
type aggregator func(mergeType mergeType, idx int, oldValue uint64) (completed bool)

// define getValue func for getting block store data
type getValue func(idx int) (value uint64)

// define mergeFunc func for merging block store value and compress value
type mergeFunc func(mergeType mergeType, idx int, oldValue uint64)

// define aggFunc for aggregating block store value and compress value
type aggFunc func(idx int, oldValue uint64) (value uint64)

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

// scan scans the block store data and compress data
func (c *container) scan(memScanCtx *memScanContext, aggregator aggregator) {
	hasOld := len(c.compress) > 0
	hasNew := c.container != 0
	switch {
	case !hasOld && hasNew: // scans current block store buffer data
		end := c.getEndTime() - c.startTime
		for i := 0; i <= end; i++ {
			if !c.hasValue(i) {
				continue
			}
			if aggregator(appendNew, i, 0) {
				return
			}
		}
	case hasOld && hasNew: // scans current buffer data and compress data, then merges them for same time slot
		tsd := memScanCtx.tsd
		tsd.Reset(c.compress)
		scanner := newMergeScanner(c, tsd)
		scanner.mergeFunc = func(mergeType mergeType, pos int, oldValue uint64) {
			if aggregator(mergeType, pos, oldValue) {
				scanner.complete = true
			}
		}
		scanner.scan()
	case hasOld: // scans compress data
		tsd := memScanCtx.tsd
		tsd.Reset(c.compress)
		for tsd.Error() == nil && tsd.Next() {
			if tsd.HasValue() {
				timeSlot := tsd.Slot()
				val := tsd.Value()
				if aggregator(appendOld, timeSlot, val) {
					return
				}
			}
		}
	}
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

// merge merges values and compress data of container based on value type nad agg func
func (c *container) merge(getValue getValue, aggFunc aggFunc) (start, end int, err error) {
	hasOld := len(c.compress) > 0
	hasNew := c.container != 0
	var encode *encoding.TSDEncoder
	switch {
	case !hasOld && !hasNew: // no data
		return 0, 0, nil
	case !hasOld: // compact current buffer data
		end = c.getEndTime()
		start = c.startTime
		encode = encoding.NewTSDEncoder(start)
		for i := start; i <= end; i++ {
			idx := i - start
			if c.hasValue(idx) {
				encode.AppendTime(bit.One)
				encode.AppendValue(getValue(idx))
			} else {
				encode.AppendTime(bit.Zero)
			}
		}
	case hasOld && !hasNew: // just decode time slot range for compress data
		start, end = encoding.DecodeTSDTime(c.compress)
		return
	default: // merge current buffer data and compress data
		tsd := encoding.GetTSDDecoder()

		tsd.Reset(c.compress)
		scanner := newMergeScanner(c, tsd)
		encode = encoding.NewTSDEncoder(scanner.start)
		scanner.mergeFunc = func(mergeType mergeType, idx int, oldValue uint64) {
			switch mergeType {
			case appendEmpty:
				encode.AppendTime(bit.Zero)
			case appendNew:
				encode.AppendTime(bit.One)
				encode.AppendValue(getValue(idx))
			case appendOld:
				encode.AppendTime(bit.One)
				encode.AppendValue(oldValue)
			case mergeType:
				encode.AppendTime(bit.One)
				encode.AppendValue(aggFunc(idx, oldValue))
			}
		}
		scanner.scan()
		encoding.ReleaseTSDDecoder(tsd)
		start = scanner.start
		end = scanner.end
	}
	// reset compress data and clear current buffer
	if encode != nil {
		data, err := encode.Bytes()
		if err != nil {
			return 0, 0, err
		}
		c.compress = data
		c.container = 0
	}
	return start, end, err
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
func (b *intBlock) compact(aggFunc field.AggFunc) (startSlot, endSlot int, err error) {
	return b.container.merge(func(idx int) (value uint64) {
		value = encoding.ZigZagEncode(b.values[idx])
		return
	}, func(idx int, oldValue uint64) (value uint64) {
		return encoding.ZigZagEncode(aggFunc.AggregateInt(b.values[idx], encoding.ZigZagDecode(oldValue)))
	})
}

// memsize returns the memory size in bytes count
func (b *intBlock) memsize() int {
	return b.container.memsize() + 24 + cap(b.values)*8
}

// scan scans block data, then aggregates the data
func (b *intBlock) scan(aggFunc field.AggFunc, agg []aggregation.PrimitiveAggregator, memScanCtx *memScanContext) {
	b.container.scan(memScanCtx, func(mergeType mergeType, idx int, oldValue uint64) (completed bool) {
		value := 0.0
		// 1. get value and time slot
		switch mergeType {
		case appendOld:
			value = float64(encoding.ZigZagDecode(oldValue))
		case appendNew:
			value = float64(b.values[idx])
			idx += b.startTime
		case merge:
			value = float64(aggFunc.AggregateInt(b.values[idx], encoding.ZigZagDecode(oldValue)))
			idx += b.startTime
		default:
			return
		}
		// 2. aggregate the value based on time slot
		for _, a := range agg {
			completed = a.Aggregate(idx, value)
		}
		return
	})
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
func (b *floatBlock) compact(aggFunc field.AggFunc) (startSlot, endSlot int, err error) {
	return b.container.merge(func(idx int) (value uint64) {
		value = math.Float64bits(b.values[idx])
		return
	}, func(idx int, oldValue uint64) (value uint64) {
		return math.Float64bits(aggFunc.AggregateFloat(b.values[idx], math.Float64frombits(oldValue)))
	})
}

// memsize returns the memory size in bytes count
func (b *floatBlock) memsize() int {
	return b.container.memsize() + 24 + cap(b.values)*8
}

// scan scans block data, then aggregates the data
func (b *floatBlock) scan(aggFunc field.AggFunc, agg []aggregation.PrimitiveAggregator, memScanCtx *memScanContext) {
	b.container.scan(memScanCtx, func(mergeType mergeType, idx int, oldValue uint64) (completed bool) {
		value := 0.0
		// 1. get value and time slot
		switch mergeType {
		case appendOld:
			value = math.Float64frombits(oldValue)
		case appendNew:
			value = b.values[idx]
			idx += b.startTime
		case merge:
			value = aggFunc.AggregateFloat(b.values[idx], math.Float64frombits(oldValue))
			idx += b.startTime
		default:
			return
		}
		// 2. aggregate the value based on time slot
		for _, a := range agg {
			completed = a.Aggregate(idx, value)
		}
		return
	})
}

// mergeScanner represents the scanner which scans the block store current buffer data and compress data
type mergeScanner struct {
	container        *container           // current buffer
	tsd              *encoding.TSDDecoder // old value
	start, end       int                  // target time slot range
	curStart, curEnd int                  // current buffer time slot range
	oldStart, oldEnd int                  // compress data time slot range

	complete  bool
	mergeFunc mergeFunc
}

// newMergeScanner creates a merge scanner
func newMergeScanner(container *container, tsd *encoding.TSDDecoder) *mergeScanner {
	scanner := &mergeScanner{
		container: container,
		tsd:       tsd,
	}
	// init scanner time slot ranges
	scanner.init()
	return scanner
}

// init initializes the scanner's time slot ranges
func (s *mergeScanner) init() {
	// start time slot
	s.curStart = s.container.startTime
	s.oldStart = s.tsd.StartTime()
	s.start = s.curStart
	if s.start > s.oldStart {
		s.start = s.oldStart
	}
	// end time slot
	s.curEnd = s.container.getEndTime()
	s.oldEnd = s.tsd.EndTime()
	s.end = s.curEnd
	if s.end < s.oldEnd {
		s.end = s.oldEnd
	}
}

// scan scans the block store current buffer data and compress data based on target time slot range
func (s *mergeScanner) scan() {
	for i := s.start; i <= s.end; i++ {
		// if scanner is completed, return it
		if s.complete {
			return
		}
		inCurrentRange := isInRange(i, s.curStart, s.curEnd)
		inOldRange := isInRange(i, s.oldStart, s.oldEnd)
		newSlot := i - s.curStart
		oldSlot := i - s.oldStart
		hasValue := s.container.hasValue(newSlot)
		hasOldValue := s.tsd.HasValueWithSlot(oldSlot)
		switch {
		case inCurrentRange && inOldRange:
			s.merge(hasValue, hasOldValue, newSlot)
		case inCurrentRange && hasValue:
			// just compress current block value with pos
			s.mergeFunc(appendNew, newSlot, 0)
		case inCurrentRange && !hasValue:
			s.mergeFunc(appendEmpty, newSlot, 0)
		case inOldRange && hasOldValue:
			// read compress data and compress it again with new pos
			s.mergeFunc(appendOld, i, s.tsd.Value())
		case inOldRange && !hasOldValue:
			s.mergeFunc(appendEmpty, i, 0)
		default:
			s.mergeFunc(appendEmpty, i, 0)
		}
	}
}

func (s *mergeScanner) merge(hasValue bool, hasOldValue bool, newSlot int) {
	// merge current block value and value in compress data with pos
	switch {
	case hasValue && hasOldValue:
		// has value both in current and old, do rollup operation with agg func
		s.mergeFunc(merge, newSlot, s.tsd.Value())
	case hasValue:
		// append current block block
		s.mergeFunc(appendNew, newSlot, 0)
	case hasOldValue:
		// read old compress value then append value with new pos
		s.mergeFunc(appendOld, newSlot, s.tsd.Value())
	default:
		// just append empty value with pos
		s.mergeFunc(appendEmpty, newSlot, 0)
	}
}

// isInRange return slot if in range, yes return true
func isInRange(slot, start, end int) bool {
	return slot >= start && slot <= end
}
