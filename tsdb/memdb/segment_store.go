package memdb

import (
	"fmt"

	"github.com/eleme/lindb/pkg/encoding"
	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/logger"
)

// segmentStore stores field data based on family start time
type segmentStore interface {
	bytes() (data []byte, startSlot, endSlot int, err error)
	writeInt(blockStore *blockStore, slotTime int, value int64)
}

// singleFieldStore stores single field
type simpleFieldStore struct {
	block   *intBlock
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

func (fs *simpleFieldStore) writeInt(blockStore *blockStore, slotTime int, value int64) {
	currentBlock := fs.block
	if currentBlock == nil {
		currentBlock = blockStore.allocIntBlock()
		currentBlock.setStartTime(slotTime)
		currentBlock.setValue(0)
		currentBlock.updateValue(0, value)
		fs.block = currentBlock
	} else {
		startTime := currentBlock.getStartTime()
		if slotTime < startTime || slotTime >= startTime+blockStore.timeWindow {
			// if current slot time out of current time window, need compress block data
			err := currentBlock.compact(fs.aggFunc)
			if err != nil {
				memDBLogger.Error("compress block data error, data will lost", logger.Error(err))
			} else {
				currentBlock.setStartTime(slotTime) // reset start time using slot time
				currentBlock.setValue(0)
				currentBlock.updateValue(0, value)
			}
		} else {
			// in current time window, do rollup value
			var pos = slotTime - startTime
			if currentBlock.hasValue(pos) {
				// do rollup using agg func
				currentBlock.updateValue(pos, fs.aggFunc.AggregateInt(currentBlock.getValue(pos), value))
			} else {
				currentBlock.setValue(pos)
				currentBlock.updateValue(pos, value)
			}
		}
	}
}
func (fs *simpleFieldStore) bytes() (data []byte, startSlot, endSlot int, err error) {
	if fs.block == nil {
		err = fmt.Errorf("block is empty")
		return
	}
	if err = fs.block.compact(fs.aggFunc); err != nil {
		err = fmt.Errorf("compact block data in simple field store error:%s", err)
		return
	}
	data = fs.block.bytes()
	startSlot, endSlot = encoding.DecodeTSDTime(data)
	return
}
