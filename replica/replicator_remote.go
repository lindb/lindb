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

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/rpc"
)

type remoteReplicator struct {
	replicator

	ctx   context.Context
	state ReplicatorState

	//inFlight *InFlightReplica

	cliFct        rpc.ClientStreamFactory
	replicaCli    protoReplicaV1.ReplicaServiceClient
	replicaStream protoReplicaV1.ReplicaService_ReplicaClient
	stateMgr      storage.StateManager

	rwMutex sync.RWMutex

	logger *logger.Logger
}

// NewRemoteReplicator creates remote replicator.
func NewRemoteReplicator(
	ctx context.Context,
	channel *ReplicatorChannel,
	stateMgr storage.StateManager,
	cliFct rpc.ClientStreamFactory,
) Replicator {
	return &remoteReplicator{
		ctx: ctx,
		replicator: replicator{
			channel: channel,
		},
		cliFct:   cliFct,
		stateMgr: stateMgr,
		state:    ReplicatorInitState,
		logger:   logger.GetLogger("replica", "RemoteReplicator"),
	}
}

// IsReady returns remote replicator channel is ready.
// 1. state == ready, return true
// 2. state != ready, do channel init like tcp three-way handshake.
//    a. next remote replica index = current node's replica index, return true.
//    b. last remote ack index < current node's smallest ack, need reset remote replica index, then return true.
//    c. last remote ack index > current node's append index,
//   	 need reset current append index/replica index, then return true.
func (r *remoteReplicator) IsReady() bool {
	r.rwMutex.Lock()
	if r.state == ReplicatorReadyState {
		r.rwMutex.Unlock()
		return true
	}

	// replicator is not ready, need do init like tcp three-way handshake
	defer r.rwMutex.Unlock()

	node, ok := r.stateMgr.GetLiveNode(r.replicator.channel.State.Follower)
	if !ok {
		r.logger.Warn("follower node is offline")
		return false
	}
	//TODO close cli/stream if re-connect???
	replicaCli, err := r.cliFct.CreateReplicaServiceClient(&node)
	if err != nil {
		//TODO add metric
		r.logger.Warn("create replica service client err", logger.Error(err))
		return false
	}
	r.replicaCli = replicaCli
	// pass metadata(database/shard state) when create rpc connection.
	replicaState := encoding.JSONMarshal(&r.channel.State)
	ctx := rpc.CreateOutgoingContextWithPairs(r.ctx,
		constants.RPCMetaReplicaState, string(replicaState))
	r.replicaStream, err = replicaCli.Replica(ctx) //TODO add timeout ??
	if err != nil {
		//TODO add metric
		r.logger.Warn("create replica service client stream err", logger.Error(err))
		return false
	}

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
			logger.String("replicator", r.String()),
			logger.Int64("lastReplicaAckIdx", lastReplicaAckIdx),
			logger.Int64("smallestAckIdx", smallestAckIdx),
			logger.Int64("resetReplicaIdx", needResetReplicaIdx))
		// send reset index request
		_, err := r.replicaCli.Reset(context.TODO(), &protoReplicaV1.ResetIndexRequest{
			Database:    r.channel.State.Database,
			Shard:       int32(r.channel.State.ShardID),
			Leader:      int32(r.channel.State.Leader),
			AppendIndex: needResetReplicaIdx,
		})
		if err != nil {
			r.logger.Warn("do reset replica append index err",
				logger.String("replicator", r.String()),
				logger.Error(err))
			return false
		}
		_ = r.ResetReplicaIndex(nextReplicaIdx)
		r.state = ReplicatorReadyState
		return true
	case lastReplicaAckIdx > appendIdx:
		// new writeTask data will be lost, because leader's lost old wal data
		r.ResetAppendIndex(nextReplicaIdx)
	}
	// remote replica ack idx > current ack idx, maybe ack request lost
	_ = r.ResetReplicaIndex(nextReplicaIdx)
	r.logger.Warn("remote replica ack idx != current replica idx, reset current replica idx",
		logger.String("replicator", r.String()),
		logger.Int64("lastReplicaAckIdx", lastReplicaAckIdx),
		logger.Int64("replicaIdx", replicaIdx),
		logger.Int64("nextReplicaIdx", nextReplicaIdx),
	)
	r.state = ReplicatorReadyState
	return true
}

// Replica sends data to remote replica node.
func (r *remoteReplicator) Replica(idx int64, msg []byte) {
	cli := r.replicaStream
	err := cli.Send(&protoReplicaV1.ReplicaRequest{
		ReplicaIndex: idx,
		Record:       msg,
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
	r.SetAckIndex(resp.AckIndex)
}

// getLastAckIdxFromReplica returns replica replica ack index.
func (r *remoteReplicator) getLastAckIdxFromReplica() (int64, error) {
	resp, err := r.replicaCli.GetReplicaAckIndex(context.TODO(), &protoReplicaV1.GetReplicaAckIndexRequest{
		Database: r.channel.State.Database,
		Shard:    int32(r.channel.State.ShardID),
		Leader:   int32(r.channel.State.Leader),
	})
	if err != nil {
		return 0, err
	}
	return resp.AckIndex, nil
}
