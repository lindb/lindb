package memdb

import (
	"fmt"
	"math"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

const (
	emptyComplexFieldStoreSize = 8 + // familyTime
		8 + // aggFunc
		8 // block pointer
)

type complexFieldStore struct {
	familyTime int64
	schema     field.Schema
	blocks     map[uint16]block
}

func newComplexFieldStore(familyTime int64, fieldType field.Type) sStoreINTF {
	return &complexFieldStore{
		familyTime: familyTime,
		schema:     fieldType.GetSchema(),
		blocks:     make(map[uint16]block),
	}
}

func (fs *complexFieldStore) GetFamilyTime() int64 {
	return fs.familyTime
}

func (fs *complexFieldStore) SlotRange() (
	startSlot,
	endSlot int,
	err error) {
	if len(fs.blocks) == 0 {
		err = fmt.Errorf("block is empty")
		return
	}
	startSlot = math.MaxInt32
	endSlot = -1
	for _, block := range fs.blocks {
		start, end := block.slotRange()
		if start < startSlot {
			startSlot = start
		}
		if end > endSlot {
			endSlot = end
		}
	}
	return
}

func (fs *complexFieldStore) FlushFieldTo(
	tableFlusher metricsdata.Flusher,
	fieldMeta field.Meta,
) (
	flushedSize int,
) {
	if len(fs.blocks) == 0 {
		return
	}
	schema := fieldMeta.Type.GetSchema()

	for fieldID, block := range fs.blocks {
		aggFunc := schema.GetAggFunc(fieldID)
		if _, _, err := block.compact(aggFunc); err != nil {
			memDBLogger.Error("flush complex segment store data err", logger.Error(err))
			return
		}
		data := block.bytes()
		tableFlusher.FlushPrimitiveField(fieldID, data)
	}

	return fs.MemSize()
}

func (fs *complexFieldStore) MemSize() int {
	if len(fs.blocks) == 0 {
		return emptyComplexFieldStoreSize
	}
	size := emptyComplexFieldStoreSize
	for _, block := range fs.blocks {
		size += 2 + // pField id size
			block.memsize() // block mem size
	}
	return size
}
