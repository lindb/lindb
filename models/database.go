package models

import "github.com/lindb/lindb/pkg/option"

// Database defines database config, database can include multi-cluster
type Database struct {
	Name     string            `json:"name"`
	Clusters []DatabaseCluster `json:"clusters"`
}

// DatabaseCluster represents database's storage cluster config
type DatabaseCluster struct {
	Name          string              `json:"name"`          // database's name
	NumOfShard    int                 `json:"numOfShard"`    // num. of shard
	ReplicaFactor int                 `json:"replicaFactor"` // replica refactor
	Engine        option.EngineOption `json:"engine"`        // time series engine option
}

// Replica defines replica list for spec shard of database
type Replica struct {
	Replicas []int `json:"replicas"`
}

// ShardAssignment defines shard assignment for database
type ShardAssignment struct {
	Config DatabaseCluster  `json:"cluster"`
	Nodes  map[int]*Node    `json:"nodes"`
	Shards map[int]*Replica `json:"shards"`
}

// NewShardAssignment returns empty shard assignment instance
func NewShardAssignment() *ShardAssignment {
	return &ShardAssignment{
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
