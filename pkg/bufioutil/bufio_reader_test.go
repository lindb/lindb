// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
	br.Next()
	assert.Equal(t, 1, len(br.content))
	assert.Equal(t, 1, cap(br.content))

	bw.Write([]byte("abcde"))
	bw.Flush()
	br.Next()
	assert.Equal(t, 5, len(br.content))
	assert.Equal(t, 5, cap(br.content))

	bw.Write([]byte("xy"))
	bw.Flush()
	br.Next()
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
		assert.True(b, br.Next())
		content, err := br.Read()
		if i < 100 || i > b.N-100 {
			assert.Equal(b, _testContent, content)
			assert.Nil(b, err)
		}
	}
	assert.False(b, br.Next())
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

	for br.Next() {
		br.Read()
	}
	assert.Equal(t, (len(_testContent)+1)*100000, int(br.Count()))

	err := br.Reset("new" + _testFile)
	assert.Nil(t, err)
	assert.Equal(t, 0, int(br.Count()))

	assert.Nil(t, br.Close())
}

func TestBufioReader_Size(t *testing.T) {
	defer os.Remove(_testFile)
	bw, _ := NewBufioWriter(_testFile)
	br, _ := NewBufioReader(_testFile)

	var (
		s0     []byte      // 0 + 1
		s64    [64]byte    // 64 + 1
		s128   [128]byte   // 128 + 2
		s16383 [16383]byte // 16383 + 2
		s16384 [16384]byte // 16384 + 3
	)

	bw.Write(s0)
	bw.Flush()
	size, err := br.Size()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), size)

	bw.Write(s64[:])
	bw.Flush()
	size, _ = br.Size()
	assert.Equal(t, int64(1+65), size)

	bw.Write(s128[:])
	bw.Flush()
	size, _ = br.Size()
	assert.Equal(t, int64(1+65+130), size)

	bw.Write(s16383[:])
	bw.Flush()
	size, _ = br.Size()
	assert.Equal(t, int64(1+65+130+16385), size)

	bw.Write(s16384[:])
	bw.Flush()
	size, _ = br.Size()
	assert.Equal(t, int64(1+65+130+16385+16387), size)
}
