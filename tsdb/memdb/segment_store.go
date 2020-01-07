package memdb

import (
	"fmt"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

//go:generate mockgen -source ./segment_store.go -destination=./segment_store_mock_test.go -package memdb

const (
	emptySimpleFieldStoreSize = 8 + // familyTime
		8 + // aggFunc
		8 // block pointer
)

// sStoreINTF represents segment-store,
// which abstracts a store for storing field data based on family start time
type sStoreINTF interface {
	GetFamilyTime() int64

	SlotRange() (
		startSlot,
		endSlot int,
		err error)

	// FlushFieldTo flushes segment's data to writer
	FlushFieldTo(
		tableFlusher metricsdata.Flusher,
		fieldMeta field.Meta,
	) (
		flushedSize int,
	)

	// WriteInt writes a int value, and returns the written length
	WriteInt(
		pFieldID uint16,
		value int64,
		writeCtx writeContext,
	) int

	// WriteFloat writes a float64 value, and returns the written length
	WriteFloat(
		pFieldID uint16,
		value float64,
		writeCtx writeContext,
	) int

	MemSize() int

	// scan scans segment store data based on query time range
	scan(agg aggregation.SeriesAggregator, memScanCtx *memScanContext)
}

// calcTimeWindow calculates time window's block for storing field data based on slot time and value type.
// return int=>pos(slot in time window), bool=>needRollup(if rollup with old value)
// 1) block is nil, create new block, return newBlock, 0, false
// 2) slot time out of current time window, need compress time window then create new one, return block, 0, false
// 3) in current time window, if has old value return pos, true, else return block, pos, false
func calcTimeWindow(block block, blockStore *blockStore, slotTime int,
	valueType field.ValueType, aggFunc field.AggFunc) (block, int, bool) {
	currentBlock := block

	// block is nil
	if currentBlock == nil {
		currentBlock = blockStore.allocBlock(valueType)
		currentBlock.setStartTime(slotTime)
		return currentBlock, 0, false
	}

	startTime := currentBlock.getStartTime()

	// if current slot time out of current time window, need compress block data, start new time window
	if slotTime < startTime || slotTime >= startTime+blockStore.timeWindow {
		_, _, err := currentBlock.compact(aggFunc)
		if err != nil {
			memDBLogger.Error("compress block data error, data will lost", logger.Error(err))
		} else {
			// reset start time using slot time
			currentBlock.setStartTime(slotTime)
		}
		return currentBlock, 0, false
	}

	// in current time window, do rollup value
	pos := slotTime - startTime
	needRollup := false
	if currentBlock.hasValue(pos) {
		// has old value, need do rollup
		needRollup = true
	}
	return currentBlock, pos, needRollup
}

// singleFieldStore stores single field
type simpleFieldStore struct {
	familyTime int64
	aggFunc    field.AggFunc
	block      block
}

// newSingleFieldStore returns a new segment store for simple field store
func newSimpleFieldStore(familyTime int64, aggFunc field.AggFunc) sStoreINTF {
	return &simpleFieldStore{
		familyTime: familyTime,
		aggFunc:    aggFunc,
	}
}

func (fs *simpleFieldStore) GetFamilyTime() int64 {
	return fs.familyTime
}

func (fs *simpleFieldStore) FlushFieldTo(
	tableFlusher metricsdata.Flusher,
	fieldMeta field.Meta,
) (
	flushedSize int,
) {
	if fs.block == nil {
		return
	}
	if _, _, err := fs.block.compact(fs.aggFunc); err != nil {
		memDBLogger.Error("flush simple segment store data err", logger.Error(err))
		return
	}
	data := fs.block.bytes()
	tableFlusher.FlushPrimitiveField(fieldMeta.Type.GetSchema().GetAllPrimitiveFields()[0], data)

	return fs.MemSize()
}

func (fs *simpleFieldStore) SlotRange() (startSlot, endSlot int, err error) {
	if fs.block == nil {
		err = fmt.Errorf("block is empty")
		return
	}
	startSlot, endSlot = fs.block.slotRange()
	return
}

func (fs *simpleFieldStore) MemSize() int {
	if fs.block == nil {
		return emptySimpleFieldStoreSize
	}
	return emptySimpleFieldStoreSize + fs.block.memsize()
}
