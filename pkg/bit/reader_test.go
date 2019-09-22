package bit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/bufioutil"
)

func Test_Reader(t *testing.T) {
	var data []byte
	buf := bufioutil.NewBuffer(data)
	reader := NewReader(buf)

	_, err := reader.ReadBit()
	assert.NotNil(t, err)
	_, err = reader.ReadByte()
	assert.NotNil(t, err)
	_, err = reader.ReadBits(10)
	assert.NotNil(t, err)
	_, err = reader.ReadBits(1)
	assert.NotNil(t, err)

	data = append(data, []byte{1, 2, 3, 4, 5, 6, 7, 8}...)
	buf.SetBuf(data)
	reader.Reset()
	_, err = reader.ReadBits(10)
	assert.Nil(t, err)
}
