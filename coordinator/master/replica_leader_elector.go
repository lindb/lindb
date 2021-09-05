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
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
)

//go:generate mockgen -source=./replica_leader_elector.go -destination=./replica_leader_elector_mock.go -package=master

type ReplicaLeaderElector interface {
	ElectLeader(shardAssignment *models.ShardAssignment,
		liveNodes map[models.NodeID]models.StatefulNode,
		shardID models.ShardID,
	) (leader models.NodeID, err error)
}

type replicaLeaderElector struct {
}

func newReplicaLeaderElector() ReplicaLeaderElector {
	return &replicaLeaderElector{}
}

func (r *replicaLeaderElector) ElectLeader(shardAssignment *models.ShardAssignment,
	liveNodes map[models.NodeID]models.StatefulNode,
	shardID models.ShardID,
) (leader models.NodeID, err error) {
	replicas, ok := shardAssignment.Shards[shardID]
	if !ok {
		// shard not exist
		err = constants.ErrShardNotFound
		return
	}
	// build live replica node
	liveReplicaNodes := models.Replica{}
	for _, replica := range replicas.Replicas {
		_, ok := liveNodes[replica]
		if ok {
			liveReplicaNodes.Replicas = append(liveReplicaNodes.Replicas, replica)
		}
	}
	if len(liveReplicaNodes.Replicas) == 0 {
		// no live replica node
		err = constants.ErrNoLiveReplica
		return
	}
	//elect leader from live replicas
	leader = liveReplicaNodes.Replicas[0]
	return
}
