package memdb

import (
	"fmt"
	"math"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/encoding"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestSimpleSegmentStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flusher := metricsdata.NewMockFlusher(ctrl)

	store := newSimpleFieldStore(0)
	assert.Equal(t, int64(0), store.GetFamilyTime())
	assert.NotNil(t, store)
	ss, ok := store.(*simpleFieldStore)
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
		writeCtx.slotIndex = time
		ss.CheckAndCompact(field.SumField, writeCtx)
		ss.Write(field.SumField, &pb.Field{Name: "sum", Field: &pb.Field_Sum{
			Sum: &pb.Sum{
				Value: points[time],
			},
		}}, writeCtx)
	}

	startSlot, endSlot := store.SlotRange()
	assert.Equal(t, uint16(5), startSlot)
	assert.Equal(t, uint16(41), endSlot)

	flusher.EXPECT().FlushPrimitiveField(gomock.Any(), gomock.Any())
	flushSize := store.FlushFieldTo(flusher, field.Meta{
		ID:   10,
		Type: field.SumField,
		Name: "f1",
	}, flushContext{start: 0, end: 100})
	assert.True(t, flushSize > 0)
	startSlot, endSlot = store.SlotRange()
	assert.Equal(t, uint16(0), startSlot)
	assert.Equal(t, uint16(100), endSlot)
	fs := store.(*simpleFieldStore)
	result := map[uint16]float64{
		5:  5.0,
		6:  6.0,
		10: 30.0,
		11: 11.0,
		40: 40.0,
		41: 41.0,
	}
	tsd := encoding.GetTSDDecoder()
	tsd.Reset(fs.compress)
	count := 0
	for i := uint16(0); i <= 100; i++ {
		if tsd.HasValueWithSlot(i) {
			count++
			assert.Equal(t, result[i], math.Float64frombits(tsd.Value()))
		}
	}
	assert.Equal(t, len(result), count)
}

func TestSimpleFieldStore_compact_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		encodeFunc = encoding.NewTSDEncoder
	}()
	encoder := encoding.NewMockTSDEncoder(ctrl)
	encoder.EXPECT().AppendTime(gomock.Any()).AnyTimes()
	encoder.EXPECT().AppendValue(gomock.Any()).AnyTimes()
	encoder.EXPECT().Bytes().Return(nil, fmt.Errorf("err")).AnyTimes()
	encodeFunc = func(startTime uint16) encoding.TSDEncoder {
		return encoder
	}

	flusher := metricsdata.NewMockFlusher(ctrl)

	store := newSimpleFieldStore(0)

	writeCtx := writeContext{
		blockStore:   newBlockStore(30),
		timeInterval: 10,
		metricID:     1,
		familyTime:   0,
	}

	writeCtx.slotIndex = 10
	store.CheckAndCompact(field.SumField, writeCtx)
	store.Write(field.SumField, &pb.Field{Name: "sum", Field: &pb.Field_Sum{
		Sum: &pb.Sum{
			Value: 10.0,
		},
	}}, writeCtx)

	// compact because slot out of current time window
	writeCtx.slotIndex = 40
	store.CheckAndCompact(field.SumField, writeCtx)
	store.Write(field.SumField, &pb.Field{Name: "sum", Field: &pb.Field_Sum{
		Sum: &pb.Sum{
			Value: 40.0,
		},
	}}, writeCtx)
	startSlot, endSlot := store.SlotRange()
	assert.Equal(t, uint16(40), startSlot)
	assert.Equal(t, uint16(40), endSlot)

	flushSize := store.FlushFieldTo(flusher, field.Meta{
		ID:   10,
		Type: field.SumField,
		Name: "f1",
	}, flushContext{start: 0, end: 100})
	assert.Equal(t, 0, flushSize)
}

func TestSimpleFieldStore_load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := newSimpleFieldStore(0)
	agg := aggregation.NewMockPrimitiveAggregator(ctrl)
	ss := store.(*simpleFieldStore)

	writeCtx := writeContext{
		blockStore:   newBlockStore(30),
		timeInterval: 10,
		metricID:     1,
		familyTime:   0,
	}

	writeCtx.slotIndex = 10
	ss.CheckAndCompact(field.SumField, writeCtx)
	ss.Write(field.SumField, &pb.Field{Name: "sum", Field: &pb.Field_Sum{
		Sum: &pb.Sum{
			Value: 10.0,
		},
	}}, writeCtx)
	gomock.InOrder(
		agg.EXPECT().Aggregate(10, 10.0),
	)
	ss.load(field.SumField, 10, 10, []aggregation.PrimitiveAggregator{agg}, &memScanContext{tsd: encoding.GetTSDDecoder()})

	// compact because slot out of current time window
	writeCtx.slotIndex = 40
	ss.CheckAndCompact(field.SumField, writeCtx)
	ss.Write(field.SumField, &pb.Field{Name: "sum", Field: &pb.Field_Sum{
		Sum: &pb.Sum{
			Value: 40.0,
		},
	}}, writeCtx)

	gomock.InOrder(
		agg.EXPECT().Aggregate(10, 10.0),
		agg.EXPECT().Aggregate(40, 40.0),
	)
	ss.load(field.SumField, 10, 40, []aggregation.PrimitiveAggregator{agg}, &memScanContext{tsd: encoding.GetTSDDecoder()})
	// compact before time window
	writeCtx.slotIndex = 10
	ss.CheckAndCompact(field.SumField, writeCtx)
	ss.Write(field.SumField, &pb.Field{Name: "sum", Field: &pb.Field_Sum{
		Sum: &pb.Sum{
			Value: 10.0,
		},
	}}, writeCtx)
	gomock.InOrder(
		agg.EXPECT().Aggregate(10, 20.0).Return(true),
	)
	ss.load(field.SumField, 10, 40, []aggregation.PrimitiveAggregator{agg}, &memScanContext{tsd: encoding.GetTSDDecoder()})
}

func TestSimpleFieldStore_getFieldValue(t *testing.T) {
	store := newSimpleFieldStore(0)
	ss := store.(*simpleFieldStore)
	value := ss.getFieldValue(field.SumField, &pb.Field{Name: "sum", Field: &pb.Field_Sum{
		Sum: &pb.Sum{
			Value: 10.0,
		},
	}})
	assert.Equal(t, 10.0, value)
	value = ss.getFieldValue(field.MinField, &pb.Field{Name: "sum", Field: &pb.Field_Min{
		Min: &pb.Min{
			Value: 10.0,
		},
	}})
	assert.Equal(t, 10.0, value)
	value = ss.getFieldValue(field.MaxField, &pb.Field{Name: "sum", Field: &pb.Field_Max{
		Max: &pb.Max{
			Value: 10.0,
		},
	}})
	assert.Equal(t, 10.0, value)
	value = ss.getFieldValue(field.GaugeField, &pb.Field{Name: "sum", Field: &pb.Field_Gauge{
		Gauge: &pb.Gauge{
			Value: 10.0,
		},
	}})
	assert.Equal(t, 10.0, value)
	value = ss.getFieldValue(field.Unknown, nil)
	assert.Equal(t, 0.0, value)
}
