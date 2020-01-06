package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series/field"
)

func TestComplexFieldStore_GetFamilyTime(t *testing.T) {
	store := newComplexFieldStore(10, field.SummaryField)
	assert.Equal(t, int64(10), store.GetFamilyTime())
}

func TestComplexFieldStore_Bytes(t *testing.T) {
	store := newComplexFieldStore(10, field.SummaryField)
	data, _, _, err := store.Bytes(false)
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func TestComplexFieldStore_MemSize(t *testing.T) {
	store := newComplexFieldStore(10, field.SummaryField)
	assert.Equal(t, emptyComplexFieldStoreSize, store.MemSize())
}

func TestComplexFieldStore_SlotRange(t *testing.T) {
	store := newComplexFieldStore(10, field.HistogramField)
	_, _, err := store.SlotRange()
	assert.NoError(t, err)
}
