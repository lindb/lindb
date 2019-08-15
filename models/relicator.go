package models

import "fmt"

// BrokerReplicaState represents the replica state list of the broker
type BrokerReplicaState struct {
	ReportTime int64          `json:"reportTime"` // broker report state's time(millisecond)
	Replicas   []ReplicaState `json:"replicas"`   // replica state list under this broker
}

// ReplicaState represents the status of replicator's channel
type ReplicaState struct {
	Database     string `json:"database"`     // database name
	ShardID      int32  `json:"shardID"`      // shard id
	Target       Node   `json:"target"`       // target storage node for database's shard
	Pending      int64  `json:"pending"`      // the num. of pending which it need replica msg
	ReplicaIndex int64  `json:"replicaIndex"` // replica index for current replicator's channel
	AckIndex     int64  `json:"ackIndex"`     // commit index
}

// ShardIndicator returns shard indicator based on database/shard id
func (r ReplicaState) ShardIndicator() string {
	return fmt.Sprintf("%s/%d", r.Database, r.ShardID)
}
