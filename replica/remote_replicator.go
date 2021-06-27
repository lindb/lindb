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

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/rpc"
	replicaRpc "github.com/lindb/lindb/rpc/proto/replica"
)

type remoteReplicator struct {
	replicator

	state ReplicatorState

	//inFlight *InFlightReplica

	cliFct        rpc.ClientStreamFactory
	replicaCli    replicaRpc.ReplicaServiceClient
	replicaStream replicaRpc.ReplicaService_ReplicaClient

	rwMutex sync.RWMutex

	logger *logger.Logger
}

// NewRemoteReplicator creates remote replicator.
func NewRemoteReplicator(channel *ReplicatorChannel,
	cliFct rpc.ClientStreamFactory,
) Replicator {
	return &remoteReplicator{
		replicator: replicator{
			channel: channel,
		},
		cliFct: cliFct,
		state:  ReplicatorInitState,
		logger: logger.GetLogger("replica", "remoteReplicator"),
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

	//TODO check node is alive
	//TODO close cli/stream if re-connect???
	replicaCli, err := r.cliFct.CreateReplicaServiceClient(models.Node{})
	if err != nil {
		//TODO add metric
		r.logger.Warn("create replica service client err", logger.Error(err))
		return false
	}
	r.replicaCli = replicaCli
	r.replicaStream, err = replicaCli.Replica(context.TODO()) //TODO add timeout ??
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
		_, err := r.replicaCli.Reset(context.TODO(), &replicaRpc.ResetIndexRequest{
			Database:    r.channel.Database,
			Shard:       int32(r.channel.ShardID),
			Leader:      int32(r.channel.From),
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
		// new write data will be lost, because leader's lost old wal data
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
	err := cli.Send(&replicaRpc.ReplicaRequest{
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
	resp, err := r.replicaCli.GetReplicaAckIndex(context.TODO(), &replicaRpc.GetReplicaAckIndexRequest{
		Database: r.channel.Database,
		Shard:    int32(r.channel.ShardID),
		Leader:   int32(r.channel.From),
	})
	if err != nil {
		return 0, err
	}
	return resp.AckIndex, nil
}
