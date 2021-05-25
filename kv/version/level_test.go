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

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_File_Level(t *testing.T) {
	level := newLevel()

	level.addFile(NewFileMeta(1, 1, 10, 1024))
	level.addFile(NewFileMeta(1, 1, 10, 1024))
	level.addFile(NewFileMeta(1, 1, 10, 1024))

	var files = level.getFiles()
	assert.Equal(t, 1, len(files), "add file wrong")
	assert.Equal(t, 1, level.numberOfFiles())

	//add file
	level.addFile(NewFileMeta(2, 1, 10, 1024))
	level.addFile(NewFileMeta(20, 1, 10, 1024))

	//delete file
	level.deleteFile(2)

	files = level.getFiles()
	assert.Equal(t, 2, len(files), "delete file wrong")
	assert.Equal(t, 2, level.numberOfFiles())
}

func Test_Add_Files(t *testing.T) {
	level := newLevel()

	level.addFiles(NewFileMeta(1, 1, 10, 1024), NewFileMeta(2, 1, 10, 1024), NewFileMeta(3, 1, 10, 1024))

	var files = level.getFiles()

	assert.Equal(t, 3, len(files), "add files wrong")
}
