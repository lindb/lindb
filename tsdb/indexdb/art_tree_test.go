package indexdb

import (
	"encoding/binary"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ARTTree_marshaller_unmarshaler(t *testing.T) {
	tree := newArtTree()
	assert.NotNil(t, tree)
	for i := 0; i < 10000; i++ {
		tree.Insert([]byte(strconv.Itoa(i)), uint32(i))
	}
	data, err := tree.MarshalBinary()
	assert.Nil(t, err)

	tree2 := newArtTree()
	err = tree2.UnmarshalBinary(data)
	assert.Nil(t, err)
	_, err = tree2.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, tree.Size(), tree2.Size())
	// mock type assertion failure
	tree.Insert([]byte("bad-value"), 1.1111)
	_, err = tree.MarshalBinary()
	assert.Nil(t, err)
}

func Test_ARTTree_error(t *testing.T) {
	tree := newArtTree()
	assert.Nil(t, tree.UnmarshalBinary(nil))

	assert.NotNil(t, tree.UnmarshalBinary([]byte{1, 2}))

	goodTree := newArtTree()
	goodTree.Insert([]byte("1"), uint32(1))
	goodData, _ := goodTree.MarshalBinary()
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
	binary.PutUvarint(buf[:], math.MaxUint32)
	badData2 = append(badData2, buf[:]...)
	assert.NotNil(t, tree.UnmarshalBinary(badData2))

}
