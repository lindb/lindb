package memdb

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/logger"
)

// segmentStore stores field data based on famliy start time
type segmentStore interface {
	bytes() ([]byte, error)
}

// singleFieldStore stores single field
type simpleFieldStore struct {
	familyStartTime int64
	block           *intBlock
	aggFunc         field.AggFunc
}

// newSingleFieldStore returns a new segment store for simple field store
func newSimpleFieldStore(familyStartTime int64, aggFunc field.AggFunc) segmentStore {
	return &simpleFieldStore{
		familyStartTime: familyStartTime,
		aggFunc:         aggFunc,
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
				logger.GetLogger().Error("compress block data error, data will lost", zap.Error(err))
			} else {
				currentBlock.setStartTime(slotTime) // reset start time using slot time
				currentBlock.setValue(0)
				currentBlock.updateValue(0, value)
			}
		} else {
			// in current time window, do rollup value
			var pos = slotTime - startTime
			if currentBlock.hasValue(pos) {
				// do rullup using agg func
				currentBlock.updateValue(pos, fs.aggFunc.AggregateInt(currentBlock.getValue(pos), value))
			} else {
				currentBlock.setValue(pos)
				currentBlock.updateValue(pos, value)
			}
		}
	}
}
func (fs *simpleFieldStore) bytes() ([]byte, error) {
	if fs.block == nil {
		return nil, fmt.Errorf("block is empty")
	}
	if err := fs.block.compact(fs.aggFunc); err != nil {
		return nil, fmt.Errorf("compact block data in simple field store error:%s", err)
	}
	return fs.block.bytes(), nil
}
