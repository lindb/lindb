package tblstore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series/field"
)

func buildMetaBlock() (data [][]byte) {
	nopFlusher := kv.NewNopFlusher()
	metaFlusher := NewMetricsMetaFlusher(nopFlusher)

	metaFlusher.FlushTagKeyID("tag1", 1)
	metaFlusher.FlushTagKeyID("tag2", 2)
	metaFlusher.FlushTagKeyID("tag3", 3)
	metaFlusher.FlushFieldMeta(field.Meta{ID: 1, Type: field.SumField, Name: "f1"})
	metaFlusher.FlushFieldMeta(field.Meta{ID: 2, Type: field.SumField, Name: "f2"})
	_ = metaFlusher.FlushMetricMeta(1)
	data = append(data, append([]byte{}, nopFlusher.Bytes()...))

	metaFlusher.FlushTagKeyID("tag4", 4)
	metaFlusher.FlushTagKeyID("tag5", 5)
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

	tagKeyItr := newTagKeyIDIterator(tagMetaBlock)
	var tagKeyIDCount = 0
	for tagKeyItr.HasNext() {
		tagKeyIDCount++
		tagKey, tagKeyID := tagKeyItr.Next()
		assert.Equal(t, fmt.Sprintf("tag%d", tagKeyIDCount), tagKey)
		assert.Equal(t, uint32(tagKeyIDCount), tagKeyID)
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
