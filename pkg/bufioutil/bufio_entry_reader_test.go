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

func Test_NewBufioEntryReader(t *testing.T) {
	defer func() {
		_ = os.Remove(_testFile)
	}()
	br, err := NewBufioEntryReader(_testFile)
	assert.NotNil(t, err)
	assert.Nil(t, br)

	_, _ = os.Create(_testFile)
	br, err = NewBufioEntryReader(_testFile)
	assert.Nil(t, err)
	assert.NotNil(t, br)
}

func TestBufioReader_content(t *testing.T) {
	defer func() {
		_ = os.Remove(_testFile)
	}()
	bw, _ := NewBufioEntryWriter(_testFile)

	f, _ := os.Open(_testFile)
	br := bufioEntryReader{
		f: f,
		r: bufio.NewReader(f)}

	_, _ = bw.Write([]byte("a"))
	_ = bw.Flush()
	br.Next()
	assert.Equal(t, 1, len(br.content))
	assert.Equal(t, 1, cap(br.content))

	_, _ = bw.Write([]byte("abcde"))
	_ = bw.Flush()
	br.Next()
	assert.Equal(t, 5, len(br.content))
	assert.Equal(t, 5, cap(br.content))

	_, _ = bw.Write([]byte("xy"))
	_ = bw.Flush()
	br.Next()
	assert.Equal(t, 2, len(br.content))
	assert.Equal(t, 5, cap(br.content))
}

func BenchmarkBufioReader_Read(b *testing.B) {
	defer func() {
		_ = os.Remove(_testFile)
	}()
	bw, _ := NewBufioEntryWriter(_testFile)
	br, _ := NewBufioEntryReader(_testFile)

	for i := 0; i < b.N; i++ {
		_, _ = bw.Write(_testContent)
	}
	_ = bw.Sync()
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
	defer func() {
		_ = os.Remove(_testFile)
	}()
	defer func() {
		_ = os.Remove("new" + _testFile)
	}()
	_, _ = os.Create("new" + _testFile)
	bw, _ := NewBufioEntryWriter(_testFile)
	br, _ := NewBufioEntryReader(_testFile)

	for i := 0; i < 100000; i++ {
		_, _ = bw.Write(_testContent)
	}
	_ = bw.Sync()

	for br.Next() {
		_, _ = br.Read()
	}
	assert.Equal(t, (len(_testContent)+1)*100000, int(br.Count()))

	err := br.Reset("new" + _testFile)
	assert.Nil(t, err)
	assert.Equal(t, 0, int(br.Count()))

	assert.Nil(t, br.Close())
}

func TestBufioReader_Size(t *testing.T) {
	defer func() {
		_ = os.Remove(_testFile)
	}()
	bw, _ := NewBufioEntryWriter(_testFile)
	br, _ := NewBufioEntryReader(_testFile)

	var (
		s0     []byte      // 0 + 1
		s64    [64]byte    // 64 + 1
		s128   [128]byte   // 128 + 2
		s16383 [16383]byte // 16383 + 2
		s16384 [16384]byte // 16384 + 3
	)

	_, _ = bw.Write(s0)
	_ = bw.Flush()
	size, err := br.Size()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), size)

	_, _ = bw.Write(s64[:])
	_ = bw.Flush()
	size, _ = br.Size()
	assert.Equal(t, int64(1+65), size)

	_, _ = bw.Write(s128[:])
	_ = bw.Flush()
	size, _ = br.Size()
	assert.Equal(t, int64(1+65+130), size)

	_, _ = bw.Write(s16383[:])
	_ = bw.Flush()
	size, _ = br.Size()
	assert.Equal(t, int64(1+65+130+16385), size)

	_, _ = bw.Write(s16384[:])
	_ = bw.Flush()
	size, _ = br.Size()
	assert.Equal(t, int64(1+65+130+16385+16387), size)
}
