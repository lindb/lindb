package tblstore

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/tsdb/field"
)

func Test_MetricsNameIDReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockReader1 := table.NewMockReader(ctrl)
	mockReader2 := table.NewMockReader(ctrl)

	metricNameIDReader := NewMetricsNameIDReader([]table.Reader{mockReader1, mockReader2})
	// mock readers return nil
	mockReader1.EXPECT().Get(uint32(1)).Return(nil)
	mockReader2.EXPECT().Get(uint32(1)).Return(nil)
	data, metricIDSeq, tagIDSeq, ok := metricNameIDReader.ReadMetricNS(1)
	assert.Nil(t, data)
	assert.Zero(t, metricIDSeq)
	assert.Zero(t, tagIDSeq)
	assert.False(t, ok)
	// mock ok
	mockReader1.EXPECT().Get(uint32(2)).Return([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	mockReader2.EXPECT().Get(uint32(2)).Return(nil)
	data, metricIDSeq, tagIDSeq, ok = metricNameIDReader.ReadMetricNS(2)
	for _, d := range data {
		assert.Len(t, d, 0)
	}
	assert.NotZero(t, metricIDSeq)
	assert.NotZero(t, tagIDSeq)
	assert.True(t, ok)
}

func prepareData(ctrl *gomock.Controller) ([]byte, []byte) {
	mockFlusher := kv.NewMockFlusher(ctrl)
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	metaFlusherINTF1 := NewMetricsMetaFlusher(mockFlusher)
	metaFlusherINTF2 := NewMetricsMetaFlusher(mockFlusher)
	metaFlusher1 := metaFlusherINTF1.(*metricsMetaFlusher)
	metaFlusher2 := metaFlusherINTF2.(*metricsMetaFlusher)

	metaFlusherINTF1.FlushFieldID("sum1", field.SumField, 1)
	metaFlusherINTF1.FlushFieldID("min1", field.MinField, 2)
	metaFlusherINTF1.FlushTagKeyID("a1", 3)
	metaFlusherINTF1.FlushTagKeyID("b1", 4)
	metaFlusher1.buildMetricMeta()

	metaFlusherINTF2.FlushFieldID("sum2", field.SumField, 5)
	metaFlusherINTF2.FlushFieldID("min2", field.MinField, 6)
	metaFlusherINTF2.FlushTagKeyID("a2", 7)
	metaFlusherINTF2.FlushTagKeyID("b2", 8)
	metaFlusher2.buildMetricMeta()

	valueBufData1, _ := metaFlusher1.valueBufWriter.Bytes()
	valueBufData2, _ := metaFlusher2.valueBufWriter.Bytes()
	return valueBufData1, valueBufData2
}

func Test_MetricsMetaReader_ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader1 := table.NewMockReader(ctrl)
	mockReader2 := table.NewMockReader(ctrl)

	metaReader := NewMetricsMetaReader([]table.Reader{mockReader1, mockReader2})
	assert.NotNil(t, metaReader)

	// mock nil
	mockReader1.EXPECT().Get(uint32(1)).Return(nil).Times(2)
	mockReader2.EXPECT().Get(uint32(1)).Return(nil).Times(2)
	metaReader.ReadTagID(1, "test-tag")
	metaReader.ReadFieldID(1, "test-field")

	// mockOK
	data1, data2 := prepareData(ctrl)
	mockReader1.EXPECT().Get(uint32(2)).Return(data1).AnyTimes()
	mockReader2.EXPECT().Get(uint32(2)).Return(data2).AnyTimes()
	// tag found
	tagID, ok := metaReader.ReadTagID(2, "a2")
	assert.Equal(t, uint32(7), tagID)
	assert.True(t, ok)
	assert.Len(t, metaReader.SuggestTagKeys(2, "a", 100), 2)
	assert.Len(t, metaReader.SuggestTagKeys(2, "a", 1), 1)
	// tag not found
	tagID, ok = metaReader.ReadTagID(2, "a3")
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
	metaReader1 := NewMetricsMetaReader(nil)
	assert.Zero(t, metaReader1.ReadMaxFieldID(1))

	// mock normal readers
	mockReader2 := table.NewMockReader(ctrl)
	metaReader := NewMetricsMetaReader([]table.Reader{mockReader2})
	_, data2 := prepareData(ctrl)
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
	metaReaderINTF := NewMetricsMetaReader(nil)
	metaReader := metaReaderINTF.(*metricsMetaReader)

	mockReader := table.NewMockReader(ctrl)

	// remainingBlock corrupt
	ret, _ := prepareData(ctrl)
	ret = append(ret, byte(3))
	mockReader.EXPECT().Get(uint32(1)).Return(ret)
	data1, data2 := metaReader.readMetasBlock(mockReader, 1)
	assert.Nil(t, data1)
	assert.Nil(t, data2)

	// block size not ok
	ret, _ = prepareData(ctrl)
	ret = ret[:5]
	mockReader.EXPECT().Get(uint32(1)).Return(ret)
	data1, data2 = metaReader.readMetasBlock(mockReader, 1)
	assert.Nil(t, data1)
	assert.Nil(t, data2)
}
