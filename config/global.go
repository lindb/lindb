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

import "sync/atomic"

var (
	// StandaloneMode represents LinDB run as standalone mode
	StandaloneMode = false

	globalBrokerCfg  atomic.Value
	globalStorageCfg atomic.Value
)

func init() {
	SetGlobalBrokerConfig(NewDefaultBrokerBase())
	SetGlobalStorageConfig(NewDefaultStorageBase())
}

// SetGlobalBrokerConfig sets global the broker config,
// this config will be triggered to reload when receiving a SIGHUP
func SetGlobalBrokerConfig(cfg *BrokerBase) {
	globalBrokerCfg.Store(cfg)
}

// SetGlobalStorageConfig sets global the storage config,
// this config will be triggered to reload when receiving a SIGHUP
func SetGlobalStorageConfig(cfg *StorageBase) {
	globalStorageCfg.Store(cfg)
}

// GlobalBrokerConfig returns the global broker config
func GlobalBrokerConfig() *BrokerBase {
	return globalBrokerCfg.Load().(*BrokerBase)
}

// GlobalStorageConfig returns the global storage config
func GlobalStorageConfig() *StorageBase {
	return globalStorageCfg.Load().(*StorageBase)
}
