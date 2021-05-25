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
	"bufio"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// EncodeToml encodes data into file using toml format,
// encode data to tmp file, if success then rename tmp => target file
func EncodeToml(fileName string, v interface{}) error {
	tmp := fmt.Sprintf("%s.tmp", fileName)
	f, _ := os.Create(tmp)
	defer func() {
		_ = f.Close()
	}()
	w := bufio.NewWriter(f)
	// write info using toml format
	if err := toml.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	if err := os.Rename(tmp, fileName); err != nil {
		return fmt.Errorf("rename tmp file[%s] name error:%s", tmp, err)
	}
	return nil
}

func WriteConfig(fileName string, content string) error {
	tmp := fmt.Sprintf("%s.tmp", fileName)
	f, _ := os.Create(tmp)
	defer func() {
		_ = f.Close()
	}()
	w := bufio.NewWriter(f)
	if _, err := w.WriteString(content); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	if err := os.Rename(tmp, fileName); err != nil {
		return fmt.Errorf("rename tmp file[%s] name error:%s", tmp, err)
	}
	return nil
}

// DecodeToml decodes data from file using toml format
func DecodeToml(fileName string, v interface{}) error {
	if _, err := toml.DecodeFile(fileName, v); err != nil {
		return err
	}
	return nil
}
