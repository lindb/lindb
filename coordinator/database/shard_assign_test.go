package database

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/models"
)

func TestShardAssign(t *testing.T) {
	storageNodeIDs := []int{0, 1, 2, 3, 4}

	_, err1 := ShardAssignment(storageNodeIDs,
		models.Database{
			Name:          "test",
			NumOfShard:    0,
			ReplicaFactor: 3,
		})
	assert.NotNil(t, err1)

	_, err2 := ShardAssignment(storageNodeIDs,
		models.Database{
			Name:          "test",
			NumOfShard:    10,
			ReplicaFactor: 6,
		})
	assert.NotNil(t, err2)

	shardAssignment, _ := ShardAssignment(storageNodeIDs,
		models.Database{
			Name:          "test",
			NumOfShard:    10,
			ReplicaFactor: 3,
		})

	shardAssignment2 := models.NewShardAssignment()
	shardAssignment2.AddReplica(0, 0)
	shardAssignment2.AddReplica(0, 1)
	shardAssignment2.AddReplica(0, 2)

	shardAssignment2.AddReplica(1, 1)
	shardAssignment2.AddReplica(1, 2)
	shardAssignment2.AddReplica(1, 3)

	shardAssignment2.AddReplica(2, 2)
	shardAssignment2.AddReplica(2, 3)
	shardAssignment2.AddReplica(2, 4)

	shardAssignment2.AddReplica(3, 3)
	shardAssignment2.AddReplica(3, 4)
	shardAssignment2.AddReplica(3, 0)

	shardAssignment2.AddReplica(4, 4)
	shardAssignment2.AddReplica(4, 0)
	shardAssignment2.AddReplica(4, 1)

	shardAssignment2.AddReplica(5, 0)
	shardAssignment2.AddReplica(5, 2)
	shardAssignment2.AddReplica(5, 3)

	shardAssignment2.AddReplica(6, 1)
	shardAssignment2.AddReplica(6, 3)
	shardAssignment2.AddReplica(6, 4)

	shardAssignment2.AddReplica(7, 2)
	shardAssignment2.AddReplica(7, 4)
	shardAssignment2.AddReplica(7, 0)

	shardAssignment2.AddReplica(8, 3)
	shardAssignment2.AddReplica(8, 0)
	shardAssignment2.AddReplica(8, 1)

	shardAssignment2.AddReplica(9, 4)
	shardAssignment2.AddReplica(9, 1)
	shardAssignment2.AddReplica(9, 2)

	assert.Equal(t, *shardAssignment, *shardAssignment2)
}
