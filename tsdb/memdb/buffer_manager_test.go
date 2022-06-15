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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestBufferManager_AllocBuffer(t *testing.T) {
	defer func() {
		mkdirFunc = fileutil.MkDirIfNotExist
	}()

	mgr := NewBufferManager(t.TempDir())
	// case 1: allocate err
	mkdirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	buf, err := mgr.AllocBuffer(timeutil.Now())
	assert.Error(t, err)
	assert.Nil(t, buf)

	// case 2: allocate ok
	mkdirFunc = fileutil.MkDirIfNotExist
	buf, err = mgr.AllocBuffer(timeutil.Now())
	assert.NoError(t, err)
	assert.NotNil(t, buf)
}

func TestBufferManager_Cleanup(t *testing.T) {
	defer func() {
		removeFunc = fileutil.RemoveDir
	}()

	mgr := NewBufferManager(t.TempDir())

	// case 1: cleanup err
	removeFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	mgr.Cleanup()
	// case 2: cleanup ok
	removeFunc = fileutil.RemoveDir
	mgr.Cleanup()
}

func TestBufferManager_GarbageCollect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	mgr := NewBufferManager(t.TempDir())
	// case 1: no buf
	mgr.GarbageCollect()
	// case 2: no dirty buf
	mgr1 := mgr.(*bufferManager)
	newSet := make([]DataPointBuffer, 0)
	buf1 := NewMockDataPointBuffer(ctrl)
	buf2 := NewMockDataPointBuffer(ctrl)
	buf3 := NewMockDataPointBuffer(ctrl)
	newSet = append(newSet, []DataPointBuffer{buf1, buf2, buf3}...)
	mgr1.value.Store(newSet)
	buf1.EXPECT().IsDirty().Return(false)
	buf2.EXPECT().IsDirty().Return(false)
	buf3.EXPECT().IsDirty().Return(false)
	mgr.GarbageCollect()
	oldSet := mgr1.value.Load().([]DataPointBuffer)
	assert.Len(t, oldSet, 3)
	// case 3: gc dirty buf
	buf1.EXPECT().IsDirty().Return(false)
	buf2.EXPECT().IsDirty().Return(true) // gc
	buf2.EXPECT().Close().Return(nil)
	buf3.EXPECT().IsDirty().Return(true) // keep it
	buf3.EXPECT().Close().Return(fmt.Errorf("err"))
	mgr.GarbageCollect()
	oldSet = mgr1.value.Load().([]DataPointBuffer)
	assert.Len(t, oldSet, 2)
	assert.Same(t, oldSet[0], buf1)
	assert.Same(t, oldSet[1], buf3)
}
