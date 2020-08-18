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
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

var bitMapMarshal = encoding.BitmapMarshal

func TestFlusher_flush_metric(t *testing.T) {
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas([]field.Meta{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}})
	// no field for series
	flusher.FlushSeries(5)

	flusher.FlushField([]byte{1, 2, 3})
	flusher.FlushField([]byte{10, 20, 30})
	flusher.FlushSeries(10)
	// flush has one field
	flusher.FlushField([]byte{10, 20, 30})
	flusher.FlushField(nil)
	flusher.FlushSeries(100)

	f := flusher.GetFieldMetas()
	assert.Equal(t, field.Metas{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}}, f)
	err := flusher.FlushMetric(39, 10, 13)
	assert.NoError(t, err)
	// field not exist not flush metric
	assert.Empty(t, flusher.GetFieldMetas())

	flusher.FlushFieldMetas([]field.Meta{{ID: 1, Type: field.SumField}})
	flusher.FlushField([]byte{1, 2, 3})
	err = flusher.FlushMetric(40, 10, 13)
	assert.NoError(t, err)

	// metric hasn't series ids
	flusher.FlushFieldMetas([]field.Meta{{ID: 1, Type: field.SumField}})
	flusher.FlushField(nil)
	err = flusher.FlushMetric(50, 10, 13)
	assert.NoError(t, err)

	err = flusher.Commit()
	assert.NoError(t, err)
}

func TestFlusher_flush_big_series_id(t *testing.T) {
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas([]field.Meta{{ID: 1, Type: field.SumField}, {ID: 2, Type: field.SumField}})
	flusher.FlushField([]byte{1, 2, 3})
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
	flusher.FlushField([]byte{1, 2, 3})
	flusher.FlushSeries(100000)
	encoding.BitmapMarshal = func(bitmap *roaring.Bitmap) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	err := flusher.FlushMetric(39, 10, 13)
	assert.Error(t, err)
	assert.Empty(t, flusher.GetFieldMetas())
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
		flusher.FlushField(data)
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
	block := series.NewMockBlock(ctrl)
	qFlow := flow.NewMockStorageQueryFlow(ctrl)
	// case 2: load data success
	cAgg := aggregation.NewMockContainerAggregator(ctrl)
	cAgg.EXPECT().GetFieldAggregates().Return(aggregation.FieldAggregates{sAgg1, nil})
	qFlow.EXPECT().GetAggregator(uint16(0)).Return(cAgg)
	sAgg1.EXPECT().GetAggregator(gomock.Any()).Return(fAgg1, true)
	fAgg1.EXPECT().GetBlock().Return(block)
	block.EXPECT().Append(gomock.Any(), gomock.Any()).AnyTimes()
	qFlow.EXPECT().Reduce("host", gomock.Any()).AnyTimes()
	r.Load(qFlow, 10, []field.ID{2}, 0, roaring.BitmapOf(1, 2, 3, 4).GetContainer(0))
}
