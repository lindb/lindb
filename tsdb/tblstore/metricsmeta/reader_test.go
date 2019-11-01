package metricsmeta

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

func prepareData() ([]byte, []byte) {
	nopKVFlusher := kv.NewNopFlusher()

	metaFlusherINTF1 := NewFlusher(nopKVFlusher)
	metaFlusherINTF2 := NewFlusher(nopKVFlusher)

	metaFlusherINTF1.FlushTagMeta(tag.Meta{Key: "a1", ID: 3})
	metaFlusherINTF1.FlushTagMeta(tag.Meta{Key: "b1", ID: 4})
	metaFlusherINTF1.FlushFieldMeta(field.Meta{ID: 1, Type: field.SumField, Name: "sum1"})
	metaFlusherINTF1.FlushFieldMeta(field.Meta{ID: 2, Type: field.MinField, Name: "min1"})
	_ = metaFlusherINTF1.FlushMetricMeta(2)
	data1 := append([]byte{}, nopKVFlusher.Bytes()...)

	metaFlusherINTF2.FlushTagMeta(tag.Meta{Key: "a2", ID: 7})
	metaFlusherINTF2.FlushTagMeta(tag.Meta{Key: "b2", ID: 8})
	metaFlusherINTF2.FlushFieldMeta(field.Meta{ID: 5, Type: field.SumField, Name: "sum2"})
	metaFlusherINTF2.FlushFieldMeta(field.Meta{ID: 6, Type: field.MinField, Name: "min2"})
	_ = metaFlusherINTF2.FlushMetricMeta(2)
	data2 := append([]byte{}, nopKVFlusher.Bytes()...)
	return data1, data2
}

func Test_MetricsMetaReader_ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader1 := table.NewMockReader(ctrl)
	mockReader2 := table.NewMockReader(ctrl)

	metaReader := NewReader([]table.Reader{mockReader1, mockReader2})
	assert.NotNil(t, metaReader)

	// mock nil
	mockReader1.EXPECT().Get(uint32(1)).Return(nil).Times(2)
	mockReader2.EXPECT().Get(uint32(1)).Return(nil).Times(2)
	metaReader.ReadTagKeyID(1, "test-tag")
	metaReader.ReadFieldID(1, "test-field")

	// mockOK
	data1, data2 := prepareData()
	mockReader1.EXPECT().Get(uint32(2)).Return(data1).AnyTimes()
	mockReader2.EXPECT().Get(uint32(2)).Return(data2).AnyTimes()
	// tag found
	tagID, ok := metaReader.ReadTagKeyID(2, "a2")
	assert.Equal(t, uint32(7), tagID)
	assert.True(t, ok)
	assert.Len(t, metaReader.SuggestTagKeys(2, "a", 100), 2)
	assert.Len(t, metaReader.SuggestTagKeys(2, "a", 1), 1)
	// tag not found
	tagID, ok = metaReader.ReadTagKeyID(2, "a3")
	assert.Zero(t, tagID)
	assert.False(t, ok)

	// field found
	fieldID, fieldType, ok := metaReader.ReadFieldID(2, "sum2")
	assert.True(t, ok)
	assert.Equal(t, uint16(5), fieldID)
	assert.Equal(t, field.SumField, fieldType)
	// field not found
	fieldID, fieldType, ok = metaReader.ReadFieldID(2, "sum3")
	assert.Equal(t, uint16(0), fieldID)
	assert.False(t, ok)
	assert.Equal(t, field.Type(0), fieldType)
}

func Test_MetricsMetaReader_ReadMaxFieldID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// empty readers
	metaReader1 := NewReader(nil)
	assert.Zero(t, metaReader1.ReadMaxFieldID(1))

	// mock normal readers
	mockReader2 := table.NewMockReader(ctrl)
	metaReader := NewReader([]table.Reader{mockReader2})
	_, data2 := prepareData()
	mockReader2.EXPECT().Get(uint32(2)).Return(data2)
	assert.Equal(t, uint16(6), metaReader.ReadMaxFieldID(2))

	// mock corrupt data
	data2 = append(data2, byte(32))
	mockReader2.EXPECT().Get(uint32(2)).Return(data2).Times(2)
	assert.Equal(t, uint16(0), metaReader.ReadMaxFieldID(2))
	assert.Nil(t, metaReader.SuggestTagKeys(2, "", 100))
}

func Test_MetricsMetaReader_readBlock_corrupt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	metaReaderINTF := NewReader(nil)
	metaReader := metaReaderINTF.(*reader)

	mockReader := table.NewMockReader(ctrl)

	// remainingBlock corrupt
	ret, _ := prepareData()
	ret = append(ret, byte(3))
	mockReader.EXPECT().Get(uint32(1)).Return(ret)
	data1, data2 := metaReader.readMetasBlock(mockReader.Get(1))
	assert.Nil(t, data1)
	assert.Nil(t, data2)

	// block size not ok
	ret, _ = prepareData()
	ret = ret[:5]
	mockReader.EXPECT().Get(uint32(1)).Return(ret)
	data1, data2 = metaReader.readMetasBlock(mockReader.Get(1))
	assert.Nil(t, data1)
	assert.Nil(t, data2)
}
