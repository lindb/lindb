package bufioutil

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewBufioReader(t *testing.T) {
	defer os.Remove(_testFile)
	br, err := NewBufioReader(_testFile)
	assert.NotNil(t, err)
	assert.Nil(t, br)

	os.Create(_testFile)
	br, err = NewBufioReader(_testFile)
	assert.Nil(t, err)
	assert.NotNil(t, br)
}

func TestBufioReader_content(t *testing.T) {
	defer os.Remove(_testFile)
	bw, _ := NewBufioWriter(_testFile)

	f, _ := os.Open(_testFile)
	br := bufioReader{
		f: f,
		r: bufio.NewReader(f)}

	bw.Write([]byte("a"))
	bw.Flush()
	br.Read()
	assert.Equal(t, 1, len(br.content))
	assert.Equal(t, 1, cap(br.content))

	bw.Write([]byte("abcde"))
	bw.Flush()
	br.Read()
	assert.Equal(t, 5, len(br.content))
	assert.Equal(t, 5, cap(br.content))

	bw.Write([]byte("xy"))
	bw.Flush()
	br.Read()
	assert.Equal(t, 2, len(br.content))
	assert.Equal(t, 5, cap(br.content))
}

func BenchmarkBufioReader_Read(b *testing.B) {
	defer os.Remove(_testFile)
	bw, _ := NewBufioWriter(_testFile)
	br, _ := NewBufioReader(_testFile)

	for i := 0; i < b.N; i++ {
		bw.Write(_testContent)
	}
	bw.Sync()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eof, content, err := br.Read()
		if i < 100 || i > b.N-100 {
			assert.False(b, eof)
			assert.Equal(b, _testContent, content)
			assert.Nil(b, err)
		}
	}
	eof, content, err := br.Read()
	assert.True(b, eof)
	assert.Nil(b, content)
	assert.Nil(b, err)
}

func TestBufioReader_Count_Reset_Close(t *testing.T) {
	defer os.Remove(_testFile)
	defer os.Remove("new" + _testFile)
	os.Create("new" + _testFile)
	bw, _ := NewBufioWriter(_testFile)
	br, _ := NewBufioReader(_testFile)

	for i := 0; i < 100000; i++ {
		bw.Write(_testContent)
	}
	bw.Sync()

	for {
		eof, _, _ := br.Read()
		if eof {
			break
		}
	}
	assert.Equal(t, int(br.Count()), (len(_testContent)+4)*100000)

	err := br.Reset("new" + _testFile)
	assert.Nil(t, err)
	assert.Equal(t, 0, int(br.Count()))

	assert.Nil(t, br.Close())
}

func TestBufioReader_Size(t *testing.T) {
	defer os.Remove(_testFile)
	bw, _ := NewBufioWriter(_testFile)
	br, _ := NewBufioReader(_testFile)

	for i := 1; i < 11; i++ {
		bw.Write(_testContent)
		bw.Flush()
		size, err := br.Size()
		assert.Nil(t, err)
		assert.Equal(t, int64(len(_testContent)*i+4*i), size)
	}
}
