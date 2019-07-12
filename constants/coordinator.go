package constants

import "github.com/eleme/lindb/coordinator/task"

// defines all path constants
const (
	// StorageClusterConfigPath represents cluster config store
	StorageClusterConfigPath = "/storage/clusters"
	// ActiveNodesPath represents active nodes prefix path for node register
	ActiveNodesPath = "/active/nodes"
	// DatabaseConfigPath represents database config path
	DatabaseConfigPath = "/database/config"
	// DatabaseAssignPath represents database shard assignment
	DatabaseAssignPath = "/database/assign"
)

// defines all task kinds
const (
	// CreateShard represents task kind which is create shard for storage node
	CreateShard task.Kind = "create-shard"
)
