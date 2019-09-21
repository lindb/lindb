package tblstore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

func buildMetaBlock() (data [][]byte) {
	nopFlusher := kv.NewNopFlusher()
	metaFlusher := NewMetricsMetaFlusher(nopFlusher)

	metaFlusher.FlushTagMeta(tag.Meta{Key: "tag1", ID: 1})
	metaFlusher.FlushTagMeta(tag.Meta{Key: "tag2", ID: 2})
	metaFlusher.FlushTagMeta(tag.Meta{Key: "tag3", ID: 3})
	metaFlusher.FlushFieldMeta(field.Meta{ID: 1, Type: field.SumField, Name: "f1"})
	metaFlusher.FlushFieldMeta(field.Meta{ID: 2, Type: field.SumField, Name: "f2"})
	_ = metaFlusher.FlushMetricMeta(1)
	data = append(data, append([]byte{}, nopFlusher.Bytes()...))

	metaFlusher.FlushTagMeta(tag.Meta{Key: "tag4", ID: 4})
	metaFlusher.FlushTagMeta(tag.Meta{Key: "tag5", ID: 5})
	metaFlusher.FlushFieldMeta(field.Meta{ID: 3, Type: field.SumField, Name: "f3"})
	metaFlusher.FlushFieldMeta(field.Meta{ID: 4, Type: field.SumField, Name: "f4"})
	_ = metaFlusher.FlushMetricMeta(1)
	data = append(data, append([]byte{}, nopFlusher.Bytes()...))

	return data
}

func Test_MetricsMetaMerger(t *testing.T) {
	m := NewMetricsMetaMerger()
	// merge unavailable block
	data, err := m.Merge(0, nil)
	assert.Nil(t, data)
	assert.Error(t, err)
	// merge normal block
	block := buildMetaBlock()
	data, err = m.Merge(1, block)
	assert.Nil(t, err)
	assert.NotNil(t, data)

	reader := NewMetricsMetaReader(nil).(*metricsMetaReader)
	tagMetaBlock, fieldMetaBlock := reader.readMetasBlock(data)

	tagMetaItr := newTagMetaIterator(tagMetaBlock)
	var tagKeyIDCount = 0
	for tagMetaItr.HasNext() {
		tagKeyIDCount++
		tagMeta := tagMetaItr.Next()
		assert.Equal(t, fmt.Sprintf("tag%d", tagKeyIDCount), tagMeta.Key)
		assert.Equal(t, uint32(tagKeyIDCount), tagMeta.ID)
	}

	fieldItr := newFieldMetaIterator(fieldMetaBlock)
	var fieldIDCount = 0
	for fieldItr.HasNext() {
		fieldIDCount++
		fieldMeta := fieldItr.Next()
		assert.Equal(t, fmt.Sprintf("f%d", fieldIDCount), fieldMeta.Name)
		assert.Equal(t, field.SumField, fieldMeta.Type)
		assert.Equal(t, uint16(fieldIDCount), fieldMeta.ID)
	}
}
