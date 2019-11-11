package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseAssignPath(t *testing.T) {
	assert.Equal(t, DatabaseAssignPath+"/name", GetDatabaseAssignPath("name"))
}

func TestGetDatabaseConfigPath(t *testing.T) {
	assert.Equal(t, DatabaseConfigPath+"/name", GetDatabaseConfigPath("name"))
}

func TestGetNodePath(t *testing.T) {
	assert.Equal(t, "prefix/name", GetNodePath("prefix", "name"))
}

func TestGetStorageClusterConfigPath(t *testing.T) {
	assert.Equal(t, StorageClusterConfigPath+"/name", GetStorageClusterConfigPath("name"))

}
func TestGetStorageClusterStatePath(t *testing.T) {
	assert.Equal(t, StorageClusterNodeStatePath+"/name", GetStorageClusterNodeStatePath("name"))
}

func TestGetStorageClusterStatPath(t *testing.T) {
	assert.Equal(t, StorageClusterStatPath+"/name", GetStorageClusterStatPath("name"))
}

func TestGetReplicaStatePath(t *testing.T) {
	assert.Equal(t, ReplicaStatePath+"/1.1.1.1:port", GetReplicaStatePath("1.1.1.1:port"))
}

func TestGetNodeMonitoringStatPath(t *testing.T) {
	assert.Equal(t, StateNodesPath+"/1.1.1.1:port", GetNodeMonitoringStatPath("1.1.1.1:port"))
}
