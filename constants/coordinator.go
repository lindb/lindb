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

package constants

import (
	"fmt"
)

// StatePathSeparator represents the separator of state store's path.
const StatePathSeparator = "/"

// defines the role type of node.
const (
	MasterRole  = "Master"
	RootRole    = "Root"
	BrokerRole  = "Broker"
	StorageRole = "Storage"
)

// defines all metadata type.
const (
	LiveNode        = "LiveNode"
	DatabaseConfig  = "DatabaseConfig"
	BrokerState     = "BrokerState"
	StorageState    = "StorageState"
	ShardAssignment = "ShardAssignment"
	Master          = "Master"
	StorageConfig   = "StorageConfig"
)

// defines common constants will be used in broker and storage.
const (
	// LiveNodesPath represents live nodes prefix path for node register.
	LiveNodesPath = "/live/nodes"
)

// defines broker level constants will be used in broker.
const (
	// MasterPath represents master elect path.
	MasterPath = "/master/node"
	// MasterElectedPath represents register path after master finished election.
	MasterElectedPath = "/master/elected"
	// DatabaseConfigPath represents database config path.
	DatabaseConfigPath = "/database/config"
	// DatabaseLimitPath represents database limit path.
	DatabaseLimitPath = "/database/limit"
	// ShardAssignmentPath represents database shard assignment.
	ShardAssignmentPath = "/database/assign"
	// StorageConfigPath represents storage cluster's config.
	StorageLiveNodesPath = "/storage/live/nodes"
	// StorageStatePath represents storage cluster's state.
	StorageStatePath = "/storage/state"
	// BrokerConfigPath represents broker cluster's config.
	BrokerConfigPath = "/broker/config"
)

// GetBrokerClusterConfigPath returns path which storing config of broker cluster.
func GetBrokerClusterConfigPath(name string) string {
	return fmt.Sprintf("%s/%s", BrokerConfigPath, name)
}

// GetDatabaseConfigPath returns path which storing config of database
func GetDatabaseConfigPath(name string) string {
	return fmt.Sprintf("%s/%s", DatabaseConfigPath, name)
}

// GetDatabaseLimitPath returns path which storing limit of database
func GetDatabaseLimitPath(name string) string {
	return fmt.Sprintf("%s/%s", DatabaseLimitPath, name)
}

// GetDatabaseAssignPath returns path which storing shard assignment of database
func GetDatabaseAssignPath(name string) string {
	return fmt.Sprintf("%s/%s", ShardAssignmentPath, name)
}

// GetShardAssignPath returns path which storing shard assignment of database.
func GetShardAssignPath(name string) string {
	return fmt.Sprintf("%s/%s", ShardAssignmentPath, name)
}

// GetStorageLiveNodePath returns live node register path for storage.
func GetStorageLiveNodePath(node string) string {
	return fmt.Sprintf("%s/%s", StorageLiveNodesPath, node)
}

// GetLiveNodePath returns live node register path for broker/root.
func GetLiveNodePath(node string) string {
	return fmt.Sprintf("%s/%s", LiveNodesPath, node)
}
