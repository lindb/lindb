package diskdb

import (
	"encoding/binary"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_nameIDCompressor(t *testing.T) {
	compressor := newNameIDCompressor()
	for i := 0; i < 10000; i++ {
		compressor.AddNameID(strconv.Itoa(i), uint32(i))
	}
	data, err := compressor.Close()
	assert.Nil(t, err)

	tree := newArtTree()
	err = tree.UnmarshalBinary(data)
	assert.Nil(t, err)
	assert.Equal(t, 10000, tree.Size())

	compressor2 := newNameIDCompressor()
	for it := tree.Iterator(); it.HasNext(); {
		item, _ := it.Next()
		compressor2.AddNameID(string(item.Key()), item.Value().(uint32))
	}
	_, err = compressor2.Close()
	assert.Nil(t, err)
	_, err = compressor2.Close()
	assert.Nil(t, err)
}

func Test_ARTTree_error(t *testing.T) {
	tree := newArtTree()
	assert.Nil(t, tree.UnmarshalBinary(nil))
	assert.NotNil(t, tree.UnmarshalBinary([]byte{1, 2}))

	compressor := newNameIDCompressor()
	compressor.AddNameID("1", 1)
	goodData, _ := compressor.Close()
	// mock bad length
	badData1 := append([]byte{}, goodData...)
	badData1 = append(badData1, byte(32))
	assert.NotNil(t, tree.UnmarshalBinary(badData1))
	// mock bad metricName
	var buf [8]byte
	binary.PutUvarint(buf[:], 3)
	badData2 := append([]byte{}, goodData...)
	badData2 = append(badData2, buf[:1]...)
	badData2 = append(badData2, []byte("abc")...)
	badData2 = append(badData2, byte(1))
	assert.NotNil(t, tree.UnmarshalBinary(badData2))

}
