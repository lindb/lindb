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

package fileutil

import (
	"os"
	"path/filepath"
)

var (
	mkdirAllFunc  = os.MkdirAll
	removeAllFunc = os.RemoveAll
	removeFunc    = os.Remove
)

// MkDirIfNotExist creates given dir if it not exist
func MkDirIfNotExist(path string) error {
	if !Exist(path) {
		if e := mkdirAllFunc(path, os.ModePerm); e != nil {
			return e
		}
	}
	return nil
}

// RemoveDir deletes dir include children if exist
func RemoveDir(path string) error {
	if Exist(path) {
		if e := removeAllFunc(path); e != nil {
			return e
		}
	}
	return nil
}

// RemoveFile removes the file if exist
func RemoveFile(file string) error {
	if Exist(file) {
		if e := removeFunc(file); e != nil {
			return e
		}
	}
	return nil
}

// MkDir create dir
func MkDir(path string) error {
	if e := mkdirAllFunc(path, os.ModePerm); e != nil {
		return e
	}
	return nil
}

// ListDir reads the directory named by dirname and returns a list of filename.
func ListDir(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, file := range files {
		result = append(result, file.Name())
	}
	return result, nil
}

// Exist check file or dir if exist
func Exist(file string) bool {
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// GetExistPath get exist path based on given path
func GetExistPath(path string) string {
	if Exist(path) {
		return path
	}
	dir, _ := filepath.Split(path)
	length := len(dir)
	if length == 0 {
		return dir
	}
	if length > 0 && os.IsPathSeparator(dir[length-1]) {
		dir = dir[:length-1]
	}
	return GetExistPath(dir)
}
