package bufioutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_testReadFile = "bufioReader.test"
)

var (
	_readContent = []byte("eleme.ci.etrace")
)

func Test_NewBufioReader(t *testing.T) {
	defer os.Remove(_testReadFile)
	br, err := NewBufioReader(_testReadFile)
	assert.Nil(t, err)
	assert.NotNil(t, br)
}

func TestBufioWriter_setContent(t *testing.T) {
	br := bufioReader{}

	br.setContent([]byte("a"))
	assert.Equal(t, 1, len(br.content))
	assert.Equal(t, 1, cap(br.content))

	br.setContent([]byte("abcde"))
	assert.Equal(t, 5, len(br.content))
	assert.Equal(t, 5, cap(br.content))

	br.setContent([]byte("xy"))
	assert.Equal(t, 2, len(br.content))
	assert.Equal(t, 5, cap(br.content))
}

func BenchmarkBufioReader_Read(b *testing.B) {
	defer os.Remove(_testReadFile)
	bw, _ := NewBufioWriter(_testReadFile)
	br, _ := NewBufioReader(_testReadFile)

	for i := 0; i < b.N; i++ {
		bw.Write(_readContent)
	}
	bw.Sync()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eof, content, err := br.Read()
		if i < 100 || i > b.N-100 {
			assert.False(b, eof)
			assert.Equal(b, _readContent, content)
			assert.Nil(b, err)
		}
	}
	eof, content, err := br.Read()
	assert.True(b, eof)
	assert.Nil(b, content)
	assert.Nil(b, err)
}

func TestBufioReader_Count_Reset_Close(t *testing.T) {
	defer os.Remove(_testReadFile)
	defer os.Remove("new" + _testReadFile)
	bw, _ := NewBufioWriter(_testReadFile)
	br, _ := NewBufioReader(_testReadFile)

	for i := 0; i < 100000; i++ {
		bw.Write(_readContent)
	}
	bw.Sync()

	for {
		eof, _, _ := br.Read()
		if eof {
			break
		}
	}
	assert.Equal(t, int(br.Count()), (len(_readContent)+4)*100000)

	err := br.Reset("new" + _testReadFile)
	assert.Nil(t, err)
	assert.Equal(t, 0, int(br.Count()))

	assert.Nil(t, br.Close())
}
