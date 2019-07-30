package constants

import "github.com/eleme/lindb/coordinator/task"

// defines common constants will be used in broker and storage
const (
	// ActiveNodesPath represents active nodes prefix path for node register
	ActiveNodesPath = "/active/nodes"
)

// defines storage level constants will be used in storage
const (
	// DatabaseAssignPath represents database shard assignment
	DatabaseAssignPath = "/database/assign"
)

// defines broker level constants will be used in broker
const (
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
