package metricsdata

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series/field"
)

func TestField_read(t *testing.T) {
	block := mockMetricMergeBlock([]uint32{1}, 5, 5)
	r, err := NewReader("1.sst", block)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	scanner := newDataScanner(r)
	seriesPos := scanner.scan(0, 1)
	fReader := newFieldReader(scanner.fieldIndexes(), block, seriesPos, 5, 5)
	start, end := fReader.slotRange()
	assert.Equal(t, uint16(5), start)
	assert.Equal(t, uint16(5), end)
	// case 1: field 1 not exist
	data := fReader.getFieldData(1)
	assert.Nil(t, data)
	// case 2: field 2 exist
	data = fReader.getFieldData(2)
	assert.True(t, len(data) > 0)
	// case 3: field 10 exist
	data = fReader.getFieldData(10)
	assert.True(t, len(data) > 0)
	// case 4: field 20 not exist
	data = fReader.getFieldData(20)
	assert.Nil(t, data)
	// case 5: complete cannot get field
	fReader.close()
	data = fReader.getFieldData(10)
	assert.Nil(t, data)
	// case 6: no fields
	fReader = newFieldReader(scanner.fieldIndexes(), []byte{0, 0, 0}, 0, 5, 5)
	data = fReader.getFieldData(10)
	assert.Nil(t, data)
}

func TestFieldReader_close(t *testing.T) {
	block := mockMetricMergeBlock([]uint32{1}, 5, 5)
	r, err := NewReader("1.sst", block)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	scanner := newDataScanner(r)
	seriesPos := scanner.scan(0, 1)
	fReader := newFieldReader(scanner.fieldIndexes(), block, seriesPos, 5, 5)
	fReader.close()
	data := fReader.getFieldData(2)
	assert.Nil(t, data)
}

func TestFieldReader_reset(t *testing.T) {
	block := mockMetricMergeBlock([]uint32{1}, 5, 5)
	r, err := NewReader("1.sst", block)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	scanner := newDataScanner(r)
	seriesPos := scanner.scan(0, 1)
	fReader := newFieldReader(scanner.fieldIndexes(), block, seriesPos, 5, 5)
	start, end := fReader.slotRange()
	assert.Equal(t, uint16(5), start)
	assert.Equal(t, uint16(5), end)
	data := fReader.getFieldData(2)
	assert.True(t, len(data) > 0)
	data = fReader.getFieldData(10)
	assert.True(t, len(data) > 0)

	// mock diff field
	nopKVFlusher := kv.NewNopFlusher()
	flusher := NewFlusher(nopKVFlusher)
	flusher.FlushFieldMetas(field.Metas{
		{ID: 10, Type: field.MinField},
	})
	flusher.FlushField([]byte{1, 2, 3})
	flusher.FlushSeries(10)
	_ = flusher.FlushMetric(uint32(10), start, end)
	block = nopKVFlusher.Bytes()

	// reset value
	fReader.reset(block, seriesPos, 15, 15)
	start, end = fReader.slotRange()
	assert.Equal(t, uint16(15), start)
	assert.Equal(t, uint16(15), end)
	data = fReader.getFieldData(10)
	assert.True(t, len(data) > 0)
}
func TestFieldReader_read_one_field(t *testing.T) {
	block := mockMetricMergeBlockOneField([]uint32{1}, 5, 5)
	r, err := NewReader("1.sst", block)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	scanner := newDataScanner(r)
	seriesPos := scanner.scan(0, 1)
	fReader := newFieldReader(scanner.fieldIndexes(), block, seriesPos, 5, 5)
	start, end := fReader.slotRange()
	assert.Equal(t, uint16(5), start)
	assert.Equal(t, uint16(5), end)
	// case 1: field 1 not exist
	data := fReader.getFieldData(1)
	assert.Nil(t, data)
	// case 2: field 2 exist
	data = fReader.getFieldData(2)
	assert.True(t, len(data) > 0)
	// case 3: close cannot reader data
	fReader.close()
	data = fReader.getFieldData(2)
	assert.Nil(t, data)
}
