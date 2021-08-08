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

package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/ltoml"
)

var testPath = "./tmp"

func TestTCP_TOML(t *testing.T) {
	tcp := &TCP{
		Port: 8080,
	}
	assert.Equal(t, "\n    port = 8080", tcp.TOML())
}

func Test_NewConfig(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testPath)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	// validate broker config
	brokerCfgPath := filepath.Join(testPath, "broker.toml")
	assert.Nil(t, ltoml.WriteConfig(brokerCfgPath, NewDefaultBrokerTOML()))
	var brokerCfg Broker
	assert.Nil(t, ltoml.DecodeToml(brokerCfgPath, &brokerCfg))
	assert.Equal(t, brokerCfg.BrokerBase, *NewDefaultBrokerBase())
	assert.Equal(t, brokerCfg.Logging, *NewDefaultLogging())
	assert.Equal(t, brokerCfg.Monitor, *NewDefaultMonitor())

	// validate storage config
	storageCfgPath := filepath.Join(testPath, "storage.toml")
	assert.Nil(t, ltoml.WriteConfig(storageCfgPath, NewDefaultStorageTOML()))
	var storageCfg Storage
	assert.Nil(t, ltoml.DecodeToml(storageCfgPath, &storageCfg))
	assert.Equal(t, storageCfg.StorageBase, *NewDefaultStorageBase())
	assert.Equal(t, storageCfg.Logging, *NewDefaultLogging())
	assert.Equal(t, storageCfg.Monitor, *NewDefaultMonitor())

	// validate standalone config
	standaloneCfgPath := filepath.Join(testPath, "standalone.toml")
	assert.Nil(t, ltoml.WriteConfig(standaloneCfgPath, NewDefaultStandaloneTOML()))
	var standaloneCfg Standalone
	assert.Nil(t, ltoml.DecodeToml(standaloneCfgPath, &standaloneCfg))
	assert.Equal(t, standaloneCfg.BrokerBase, *NewDefaultBrokerBase())
	assert.Equal(t, standaloneCfg.StorageBase, *NewDefaultStorageBase())
	assert.Equal(t, standaloneCfg.Logging, *NewDefaultLogging())
	assert.Equal(t, standaloneCfg.Monitor, *NewDefaultMonitor())
}

func Test_ReplicationChannel_SegmentFileSizeInBytes(t *testing.T) {
	var rc ReplicationChannel
	assert.Equal(t, int64(1024*1024), rc.GetDataSizeLimit())
	assert.Zero(t, rc.BufferSizeInBytes())
	rc.DataSizeLimit = 10
	assert.Equal(t, int64(10*1024*1024), rc.GetDataSizeLimit())
	rc.DataSizeLimit = 10000
	assert.Equal(t, int64(1024*1024*1024), rc.GetDataSizeLimit())
}

func Test_Global(t *testing.T) {
	assert.NotNil(t, GlobalBrokerConfig())
	assert.NotNil(t, GlobalStorageConfig())
}
