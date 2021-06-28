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

	"github.com/lindb/lindb/coordinator/task"
)

// defines common constants will be used in broker and storage
const (
	// ActiveNodesPath represents active nodes prefix path for node register
	ActiveNodesPath = "/active/nodes"
	// StateNodesPath represents the state of node that node will report runtime status
	StateNodesPath = "/state/nodes"
)

// defines storage level constants will be used in storage
//const ()

// defines broker level constants will be used in broker
const (
	// MasterPath represents master elect path
	MasterPath = "/master/node"
	// DatabaseAssignPath represents database shard assignment
	DatabaseAssignPath = "/database/assign"
	// StorageClusterConfigPath represents cluster config store
	StorageClusterConfigPath = "/storage/cluster/config"
	// DatabaseConfigPath represents database config path
	DatabaseConfigPath = "/database/config"

	// StorageClusterNodeStatePath represents storage cluster's node state
	StorageClusterNodeStatePath = "/state/storage/nodes/cluster"
	// ReplicaStatePath represents the replica's state
	ReplicaStatePath = "/state/replica"
	// StorageClusterStatPath represents storage cluster's node monitoring stat
	StorageClusterStatPath = "/state/storage/stat/cluster"
)

// defines all task kinds
const (
	// CreateShard represents task kind which is create shard for storage node
	CreateShard task.Kind = "create-shard"
	// FlushDatabase represents task kind which is flush memory database for storage node
	FlushDatabase task.Kind = "flush-database"
)

// GetStorageClusterConfigPath returns path which storing config of storage cluster
func GetStorageClusterConfigPath(name string) string {
	return fmt.Sprintf("%s/%s", StorageClusterConfigPath, name)
}

// GetStorageClusterNodeStatePath returns path whine storing state of storage cluster
func GetStorageClusterNodeStatePath(name string) string {
	return fmt.Sprintf("%s/%s", StorageClusterNodeStatePath, name)
}

// GetStorageClusterStatPath returns path whine storing monitoring stat of storage cluster
func GetStorageClusterStatPath(name string) string {
	return fmt.Sprintf("%s/%s", StorageClusterStatPath, name)
}

// GetDatabaseConfigPath returns path which storing config of database
func GetDatabaseConfigPath(name string) string {
	return fmt.Sprintf("%s/%s", DatabaseConfigPath, name)
}

// GetDatabaseAssignPath returns path which storing shard assignment of database
func GetDatabaseAssignPath(name string) string {
	return fmt.Sprintf("%s/%s", DatabaseAssignPath, name)
}

// GetNodePath returns node register path
func GetNodePath(prefix, node string) string {
	return fmt.Sprintf("%s/data/%s", prefix, node)
}

// GetNodeIDPath returns node id register path
func GetNodeIDPath(prefix, node string) string {
	return fmt.Sprintf("%s/ids/%s", prefix, node)
}

// GetNodeSeqPath returns node id's generate path
func GetNodeSeqPath(prefix string) string {
	return fmt.Sprintf("%s/seq", prefix)
}

// GetReplicaStatePath returns replica's state path
func GetReplicaStatePath(node string) string {
	return fmt.Sprintf("%s/%s", ReplicaStatePath, node)
}

// GetNodeMonitoringStatPath returns the node monitoring stat's path
func GetNodeMonitoringStatPath(node string) string {
	return fmt.Sprintf("%s/%s", StateNodesPath, node)
}
