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
	assert.Equal(t, int64(5), bw.Size())

	err := bw.Reset("new" + _testFile)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), bw.Size())
	bw.Write([]byte("abcd"))

	stat, _ := os.Stat(_testFile)
	assert.Equal(t, int64(5), stat.Size())
}

func TestBufioWriter_Write_Size(t *testing.T) {
	defer os.Remove(_testFile)
	bw, _ := NewBufioWriter(_testFile)
	assert.Equal(t, int64(0), bw.Size())
	n, err := bw.Write([]byte(""))
	assert.Equal(t, 1, n)
	assert.Equal(t, int64(1), bw.Size())
	assert.Nil(t, err)

	n, _ = bw.Write([]byte("xyz"))
	assert.Equal(t, 4, n)
	assert.Equal(t, int64(5), bw.Size())

	var s [128]byte
	n, _ = bw.Write(s[:])
	assert.Equal(t, 130, n)
	assert.Equal(t, int64(135), bw.Size())
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

	expectedLength := (len(_testContent) + 1) * 100000
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
