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

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lindb/common/models"

	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/validate"
)

// DatabaseNames represents the database name list.
type DatabaseNames []string

// ToTable returns database name list as table if it has value, else return empty string.
func (dbs DatabaseNames) ToTable() (rows int, tableStr string) {
	if len(dbs) == 0 {
		return 0, ""
	}
	writer := models.NewTableFormatter()
	writer.AppendHeader(table.Row{"Database"})
	for i := range dbs {
		r := dbs[i]
		writer.AppendRow(table.Row{r})
	}
	return len(dbs), writer.Render()
}

// Databases represents the database list.
type Databases []Database

// ToTable returns database list as table if it has value, else return empty string.
func (dbs Databases) ToTable() (rows int, tableStr string) {
	if len(dbs) == 0 {
		return 0, ""
	}
	writer := models.NewTableFormatter()
	writer.AppendHeader(table.Row{"Name", "Desc"})
	for i := range dbs {
		r := dbs[i]
		writer.AppendRow(table.Row{r.Name, r.Desc})
	}
	return len(dbs), writer.Render()
}

// ShardID represents type for shard id.
type ShardID int

func (s ShardID) String() string { return strconv.Itoa(int(s)) }

func (s ShardID) Int() int { return int(s) }

// ParseShardID returns ShardID by given shard string value.
func ParseShardID(shard string) ShardID {
	shardID, _ := strconv.Atoi(shard)
	return ShardID(shardID)
}

// DatabaseConfig represents a database configuration about config and families
type DatabaseConfig struct {
	Option   *option.DatabaseOption `toml:"option" json:"option"`
	Name     string                 `toml:"name" json:"name"`
	ShardIDs []ShardID              `toml:"shardIDs" json:"shardIDs"`
}

// Router represents the router of database.
type Router struct {
	Key      string   `json:"key" validate:"required"`    // routing key
	Broker   string   `json:"broker" validate:"required"` // target broker
	Database string   `json:"database,omitempty"`         // target database
	Values   []string `json:"values" validate:"required"` // routing values
}

// LogicDatabase defines database logic config, database can include multi-cluster.
type LogicDatabase struct {
	Name    string   `json:"name" validate:"required"` // database's name
	Desc    string   `json:"desc,omitempty"`
	Routers []Router `json:"routers" validate:"required"` // database router
}

// Database defines database config.
type Database struct {
	Option *option.DatabaseOption `json:"option"`                   // time series database option
	Name   string                 `json:"name" validate:"required"` // database's name
	Desc   string                 `json:"desc,omitempty"`
}

func (db *Database) Default() {
	if db.Option == nil {
		db.Option = &option.DatabaseOption{}
	}
	db.Option.Default()
}

func (db *Database) Validate() error {
	err := validate.Validator.Struct(db)
	if err != nil {
		return err
	}
	// validate time series engine option
	err = db.Option.Validate()
	if err != nil {
		return err
	}
	return nil
}

// String returns the database's description.
func (db *Database) String() string {
	result := "create database " + db.Name + " with ("
	result += "numOfShard=" + fmt.Sprintf("%d", db.Option.NumOfShard) + ", replicaRactor=" + fmt.Sprintf("%d", db.Option.ReplicaFactor)
	result += ", intervals " + db.Option.Intervals.String()
	return result
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
	Shards map[ShardID]*Replica `json:"shards"`
	Name   string               `json:"name"` // database's name

	replicaFactor int // for storage recover
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
	if !replica.Contain(replicaID) {
		replica.Replicas = append(replica.Replicas, replicaID)
	}

	if len(replica.Replicas) > s.replicaFactor {
		s.replicaFactor = len(replica.Replicas)
	}
}

// GetReplicaFactor returns the factor of replica.
func (s *ShardAssignment) GetReplicaFactor() int {
	return s.replicaFactor
}
