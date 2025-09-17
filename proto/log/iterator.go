package log

import (
	"github.com/lindb/common/proto/gen/v1/flatLogV1"
)

type FieldIterator struct {
	l     *flatLogV1.Log
	field flatLogV1.Field

	idx int
	num int
}

func NewFieldIterator(l *flatLogV1.Log) *FieldIterator {
	return &FieldIterator{
		l:   l,
		idx: -1,
		num: l.FieldsLength(),
	}
}

func (it *FieldIterator) HasNext() bool {
	it.idx++
	if it.idx >= it.num {
		return false
	}
	return it.l.Fields(&it.field, it.idx)
}

func (it *FieldIterator) NextName() []byte { return it.field.Name() }

func (it *FieldIterator) NextValue() []byte { return it.field.Value() }

func (it *FieldIterator) Len() int { return it.num }

func (it *FieldIterator) Reset() { it.idx = -1 }
