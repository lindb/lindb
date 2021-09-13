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

func Test_NewConfig(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testPath)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	// validate broker config
	brokerCfgPath := filepath.Join(testPath, "broker.toml")
	var brokerCfg Broker
	// not-exist
	assert.NotNil(t, LoadAndSetBrokerConfig("not-exist-path", "broker.toml", &brokerCfg))
	// bad broker config
	assert.Nil(t, ltoml.WriteConfig(brokerCfgPath, ""))
	assert.Error(t, LoadAndSetBrokerConfig(brokerCfgPath, "broker.toml", &brokerCfg))

	// ok
	assert.Nil(t, ltoml.WriteConfig(brokerCfgPath, NewDefaultBrokerTOML()))
	assert.Nil(t, LoadAndSetBrokerConfig(brokerCfgPath, "broker.toml", &brokerCfg))
	assert.Nil(t, ltoml.DecodeToml(brokerCfgPath, &brokerCfg))
	assert.Equal(t, brokerCfg.BrokerBase, *NewDefaultBrokerBase())
	assert.Equal(t, brokerCfg.Logging, *NewDefaultLogging())
	assert.Equal(t, brokerCfg.Monitor, *NewDefaultMonitor())

	// validate storage config
	storageCfgPath := filepath.Join(testPath, "storage.toml")
	var storageCfg Storage
	// not exist
	assert.Error(t, LoadAndSetStorageConfig("not-exist-path", "storage.toml", &storageCfg))
	// bad storage config
	assert.Nil(t, ltoml.WriteConfig(storageCfgPath, ""))
	assert.Error(t, LoadAndSetStorageConfig(storageCfgPath, "storage.toml", &storageCfg))

	// ok
	assert.Nil(t, ltoml.WriteConfig(storageCfgPath, NewDefaultStorageTOML()))
	assert.Nil(t, LoadAndSetStorageConfig(storageCfgPath, "storage.toml", &storageCfg))
	assert.Nil(t, ltoml.DecodeToml(storageCfgPath, &storageCfg))
	assert.Equal(t, storageCfg.StorageBase, *NewDefaultStorageBase())
	assert.Equal(t, storageCfg.Logging, *NewDefaultLogging())
	assert.Equal(t, storageCfg.Monitor, *NewDefaultMonitor())

	// validate standalone config
	standaloneCfgPath := filepath.Join(testPath, "standalone.toml")
	var standaloneCfg Standalone
	// not-exist
	assert.Error(t, LoadAndSetStandAloneConfig("not-exist-path", "standalone.toml", &standaloneCfg))
	// bad broker config
	assert.Nil(t, ltoml.WriteConfig(standaloneCfgPath, ""))
	assert.Error(t, LoadAndSetStandAloneConfig(standaloneCfgPath, "standalone.toml", &standaloneCfg))

	// ok
	assert.Nil(t, ltoml.WriteConfig(standaloneCfgPath, NewDefaultStandaloneTOML()))
	assert.Nil(t, LoadAndSetStandAloneConfig(standaloneCfgPath, "standalone.toml", &standaloneCfg))
	assert.Nil(t, ltoml.DecodeToml(standaloneCfgPath, &standaloneCfg))
	assert.Equal(t, standaloneCfg.BrokerBase, *NewDefaultBrokerBase())
	assert.Equal(t, standaloneCfg.StorageBase, *NewDefaultStorageBase())
	assert.Equal(t, standaloneCfg.Logging, *NewDefaultLogging())
	assert.Equal(t, standaloneCfg.Monitor, *NewDefaultMonitor())
}

func Test_Global(t *testing.T) {
	assert.NotNil(t, GlobalBrokerConfig())
	assert.NotNil(t, GlobalStorageConfig())
}

func Test_checkBrokerBaseCfg(t *testing.T) {
	emptyBrokerBase := &BrokerBase{}
	assert.Error(t, checkBrokerBaseCfg(emptyBrokerBase))

	// grpc failure
	brokerCfg1 := &BrokerBase{}
	assert.Error(t, checkBrokerBaseCfg(brokerCfg1))

	// http port failure
	brokerCfg2 := &BrokerBase{
		GRPC: GRPC{Port: 2379},
	}
	assert.Error(t, checkBrokerBaseCfg(brokerCfg2))

	// ok
	brokerCfg3 := &BrokerBase{
		GRPC: GRPC{Port: 2379},
		HTTP: HTTP{Port: 9000},
	}
	assert.NoError(t, checkBrokerBaseCfg(brokerCfg3))
	assert.NotZero(t, brokerCfg3.HTTP.ReadTimeout)
	assert.NotZero(t, brokerCfg3.HTTP.IdleTimeout)
	assert.NotZero(t, brokerCfg3.HTTP.WriteTimeout)
	assert.NotZero(t, brokerCfg3.Ingestion.IngestTimeout)
}

func Test_checkStorageBaseCfg(t *testing.T) {
	emptyStorageBase := &StorageBase{}
	assert.Error(t, checkStorageBaseCfg(emptyStorageBase))

	// grpc failure
	storageCfg1 := &StorageBase{Indicator: 1}
	assert.Error(t, checkStorageBaseCfg(storageCfg1))

	// http port failure
	storageCfg2 := &StorageBase{
		Indicator: 1,
		GRPC:      GRPC{Port: 2379},
	}
	assert.Error(t, checkStorageBaseCfg(storageCfg2))

	// tsdb error
	storageCfg3 := &StorageBase{
		Indicator: 1,
		GRPC:      GRPC{Port: 2379},
	}
	assert.Error(t, checkStorageBaseCfg(storageCfg3))

	// ok
	storageCfg4 := &StorageBase{
		Indicator: 1,
		GRPC:      GRPC{Port: 2379},
		TSDB:      TSDB{Dir: "/tmp/lindb"},
	}
	assert.NoError(t, checkStorageBaseCfg(storageCfg4))
	assert.NotZero(t, storageCfg4.TSDB.MaxMemDBSize)
	assert.NotZero(t, storageCfg4.TSDB.MaxMemDBTotalSize)
	assert.NotZero(t, storageCfg4.TSDB.MaxMemDBNumber)
	assert.NotZero(t, storageCfg4.TSDB.MutableMemDBTTL)
	assert.NotZero(t, storageCfg4.TSDB.MaxMemUsageBeforeFlush)
	assert.NotZero(t, storageCfg4.TSDB.TargetMemUsageAfterFlush)
	assert.NotZero(t, storageCfg4.TSDB.FlushConcurrency)
	assert.NotZero(t, storageCfg4.TSDB.MaxSeriesIDsNumber)
	assert.NotZero(t, storageCfg4.TSDB.MaxTagKeysNumber)
}

func Test_checkCoordinatorCfg(t *testing.T) {
	var repo RepoState
	assert.Error(t, checkCoordinatorCfg(&repo))

	repo = RepoState{Namespace: "/1"}
	assert.Error(t, checkCoordinatorCfg(&repo))

	repo = RepoState{Namespace: "/1", Endpoints: []string{"http://localhost:2379"}}
	assert.NoError(t, checkCoordinatorCfg(&repo))

	assert.Equal(t, "/1/2", repo.WithSubNamespace("2").Namespace)
}
