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

package memdb

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestDataPointBuffer_New_err(t *testing.T) {
	defer func() {
		mkdirFunc = fileutil.MkDirIfNotExist
	}()
	mkdirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	buf, err := newDataPointBuffer(t.TempDir())
	assert.Error(t, err)
	assert.Nil(t, buf)
}

func TestDataPointBuffer_AllocPage(t *testing.T) {
	path := "buf_alloc_test"
	defer func() {
		assert.NoError(t, fileutil.RemoveDir(path))
	}()
	buf, err := newDataPointBuffer(path)
	assert.NoError(t, err)
	for i := 0; i < 10000; i++ {
		var b []byte
		b, err = buf.AllocPage()
		assert.NoError(t, err)
		assert.NotNil(t, b)
	}
	assert.NoError(t, buf.Close())
	assert.False(t, buf.IsDirty())
	buf.Release()
	assert.True(t, buf.IsDirty())
	assert.NoError(t, buf.Close())
}

func TestDataPointBuffer_AllocPage_err(t *testing.T) {
	defer func() {
		mkdirFunc = fileutil.MkDirIfNotExist
		mapFunc = fileutil.RWMap
		openFileFunc = os.OpenFile
	}()
	buf, err := newDataPointBuffer(t.TempDir())
	assert.NoError(t, err)
	mkdirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	// case 1: make file path err
	b, err := buf.AllocPage()
	assert.Error(t, err)
	assert.Nil(t, b)
	mkdirFunc = fileutil.MkDirIfNotExist

	// case 1: open file err
	buf, err = newDataPointBuffer(t.TempDir())
	assert.NoError(t, err)
	openFileFunc = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return nil, fmt.Errorf("err")
	}
	b, err = buf.AllocPage()
	assert.Error(t, err)
	assert.Nil(t, b)
	openFileFunc = os.OpenFile

	// case 2: wrong region
	b, err = buf.AllocPage()
	assert.Error(t, err)
	assert.Nil(t, b)
	buf.Release()
	err = buf.Close()
	assert.NoError(t, err)

	mapFunc = func(file *os.File, size int) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	// case 3: map file err
	buf, err = newDataPointBuffer(t.TempDir())
	assert.NoError(t, err)
	b, err = buf.AllocPage()
	assert.Error(t, err)
	assert.Nil(t, b)
	buf.Release()
	err = buf.Close()
	assert.NoError(t, err)
}

func TestDataPointBuffer_Close_err(t *testing.T) {
	path := "buf_close_err_test"
	defer func() {
		removeFunc = fileutil.RemoveDir
		assert.NoError(t, fileutil.RemoveDir(path))
	}()
	buf, err := newDataPointBuffer(path)
	assert.NoError(t, err)
	b, err := buf.AllocPage()
	assert.NoError(t, err)
	assert.NotNil(t, b)
	buf.Release()
	// case 1: remove dir err
	removeFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	assert.NoError(t, buf.Close())

	// case 2: unmap err
	buf, err = newDataPointBuffer(path)
	assert.NoError(t, err)
	b, err = buf.AllocPage()
	assert.NoError(t, err)
	assert.NotNil(t, b)
	buf.Release()
	removeFunc = fileutil.RemoveDir
	assert.NoError(t, buf.Close())
}
