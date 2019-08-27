package constants

import (
	"fmt"

	"github.com/lindb/lindb/coordinator/task"
)

// defines common constants will be used in broker and storage
const (
	// ActiveNodesPath represents active nodes prefix path for node register
	ActiveNodesPath = "/active/nodes"
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

	// StorageClusterStatePath represents storage cluster state
	StorageClusterStatePath = "/state/storage/cluster"
	// ReplicaStatePath represents the replica's state
	ReplicaStatePath = "/state/replica"
)

// defines all task kinds
const (
	// CreateShard represents task kind which is create shard for storage node
	CreateShard task.Kind = "create-shard"
)

// GetStorageClusterConfigPath returns path which storing config of storage cluster
func GetStorageClusterConfigPath(name string) string {
	return fmt.Sprintf("%s/%s", StorageClusterConfigPath, name)
}

// GetStorageClusterStatePath returns path whine storing state of storage cluster
func GetStorageClusterStatePath(name string) string {
	return fmt.Sprintf("%s/%s", StorageClusterStatePath, name)
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
