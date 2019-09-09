package tblstore

import (
	"testing"

	"github.com/lindb/lindb/kv"

	art "github.com/plar/go-adaptive-radix-tree"
	"github.com/stretchr/testify/assert"
)

func buildNameIDBlocks() (data [][]byte) {
	nopFlusher := kv.NewNopFlusher()
	nameIDFlusher := NewMetricsNameIDFlusher(nopFlusher)

	nameIDFlusher.FlushNameID("1", 1)
	nameIDFlusher.FlushNameID("2", 2)
	nameIDFlusher.FlushNameID("3", 3)
	nameIDFlusher.FlushNameID("4", 4)
	nameIDFlusher.FlushMetricsNS(1, 4, 8)

	data = append(data, append([]byte{}, nopFlusher.Bytes()...))

	nameIDFlusher.FlushNameID("5", 5)
	nameIDFlusher.FlushNameID("6", 6)
	nameIDFlusher.FlushMetricsNS(1, 6, 7)
	data = append(data, append([]byte{}, nopFlusher.Bytes()...))
	return data
}

func Test_MetricsNameIDMerger(t *testing.T) {
	reader := NewMetricsNameIDReader(nil).(*metricsNameIDReader)
	m := NewMetricsNameIDMerger()
	// empty value
	data, err := m.Merge(0, nil)
	assert.Nil(t, data)
	assert.NotNil(t, err)

	// invalid block
	data, err = m.Merge(0, [][]byte{{1, 2, 3, 4}})
	assert.Nil(t, data)
	assert.NotNil(t, err)

	// build test data
	blocks := buildNameIDBlocks()
	data, err = m.Merge(1, blocks)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	content, metricIDSeq, tagKeyIDSeq, _ := reader.ReadBlock(data)
	assert.NotNil(t, content)
	assert.Equal(t, uint32(6), metricIDSeq)
	assert.Equal(t, uint32(8), tagKeyIDSeq)

	tree := art.New()
	assert.Nil(t, reader.UnmarshalBinaryToART(tree, content))
	assert.Equal(t, 6, tree.Size())
}

func Test_MetricsNameIDMerger_error(t *testing.T) {
	m := NewMetricsNameIDMerger()

	goodData := []byte{31, 139, 8, 0, 0, 0, 0, 0, 4, 255, 0, 6, 0, 249, 255, 1, 49, 1, 0, 0, 0, 1, 0, 0,
		255, 255, 85, 132, 99, 94, 6, 0, 0, 0,
		1, 0, 0, 0, 1, 0, 0, 0}
	offset := len(goodData) - metricNameIDSequenceSize

	// mock bad length
	badData1 := append([]byte{}, goodData[:offset]...)
	badData1 = append(badData1, byte(32))
	badData1 = append(badData1, goodData[offset:]...)
	_, err := m.Merge(0, [][]byte{badData1})
	assert.NotNil(t, err)
}
