package tblstore

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

func Test_MetricsDataFlusher(t *testing.T) {
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewMetricsDataFlusher(nopKVFlusher, 10)

	flush := func() []byte {
		for version := 0; version < 10; version++ {
			flusher.FlushFieldMetas([]field.Meta{
				{ID: 1, Type: field.SumField, Name: "sum1"},
				{ID: 2, Type: field.SumField, Name: "sum2"},
				{ID: 3, Type: field.SumField, Name: "sum3"},
				{ID: 4, Type: field.SumField, Name: "sum4"},
			})

			for seriesID := 0; seriesID < 100; seriesID++ {
				flusher.FlushField(1, []byte{1, 2}, 0, 10)
				flusher.FlushField(2, []byte{2, 3}, 10, 20)
				flusher.FlushField(3, []byte{3, 4}, 20, 30)
				flusher.FlushSeries(uint32(seriesID))
			}
			flusher.FlushVersion(series.Version(version))
		}
		assert.Nil(t, flusher.FlushMetric(1))
		return append([]byte{}, nopKVFlusher.Bytes()...)
	}
	// assert resettable
	data1 := flush()
	data2 := flush()
	assert.Equal(t, data1, data2)
}

func Test_MetricsDataFlusher_Commit(t *testing.T) {
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewMetricsDataFlusher(nopKVFlusher, 10)
	assert.Nil(t, flusher.Commit())

	assert.Nil(t, flusher.FlushMetric(1))
}
