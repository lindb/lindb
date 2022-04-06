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

package master

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
)

func TestReplicaLeaderElector_ElectLeader(t *testing.T) {
	elect := newReplicaLeaderElector()
	_, err := elect.ElectLeader(models.NewShardAssignment("test"), nil, models.ShardID(1))
	assert.Equal(t, constants.ErrShardNotFound, err)

	shardAssignment := models.NewShardAssignment("test")
	shardAssignment.AddReplica(models.ShardID(1), models.NodeID(1))
	liveNodes := make(map[models.NodeID]models.StatefulNode)

	_, err = elect.ElectLeader(shardAssignment, liveNodes, models.ShardID(1))
	assert.Equal(t, constants.ErrNoLiveReplica, err)
	liveNodes[models.NodeID(1)] = models.StatefulNode{}

	leader, err := elect.ElectLeader(shardAssignment, liveNodes, models.ShardID(1))
	assert.NoError(t, err)
	assert.Equal(t, models.NodeID(1), leader)
}
