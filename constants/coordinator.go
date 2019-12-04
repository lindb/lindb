package constants

import (
	"fmt"

	"github.com/lindb/lindb/coordinator/task"
)

// defines common constants will be used in broker and storage
const (
	// ActiveNodesPath represents active nodes prefix path for node register
	ActiveNodesPath = "/active/nodes"
	// StateNodesPath represents the state of node that node will report runtime status
	StateNodesPath = "/state/nodes"
)

// defines storage level constants will be used in storage
//const ()

// defines broker level constants will be used in broker
const (
	// MasterPath represents master elect path
	MasterPath = "/master/node"
	// DatabaseAssignPath represents database shard assignment
	DatabaseAssignPath = "/database/assign"
	// StorageClusterConfigPath represents cluster config store
	StorageClusterConfigPath = "/storage/cluster/config"
	// DatabaseConfigPath represents database config path
	DatabaseConfigPath = "/database/config"

	// StorageClusterNodeStatePath represents storage cluster's node state
	StorageClusterNodeStatePath = "/state/storage/nodes/cluster"
	// ReplicaStatePath represents the replica's state
	ReplicaStatePath = "/state/replica"
	// StorageClusterStatPath represents storage cluster's node monitoring stat
	StorageClusterStatPath = "/state/storage/stat/cluster"
)

// defines all task kinds
const (
	// CreateShard represents task kind which is create shard for storage node
	CreateShard task.Kind = "create-shard"
	// FlushDatabase represents task kind which is flush memory database for storage node
	FlushDatabase task.Kind = "flush-database"
)

// GetStorageClusterConfigPath returns path which storing config of storage cluster
func GetStorageClusterConfigPath(name string) string {
	return fmt.Sprintf("%s/%s", StorageClusterConfigPath, name)
}

// GetStorageClusterNodeStatePath returns path whine storing state of storage cluster
func GetStorageClusterNodeStatePath(name string) string {
	return fmt.Sprintf("%s/%s", StorageClusterNodeStatePath, name)
}

// GetStorageClusterStatPath returns path whine storing monitoring stat of storage cluster
func GetStorageClusterStatPath(name string) string {
	return fmt.Sprintf("%s/%s", StorageClusterStatPath, name)
}

// GetDatabaseConfigPath returns path which storing config of database
func GetDatabaseConfigPath(name string) string {
	return fmt.Sprintf("%s/%s", DatabaseConfigPath, name)
}

// GetDatabaseAssignPath returns path which storing shard assignment of database
func GetDatabaseAssignPath(name string) string {
	return fmt.Sprintf("%s/%s", DatabaseAssignPath, name)
}

// GetNodePath returns node register path
func GetNodePath(prefix, node string) string {
	return fmt.Sprintf("%s/%s", prefix, node)
}

// GetReplicaStatePath returns replica's state path
func GetReplicaStatePath(node string) string {
	return fmt.Sprintf("%s/%s", ReplicaStatePath, node)
}

// GetNodeMonitoringStatPath returns the node monitoring stat's path
func GetNodeMonitoringStatPath(node string) string {
	return fmt.Sprintf("%s/%s", StateNodesPath, node)
}
