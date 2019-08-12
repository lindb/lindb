package indextbl

import (
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/field"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_MetricsNameIDReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSnapShot := kv.NewMockSnapshot(ctrl)
	mockReader1 := table.NewMockReader(ctrl)
	mockReader2 := table.NewMockReader(ctrl)
	mockSnapShot.EXPECT().Readers().Return([]table.Reader{mockReader1, mockReader2}).AnyTimes()

	metricNameIDReader := NewMetricsNameIDReader(mockSnapShot)
	// mock readers return nil
	mockReader1.EXPECT().Get(uint32(1)).Return(nil)
	mockReader2.EXPECT().Get(uint32(1)).Return(nil)
	data, metricIDSeq, tagIDSeq := metricNameIDReader.ReadMetricNS(1)
	assert.Nil(t, data)
	assert.Zero(t, metricIDSeq)
	assert.Zero(t, tagIDSeq)
	// mock ok
	mockReader1.EXPECT().Get(uint32(2)).Return([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	mockReader2.EXPECT().Get(uint32(2)).Return(nil)
	data, metricIDSeq, tagIDSeq = metricNameIDReader.ReadMetricNS(2)
	assert.Nil(t, data)
	assert.NotZero(t, metricIDSeq)
	assert.NotZero(t, tagIDSeq)
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

	return metaFlusher1.valueBuf.Bytes(), metaFlusher2.valueBuf.Bytes()
}

func Test_MetricsMetaReader_ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSnapShot := kv.NewMockSnapshot(ctrl)
	mockReader1 := table.NewMockReader(ctrl)
	mockReader2 := table.NewMockReader(ctrl)
	mockSnapShot.EXPECT().Readers().Return([]table.Reader{mockReader1, mockReader2}).AnyTimes()

	metaReader := NewMetricsMetaReader(mockSnapShot)
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
	tagID := metaReader.ReadTagID(2, "a2")
	assert.Equal(t, uint32(7), tagID)
	// tag not found
	tagID = metaReader.ReadTagID(2, "a3")
	assert.Zero(t, tagID)

	// field found
	fieldID, fieldType := metaReader.ReadFieldID(2, "sum2")
	assert.Equal(t, uint16(5), fieldID)
	assert.Equal(t, field.SumField, fieldType)
	// field not found
	fieldID, fieldType = metaReader.ReadFieldID(2, "sum3")
	assert.Equal(t, uint16(0), fieldID)
	assert.Equal(t, field.Type(0), fieldType)
}

func Test_MetricsMetaReader_readBlock_corrupt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSnapShot := kv.NewMockSnapshot(ctrl)
	metaReaderINTF := NewMetricsMetaReader(mockSnapShot)
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
