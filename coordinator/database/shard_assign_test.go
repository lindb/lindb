package database

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
)

func TestShardAssign(t *testing.T) {
	storageNodeIDs := []int{0, 1, 2, 3, 4}

	_, err1 := ShardAssignment(storageNodeIDs,
		&models.Database{
			Name:          "test",
			NumOfShard:    0,
			ReplicaFactor: 3,
		})
	assert.NotNil(t, err1)

	_, err1 = ShardAssignment(storageNodeIDs,
		&models.Database{
			Name:          "test",
			NumOfShard:    3,
			ReplicaFactor: 0,
		})
	assert.NotNil(t, err1)

	_, err2 := ShardAssignment(storageNodeIDs,
		&models.Database{
			Name:          "test",
			NumOfShard:    10,
			ReplicaFactor: 6,
		})
	assert.NotNil(t, err2)

	shardAssignment, _ := ShardAssignment(storageNodeIDs,
		&models.Database{
			Name:          "test",
			NumOfShard:    10,
			ReplicaFactor: 3,
		})
	checkShardAssignResult(shardAssignment, t)
}

func checkShardAssignResult(shardAssignment *models.ShardAssignment, t *testing.T) {
	assert.Equal(t, 10, len(shardAssignment.Shards))
	var nodes = make(map[int]map[int]int)
	for shardID, replica := range shardAssignment.Shards {
		for _, nodeID := range replica.Replicas {
			node, ok := nodes[nodeID]
			if !ok {
				node = make(map[int]int)
				nodes[nodeID] = node
			}
			node[shardID] = shardID
		}
	}
	assert.Equal(t, 5, len(nodes))
	for _, replicas := range nodes {
		assert.Equal(t, 6, len(replicas))
	}
}
