package tblstore

import (
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"

	"github.com/stretchr/testify/assert"
)

func Test_MetricsMetaFlusher(t *testing.T) {
	nopKVFlusher := kv.NewNopFlusher()

	metaFlusher := NewMetricsMetaFlusher(nopKVFlusher)
	// write only tags
	metaFlusher.FlushTagMeta(tag.Meta{Key: "k1", ID: 1})
	metaFlusher.FlushTagMeta(tag.Meta{Key: "k2", ID: 2})
	metaFlusher.FlushMetricMeta(1)
	assert.Equal(t, []byte{
		0x2, 0x6b, 0x31, 0x1, 0x0, 0x0, 0x0, 0x2, 0x6b, 0x32, 0x2, 0x0, 0x0, 0x0, 0xe, 0x0, 0x0, 0x0},
		nopKVFlusher.Bytes())

	// write only fields
	metaFlusher.FlushFieldMeta(field.Meta{ID: 3, Type: field.SumField, Name: "f3"})
	metaFlusher.FlushFieldMeta(field.Meta{ID: 4, Type: field.MinField, Name: "f4"})
	metaFlusher.FlushMetricMeta(2)
	assert.Equal(t, []byte{
		0x3, 0x0, 0x1, 0x2, 0x66, 0x33, 0x4, 0x0, 0x2, 0x2, 0x66, 0x34, 0x0, 0x0, 0x0, 0x0},
		nopKVFlusher.Bytes())

	// write tags fields
	metaFlusher.FlushTagMeta(tag.Meta{Key: "k1", ID: 1})
	metaFlusher.FlushFieldMeta(field.Meta{ID: 3, Type: field.SumField, Name: "f3"})
	metaFlusher.FlushMetricMeta(3)
	assert.Equal(t, []byte{
		0x2, 0x6b, 0x31, 0x1, 0x0, 0x0, 0x0, 0x3, 0x0, 0x1, 0x2, 0x66, 0x33, 0x7, 0x0, 0x0, 0x0},
		nopKVFlusher.Bytes())
	assert.Nil(t, metaFlusher.Commit())
}

func Test_flusher_invalid_input(t *testing.T) {
	badKey := ""
	for i := 0; i < 1000; i++ {
		badKey += "X"
	}

	metaFlusher := NewMetricsMetaFlusher(nil)
	metaFlusher.FlushTagMeta(tag.Meta{Key: "", ID: 1})
	metaFlusher.FlushTagMeta(tag.Meta{Key: badKey, ID: 1})
	metaFlusher.FlushFieldMeta(field.Meta{ID: 1, Type: field.SumField, Name: ""})
	metaFlusher.FlushFieldMeta(field.Meta{ID: 1, Type: field.SumField, Name: badKey})
}
