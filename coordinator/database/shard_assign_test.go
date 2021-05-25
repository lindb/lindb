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
		}, -1, -1)
	assert.NotNil(t, err1)

	_, err1 = ShardAssignment(storageNodeIDs,
		&models.Database{
			Name:          "test",
			NumOfShard:    3,
			ReplicaFactor: 0,
		}, -1, -1)
	assert.NotNil(t, err1)

	_, err2 := ShardAssignment(storageNodeIDs,
		&models.Database{
			Name:          "test",
			NumOfShard:    10,
			ReplicaFactor: 6,
		}, -1, -1)
	assert.NotNil(t, err2)

	shardAssignment, _ := ShardAssignment(storageNodeIDs,
		&models.Database{
			Name:          "test",
			NumOfShard:    10,
			ReplicaFactor: 3,
		}, -1, -1)
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
