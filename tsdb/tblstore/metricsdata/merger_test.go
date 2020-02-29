package metricsdata

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

func TestMerger_merge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	flusher := NewMockFlusher(ctrl)
	seriesMerger := NewMockSeriesMerger(ctrl)
	merge := NewMerger()
	m := merge.(*merger)
	m.dataFlusher = flusher
	m.seriesMerger = seriesMerger
	// case 1: new reader err
	data, err := merge.Merge(1, [][]byte{{1, 2, 3}})
	assert.Error(t, err)
	assert.Nil(t, data)
	// case 2: series merge err
	flusher.EXPECT().FlushFieldMetas(field.Metas{{ID: 2, Type: field.SumField}, {ID: 10, Type: field.MinField}}).AnyTimes()
	seriesMerger.EXPECT().merge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), uint16(5), uint16(15)).
		Return(fmt.Errorf("err"))
	data, err = merge.Merge(
		1,
		[][]byte{
			mockMetricMergeBlock([]uint32{1, 2, 4}, 10, 10),
			mockMetricMergeBlock([]uint32{2, 20}, 15, 15),
			mockMetricMergeBlock([]uint32{2, 30}, 5, 5),
		})
	assert.Error(t, err)
	assert.Nil(t, data)
	// case 3: merge success
	seriesMerger.EXPECT().merge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), uint16(10), uint16(15)).
		Return(nil).Times(4)
	gomock.InOrder(
		flusher.EXPECT().FlushSeries(uint32(1)),
		flusher.EXPECT().FlushSeries(uint32(2)),
		flusher.EXPECT().FlushSeries(uint32(4)),
		flusher.EXPECT().FlushSeries(uint32(20)),
		flusher.EXPECT().FlushMetric(uint32(1), uint16(10), uint16(15)).Return(nil),
	)
	data, err = merge.Merge(
		1,
		[][]byte{
			mockMetricMergeBlock([]uint32{1, 2, 4}, 10, 10),
			mockMetricMergeBlock([]uint32{2, 20}, 15, 15),
		})
	assert.NoError(t, err)
	assert.False(t, len(data) > 0) // data flush is mock
	// case 4: flush metric err
	seriesMerger.EXPECT().merge(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), uint16(10), uint16(10)).
		Return(nil)
	gomock.InOrder(
		flusher.EXPECT().FlushSeries(uint32(1)),
		flusher.EXPECT().FlushMetric(uint32(1), uint16(10), uint16(10)).Return(fmt.Errorf("err")),
	)
	data, err = merge.Merge(
		1,
		[][]byte{
			mockMetricMergeBlock([]uint32{1}, 10, 10),
		})
	assert.Error(t, err)
	assert.Nil(t, data)
}

func mockMetricMergeBlock(seriesIDs []uint32, start, end uint16) []byte {
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas(field.Metas{
		{ID: 2, Type: field.SumField},
		{ID: 10, Type: field.MinField},
	})
	for _, seriesID := range seriesIDs {
		flusher.FlushField(field.Key(stream.ReadUint16([]byte{2, byte(1)}, 0)), []byte{1, 2, 3})
		flusher.FlushField(field.Key(stream.ReadUint16([]byte{10, byte(2)}, 0)), []byte{1, 2, 3})
		flusher.FlushSeries(seriesID)
	}
	_ = flusher.FlushMetric(uint32(10), start, end)
	return nopKVFlusher.Bytes()
}
