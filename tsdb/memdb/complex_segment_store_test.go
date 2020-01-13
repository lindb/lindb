package memdb

import (
	"math"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/encoding"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestComplexFieldStore_FlushFieldTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flusher := metricsdata.NewMockFlusher(ctrl)

	store := newComplexFieldStore(0)
	assert.Equal(t, int64(0), store.GetFamilyTime())
	assert.NotNil(t, store)
	ss, ok := store.(*complexFieldStore)
	assert.True(t, ok)

	writeCtx := writeContext{
		blockStore:   newBlockStore(30),
		timeInterval: 10,
		metricID:     1,
		familyTime:   0,
	}
	times := []uint16{
		10,
		11,
		10, // memory auto rollup
		40, // compact because slot out of current time window
		10, // compact before time window
		41, // compact because slot out of current time window
		5,  // compact because slot out of current time window
		6}
	points := map[uint16]float64{
		5:  5.0,
		6:  6.0,
		10: 10.0,
		11: 11.0,
		40: 40.0,
		41: 41.0,
	}
	for _, time := range times {
		value := points[time]
		writeCtx.slotIndex = time
		ss.CheckAndCompact(field.SummaryField, writeCtx)
		ss.Write(field.SummaryField, &pb.Field{Name: "summary", Field: &pb.Field_Summary{
			Summary: &pb.Summary{
				Sum:   value,
				Count: 2 * value,
			},
		}}, writeCtx)
	}

	startSlot, endSlot := store.SlotRange()
	assert.Equal(t, uint16(5), startSlot)
	assert.Equal(t, uint16(41), endSlot)

	flusher.EXPECT().FlushPrimitiveField(gomock.Any(), gomock.Any())
	flushSize := store.FlushFieldTo(flusher, field.Meta{
		ID:   10,
		Type: field.SummaryField,
		Name: "f1",
	}, flushContext{start: 0, end: 100})
	assert.True(t, flushSize > 0)
	startSlot, endSlot = store.SlotRange()
	assert.Equal(t, uint16(0), startSlot)
	assert.Equal(t, uint16(100), endSlot)
	fs := store.(*complexFieldStore)
	result := map[uint16]float64{
		5:  5.0,
		6:  6.0,
		10: 30.0,
		11: 11.0,
		40: 40.0,
		41: 41.0,
	}
	reader := encoding.NewTSDStreamReader(fs.compress)
	defer reader.Close()
	startSlot, endSlot = reader.TimeRange()
	assert.Equal(t, uint16(0), startSlot)
	assert.Equal(t, uint16(100), endSlot)
	// test sum values
	assert.True(t, reader.HasNext())
	fieldID, tsd := reader.Next()
	assert.Equal(t, uint16(1), fieldID)
	count := 0
	for i := uint16(0); i <= 100; i++ {
		if tsd.HasValueWithSlot(i) {
			count++
			assert.Equal(t, result[i], math.Float64frombits(tsd.Value()))
		}
	}
	assert.Equal(t, len(result), count)

	// test count values
	assert.True(t, reader.HasNext())
	fieldID, tsd = reader.Next()
	assert.Equal(t, uint16(2), fieldID)
	count = 0
	for i := uint16(0); i <= 100; i++ {
		if tsd.HasValueWithSlot(i) {
			count++
			assert.Equal(t, result[i]*2, math.Float64frombits(tsd.Value()))
		}
	}
	assert.Equal(t, len(result), count)

	assert.False(t, reader.HasNext())
}

func TestComplexFieldStore_GetFamilyTime(t *testing.T) {
	store := newComplexFieldStore(10)
	writeCtx := writeContext{
		blockStore:   newBlockStore(30),
		timeInterval: 10,
		metricID:     1,
		familyTime:   0,
	}
	store.Write(field.HistogramField, &pb.Field{Name: "test", Field: &pb.Field_Histogram{
		Histogram: &pb.Histogram{
			Sum:   0,
			Count: 0,
		},
	}}, writeCtx)
	assert.Equal(t, int64(10), store.GetFamilyTime())
}

//
//func TestComplexFieldStore_FlushFieldTo(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// test store is empty
//	flusher := metricsdata.NewMockFlusher(ctrl)
//	store := newComplexFieldStore(10, field.SummaryField)
//	flushSize := store.FlushFieldTo(flusher, field.Meta{
//		ID:   10,
//		Type: field.SummaryField,
//		Name: "f1",
//	}, flushContext{start: 0, end: 100})
//	assert.True(t, flushSize == 0)
//
//	// test normal case
//	writeCtx := writeContext{
//		blockStore:   newBlockStore(30),
//		timeInterval: 10,
//		metricID:     1,
//		slotIndex:    1,
//		familyTime:   0,
//	}
//	store.WriteInt(uint16(1), int64(10), writeCtx)
//	flusher.EXPECT().FlushPrimitiveField(gomock.Any(), gomock.Any())
//	flushSize = store.FlushFieldTo(flusher, field.Meta{
//		ID:   10,
//		Type: field.SummaryField,
//		Name: "f1",
//	}, flushContext{start: 0, end: 100})
//	assert.True(t, flushSize > 0)
//
//	// test block compact err
//	s := store.(*complexFieldStore)
//	block := NewMockblock(ctrl)
//	block.EXPECT().compact(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
//	s.blocks[uint16(3)] = block
//	flusher.EXPECT().FlushPrimitiveField(gomock.Any(), gomock.Any())
//	flushSize = store.FlushFieldTo(flusher, field.Meta{
//		ID:   10,
//		Type: field.SummaryField,
//		Name: "f1",
//	}, flushContext{start: 0, end: 100})
//	assert.True(t, flushSize == 0)
//}
//
//func TestComplexFieldStore_MemSize(t *testing.T) {
//	store := newComplexFieldStore(10, field.SummaryField)
//	assert.Equal(t, emptyComplexFieldStoreSize, store.MemSize())
//}
//
//func TestComplexFieldStore_SlotRange(t *testing.T) {
//	store := newComplexFieldStore(10, field.SummaryField)
//	_, _, err := store.SlotRange()
//	assert.Error(t, err)
//	writeCtx := writeContext{
//		blockStore:   newBlockStore(30),
//		timeInterval: 10,
//		metricID:     1,
//		slotIndex:    1,
//		familyTime:   0,
//	}
//	store.WriteInt(uint16(1), int64(10), writeCtx)
//	start, end, err := store.SlotRange()
//	assert.NoError(t, err)
//	assert.Equal(t, 1, start)
//	assert.Equal(t, 1, end)
//}
