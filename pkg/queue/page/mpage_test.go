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

package page

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestMappedPage_err(t *testing.T) {
	defer func() {
		mapFileFunc = fileutil.RWMap
	}()

	mapFileFunc = func(file *os.File, size int) ([]byte, error) {
		return nil, fmt.Errorf("err")
	}
	mp, err := NewMappedPage(filepath.Join(t.TempDir(), t.Name()), 128)
	assert.Error(t, err)
	assert.Nil(t, mp)
}

func TestMappedPage(t *testing.T) {
	bytes := []byte("12345")

	tmpDir := t.TempDir()
	mp, err := NewMappedPage(filepath.Join(tmpDir, t.Name()), 128)
	assert.NoError(t, err)

	// copy data
	mp.WriteBytes(bytes, 0)

	assert.NoError(t, mp.Sync())
	assert.Equal(t, filepath.Join(tmpDir, t.Name()), mp.FilePath())
	assert.NotNil(t, 128, mp.Size())
	assert.Equal(t, bytes, mp.ReadBytes(0, 5))
	assert.False(t, mp.Closed())
	assert.NoError(t, mp.Close())
	assert.True(t, mp.Closed())
	assert.NoError(t, mp.Close())
}

func TestMappedPage_Write_number(t *testing.T) {
	mp, err := NewMappedPage(filepath.Join(t.TempDir(), t.Name()), 128)
	assert.NoError(t, err)
	mp.PutUint32(10, 0)
	mp.PutUint64(999, 8)
	mp.PutUint8(50, 16)
	assert.Equal(t, uint32(999), mp.ReadUint32(8))
	assert.Equal(t, uint64(10), mp.ReadUint64(0))
	assert.Equal(t, uint8(50), mp.ReadUint8(16))

	err = mp.Close()
	assert.NoError(t, err)
}
