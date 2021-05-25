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

import "github.com/lindb/lindb/kv/table"

// level stores sst files of level
type level struct {
	files map[table.FileNumber]*FileMeta
}

// newLevel new level instance
func newLevel() *level {
	return &level{
		files: make(map[table.FileNumber]*FileMeta),
	}
}

// addFile adds new file into file list
func (l *level) addFile(file *FileMeta) {
	l.files[file.GetFileNumber()] = file
}

// addFiles adds new files into file list
func (l *level) addFiles(files ...*FileMeta) {
	for _, file := range files {
		l.addFile(file)
	}
}

// deleteFile removes file from file list
func (l *level) deleteFile(fileNumber table.FileNumber) {
	delete(l.files, fileNumber)
}

// getFiles returns all files in current level
func (l *level) getFiles() []*FileMeta {
	var values []*FileMeta
	for _, v := range l.files {
		values = append(values, v)
	}
	return values
}

// numberOfFiles returns the number of files in current level
func (l *level) numberOfFiles() int {
	return len(l.files)
}
