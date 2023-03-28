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

package ltoml

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

type User struct {
	Name string
}

func Test_Encode(t *testing.T) {
	testPath := t.TempDir()
	user := User{Name: "LinDB"}
	file := filepath.Join(testPath, "toml")
	err := EncodeToml(file, &user)
	if err != nil {
		t.Fatal(err)
	}
	user2 := User{}
	err = DecodeToml(file, &user2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, user, user2)

	files, _ := fileutil.ListDir(testPath)
	assert.Equal(t, "toml", files[0])

	assert.NotNil(t, EncodeToml(filepath.Join(os.TempDir(), "tmp", "test.toml"), []byte{}))
}

func Test_WriteConfig(t *testing.T) {
	testPath := t.TempDir()
	assert.Nil(t, WriteConfig(filepath.Join(testPath, "toml"), ""))
}
