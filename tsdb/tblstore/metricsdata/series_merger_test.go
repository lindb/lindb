package metricsdata

import (
	"fmt"
	"math"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
)

func TestSeriesMerger_merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	flusher := NewMockFlusher(ctrl)
	merger := newSeriesMerger(flusher)
	decodeStreams := make([]*encoding.TSDDecoder, 2)
	reader1 := NewMockFieldReader(ctrl)
	reader2 := NewMockFieldReader(ctrl)
	readers := []FieldReader{reader1, reader2}

	encodeStream := encoding.NewTSDEncoder(5)
	// case 1: merge success and rollup
	reader1.EXPECT().getPrimitiveData(gomock.Any(), gomock.Any()).Return(mockPrimitiveField(10))
	reader1.EXPECT().slotRange().Return(uint16(10), uint16(10))
	reader2.EXPECT().getPrimitiveData(gomock.Any(), gomock.Any()).Return(mockPrimitiveField(10))
	reader2.EXPECT().slotRange().Return(uint16(10), uint16(10))
	var result []byte
	flusher.EXPECT().FlushField(gomock.Any(), gomock.Any()).DoAndReturn(func(key field.Key, data []byte) {
		result = data
	})
	err := merger.merge(field.Metas{{ID: 1, Type: field.SumField}}, decodeStreams, encodeStream, readers, 5, 15)
	assert.NoError(t, err)
	tsd := encoding.GetTSDDecoder()
	tsd.ResetWithTimeRange(result, 5, 15)
	slot := uint16(0)
	for i := uint16(5); i <= 15; i++ {
		if tsd.HasValueWithSlot(i) {
			slot = i
			assert.Equal(t, 20.0, math.Float64frombits(tsd.Value()))
		}
	}
	assert.Equal(t, uint16(10), slot)
	// case 2: merge success with diff slot range
	reader1.EXPECT().getPrimitiveData(gomock.Any(), gomock.Any()).Return(mockPrimitiveField(10))
	reader1.EXPECT().slotRange().Return(uint16(10), uint16(10))
	reader2.EXPECT().getPrimitiveData(gomock.Any(), gomock.Any()).Return(mockPrimitiveField(12))
	reader2.EXPECT().slotRange().Return(uint16(12), uint16(12))
	flusher.EXPECT().FlushField(gomock.Any(), gomock.Any()).DoAndReturn(func(key field.Key, data []byte) {
		result = data
	})
	err = merger.merge(field.Metas{{ID: 1, Type: field.SumField}}, decodeStreams, encodeStream, readers, 5, 15)
	assert.NoError(t, err)
	tsd.ResetWithTimeRange(result, 5, 15)
	c := 0
	for i := uint16(5); i <= 15; i++ {
		if tsd.HasValueWithSlot(i) && (i == 10 || i == 12) {
			c++
			assert.Equal(t, 10.0, math.Float64frombits(tsd.Value()))
		}
	}
	assert.Equal(t, 2, c)
	// case 3: encode stream err
	encodeStream2 := encoding.NewMockTSDEncoder(ctrl)
	reader1.EXPECT().getPrimitiveData(gomock.Any(), gomock.Any()).Return(mockPrimitiveField(10))
	reader1.EXPECT().slotRange().Return(uint16(10), uint16(10))
	reader2.EXPECT().getPrimitiveData(gomock.Any(), gomock.Any()).Return(mockPrimitiveField(12))
	reader2.EXPECT().slotRange().Return(uint16(12), uint16(12))
	encodeStream2.EXPECT().AppendTime(gomock.Any()).AnyTimes()
	encodeStream2.EXPECT().AppendValue(gomock.Any()).AnyTimes()
	encodeStream2.EXPECT().BytesWithoutTime().Return(nil, fmt.Errorf("err"))
	err = merger.merge(field.Metas{{ID: 1, Type: field.SumField}}, decodeStreams, encodeStream2, readers, 5, 15)
	assert.Error(t, err)
}

func mockPrimitiveField(start uint16) []byte {
	encoder := encoding.NewTSDEncoder(start)
	encoder.AppendTime(bit.One)
	encoder.AppendValue(math.Float64bits(10.0))
	data, _ := encoder.BytesWithoutTime()
	return data
}
