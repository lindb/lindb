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

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageState(t *testing.T) {
	storageState := NewStorageState("test")
	storageState.NodeOnline(StatefulNode{
		StatelessNode: StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000},
		ID:            1,
	})
	storageState.NodeOnline(StatefulNode{
		StatelessNode: StatelessNode{HostIP: "1.1.1.2", GRPCPort: 9000},
		ID:            2,
	})
	storageState.NodeOnline(StatefulNode{
		StatelessNode: StatelessNode{HostIP: "1.1.1.3", GRPCPort: 9000},
		ID:            3,
	})
	storageState.NodeOffline(2)
	assert.Len(t, storageState.LiveNodes, 2)
	storageState.ShardAssignments["test"] = &ShardAssignment{
		Name:   "test",
		Shards: map[ShardID]*Replica{1: {Replicas: []NodeID{1, 2, 3}}},
	}
	rs := storageState.ReplicasOnNode(3)
	assert.Len(t, rs, 1)
	assert.Equal(t, rs["test"], []ShardID{1})

	storageState.ShardStates["test"] = map[ShardID]ShardState{1: {
		ID:     1,
		Leader: 2,
	}}
	rs1 := storageState.LeadersOnNode(2)
	assert.Len(t, rs1, 1)
	assert.Equal(t, rs1["test"], []ShardID{1})
}
