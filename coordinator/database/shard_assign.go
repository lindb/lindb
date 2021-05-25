// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package database

import (
	"fmt"
	"math/rand"

	"github.com/lindb/lindb/models"
)

// Shard assigment reference kafka partition assigment
// kafka implement => (https://github.com/apache/kafka/blob/2.3/core/src/main/scala/kafka/admin/AdminUtils.scala)

// ShardAssignment assigns replica list for storage cluster
// which database's each shard based on selected node list in cluster.
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
func ShardAssignment(storageNodeIDs []int, cfg *models.Database, fixedStartIndex, startShardID int) (*models.ShardAssignment, error) {
	numOfShard := cfg.NumOfShard
	replicaFactor := cfg.ReplicaFactor
	if numOfShard <= 0 {
		return nil, fmt.Errorf("shard assign error for databaes[%s], because num. of shard <=0", cfg.Name)
	}
	if replicaFactor <= 0 {
		return nil, fmt.Errorf("shard assign error for databaes[%s], bacause replica factor <=0", cfg.Name)
	}
	if replicaFactor > len(storageNodeIDs) {
		return nil,
			fmt.Errorf("shard assign error for databaes[%s], bacause replica factor > num. of storage nodes",
				cfg.Name)
	}

	shardAssignment := models.NewShardAssignment(cfg.Name)
	assignReplicasToStorageNodes(storageNodeIDs, numOfShard, replicaFactor, fixedStartIndex, startShardID, shardAssignment)

	return shardAssignment, nil
}

func ModifyShardAssignment(storageNodeIDs []int, cfg *models.Database, shardAssignment *models.ShardAssignment,
	fixedStartIndex, startShardID int) error {
	numOfShard := cfg.NumOfShard - len(shardAssignment.Shards)
	replicaFactor := cfg.ReplicaFactor
	if numOfShard <= 0 {
		return fmt.Errorf("shard assign error for databaes[%s], because add num. of shard <=0", cfg.Name)
	}
	if replicaFactor <= 0 {
		return fmt.Errorf("shard assign error for databaes[%s], bacause replica factor <=0", cfg.Name)
	}
	if replicaFactor > len(storageNodeIDs) {
		return fmt.Errorf("shard assign error for databaes[%s], bacause replica factor > num. of storage nodes",
			cfg.Name)
	}

	assignReplicasToStorageNodes(storageNodeIDs, numOfShard, replicaFactor, fixedStartIndex, startShardID, shardAssignment)

	return nil
}

// assignReplicasToStorageNodes assigns replica list for storage cluster
// which database's each shard based on selected node list in cluster.
func assignReplicasToStorageNodes(storageNodeIDs []int,
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
		shardAssignment.AddReplica(currentShardID, leader)

		// assign other replica
		for j := 0; j < replicaFactor-1; j++ {
			idx := replicaIndex(firstReplicaIndex, nextReplicaShift, j, numOfNode)
			shardAssignment.AddReplica(currentShardID, storageNodeIDs[idx])
		}

		// do next shard assign
		currentShardID++
	}

}

// replicaIndex calculates replica index based on first replica index and shift
func replicaIndex(firstReplicaIndex, secondReplicaShift, replicaIndex, numOfNode int) int {
	shift := 1 + (secondReplicaShift+replicaIndex)%(numOfNode-1)
	return (firstReplicaIndex + shift) % numOfNode
}
