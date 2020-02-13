package metricsdata

import (
	"fmt"
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
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
	assert.Equal(t, uint16(2), f.ID)
	_, ok = flusher.GetFieldMeta(field.ID(20))
	assert.False(t, ok)

	err := flusher.FlushMetric(39)
	assert.NoError(t, err)

	// metric hasn't series ids
	err = flusher.FlushMetric(40)
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
	err := flusher.FlushMetric(39)
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
	err := flusher.FlushMetric(39)
	assert.Error(t, err)
	_, ok := flusher.GetFieldMeta(field.ID(2))
	assert.False(t, ok)
}
