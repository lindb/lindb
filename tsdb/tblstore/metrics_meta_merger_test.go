package tblstore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/tsdb/field"
)

func buildMetaBlock() (data [][]byte) {
	nopFlusher := kv.NewNopFlusher()
	metaFlusher := NewMetricsMetaFlusher(nopFlusher)

	metaFlusher.FlushTagKeyID("tag1", 1)
	metaFlusher.FlushTagKeyID("tag2", 2)
	metaFlusher.FlushTagKeyID("tag3", 3)
	metaFlusher.FlushFieldID("f1", field.SumField, 1)
	metaFlusher.FlushFieldID("f2", field.SumField, 2)
	_ = metaFlusher.FlushMetricMeta(1)
	data = append(data, append([]byte{}, nopFlusher.Bytes()...))

	metaFlusher.FlushTagKeyID("tag4", 4)
	metaFlusher.FlushTagKeyID("tag5", 5)
	metaFlusher.FlushFieldID("f3", field.SumField, 3)
	metaFlusher.FlushFieldID("f4", field.SumField, 4)
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

	fieldItr := newFieldIDIterator(fieldMetaBlock)
	var fieldIDCount = 0
	for fieldItr.HasNext() {
		fieldIDCount++
		fieldName, fieldType, fieldID := fieldItr.Next()
		assert.Equal(t, fmt.Sprintf("f%d", fieldIDCount), fieldName)
		assert.Equal(t, field.SumField, fieldType)
		assert.Equal(t, uint16(fieldIDCount), fieldID)
	}
}
