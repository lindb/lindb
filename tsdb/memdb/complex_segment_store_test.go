package memdb

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestComplexFieldStore_GetFamilyTime(t *testing.T) {
	store := newComplexFieldStore(10, field.HistogramField)
	assert.Equal(t, int64(10), store.GetFamilyTime())
}

func TestComplexFieldStore_FlushFieldTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// test store is empty
	flusher := metricsdata.NewMockFlusher(ctrl)
	store := newComplexFieldStore(10, field.SummaryField)
	flushSize := store.FlushFieldTo(flusher, field.Meta{
		ID:   10,
		Type: field.SummaryField,
		Name: "f1",
	})
	assert.True(t, flushSize == 0)

	// test normal case
	writeCtx := writeContext{
		blockStore:   newBlockStore(30),
		timeInterval: 10,
		metricID:     1,
		slotIndex:    1,
		familyTime:   0,
	}
	store.WriteInt(uint16(1), int64(10), writeCtx)
	flusher.EXPECT().FlushPrimitiveField(gomock.Any(), gomock.Any())
	flushSize = store.FlushFieldTo(flusher, field.Meta{
		ID:   10,
		Type: field.SummaryField,
		Name: "f1",
	})
	assert.True(t, flushSize > 0)

	// test block compact err
	s := store.(*complexFieldStore)
	block := NewMockblock(ctrl)
	block.EXPECT().compact(gomock.Any()).Return(0, 0, fmt.Errorf("err"))
	s.blocks[uint16(3)] = block
	flusher.EXPECT().FlushPrimitiveField(gomock.Any(), gomock.Any())
	flushSize = store.FlushFieldTo(flusher, field.Meta{
		ID:   10,
		Type: field.SummaryField,
		Name: "f1",
	})
	assert.True(t, flushSize == 0)
}

func TestComplexFieldStore_MemSize(t *testing.T) {
	store := newComplexFieldStore(10, field.SummaryField)
	assert.Equal(t, emptyComplexFieldStoreSize, store.MemSize())
}

func TestComplexFieldStore_SlotRange(t *testing.T) {
	store := newComplexFieldStore(10, field.SummaryField)
	_, _, err := store.SlotRange()
	assert.Error(t, err)
	writeCtx := writeContext{
		blockStore:   newBlockStore(30),
		timeInterval: 10,
		metricID:     1,
		slotIndex:    1,
		familyTime:   0,
	}
	store.WriteInt(uint16(1), int64(10), writeCtx)
	start, end, err := store.SlotRange()
	assert.NoError(t, err)
	assert.Equal(t, 1, start)
	assert.Equal(t, 1, end)
}
