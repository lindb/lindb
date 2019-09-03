package tblstore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/series/field"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_tsEntryBuilder(t *testing.T) {
	entry := newSeriesEntryBuilder()
	assert.NotNil(t, entry)

	assert.Nil(t, entry.bytes(nil))
	assert.Len(t, entry.bytes(nil), 0)

	entry.addField(uint16(1), []byte("abcd"), 15, 31)
	entry.addField(uint16(2), []byte("efgh"), 17, 19)
	entry.addField(uint16(4), []byte("ijk"), 14, 30)

	metaFieldsID := []uint16{1, 2, 3, 4}
	var copyData []byte
	data := entry.bytes(metaFieldsID)
	copyData = make([]byte, len(data))
	copy(copyData, data)

	r := bytes.NewReader(copyData)
	minTime, _ := binary.ReadUvarint(r)
	assert.Equal(t, 14, int(minTime)) // minTime
	maxTime, _ := binary.ReadUvarint(r)
	assert.Equal(t, 31, int(maxTime)) // maxTime
	size, _ := binary.ReadUvarint(r)  // size of bit-array
	assert.Equal(t, 1, int(size))
	// bit-array
	assert.Equal(t, 1<<0+1<<1+1<<3, int(copyData[3]))

	expected := []uint8{14, 31, 1, 11, 4, 4, 3, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107}
	assert.Equal(t, len(expected), len(copyData))
	for idx := range expected {
		assert.Equal(t, expected[idx], copyData[idx])
	}
	assert.Equal(t, copyData, entry.bytes(metaFieldsID))
}

func Test_metricBlockBuilder_addTSEntry(t *testing.T) {
	block := newBlockBuilder()
	assert.NotNil(t, block)

	block.addSeries(uint32(1), []byte("a"))
	block.addSeries(uint32(2), []byte("bc"))
	block.addSeries(uint32(3), []byte("def"))
	assert.Equal(t, "abcdef", string(block.bytes()))

	err := block.finish()
	assert.Nil(t, err)
	assert.NotEqual(t, "abcdef", string(block.bytes()))

	block.reset()
	assert.Empty(t, block.bytes())
}

func Test_metricBlockBuilder_finish(t *testing.T) {
	block := newBlockBuilder()

	block.appendFieldMeta(uint16(10), field.SumField)
	block.appendFieldMeta(uint16(11), field.SumField)
	block.appendFieldMeta(uint16(12), field.SumField)
	block.addStartEndTime(1, 2)
	block.addStartEndTime(3, 7)
	block.addStartEndTime(4, 5)

	block.addSeries(uint32(1), []byte("a"))
	block.appendFieldMeta(uint16(20), field.SumField)
	block.appendFieldMeta(uint16(21), field.SumField)
	block.appendFieldMeta(uint16(22), field.SumField)
	block.addStartEndTime(1, 9)
	block.addStartEndTime(2, 3)
	block.addStartEndTime(8, 10)
	block.addSeries(uint32(2), []byte("bc"))
	block.appendFieldMeta(uint16(30), field.SumField)
	block.appendFieldMeta(uint16(31), field.SumField)
	block.appendFieldMeta(uint16(32), field.SumField)
	block.addStartEndTime(1, 9)
	block.addStartEndTime(2, 3)
	block.addStartEndTime(8, 10)
	block.addSeries(uint32(3), []byte("def"))

	assert.Nil(t, block.finish())
	data := block.bytes()

	assert.Equal(t, "abcdef", string(data[:6]))
	// validate pos of meta
	footer := data[len(data)-16:]
	assert.Equal(t, uint32(6), binary.LittleEndian.Uint32(footer[:4]))
	posOfMeta := binary.LittleEndian.Uint32(footer[8:])
	// validate start-Time
	r := bytes.NewReader(data[posOfMeta:])
	startTime, _ := binary.ReadUvarint(r)
	assert.Equal(t, uint64(1), startTime)
	// validate end-time
	endTime, _ := binary.ReadUvarint(r)
	assert.Equal(t, uint64(10), endTime)
	// validate fields-id count
	count, _ := binary.ReadUvarint(r)
	assert.Equal(t, uint64(9), count)

	// validate reset
	block.reset()
	assert.Zero(t, block.minStartTime)
	assert.Zero(t, block.maxEndTime)
	assert.Len(t, block.metaFieldsIDMap, 0)
	assert.Len(t, block.metaFieldsID, 0)
}

func Test_TableFlusher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFlusher := kv.NewMockFlusher(ctrl)

	// add error
	tw := NewMetricsDataFlusher(mockFlusher, 10)
	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(fmt.Errorf("test error"))
	err := tw.FlushMetric(uint32(1))
	assert.NotNil(t, err)
	// close error
	mockFlusher.EXPECT().Commit().Return(fmt.Errorf("close error"))
	assert.NotNil(t, tw.Commit())
	// common write
	tw.FlushFieldMeta(1, field.SumField)
	tw.FlushField(1, []byte("test-field"), 1, 1)
	tw.FlushSeries(uint32(2))

	mockFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	assert.Nil(t, tw.FlushMetric(uint32(3)))

	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			for z := 0; z < 100; z++ {
				tw.FlushFieldMeta(uint16(z), field.SumField)
				tw.FlushField(uint16(z), []byte("test-field"), 1, 2)
				tw.FlushField(uint16(z), []byte("test-field"), 0, 2)
			}
			tw.FlushSeries(uint32(y))
		}
		assert.Nil(t, tw.FlushMetric(uint32(x)))
	}
	mockFlusher.EXPECT().Commit().Return(nil)
	assert.Nil(t, tw.Commit())
}
