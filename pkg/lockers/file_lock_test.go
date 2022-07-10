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

package lockers

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLock(t *testing.T) {
	defer func() {
		openFileFn = os.OpenFile
	}()
	lock, err := NewFileLock("t.lock")
	assert.NoError(t, err)
	assert.NotNil(t, lock)
	err = lock.Lock()
	assert.Nil(t, err, "lock error")

	fileInfo, _ := os.Stat("t.lock")
	assert.NotNil(t, fileInfo, "lock file not exist")

	err = lock.Unlock()
	assert.NoError(t, err)

	fileInfo, _ = os.Stat("t.lock")
	assert.Nil(t, fileInfo, "lock file exist")

	openFileFn = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return nil, fmt.Errorf("err")
	}
	lock, err = NewFileLock("/tmp/not_dir/t.lock")
	assert.Error(t, err)
	assert.Nil(t, lock)
}
