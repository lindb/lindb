package metricsnameid

import (
	"encoding/binary"
	"strconv"
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"

	"github.com/golang/mock/gomock"
	art "github.com/plar/go-adaptive-radix-tree"
	"github.com/stretchr/testify/assert"
)

func Test_MetricsNameIDReader_ReadMetricNS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockReader1 := table.NewMockReader(ctrl)
	mockReader2 := table.NewMockReader(ctrl)

	metricNameIDReader := NewReader([]table.Reader{mockReader1, mockReader2})
	// mock readers return nil
	mockReader1.EXPECT().Get(uint32(1)).Return(nil, false)
	mockReader2.EXPECT().Get(uint32(1)).Return(nil, true)
	data, metricIDSeq, tagKeyIDSeq, ok := metricNameIDReader.ReadMetricNS(1)
	assert.Nil(t, data)
	assert.Zero(t, metricIDSeq)
	assert.Zero(t, tagKeyIDSeq)
	assert.False(t, ok)
	// mock ok
	mockReader1.EXPECT().Get(uint32(2)).Return([]byte{1, 2, 3, 4, 5, 6, 7, 8}, true)
	mockReader2.EXPECT().Get(uint32(2)).Return(nil, false)
	data, metricIDSeq, tagKeyIDSeq, ok = metricNameIDReader.ReadMetricNS(2)
	for _, d := range data {
		assert.Len(t, d, 0)
	}
	assert.NotZero(t, metricIDSeq)
	assert.NotZero(t, tagKeyIDSeq)
	assert.True(t, ok)
}

func Test_MetricsNameIDReader_UnmarshalBinaryToART(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nopKVFlusher := kv.NewNopFlusher()
	nameIDFlusher := NewFlusher(nopKVFlusher).(*flusher)
	for i := 0; i < 10000; i++ {
		nameIDFlusher.FlushNameID(strconv.Itoa(i), uint32(i))
	}
	_ = nameIDFlusher.FlushMetricsNS(1, 1, 1)
	data := nopKVFlusher.Bytes()

	nameIDReader := NewReader(nil).(*reader)
	content, _, _, ok := nameIDReader.ReadBlock(data)
	assert.True(t, ok)

	tree := art.New()
	err := nameIDReader.UnmarshalBinaryToART(tree, content)
	assert.Nil(t, err)
	assert.Equal(t, 10000, tree.Size())
}

func Test_ARTTree_error(t *testing.T) {
	tree := art.New()
	nameIDReader := NewReader(nil).(*reader)

	assert.Nil(t, nameIDReader.UnmarshalBinaryToART(tree, nil))
	assert.NotNil(t, nameIDReader.UnmarshalBinaryToART(tree, []byte{1, 2}))

	goodData := []byte{31, 139, 8, 0, 0, 0, 0, 0, 4, 255, 0, 6, 0, 249, 255, 1, 49, 1, 0, 0, 0, 1, 0, 0,
		255, 255, 85, 132, 99, 94, 6, 0, 0, 0,
		1, 0, 0, 0, 1, 0, 0, 0}
	offset := len(goodData) - metricNameIDSequenceSize

	// mock bad length
	badData1 := append([]byte{}, goodData[:offset]...)
	badData1 = append(badData1, byte(32))
	badData1 = append(badData1, goodData[offset:]...)
	assert.NotNil(t, nameIDReader.UnmarshalBinaryToART(tree, badData1))

	// mock bad metricName
	var buf [8]byte
	binary.PutUvarint(buf[:], 3)
	badData2 := append([]byte{}, goodData[:offset]...)
	badData2 = append(badData2, buf[:1]...)
	badData2 = append(badData2, []byte("abc")...)
	badData2 = append(badData2, byte(1))
	badData2 = append(badData2, goodData[offset:]...)
	assert.NotNil(t, nameIDReader.UnmarshalBinaryToART(tree, badData2))

	// reset failure
	assert.NotNil(t, nameIDReader.UnmarshalBinaryToART(tree, []byte{1}))
}
