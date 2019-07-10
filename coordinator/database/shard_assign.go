package database

import (
	"fmt"
	"math/rand"

	"github.com/eleme/lindb/models"
)

// Shard assigment reference kafka partition assigment
// kafka implement => (https://github.com/apache/kafka/blob/2.3/core/src/main/scala/kafka/admin/AdminUtils.scala)

// ShardAssignment assigns replica list for database's each shard based on selected node list
// There are 2 goals of replica assignment:
// 1. Spread the replicas evenly among storage nodes for currently cluster state.
// 2. For shards assigned to a particular storage node, their other replicas are spread over the other storage nodes.
//
// TO achieve this goal, we:
// 1. Assign the first replica of each shard by round-robin, starting from a random position in the storage node list.
// 2. Assign the remaining replicas of each shard with an increasing shift.
//
// Here is an example of assigning, (num. of nodes = 5, num of shards = 10, replica factor = 3)
// node-0	node-1	node-2	node-3	node-4
// s0		s1		s2		s3		s4		(1st replica)
// s5		s6		s7		s8		s9		(1st replica)
// s4		s0		s1		s2		s3		(2st replica)
// s8		s9		s5		s6		s7		(2st replica)
// s3		s4		s0		s1		s2		(3st replica)
// s7		s8		s9		s5		s6		(3st replica)
func ShardAssignment(storageNodeIDs []int, database models.Database) (*models.ShardAssignment, error) {
	numOfShard := database.NumOfShard
	replicaFactor := database.ReplicaFactor
	if numOfShard <= 0 {
		return nil, fmt.Errorf("shard assign error for database[%s], because num. of shard <=0", database.Name)
	}
	if replicaFactor <= 0 {
		return nil, fmt.Errorf("shard assign error for database[%s], bacause replica factor <=0", database.Name)
	}
	if replicaFactor > len(storageNodeIDs) {
		return nil,
			fmt.Errorf("shard assign error for database[%s], bacause replica factor > num. of storage nodes",
				database.Name)
	}

	shardAssignment := models.NewShardAssignment()
	assignReplicasToStorages(storageNodeIDs, numOfShard, replicaFactor, -1, -1, shardAssignment)

	return shardAssignment, nil
}

// assignReplicasToStorages assigns replica list for database's each shard, return shard assignment result
func assignReplicasToStorages(storageNodeIDs []int,
	numOfShard, replicaFactor, fixedStartIndex, startShardID int,
	shardAssignment *models.ShardAssignment) {
	numOfNode := len(storageNodeIDs)

	// init start index/shift/current shard
	startIndex := fixedStartIndex
	nextReplicaShift := fixedStartIndex
	if fixedStartIndex < 0 {
		startIndex = rand.Intn(numOfNode)
		nextReplicaShift = rand.Intn(numOfNode)
	}
	currentShardID := 0
	if startShardID >= 0 {
		currentShardID = startShardID
	}

	// assign replica list for each shard
	for i := 0; i < numOfShard; i++ {
		if currentShardID > 0 && (currentShardID%numOfNode == 0) {
			nextReplicaShift++
		}
		firstReplicaIndex := (currentShardID + startIndex) % numOfNode

		// elect first replica as leader
		leader := storageNodeIDs[firstReplicaIndex]
		shardAssignment.AddReplica(int32(currentShardID), leader)

		// assign other replica
		for j := 0; j < replicaFactor-1; j++ {
			idx := replicaIndex(firstReplicaIndex, nextReplicaShift, j, numOfNode)
			shardAssignment.AddReplica(int32(currentShardID), storageNodeIDs[idx])
		}

		// do next shard assign
		currentShardID++
	}

}

// replicaIndex calculates replica index based on firist replica index and shift
func replicaIndex(firstReplicaIndex, secondReplicaShift, replicaIndex, numOfNode int) int {
	shift := 1 + (secondReplicaShift+replicaIndex)%(numOfNode-1)
	return (firstReplicaIndex + shift) % numOfNode
}
