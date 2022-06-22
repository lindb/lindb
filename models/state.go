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
	"encoding/json"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
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

// NodeStateType represents node state type
type NodeStateType int

const (
	NodeOnline NodeStateType = iota + 1
	NodeOffline
)

// StorageStatus represents current storage config status.
type StorageStatus int

const (
	StorageStatusUnknown StorageStatus = iota
	StorageStatusInitialize
	StorageStatusReady
)

func (s StorageStatus) String() string {
	val := "Unknown"
	switch s {
	case StorageStatusInitialize:
		val = "Initialize"
	case StorageStatusReady:
		val = "Ready"
	}
	return val
}

// MarshalJSON encodes storage status.
func (s StorageStatus) MarshalJSON() ([]byte, error) {
	val := s.String()
	return json.Marshal(&val)
}

// UnmarshalJSON decodes storage status.
func (s *StorageStatus) UnmarshalJSON(value []byte) error {
	switch string(value) {
	case `"Initialize"`:
		*s = StorageStatusInitialize
		return nil
	case `"Ready"`:
		*s = StorageStatusReady
		return nil
	default:
		*s = StorageStatusUnknown
		return nil
	}
}

// Storages represents the storage list.
type Storages []Storage

// ToTable returns storage list as table if it has value, else return empty string.
func (s Storages) ToTable() (rows int, tableStr string) {
	if len(s) == 0 {
		return 0, ""
	}
	writer := NewTableFormatter()
	writer.AppendHeader(table.Row{"Namespace", "Status", "Configuration"})
	for i := range s {
		r := s[i]
		writer.AppendRow(table.Row{
			r.Config.Namespace,
			r.Status.String(),
			r.Config.String(),
		})
	}
	return len(s), writer.Render()
}

// Storage represents storage config and state.
type Storage struct {
	config.StorageCluster
	Status StorageStatus `json:"status"`
}

// ReplicaState represents the relationship for a replica.
type ReplicaState struct {
	Database   string  `json:"database"`
	ShardID    ShardID `json:"shardId"`
	Leader     NodeID  `json:"leader"`
	Follower   NodeID  `json:"follower"`
	FamilyTime int64   `json:"familyTime"`
}

func (r ReplicaState) String() string {
	return "[" +
		"database:" + r.Database +
		",shard:" + strconv.Itoa(int(r.ShardID)) +
		",family:" + timeutil.FormatTimestamp(r.FamilyTime, timeutil.DataTimeFormat4) +
		",from(leader):" + strconv.Itoa(int(r.Leader)) +
		",to(follower):" + strconv.Itoa(int(r.Follower)) +
		"]"
}

// ShardState represents current state of shard.
type ShardState struct {
	ID      ShardID        `json:"id"`
	State   ShardStateType `json:"state"`
	Leader  NodeID         `json:"leader"`
	Replica Replica        `json:"replica"`
}

// FamilyState represents current state of shard's family.
type FamilyState struct {
	Database   string     `json:"database"`
	Shard      ShardState `json:"shard"`
	FamilyTime int64      `json:"familyTime"`
}

// StorageState represents storage cluster state.
// NOTICE: it is not safe for concurrent use. //TODO need concurrent safe????
type StorageState struct {
	Name string `json:"name"` // ref Namespace

	LiveNodes map[NodeID]StatefulNode `json:"liveNodes"`

	// TODO remove??
	ShardAssignments map[string]*ShardAssignment       `json:"shardAssignments"` // database's name => shard assignment
	ShardStates      map[string]map[ShardID]ShardState `json:"shardStates"`      // database's name => shard state
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

// LeadersOnNode returns leaders on this node.
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

// ReplicasOnNode returns replicas on this node.
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

// DropDatabase drops shard state/assignment by database's name.
func (s *StorageState) DropDatabase(name string) {
	delete(s.ShardStates, name)
	delete(s.ShardAssignments, name)
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
	return string(encoding.JSONMarshal(s))
}

// StateMachineInfo represents state machine register info.
type StateMachineInfo struct {
	Path        string             `json:"path"`
	CreateState func() interface{} `json:"-"`
}

// StateMetric represents internal state metric.
type StateMetric struct {
	Tags   map[string]string `json:"tags"`
	Fields []StateField      `json:"fields"`
}

// StateField represents internal state value.
type StateField struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}
