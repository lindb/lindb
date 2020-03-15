package metricsdata

import (
	"fmt"
	"math"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

var bitMapMarshal = encoding.BitmapMarshal

func TestFlusher_flush_metric(t *testing.T) {
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas([]field.Meta{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}})
	// no field for series
	flusher.FlushSeries(5)

	flusher.FlushField(field.Key(10), []byte{1, 2, 3})
	flusher.FlushField(field.Key(11), []byte{10, 20, 30})
	flusher.FlushSeries(10)

	f, ok := flusher.GetFieldMeta(field.ID(2))
	assert.True(t, ok)
	assert.Equal(t, field.ID(2), f.ID)
	_, ok = flusher.GetFieldMeta(field.ID(20))
	assert.False(t, ok)

	err := flusher.FlushMetric(39, 10, 13)
	assert.NoError(t, err)

	// metric hasn't series ids
	err = flusher.FlushMetric(40, 10, 13)
	assert.NoError(t, err)

	// field not exist not flush metric
	_, ok = flusher.GetFieldMeta(field.ID(2))
	assert.False(t, ok)

	err = flusher.Commit()
	assert.NoError(t, err)
}

func TestFlusher_flush_big_series_id(t *testing.T) {
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas([]field.Meta{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}})
	flusher.FlushField(field.Key(10), []byte{1, 2, 3})
	flusher.FlushSeries(100000)
	err := flusher.FlushMetric(39, 10, 13)
	assert.NoError(t, err)
	err = flusher.Commit()
	assert.NoError(t, err)
}

func TestFlusher_flush_err(t *testing.T) {
	defer func() {
		encoding.BitmapMarshal = bitMapMarshal
	}()
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas([]field.Meta{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}})
	flusher.FlushField(field.Key(10), []byte{1, 2, 3})
	flusher.FlushSeries(100000)
	encoding.BitmapMarshal = func(bitmap *roaring.Bitmap) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	err := flusher.FlushMetric(39, 10, 13)
	assert.Error(t, err)
	_, ok := flusher.GetFieldMeta(field.ID(2))
	assert.False(t, ok)
}

func TestFlusher_TooMany_Data(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	encoder := encoding.NewTSDEncoder(5)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(10.0))
	data, _ := encoder.BytesWithoutTime()

	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas([]field.Meta{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}})
	for i := 0; i < 80000; i++ {
		flusher.FlushField(field.Key(10), data)
		flusher.FlushSeries(uint32(i))
	}
	err := flusher.FlushMetric(39, 5, 5)
	assert.NoError(t, err)
	data = nopKVFlusher.Bytes()
	r, err := NewReader("1.sst", data)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	sAgg1 := aggregation.NewMockSeriesAggregator(ctrl)
	fAgg1 := aggregation.NewMockFieldAggregator(ctrl)
	pAgg1 := aggregation.NewMockPrimitiveAggregator(ctrl)
	qFlow := flow.NewMockStorageQueryFlow(ctrl)
	// case 2: load data success
	qFlow.EXPECT().GetAggregator().Return(aggregation.FieldAggregates{sAgg1, nil})
	sAgg1.EXPECT().GetAggregator(gomock.Any()).Return(fAgg1, true)
	fAgg1.EXPECT().GetAllAggregators().Return([]aggregation.PrimitiveAggregator{pAgg1})
	pAgg1.EXPECT().FieldID().Return(field.PrimitiveID(5))
	pAgg1.EXPECT().Aggregate(gomock.Any(), gomock.Any()).AnyTimes()
	qFlow.EXPECT().Reduce("host", gomock.Any()).AnyTimes()
	r.Load(qFlow, 10, []field.ID{2}, 1, map[string][]uint16{"host": {1, 2, 3, 4}})

}
