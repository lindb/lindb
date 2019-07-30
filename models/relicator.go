package models

import "fmt"

// BrokerReplicaState represents the replica state list of the broker
type BrokerReplicaState struct {
	ReportTime int64          `json:"reportTime"` // broker report state's time(millisecond)
	Replicas   []ReplicaState `json:"replicas"`   // replica state list under this broker
}

// ReplicaState represents the status of replicator's channel
type ReplicaState struct {
	Cluster      string `json:"cluster"`      // cluster which storing database
	Database     string `json:"database"`     // database name
	ShardID      int32  `json:"shardID"`      // shard id
	TO           Node   `json:"to"`           // target storage node for database's shard
	WriteIndex   int64  `json:"writeIndex"`   // wal write index
	ReplicaIndex int64  `json:"replicaIndex"` // replica index for current replicator's channel
	CommitIndex  int64  `json:"commitIndex"`  // commit index
}

// ShardIndicator returns shard indicator based on cluster/database/shard id
func (r ReplicaState) ShardIndicator() string {
	return fmt.Sprintf("%s/%s/%d", r.Cluster, r.Database, r.ShardID)
}

// Pending returns the num. of pending which it need replica msg
func (r ReplicaState) Pending() int64 {
	return r.WriteIndex - r.ReplicaIndex
}
