package memdb

import (
	"fmt"

	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/logger"
)

// segmentStore stores field data based on family start time
type segmentStore interface {
	bytes() (data []byte, startSlot, endSlot int, err error)
	writeInt(blockStore *blockStore, slotTime int, value int64)
	writeFloat(blockStore *blockStore, slotTime int, value float64)
}

// singleFieldStore stores single field
type simpleFieldStore struct {
	block   block
	aggFunc field.AggFunc
}

// newSingleFieldStore returns a new segment store for simple field store
func newSimpleFieldStore(aggFunc field.AggFunc) segmentStore {
	return &simpleFieldStore{
		aggFunc: aggFunc,
	}
}

func (fs *simpleFieldStore) AggFunc() field.AggFunc {
	//TODO using type????
	return fs.aggFunc
}
func (fs *simpleFieldStore) writeFloat(blockStore *blockStore, slotTime int, value float64) {
	pos, hasValue := fs.calcTimeWindow(blockStore, slotTime, field.Float)
	currentBlock := fs.block
	if hasValue {
		// do rollup using agg func
		currentBlock.setFloatValue(pos, fs.aggFunc.AggregateFloat(currentBlock.getFloatValue(pos), value))
	} else {
		currentBlock.setFloatValue(pos, value)
	}
}

func (fs *simpleFieldStore) writeInt(blockStore *blockStore, slotTime int, value int64) {
	pos, hasValue := fs.calcTimeWindow(blockStore, slotTime, field.Integer)
	currentBlock := fs.block
	if hasValue {
		// do rollup using agg func
		currentBlock.setIntValue(pos, fs.aggFunc.AggregateInt(currentBlock.getIntValue(pos), value))
	} else {
		currentBlock.setIntValue(pos, value)
	}
}

// calcTimeWindow calculates time window's block for storing field data based on slot time and value type.
// return int=>pos(slot in time window), bool=>needRollup(if rollup with old value)
// 1) block is nil, create new block, return 0, false
// 2) slot time out of current time window, need compress time window then create new one, return 0, false
// 3) in current time window, if has old value return pos, true, else return pos, false
func (fs *simpleFieldStore) calcTimeWindow(blockStore *blockStore, slotTime int,
	valueType field.ValueType) (int, bool) {
	currentBlock := fs.block

	// block is nil
	if currentBlock == nil {
		currentBlock = blockStore.allocBlock(valueType)
		currentBlock.setStartTime(slotTime)
		fs.block = currentBlock

		return 0, false
	}

	startTime := currentBlock.getStartTime()

	// if current slot time out of current time window, need compress block data, start new time window
	if slotTime < startTime || slotTime >= startTime+blockStore.timeWindow {
		_, _, err := currentBlock.compact(fs.aggFunc)
		if err != nil {
			memDBLogger.Error("compress block data error, data will lost", logger.Error(err))
		} else {
			// reset start time using slot time
			currentBlock.setStartTime(slotTime)
		}
		return 0, false
	}

	// in current time window, do rollup value
	pos := slotTime - startTime
	needRollup := false
	if currentBlock.hasValue(pos) {
		// has old value, need do rollup
		needRollup = true
	}
	return pos, needRollup
}

func (fs *simpleFieldStore) bytes() (data []byte, startSlot, endSlot int, err error) {
	if fs.block == nil {
		err = fmt.Errorf("block is empty")
		return
	}
	if startSlot, endSlot, err = fs.block.compact(fs.aggFunc); err != nil {
		err = fmt.Errorf("compact block data in simple field store error:%s", err)
		return
	}
	data = fs.block.bytes()
	return
}
