package bufioutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_testFile = "bufioio.test"
)

var (
	_testContent = []byte("eleme.ci.etrace")
)

func TestNewBufioWriter(t *testing.T) {
	defer os.Remove(_testFile)
	bw, err := NewBufioWriter(_testFile)

	assert.Nil(t, err)
	assert.NotNil(t, bw)
}

func TestBufioWriter_Reset(t *testing.T) {
	defer os.Remove(_testFile)
	defer os.Remove("new" + _testFile)

	bw, _ := NewBufioWriter(_testFile)
	bw.Write([]byte("test"))
	bw.Flush()
	assert.Equal(t, int64(8), bw.Size())

	err := bw.Reset("new" + _testFile)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), bw.Size())
	bw.Write([]byte("abcd"))

	stat, _ := os.Stat(_testFile)
	assert.Equal(t, int64(8), stat.Size())
}

func TestBufioWriter_Write_Size(t *testing.T) {
	defer os.Remove(_testFile)
	bw, _ := NewBufioWriter(_testFile)

	for i := 0; i < 30; i++ {
		bw.Write(_testContent)
	}
	assert.Equal(t, len(_testContent)*30+120, int(bw.Size()))
}

func BenchmarkBufioWriter_Write(b *testing.B) {
	defer os.Remove(_testFile)
	bw, _ := NewBufioWriter(_testFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := bw.Write(_testContent); err != nil {
			assert.Nil(b, err)
		}
	}
}

func TestBufioWriter_Close(t *testing.T) {
	defer os.Remove(_testFile)
	bw, _ := NewBufioWriter(_testFile)

	expectedLength := (len(_testContent) + 4) * 100000
	for i := 0; i < 100000; i++ {
		bw.Write(_testContent)
	}
	bw.Sync()
	assert.Nil(t, bw.Sync())
	assert.Nil(t, bw.Close())

	data, err := ioutil.ReadFile(_testFile)
	assert.Nil(t, err)
	assert.Len(t, data, expectedLength)
}
