package metricsdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField_read(t *testing.T) {
	block := mockMetricMergeBlock([]uint32{1}, 5, 5)
	r, err := NewReader("1.sst", block)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	scanner := newDataScanner(r)
	seriesPos := scanner.scan(0, 1)
	fReader := newFieldReader(block, seriesPos, 5, 5)
	start, end := fReader.slotRange()
	assert.Equal(t, uint16(5), start)
	assert.Equal(t, uint16(5), end)
	// case 1: field(2) > field(0), not exist
	data := fReader.getPrimitiveData(1, 0)
	assert.Nil(t, data)
	// case 2: field(2) = field(2) but pID(0)<pID(1), not exist
	data = fReader.getPrimitiveData(2, 0)
	assert.Nil(t, data)
	// case 3: field(2) = field(2) and pID(1)=pID(1), found
	data = fReader.getPrimitiveData(2, 1)
	assert.True(t, len(data) > 0)
	// case 4: field(2) = field(2) and pID(3)>pID(1), not exist, go next field
	data = fReader.getPrimitiveData(2, 3)
	assert.Nil(t, data)
	// case 5: field(10) = field(10) and pID(2)=pID(2), found
	data = fReader.getPrimitiveData(10, 2)
	assert.True(t, len(data) > 0)
	// case 5: field(10) = field(10) and pID(3)>pID(2), completed
	data = fReader.getPrimitiveData(10, 3)
	assert.Nil(t, data)
	// case 6: after completed return nil
	data = fReader.getPrimitiveData(10, 2)
	assert.Nil(t, data)
	// case 7: no fields
	fReader = newFieldReader([]byte{0, 0, 0}, 0, 5, 5)
	data = fReader.getPrimitiveData(10, 2)
	assert.Nil(t, data)
	// case 8: reset, field(100) > field(10) , completed
	fReader.reset(block, seriesPos, 5, 5)
	data = fReader.getPrimitiveData(2, 1)
	assert.True(t, len(data) > 0)
	data = fReader.getPrimitiveData(10, 2)
	assert.True(t, len(data) > 0)
	data = fReader.getPrimitiveData(100, 2)
	assert.Nil(t, data)
	data = fReader.getPrimitiveData(10, 2)
	assert.Nil(t, data)
}
