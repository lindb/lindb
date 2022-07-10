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
)

type TestCfg struct {
	Path string `toml:"path"`
}

func TestLoadConfig(t *testing.T) {
	cfgFile := filepath.Join(t.TempDir(), "cfg")
	assert.NotNil(t, LoadConfig(cfgFile, cfgFile, &TestCfg{}))

	f, err := os.Create(cfgFile)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	_, _ = f.WriteString("Hello World")
	assert.NotNil(t, LoadConfig(cfgFile, cfgFile, &TestCfg{}))
	_ = f.Close()

	_ = EncodeToml(cfgFile, &TestCfg{Path: "/data/path"})
	cfg := TestCfg{}
	err = LoadConfig(cfgFile, cfgFile, &cfg)
	assert.NoError(t, err)
	assert.Equal(t, TestCfg{Path: "/data/path"}, cfg)

	err = LoadConfig("", cfgFile, &cfg)
	assert.NoError(t, err)
	assert.Equal(t, TestCfg{Path: "/data/path"}, cfg)
}
