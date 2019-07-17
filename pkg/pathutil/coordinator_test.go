package pathutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/constants"
)

func TestGetName(t *testing.T) {
	assert.Equal(t, "name", GetName("/test/name"))
	assert.Equal(t, "name", GetName("name"))
}

func TestGetDatabaseAssignPath(t *testing.T) {
	assert.Equal(t, constants.DatabaseAssignPath+"/name", GetDatabaseAssignPath("name"))
}

func TestGetDatabaseConfigPath(t *testing.T) {
	assert.Equal(t, constants.DatabaseConfigPath+"/name", GetDatabaseConfigPath("name"))
}
func TestGetNodePath(t *testing.T) {
	assert.Equal(t, "prefix/name", GetNodePath("prefix", "name"))
}

func TestGetStorageClusterConfigPath(t *testing.T) {
	assert.Equal(t, constants.StorageClusterConfigPath+"/name", GetStorageClusterConfigPath("name"))

}
func TestGetStorageClusterStatePath(t *testing.T) {
	assert.Equal(t, constants.StorageClusterStatePath+"/name", GetStorageClusterStatePath("name"))
}
