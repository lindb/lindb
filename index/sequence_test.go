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

package index

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestSequence(t *testing.T) {
	name := "./sequence"
	defer func() {
		_ = os.RemoveAll(name)
	}()

	seq, err := NewSequence(name)
	assert.NotNil(t, seq)
	assert.NoError(t, err)

	test := func(s, c uint32, get func() uint32) {
		for ; s < c; s++ {
			assert.Equal(t, s, get())
		}
	}
	test(0, 4, seq.GenNamespaceSeq)
	test(0, 8, seq.GenMetricNameSeq)
	test(0, 16, seq.GenTagKeySeq)
	test(0, 32, seq.GenTagValueSeq)

	assert.NoError(t, seq.Sync())
	assert.NoError(t, seq.Close())

	seq, err = NewSequence(name)
	assert.NotNil(t, seq)
	assert.NoError(t, err)
	assert.Equal(t, uint32(4), seq.GenNamespaceSeq())
	assert.Equal(t, uint32(5), seq.GetNamespaceSeq())
	assert.Equal(t, uint32(8), seq.GenMetricNameSeq())
	assert.Equal(t, uint32(9), seq.GetMetricNameSeq())
	assert.Equal(t, uint32(16), seq.GenTagKeySeq())
	assert.Equal(t, uint32(32), seq.GenTagValueSeq())
	assert.NoError(t, seq.Close())
}

func TestSequence_New_Error(t *testing.T) {
	name := "./sequence_new_error"
	defer func() {
		openFileFn = os.OpenFile
		rwMapFn = fileutil.RWMap
		_ = os.RemoveAll(name)
	}()

	t.Run("create file error", func(t *testing.T) {
		openFileFn = func(_ string, _ int, _ os.FileMode) (*os.File, error) {
			return nil, fmt.Errorf("err")
		}
		seq, err := NewSequence(name)
		assert.Error(t, err)
		assert.Nil(t, seq)
	})

	t.Run("map file error", func(t *testing.T) {
		openFileFn = os.OpenFile
		rwMapFn = func(_ *os.File, _ int) (data []byte, err error) {
			return nil, fmt.Errorf("err")
		}
		seq, err := NewSequence(name)
		assert.Error(t, err)
		assert.Nil(t, seq)
	})
}
