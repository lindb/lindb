package bufioutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_testWriteFile = "bufioWriter.test"
)

var (
	_writeContent = []byte("TEST LINDB WRITE")
)

func TestNewBufioWriter(t *testing.T) {
	defer os.Remove(_testWriteFile)
	bw, err := NewBufioWriter(_testWriteFile)

	assert.Nil(t, err)
	assert.NotNil(t, bw)
}

func TestBufioWriter_Write_Size(t *testing.T) {
	defer os.Remove(_testWriteFile)
	bw, _ := NewBufioWriter(_testWriteFile)

	for i := 0; i < 30; i++ {
		bw.Write(_writeContent)
	}
	assert.Equal(t, len(_writeContent)*30+120, int(bw.Size()))
}

func BenchmarkBufioWriter_Write(b *testing.B) {
	defer os.Remove(_testWriteFile)
	bw, _ := NewBufioWriter(_testWriteFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := bw.Write(_writeContent); err != nil {
			assert.Nil(b, err)
		}
	}
}

func TestBufioWriter_Close(t *testing.T) {
	defer os.Remove(_testWriteFile)
	bw, _ := NewBufioWriter(_testWriteFile)

	expectedLength := (len(_writeContent) + 4) * 100000
	for i := 0; i < 100000; i++ {
		bw.Write(_writeContent)
	}
	bw.Sync()
	assert.Nil(t, bw.Sync())
	assert.Nil(t, bw.Close())

	data, err := ioutil.ReadFile(_testWriteFile)
	assert.Nil(t, err)
	assert.Len(t, data, expectedLength)
}
