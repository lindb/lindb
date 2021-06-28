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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseAssignPath(t *testing.T) {
	assert.Equal(t, DatabaseAssignPath+"/name", GetDatabaseAssignPath("name"))
}

func TestGetDatabaseConfigPath(t *testing.T) {
	assert.Equal(t, DatabaseConfigPath+"/name", GetDatabaseConfigPath("name"))
}

func TestGetNodePath(t *testing.T) {
	assert.Equal(t, "prefix/data/name", GetNodePath("prefix", "name"))
}

func TestGetStorageClusterConfigPath(t *testing.T) {
	assert.Equal(t, StorageClusterConfigPath+"/name", GetStorageClusterConfigPath("name"))

}
func TestGetStorageClusterStatePath(t *testing.T) {
	assert.Equal(t, StorageClusterNodeStatePath+"/name", GetStorageClusterNodeStatePath("name"))
}

func TestGetStorageClusterStatPath(t *testing.T) {
	assert.Equal(t, StorageClusterStatPath+"/name", GetStorageClusterStatPath("name"))
}

func TestGetReplicaStatePath(t *testing.T) {
	assert.Equal(t, ReplicaStatePath+"/1.1.1.1:port", GetReplicaStatePath("1.1.1.1:port"))
}

func TestGetNodeMonitoringStatPath(t *testing.T) {
	assert.Equal(t, StateNodesPath+"/1.1.1.1:port", GetNodeMonitoringStatPath("1.1.1.1:port"))
}

func TestGetNodeIDPath(t *testing.T) {
	assert.Equal(t, StateNodesPath+"/ids/1.1.1.1:port", GetNodeIDPath(StateNodesPath, "1.1.1.1:port"))
	assert.Equal(t, StateNodesPath+"/seq", GetNodeSeqPath(StateNodesPath))
}
