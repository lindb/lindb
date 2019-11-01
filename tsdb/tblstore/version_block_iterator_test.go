package tblstore

import (
	"testing"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"

	"github.com/stretchr/testify/assert"
)

func Test_VersionBlockIterator_error(t *testing.T) {
	itr, err := NewVersionBlockIterator(nil)
	assert.Nil(t, itr)
	assert.NotNil(t, err)
}

func Test_VersionBlockIterator(t *testing.T) {
	encoder := encoding.NewDeltaBitPackingEncoder()

	sw := stream.NewBufferWriter(nil)
	// write 2 versions
	sw.PutBytes([]byte{1, 1})
	encoder.Add(int32(sw.Len()))
	sw.PutBytes([]byte{2, 2})
	encoder.Add(int32(sw.Len()))
	pos := sw.Len()
	sw.PutBytes(encoder.Bytes())
	// write footer
	sw.PutUint32(uint32(pos))
	sw.PutUint32(uint32(1))

	data, _ := sw.Bytes()
	itr, err := NewVersionBlockIterator(data)
	assert.Nil(t, err)
	assert.NotNil(t, itr)

	assert.True(t, itr.HasNext())
	_, _ = itr.Next()
	assert.True(t, itr.HasNext())
	_, _ = itr.Next()
	assert.False(t, itr.HasNext())
}
