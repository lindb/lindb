package memdb

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSimpleSegmentStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flusher := metricsdata.NewMockFlusher(ctrl)

	aggFunc := field.Sum.AggFunc()
	store := newSimpleFieldStore(0, aggFunc)
	assert.Equal(t, int64(0), store.GetFamilyTime())
	assert.NotNil(t, store)
	ss, ok := store.(*simpleFieldStore)
	assert.True(t, ok)

	_, _, err := ss.SlotRange()
	assert.NotNil(t, err)

	flushSize := store.FlushFieldTo(flusher, field.Meta{
		ID:   10,
		Type: field.SumField,
		Name: "f1",
	})
	assert.True(t, flushSize == 0)

	writeCtx := writeContext{
		blockStore:   newBlockStore(30),
		timeInterval: 10,
		metricID:     1,
		familyTime:   0,
	}

	writeCtx.slotIndex = 10
	ss.WriteInt(uint16(1), 100, writeCtx)
	// memory auto rollup
	writeCtx.slotIndex = 11
	ss.WriteInt(uint16(1), 110, writeCtx)
	// memory auto rollup
	writeCtx.slotIndex = 10
	ss.WriteInt(uint16(1), 100, writeCtx)
	// compact because slot out of current time window
	writeCtx.slotIndex = 40
	ss.WriteInt(uint16(1), 20, writeCtx)
	// compact before time window
	writeCtx.slotIndex = 10
	ss.WriteInt(uint16(1), 100, writeCtx)
	// compact because slot out of current time window
	writeCtx.slotIndex = 41
	ss.WriteInt(uint16(1), 50, writeCtx)

	flusher.EXPECT().FlushPrimitiveField(gomock.Any(), gomock.Any())
	flushSize = store.FlushFieldTo(flusher, field.Meta{
		ID:   10,
		Type: field.SumField,
		Name: "f1",
	})
	assert.True(t, flushSize > 0)

	startSlot, endSlot, err := store.SlotRange()
	assert.Nil(t, err)
	assert.Equal(t, 10, startSlot)
	assert.Equal(t, 41, endSlot)
}

func TestSimpleSegmentStore_float(t *testing.T) {
	writeCtx := writeContext{
		blockStore:   newBlockStore(30),
		timeInterval: 10,
		metricID:     1,
		familyTime:   0,
	}

	aggFunc := field.Sum.AggFunc()
	store := newSimpleFieldStore(0, aggFunc)
	assert.Equal(t, int64(0), store.GetFamilyTime())
	assert.NotNil(t, store)
	ss, ok := store.(*simpleFieldStore)
	assert.True(t, ok)
	// write float test
	writeCtx.slotIndex = 10
	ss.WriteFloat(uint16(1), 10, writeCtx)
	// auto rollup
	writeCtx.slotIndex = 10
	ss.WriteFloat(uint16(1), 10, writeCtx)
}

func Test_sStore_error(t *testing.T) {
	store := newSimpleFieldStore(0, field.Sum.AggFunc())
	ss, _ := store.(*simpleFieldStore)
	// compact error test
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	flusher := metricsdata.NewMockFlusher(ctrl)

	mockBlock := NewMockblock(ctrl)
	mockBlock.EXPECT().compact(gomock.Any()).Return(0, 0, fmt.Errorf("compat error")).AnyTimes()
	mockBlock.EXPECT().setIntValue(gomock.Any(), gomock.Any()).Return().AnyTimes()
	mockBlock.EXPECT().getStartTime().Return(12).AnyTimes()
	mockBlock.EXPECT().getEndTime().Return(40).AnyTimes()
	mockBlock.EXPECT().memsize().Return(300).AnyTimes()
	ss.block = mockBlock
	flushSize := store.FlushFieldTo(flusher, field.Meta{
		ID:   10,
		Type: field.SumField,
		Name: "f1",
	})
	assert.True(t, flushSize == 0)

	writeCtx := writeContext{
		blockStore:   newBlockStore(30),
		timeInterval: 10,
		metricID:     1,
		familyTime:   0,
	}

	writeCtx.slotIndex = 10
	ss.WriteInt(uint16(1), 100, writeCtx)
	// memory auto rollup
	writeCtx.slotIndex = 11
	ss.WriteInt(uint16(1), 110, writeCtx)
}

func BenchmarkSimpleSegmentStore(b *testing.B) {
	aggFunc := field.Sum.AggFunc()
	store := newSimpleFieldStore(0, aggFunc)
	ss, _ := store.(*simpleFieldStore)

	writeCtx := writeContext{
		blockStore:   newBlockStore(30),
		timeInterval: 10,
		metricID:     1,
		familyTime:   0,
	}

	writeCtx.slotIndex = 10
	ss.WriteInt(uint16(1), 100, writeCtx)
	// memory auto rollup
	writeCtx.slotIndex = 11
	ss.WriteInt(uint16(1), 110, writeCtx)
	// memory auto rollup
	writeCtx.slotIndex = 10
	ss.WriteInt(uint16(1), 100, writeCtx)
	// compact because slot out of current time window
	writeCtx.slotIndex = 40
	ss.WriteInt(uint16(1), 20, writeCtx)
	// compact before time window
	writeCtx.slotIndex = 10
	ss.WriteInt(uint16(1), 100, writeCtx)
	// compact because slot out of current time window
	writeCtx.slotIndex = 41
	ss.WriteInt(uint16(1), 50, writeCtx)
}
