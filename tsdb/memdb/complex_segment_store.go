package memdb

import (
	"github.com/lindb/lindb/series/field"
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
	return
}

func (fs *complexFieldStore) Bytes(
	needSlotRange bool,
) (
	data []byte,
	startSlot,
	endSlot int,
	err error) {
	return
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
