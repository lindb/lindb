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

package replica

import (
	"context"
	"sync"

	"github.com/lindb/lindb/pkg/logger"
	replicaRpc "github.com/lindb/lindb/rpc/proto/replica"
)

type remoteReplicator struct {
	replicator

	state ReplicatorState

	//inFlight *InFlightReplica

	rwMutex    sync.RWMutex
	replicaCli replicaRpc.ReplicaServiceClient

	logger *logger.Logger
}

func NewRemoteReplicator() Replicator {
	return &remoteReplicator{
		state:  ReplicatorInitState,
		logger: logger.GetLogger("replica", "remoteReplicator"),
	}
}

func (r *remoteReplicator) IsReady() bool {
	//TODO check node is alive
	r.rwMutex.RLock()
	if r.state == ReplicatorReadyState {
		r.rwMutex.RUnlock()
		return true
	}

	// replicator is not ready, need do init like tcp three-way handshake
	defer r.rwMutex.RUnlock()

	lastReplicaAckIdx, err := r.getLastAckIdxFromReplica() // last ack index remote replica node
	if err != nil {
		r.logger.Warn("do get replica ack index err", logger.Error(err))
		return false
	}
	replicaIdx := r.ReplicaIndex() // current need replica index from current node
	nextReplicaIdx := lastReplicaAckIdx + 1
	if nextReplicaIdx == replicaIdx {
		// replica index == remote replica append index, can do replicator
		r.state = ReplicatorReadyState
		return true
	}

	// replica index != remote replica append index, need reset index
	appendIdx := r.AppendIndex()
	smallestAckIdx := r.AckIndex()
	switch {
	case lastReplicaAckIdx < smallestAckIdx:
		// maybe new remote replica node add in cluster or remote replica data lost.
		needResetReplicaIdx := smallestAckIdx + 1
		r.logger.Warn("replica node ack < current node ack, need reset remote replica node's append index",
			logger.Int64("lastReplicaAckIdx", lastReplicaAckIdx),
			logger.Int64("smallestAckIdx", smallestAckIdx),
			logger.Int64("resetReplicaIdx", needResetReplicaIdx))
		// send resent index request
		_, err := r.replicaCli.Reset(context.TODO(), &replicaRpc.ResetIndexRequest{
			Database:    "",
			Shard:       0,
			Leader:      0,
			AppendIndex: needResetReplicaIdx,
		})
		if err != nil {
			r.logger.Warn("do reset replica append index err", logger.Error(err))
			return false
		}
		r.ResetReplicaIndex(nextReplicaIdx)
		r.state = ReplicatorReadyState
		return true
	case lastReplicaAckIdx > appendIdx:
		// new write data will be lost, because leader's lost old wal data
		r.ResetAppendIndex(nextReplicaIdx)
	}
	// remote replica ack idx > current ack idx, maybe ack request lost
	r.ResetReplicaIndex(nextReplicaIdx)
	r.logger.Warn("remote replica ack idx != current replica idx, reset current replica idx",
		logger.Int64("lastReplicaAckIdx", lastReplicaAckIdx),
		logger.Int64("replicaIdx", replicaIdx),
		logger.Int64("nextReplicaIdx", nextReplicaIdx),
	)
	r.state = ReplicatorReadyState
	return true
}

func (r *remoteReplicator) Replica(idx int64, msg []byte) {
	cli, err := r.replicaCli.Replica(context.TODO())
	if err != nil {
		r.state = ReplicatorFailureState
		return
	}

	err = cli.Send(&replicaRpc.ReplicaRequest{
		Database:     "",
		Shard:        0,
		Leader:       0,
		ReplicaIndex: 0,
	})
	if err != nil {
		r.state = ReplicatorFailureState
		return
	}
	resp, err := cli.Recv()
	if err != nil {
		r.state = ReplicatorFailureState
		return
	}
	r.SetAckIndex(resp.ReplicaIndex)
}

func (r *remoteReplicator) getLastAckIdxFromReplica() (int64, error) {
	resp, err := r.replicaCli.GetReplicaAckIndex(context.TODO(), &replicaRpc.GetReplicaAckIndexRequest{
		Database: "",
		Shard:    0,
		Leader:   0,
	})
	if err != nil {
		return 0, nil
	}
	return resp.AckIndex, nil
}
