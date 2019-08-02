package pathutil

import (
	"fmt"
	"path/filepath"

	"github.com/lindb/lindb/constants"
)

// GetStorageClusterConfigPath returns path which storing config of storage cluster
func GetStorageClusterConfigPath(name string) string {
	return fmt.Sprintf("%s/%s", constants.StorageClusterConfigPath, name)
}

// GetStorageClusterStatePath returns path whine storing state of storage cluster
func GetStorageClusterStatePath(name string) string {
	return fmt.Sprintf("%s/%s", constants.StorageClusterStatePath, name)
}

// GetDatabaseConfigPath returns path which storing config of database
func GetDatabaseConfigPath(name string) string {
	return fmt.Sprintf("%s/%s", constants.DatabaseConfigPath, name)
}

// GetDatabaseAssignPath returns path which storing shard assignment of database
func GetDatabaseAssignPath(name string) string {
	return fmt.Sprintf("%s/%s", constants.DatabaseAssignPath, name)
}

// GetNodePath returns node register path
func GetNodePath(prefix, node string) string {
	return fmt.Sprintf("%s/%s", prefix, node)
}

// GetReplicaStatePath returns replica's state path
func GetReplicaStatePath(node string) string {
	return fmt.Sprintf("%s/%s", constants.ReplicaStatePath, node)
}

// GetName returns name, splits path and gets last path
func GetName(path string) string {
	_, name := filepath.Split(path)
	return name
}
