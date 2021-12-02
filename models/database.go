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
	"fmt"
	"strconv"

	"github.com/lindb/lindb/pkg/option"
)

// ShardID represents type for shard id.
type ShardID int

func (s ShardID) String() string { return strconv.Itoa(int(s)) }
func (s ShardID) Int() int       { return int(s) }

// ParseShardID returns ShardID by given shard string value.
func ParseShardID(shard string) ShardID {
	shardID, _ := strconv.Atoi(shard)
	return ShardID(shardID)
}

// Database defines database config, database can include multi-cluster.
type Database struct {
	Name          string                `json:"name" binding:"required"` // database's name
	Storage       string                `json:"storage"`                 // storage cluster's name
	NumOfShard    int                   `json:"numOfShard"`              // num. of shard
	ReplicaFactor int                   `json:"replicaFactor"`           // replica refactor
	Option        option.DatabaseOption `json:"option"`                  // time series database option
	Desc          string                `json:"desc,omitempty"`
}

// String returns the database's description.
func (db Database) String() string {
	result := "create database " + db.Name + " with "
	result += "shard " + fmt.Sprintf("%d", db.NumOfShard) + ", replica " + fmt.Sprintf("%d", db.ReplicaFactor)
	result += ", intervals " + db.Option.Intervals.String()
	return result
}

type DatabaseAssignment struct {
	ShardAssignment *ShardAssignment      `json:"shardAssignment"`
	Option          option.DatabaseOption `json:"option"`
}

// Replica defines replica list for spec shard of database.
type Replica struct {
	Replicas []NodeID `json:"replicas"`
}

// Contain returns if replica include node id.
func (r Replica) Contain(nodeID NodeID) bool {
	for _, id := range r.Replicas {
		if id == nodeID {
			return true
		}
	}
	return false
}

// ShardAssignment defines shard assignment for database.
type ShardAssignment struct {
	Name   string               `json:"name"` // database's name
	Shards map[ShardID]*Replica `json:"shards"`
}

// NewShardAssignment returns empty shard assignment instance.
func NewShardAssignment(name string) *ShardAssignment {
	return &ShardAssignment{
		Name:   name,
		Shards: make(map[ShardID]*Replica),
	}
}

// AddReplica adds replica id to replica list of spec shard.
func (s *ShardAssignment) AddReplica(shardID ShardID, replicaID NodeID) {
	replica, ok := s.Shards[shardID]
	if !ok {
		replica = &Replica{}
		s.Shards[shardID] = replica
	}
	replica.Replicas = append(replica.Replicas, replicaID)
}
