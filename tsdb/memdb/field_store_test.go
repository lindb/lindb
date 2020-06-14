package memdb

import (
	"fmt"
	"math"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

var encodeFunc = encoding.NewTSDEncoder

func TestFieldStore_New(t *testing.T) {
	buf := make([]byte, pageSize)

	store := newFieldStore(buf, familyID(12), field.ID(1))
	assert.NotNil(t, store)
	assert.Equal(t, familyID(12), store.GetFamilyID())
	assert.Equal(t, field.ID(1), store.GetFieldID())
	s := store.(*fieldStore)
	assert.Equal(t, uint16(0), s.getStart())
	assert.Equal(t, uint16(15), s.timeWindow())
	assert.Equal(t, field.ID(1), s.GetFieldID())
	key := uint32(1) | uint32(0)<<8 | uint32(12)<<16
	assert.Equal(t, key, s.GetKey())
}

func TestFieldStore_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	buf := make([]byte, pageSize)
	store := newFieldStore(buf, familyID(12), field.ID(1))
	assert.NotNil(t, store)
	s := store.(*fieldStore)

	writtenSize := store.Write(field.SumField, 10, 10.1)
	assert.Equal(t, valueSize+headLen, writtenSize)
	// case 1: get write value
	value, ok := s.getCurrentValue(10, 10)
	assert.True(t, ok)
	assert.InDelta(t, 10.1, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
	// case 2: get not exist value, out of time slot range
	value, ok = s.getCurrentValue(10, 12)
	assert.False(t, ok)
	assert.Equal(t, 0.0, value)
	value, ok = s.getCurrentValue(10, 0)
	assert.False(t, ok)
	assert.Equal(t, 0.0, value)
	// case 3: write exist value, need rollup
	writtenSize = store.Write(field.SumField, 10, 10.1)
	assert.Zero(t, writtenSize)
	value, ok = s.getCurrentValue(10, 10)
	assert.True(t, ok)
	assert.InDelta(t, 20.2, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
	// case 3: write new value
	writtenSize = store.Write(field.SumField, 12, 12.1)
	assert.Equal(t, valueSize, writtenSize)
	value, ok = s.getCurrentValue(10, 12)
	assert.True(t, ok)
	assert.InDelta(t, 12.1, value, 0)
	assert.Equal(t, uint16(2), s.getEnd())
	// case 4: get value in time slot range
	value, ok = s.getCurrentValue(10, 11)
	assert.False(t, ok)
	assert.Equal(t, 0.0, value)
	// case 5: test slot range [10,12]
	thisSlotRange := s.slotRange(s.getStart())
	assert.Equal(t, uint16(10), thisSlotRange.start)
	assert.Equal(t, uint16(12), thisSlotRange.end)
	// case 6: compact for slot < start time, time range[5,12]
	writtenSize = store.Write(field.SumField, 5, 5.3)
	assert.True(t, valueSize < writtenSize)
	thisSlotRange = s.slotRange(s.getStart())
	assert.Equal(t, uint16(5), thisSlotRange.start)
	assert.Equal(t, uint16(12), thisSlotRange.end)
	value, ok = s.getCurrentValue(5, 5)
	assert.True(t, ok)
	assert.InDelta(t, 5.3, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
	// case 7: write old value
	writtenSize = store.Write(field.SumField, 10, 10.1)
	assert.Equal(t, valueSize, writtenSize)
	assert.Equal(t, uint16(5), s.getEnd())
	// case 8: compact for slot > end time, time range[5,12]
	writtenSize = store.Write(field.SumField, 50, 50.1)
	assert.True(t, valueSize < writtenSize)
	thisSlotRange = s.slotRange(s.getStart())
	assert.Equal(t, uint16(5), thisSlotRange.start)
	assert.Equal(t, uint16(50), thisSlotRange.end)
	value, ok = s.getCurrentValue(50, 50)
	assert.True(t, ok)
	assert.InDelta(t, 50.1, value, 0.0)
	assert.Equal(t, uint16(0), s.getEnd())
	// case 9: write 10 slot, compact old value
	writtenSize = store.Write(field.SumField, 10, 10.1)
	assert.True(t, valueSize < writtenSize)
	assert.Equal(t, uint16(0), s.getEnd())
	value, ok = s.getCurrentValue(10, 10)
	assert.True(t, ok)
	assert.InDelta(t, 10.1, value, 0)

	// case 10: test final data by load
	writtenSize = store.Write(field.SumField, 15, 15.1)
	assert.Equal(t, valueSize, writtenSize)
	pAgg := aggregation.NewMockPrimitiveAggregator(ctrl)
	gomock.InOrder(
		pAgg.EXPECT().Aggregate(5, 5.3).Return(false),
		pAgg.EXPECT().Aggregate(10, 40.4).Return(false),
		pAgg.EXPECT().Aggregate(12, 12.1).Return(false),
		pAgg.EXPECT().Aggregate(15, 15.1).Return(false),
		pAgg.EXPECT().Aggregate(50, 50.1).Return(true),
	)
	s.Load(field.SumField,
		pAgg,
		&memScanContext{
			tsd: encoding.GetTSDDecoder(),
		},
	)
}

func TestFieldStore_Write2(t *testing.T) {
	buf := make([]byte, pageSize)
	store := newFieldStore(buf, familyID(12), field.ID(1))
	s := store.(*fieldStore)
	writtenSize := store.Write(field.SumField, 10, 178)
	assert.Equal(t, valueSize+headLen, writtenSize)
	value, ok := s.getCurrentValue(10, 10)
	assert.True(t, ok)
	assert.InDelta(t, 178.0, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
	// write with old slot
	writtenSize = store.Write(field.SumField, 10, 178)
	assert.Equal(t, 0, writtenSize)
	value, ok = s.getCurrentValue(10, 10)
	assert.True(t, ok)
	assert.InDelta(t, 178.0*2, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
}

func TestFieldStore_Write_Compact_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.TSDEncodeFunc = encodeFunc
		ctrl.Finish()
	}()

	encode := encoding.NewMockTSDEncoder(ctrl)
	encoding.TSDEncodeFunc = func(startTime uint16) encoding.TSDEncoder {
		return encode
	}

	buf := make([]byte, pageSize)
	store := newFieldStore(buf, familyID(12), field.ID(1))
	assert.NotNil(t, store)
	s := store.(*fieldStore)

	writtenSize := store.Write(field.SumField, 10, 10.1)
	assert.Equal(t, valueSize+headLen, writtenSize)
	// test compress err
	encode.EXPECT().AppendTime(gomock.Any())
	encode.EXPECT().AppendValue(gomock.Any())
	encode.EXPECT().Bytes().Return(nil, fmt.Errorf("err"))
	writtenSize = store.Write(field.SumField, 100, 100.1)
	assert.Equal(t, valueSize+headLen, writtenSize)
	value, ok := s.getCurrentValue(100, 100)
	assert.True(t, ok)
	assert.InDelta(t, 100.1, value, 0)
	assert.Equal(t, uint16(0), s.getEnd())
	// compress data is nil
	assert.Nil(t, s.compress)
}

func TestFieldStore_FlushFieldTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		encoding.TSDEncodeFunc = encodeFunc
		ctrl.Finish()
	}()

	flusher := metricsdata.NewMockFlusher(ctrl)

	buf := make([]byte, pageSize)
	store := newFieldStore(buf, familyID(12), field.ID(2))
	_ = store.Write(field.SumField, 10, 10.1)
	_ = store.Write(field.SumField, 5, 5.1)

	assert.NotNil(t, store)
	// case 1: flush success
	flusher.EXPECT().FlushField(mockFlushData())
	store.FlushFieldTo(flusher, field.Meta{Type: field.SumField}, flushContext{slotRange: slotRange{start: 2, end: 20}})
	// case 3: flush err
	encode := encoding.NewMockTSDEncoder(ctrl)
	encoding.TSDEncodeFunc = func(startTime uint16) encoding.TSDEncoder {
		return encode
	}
	encode.EXPECT().AppendTime(gomock.Any()).AnyTimes()
	encode.EXPECT().AppendValue(gomock.Any()).AnyTimes()
	encode.EXPECT().BytesWithoutTime().Return(nil, fmt.Errorf("err"))
	store.FlushFieldTo(flusher, field.Meta{Type: field.SumField}, flushContext{slotRange: slotRange{start: 2, end: 20}})
}

func mockFlushData() []byte {
	encode := encoding.NewTSDEncoder(2)
	for i := 2; i <= 20; i++ {
		if i == 5 || i == 10 {
			encode.AppendTime(bit.One)
			if i == 5 {
				encode.AppendValue(math.Float64bits(5.1))
			} else {
				encode.AppendValue(math.Float64bits(10.1))
			}
		} else {
			encode.AppendTime(bit.Zero)
		}
	}
	d, _ := encode.BytesWithoutTime()
	return d
}
