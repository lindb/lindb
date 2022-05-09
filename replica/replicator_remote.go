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

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	protoReplicaV1 "github.com/lindb/lindb/proto/gen/v1/replica"
	"github.com/lindb/lindb/rpc"
)

// remoteReplicator implements Replicator interface, do remote wal replica.
type remoteReplicator struct {
	replicator

	ctx   context.Context
	state models.ReplicatorState

	cliFct        rpc.ClientStreamFactory
	replicaCli    protoReplicaV1.ReplicaServiceClient
	replicaStream protoReplicaV1.ReplicaService_ReplicaClient
	stateMgr      storage.StateManager

	isSuspend *atomic.Bool
	suspend   chan struct{}

	rwMutex sync.RWMutex

	statistics *metrics.StorageRemoteReplicatorStatistics
	logger     *logger.Logger
}

// NewRemoteReplicator creates remote replicator.
func NewRemoteReplicator(
	ctx context.Context,
	channel *ReplicatorChannel,
	stateMgr storage.StateManager,
	cliFct rpc.ClientStreamFactory,
) Replicator {
	r := &remoteReplicator{
		ctx: ctx,
		replicator: replicator{
			channel: channel,
		},
		cliFct:     cliFct,
		stateMgr:   stateMgr,
		state:      models.ReplicatorInitState,
		isSuspend:  atomic.NewBool(false),
		suspend:    make(chan struct{}),
		statistics: metrics.NewStorageRemoteReplicatorStatistics(channel.State.Database, channel.State.ShardID.String()),
		logger:     logger.GetLogger("replica", "RemoteReplicator"),
	}

	// watch follower node state change
	stateMgr.WatchNodeStateChangeEvent(channel.State.Follower, r.handleNodeStateChangeEvent)

	r.logger.Info("start remote replicator", logger.String("replica", r.String()))
	return r
}

