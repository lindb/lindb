package option

// Database defines database config
type Database struct {
	Name          string `json:"name"`
	NumOfShard    int    `json:"numOfShard"`
	ReplicaFactor int    `json:"replicaFactor"`
}

// Replica defines replica list for spec shard of database
type Replica struct {
	Replicas []int `json:"replicas"`
}

// ShardAssignment defines shard assignment for database
type ShardAssignment struct {
	Shards map[int32]Replica `json:"shards"`
}

// NewShardAssignment returns empty shard assignment instance
func NewShardAssignment() *ShardAssignment {
	return &ShardAssignment{
		Shards: make(map[int32]Replica),
	}
}

// AddReplica adds replica id to replica list of spec shard
func (s *ShardAssignment) AddReplica(shardID int32, replicaID int) {
	replica, ok := s.Shards[shardID]
	if !ok {
		replica = Replica{}
		s.Shards[shardID] = replica
	}
	replica.Replicas = append(replica.Replicas, replicaID)
}
