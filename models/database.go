package models

import "github.com/lindb/lindb/pkg/option"

// Database defines database config, database can include multi-cluster
type Database struct {
	Name          string                `json:"name"`          // database's name
	Cluster       string                `json:"cluster"`       // storage cluster's name
	NumOfShard    int                   `json:"numOfShard"`    // num. of shard
	ReplicaFactor int                   `json:"replicaFactor"` // replica refactor
	Option        option.DatabaseOption `json:"option"`        // time series databae option
}

// Replica defines replica list for spec shard of database
type Replica struct {
	Replicas []int `json:"replicas"`
}

// ShardAssignment defines shard assignment for database
type ShardAssignment struct {
	Name   string           `json:"name"` // database's name
	Nodes  map[int]*Node    `json:"nodes"`
	Shards map[int]*Replica `json:"shards"`
}

// NewShardAssignment returns empty shard assignment instance
func NewShardAssignment(name string) *ShardAssignment {
	return &ShardAssignment{
		Name:   name,
		Nodes:  make(map[int]*Node),
		Shards: make(map[int]*Replica),
	}
}

// AddReplica adds replica id to replica list of spec shard
func (s *ShardAssignment) AddReplica(shardID int, replicaID int) {
	replica, ok := s.Shards[shardID]
	if !ok {
		replica = &Replica{}
		s.Shards[shardID] = replica
	}
	replica.Replicas = append(replica.Replicas, replicaID)
}
