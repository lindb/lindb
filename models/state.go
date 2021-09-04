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
	"github.com/lindb/lindb/pkg/encoding"
)

type ShardStateType int

const (
	UnknownShard ShardStateType = iota
	NewShard
	OnlineShard
	OfflineShard
	NonExistentShard
)

const NoLeader NodeID = -1

// ReplicaState represents the relationship for a replica.
type ReplicaState struct {
	Database string  `json:"database"`
	ShardID  ShardID `json:"shardId"`
	Leader   NodeID  `json:"leader"`
	Follower NodeID  `json:"follower"`
}

// ShardState represents current state of shard.
type ShardState struct {
	ID      ShardID        `json:"id"`
	State   ShardStateType `json:"state"`
	Leader  NodeID         `json:"leader"`
	Replica Replica        `json:"replica"`
}

// StorageState represents storage cluster state.
// NOTICE: it is not safe for concurrent use. //TODO need concurrent safe????
type StorageState struct {
	Name string `json:"name"`

	LiveNodes map[NodeID]StatefulNode

	//TODO remove??
	ShardAssignments map[string]*ShardAssignment       // database's name => shard assignment
	ShardStates      map[string]map[ShardID]ShardState // database's name => shard state
}

// NewStorageState creates storage cluster state
func NewStorageState(name string) *StorageState {
	return &StorageState{
		Name:             name,
		LiveNodes:        make(map[NodeID]StatefulNode),
		ShardAssignments: make(map[string]*ShardAssignment),
		ShardStates:      make(map[string]map[ShardID]ShardState),
	}
}

func (s *StorageState) LeadersOnNode(nodeID NodeID) map[string][]ShardID {
	result := make(map[string][]ShardID)
	for name, shards := range s.ShardStates {
		for shardID, shard := range shards {
			if shard.Leader == nodeID {
				result[name] = append(result[name], shardID)
			}
		}
	}
	return result
}

func (s *StorageState) ReplicasOnNode(nodeID NodeID) map[string][]ShardID {
	result := make(map[string][]ShardID)
	for name, shardAssignment := range s.ShardAssignments {
		shards := shardAssignment.Shards
		for shardID, replicas := range shards {
			if replicas.Contain(nodeID) {
				result[name] = append(result[name], shardID)
			}
		}
	}
	return result
}

// NodeOnline adds a live node into node list.
func (s *StorageState) NodeOnline(node StatefulNode) {
	s.LiveNodes[node.ID] = node
}

// NodeOffline removes a offline node from live node list.
func (s *StorageState) NodeOffline(nodeID NodeID) {
	delete(s.LiveNodes, nodeID)
}

// Stringer returns a human readable string
func (s *StorageState) String() string {
	content := encoding.JSONMarshal(s)
	return string(content)
}
