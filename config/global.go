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
	"fmt"
	"sync/atomic"

	"github.com/lindb/lindb/pkg/ltoml"
)

var (
	// StandaloneMode represents LinDB run as standalone mode
	StandaloneMode = false

	globalBrokerCfg  atomic.Value
	globalStorageCfg atomic.Value
)

func init() {
	globalBrokerCfg.Store(NewDefaultBrokerBase())
	globalStorageCfg.Store(NewDefaultStorageBase())
}

// GlobalBrokerConfig returns the global broker config
func GlobalBrokerConfig() *BrokerBase {
	return globalBrokerCfg.Load().(*BrokerBase)
}

// SetGlobalBrokerConfig sets global broker configuration.
func SetGlobalBrokerConfig(brokerCfg *BrokerBase) {
	globalBrokerCfg.Store(brokerCfg)
}

// GlobalStorageConfig returns the global storage config
func GlobalStorageConfig() *StorageBase {
	return globalStorageCfg.Load().(*StorageBase)
}

// SetGlobalStorageConfig sets global storage configuration.
func SetGlobalStorageConfig(storageCfg *StorageBase) {
	globalStorageCfg.Store(storageCfg)
}

// LoadAndSetBrokerConfig parses the broker config file
// this config will be triggered to reload when receiving a SIGHUP signal
func LoadAndSetBrokerConfig(cfgName, defaultPath string, brokerCfg *Broker) error {
	if err := ltoml.LoadConfig(cfgName, defaultPath, &brokerCfg); err != nil {
		return fmt.Errorf("decode broker config file error: %s", err)
	}
	checkQueryCfg(&brokerCfg.Query)
	if err := checkCoordinatorCfg(&brokerCfg.Coordinator); err != nil {
		return fmt.Errorf("failed check coordinator config: %s", err)
	}
	if err := checkBrokerBaseCfg(&brokerCfg.BrokerBase); err != nil {
		return fmt.Errorf("failed checking broker config: %s", err)
	}
	globalBrokerCfg.Store(&brokerCfg.BrokerBase)
	return nil
}

// LoadAndSetStorageConfig parses the storage config file
// this config will be triggered to reload when receiving a SIGHUP signal
func LoadAndSetStorageConfig(cfgName, defaultPath string, storageCfg *Storage) error {
	if err := ltoml.LoadConfig(cfgName, defaultPath, &storageCfg); err != nil {
		return fmt.Errorf("decode storage config file error: %s", err)
	}
	checkQueryCfg(&storageCfg.Query)
	if err := checkCoordinatorCfg(&storageCfg.Coordinator); err != nil {
		return fmt.Errorf("failed check coordinator config: %s", err)
	}
	if err := checkStorageBaseCfg(&storageCfg.StorageBase); err != nil {
		return fmt.Errorf("failed checking storage config: %s", err)
	}
	globalStorageCfg.Store(&storageCfg.StorageBase)
	return nil
}

// LoadAndSetStandAloneConfig parses the standalone config file
// then sets the global broker and storage config
// this config will be triggered to reload when receiving a SIGHUP signal
func LoadAndSetStandAloneConfig(cfgName, defaultPath string, standaloneCfg *Standalone) error {
	if err := ltoml.LoadConfig(cfgName, defaultPath, &standaloneCfg); err != nil {
		return fmt.Errorf("decode standalone config file error: %s", err)
	}
	checkQueryCfg(&standaloneCfg.Query)
	if err := checkCoordinatorCfg(&standaloneCfg.Coordinator); err != nil {
		return fmt.Errorf("failed check coordinator config: %s", err)
	}
	if err := checkBrokerBaseCfg(&standaloneCfg.BrokerBase); err != nil {
		return fmt.Errorf("failed checking broker config: %s", err)
	}
	if err := checkStorageBaseCfg(&standaloneCfg.StorageBase); err != nil {
		return fmt.Errorf("failed checking storage config: %s", err)
	}
	globalBrokerCfg.Store(&standaloneCfg.BrokerBase)
	globalStorageCfg.Store(&standaloneCfg.StorageBase)
	return nil
}