func (r *remoteReplicator) handleNodeStateChangeEvent(state models.NodeStateType) {
	if state == models.NodeOnline {
		if r.isSuspend.CAS(true, false) {
			r.logger.Info("notify replicator follower node is online", logger.String("replicator", r.String()))
			r.suspend <- struct{}{} // notify follower node online
		}
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
	if r.state == models.ReplicatorReadyState {
		r.rwMutex.Unlock()
		return true
	}

	r.statistics.NotReady.Incr()

	// replicator is not ready, need do init like tcp three-way handshake
	follower := r.replicator.channel.State.Follower
	node, ok := r.stateMgr.GetLiveNode(follower)
	if !ok {
		r.logger.Warn("follower node is offline, need suspend replicator", logger.String("replicator", r.String()))

		r.rwMutex.Unlock() // unlock
		if r.isSuspend.CAS(false, true) {
			r.statistics.FollowerOffline.Incr()
			<-r.suspend // wait follower node online
		}
		return r.IsReady() // check replicator is ready now
	}

	defer r.rwMutex.Unlock()
	if r.replicaStream != nil {
		r.statistics.NeedCloseLastStream.Incr()
		if err := r.replicaStream.CloseSend(); err != nil {
			r.statistics.CloseLastStreamFailures.Incr()
			r.logger.Warn("close replica service client stream err, when reconnection",
				logger.String("replicator", r.String()),
				logger.Error(err))
		}
	}
	replicaCli, err := r.cliFct.CreateReplicaServiceClient(&node)
	if err != nil {
		r.statistics.CreateReplicaCliFailures.Incr()
		r.logger.Warn("create replica service client err",
			logger.String("replicator", r.String()),
			logger.Error(err))
		return false
	}
	r.replicaCli = replicaCli
	r.statistics.CreateReplicaCli.Incr()

	// pass metadata(database/shard state) when create rpc connection.
	replicaState := encoding.JSONMarshal(&r.channel.State)
	ctx := rpc.CreateOutgoingContextWithPairs(r.ctx,
		constants.RPCMetaReplicaState, string(replicaState))
	r.replicaStream, err = replicaCli.Replica(ctx) // TODO add timeout ??
	if err != nil {
		r.statistics.CloseLastStreamFailures.Incr()
		r.logger.Warn("create replica service client stream err",
			logger.String("replicator", r.String()),
			logger.Error(err))
		return false
	}
	r.statistics.CreateReplicaStream.Incr()

	remoteLastReplicaAckIdx, err := r.getLastAckIdxFromReplica() // last ack index remote replica node
	if err != nil {
		r.statistics.GetLastAckFailures.Incr()
		r.logger.Warn("do get replica ack index err",
			logger.String("replicator", r.String()),
			logger.Error(err))
		return false
	}
	localReplicaIdx := r.ReplicaIndex() // current need replica index from current node
	nextReplicaIdx := remoteLastReplicaAckIdx + 1
	if nextReplicaIdx == localReplicaIdx {
		// replica index == remote replica append index, can do replicator
		r.state = models.ReplicatorReadyState
		return true
	}

	// replica index != remote replica append index, need reset index
	appendIdx := r.AppendIndex()
	smallestAckIdx := r.AckIndex()
	switch {
	case remoteLastReplicaAckIdx < smallestAckIdx:
		// maybe new remote replica node add in cluster or remote replica data lost.
		needResetReplicaIdx := smallestAckIdx + 1
		r.logger.Warn("replica node ack < current node ack, need reset remote replica node's append index",
			logger.String("replicator", r.String()),
			logger.Int64("remoteLastReplicaAckIdx", remoteLastReplicaAckIdx),
			logger.Int64("smallestAckIdx", smallestAckIdx),
			logger.Int64("resetReplicaIdx", needResetReplicaIdx))
		// send reset index request
		_, err := r.replicaCli.Reset(context.TODO(), &protoReplicaV1.ResetIndexRequest{
			Database:    r.channel.State.Database,
			Shard:       int32(r.channel.State.ShardID),
			Leader:      int32(r.channel.State.Leader),
			FamilyTime:  r.channel.State.FamilyTime,
			AppendIndex: needResetReplicaIdx,
		})
		if err != nil {
			r.statistics.ResetFollowerAppendIdxFailures.Incr()
			r.logger.Warn("do reset replica append index err",
				logger.String("replicator", r.String()),
				logger.Error(err))
			return false
		}
		r.statistics.ResetFollowerAppendIdx.Incr()
		_ = r.ResetReplicaIndex(nextReplicaIdx)
		r.state = models.ReplicatorReadyState
		return true
	case remoteLastReplicaAckIdx > appendIdx:
		// new write data will be lost, because leader's lost old wal data
		r.ResetAppendIndex(nextReplicaIdx)
		r.statistics.ResetAppendIdx.Incr()
	}
	// remote replica ack idx > current ack idx, maybe ack request lost
	_ = r.ResetReplicaIndex(nextReplicaIdx - 1)
	r.SetAckIndex(remoteLastReplicaAckIdx)
	// get new local replica idx, double check if reset replica index successfully.
	newLocalReplicaIdx := r.ReplicaIndex()
	r.logger.Warn("remote replica ack idx != current replica idx, reset current replica idx",
		logger.String("replicator", r.String()),
		logger.Int64("remoteLastReplicaAckIdx", remoteLastReplicaAckIdx),
		logger.Int64("oldLocalReplicaIdx", localReplicaIdx),
		logger.Int64("newLocalReplicaIdx", newLocalReplicaIdx),
		logger.Int64("nextReplicaIdx", nextReplicaIdx),
	)
	if newLocalReplicaIdx == nextReplicaIdx {
		// replica index == remote replica append index, can do replicator
		r.statistics.ResetReplicaIdx.Incr()
		r.state = models.ReplicatorReadyState
		r.logger.Info("remote replica ack idx != current replica idx, reset current replica idx successfully",
			logger.String("replicator", r.String()))
		return true
	}
	r.statistics.ResetReplicaIdxFailures.Incr()
	return false
}

// Replica sends data to remote replica node.
func (r *remoteReplicator) Replica(idx int64, msg []byte) {
	cli := r.replicaStream
	err := cli.Send(&protoReplicaV1.ReplicaRequest{
		ReplicaIndex: idx,
		Record:       msg,
	})
	if err != nil {
		r.state = models.ReplicatorFailureState
		r.statistics.SendMsgFailures.Incr()
		return
	}
	r.statistics.SendMsg.Incr()
	resp, err := cli.Recv()
	if err != nil {
		r.state = models.ReplicatorFailureState
		r.statistics.ReceiveMsgFailures.Incr()
		return
	}
	r.statistics.ReceiveMsg.Incr()
	r.logger.Debug("receive replica response",
		logger.String("replicator", r.String()),
		logger.Int64("replicaIdx", resp.ReplicaIndex),
		logger.Int64("ackIdx", resp.AckIndex))
	if resp.AckIndex == resp.ReplicaIndex {
		// if ack index = replica, need ack wal
		r.SetAckIndex(resp.AckIndex)
		r.statistics.AckSequence.Incr()
	} else {
		r.statistics.InvalidAckSequence.Incr()
	}
}

// getLastAckIdxFromReplica returns replica replica ack index.
func (r *remoteReplicator) getLastAckIdxFromReplica() (int64, error) {
	resp, err := r.replicaCli.GetReplicaAckIndex(context.TODO(), &protoReplicaV1.GetReplicaAckIndexRequest{
		Database:   r.channel.State.Database,
		Shard:      int32(r.channel.State.ShardID),
		Leader:     int32(r.channel.State.Leader),
		FamilyTime: r.channel.State.FamilyTime,
	})
	if err != nil {
		return 0, err
	}
	return resp.AckIndex, nil
}
