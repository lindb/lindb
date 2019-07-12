package pathutil

import (
	"fmt"
	"path/filepath"

	"github.com/eleme/lindb/constants"
)

// GetStorageClusterPath returns path which storing config of storage cluster
func GetStorageClusterPath(name string) string {
	return fmt.Sprintf("%s/%s", constants.StorageClusterConfigPath, name)
}

// GetDatabaseConfigPath returns path which storing config of database
func GetDatabaseConfigPath(name string) string {
	return constants.DatabaseConfigPath + "/" + name
}

// GetDatabaseAssignPath returns path which storing shard assignment of database
func GetDatabaseAssignPath(name string) string {
	return constants.DatabaseAssignPath + "/" + name
}

// GetNodePath returns node register path
func GetNodePath(prefix, node string) string {
	return fmt.Sprintf("%s/%s", prefix, node)
}

// GetName returns name, splits path and gets last path
func GetName(path string) string {
	_, name := filepath.Split(path)
	return name
}
