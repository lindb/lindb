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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_testFile = "bufioio.test"
)

var (
	_testContent = []byte("lindb.ci.metric")
)

func TestNewBufioEntryWriter(t *testing.T) {
	defer func() {
		_ = os.Remove(_testFile)
	}()
	bw, err := NewBufioEntryWriter(_testFile)

	assert.Nil(t, err)
	assert.NotNil(t, bw)
	err = bw.Close()
	assert.NoError(t, err)
}

func TestBufioWriter_Reset(t *testing.T) {
	defer func() {
		_ = os.Remove(_testFile)
		_ = os.Remove("new" + _testFile)
	}()

	bw, _ := NewBufioEntryWriter(_testFile)
	_, _ = bw.Write([]byte("test"))
	_ = bw.Flush()
	assert.Equal(t, int64(5), bw.Size())

	err := bw.Reset("new" + _testFile)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), bw.Size())
	_, err = bw.Write([]byte("abcd"))
	assert.NoError(t, err)

	stat, _ := os.Stat(_testFile)
	assert.Equal(t, int64(5), stat.Size())

	err = bw.Close()
	assert.NoError(t, err)
}

func TestBufioWriter_Write_Size(t *testing.T) {
	defer func() {
		_ = os.Remove(_testFile)
	}()
	bw, _ := NewBufioEntryWriter(_testFile)
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

	err = bw.Close()
	assert.NoError(t, err)
}

func TestBufioStreamWriter_Write_Size(t *testing.T) {
	defer func() {
		_ = os.Remove(_testFile)
	}()
	bw, _ := NewBufioStreamWriter(_testFile)
	assert.Equal(t, int64(0), bw.Size())
	n, err := bw.Write([]byte(""))
	assert.Equal(t, 0, n)
	assert.Equal(t, int64(0), bw.Size())
	assert.Nil(t, err)

	n, _ = bw.Write([]byte("xyz"))
	assert.Equal(t, 3, n)
	assert.Equal(t, int64(3), bw.Size())

	var s [128]byte
	n, _ = bw.Write(s[:])
	assert.Equal(t, 128, n)
	assert.Equal(t, int64(131), bw.Size())
	err = bw.Close()
	assert.NoError(t, err)
}

func BenchmarkBufioWriter_Write(b *testing.B) {
	defer func() {
		_ = os.Remove(_testFile)
	}()
	bw, _ := NewBufioEntryWriter(_testFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := bw.Write(_testContent); err != nil {
			assert.Nil(b, err)
		}
	}
	err := bw.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func TestBufioWriter_Close(t *testing.T) {
	defer func() {
		_ = os.Remove(_testFile)
	}()
	bw, _ := NewBufioEntryWriter(_testFile)

	expectedLength := (len(_testContent) + 1) * 100000
	for i := 0; i < 100000; i++ {
		_, _ = bw.Write(_testContent)
	}
	_ = bw.Sync()
	assert.Nil(t, bw.Sync())
	assert.Nil(t, bw.Close())

	data, err := os.ReadFile(_testFile)
	assert.Nil(t, err)
	assert.Len(t, data, expectedLength)
}
