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
	"fmt"
	"io"
	"sync"

	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/rpc"
	"github.com/lindb/lindb/tsdb"
)

//go:generate mockgen -source=./partition.go -destination=./partition_mock.go -package=replica

var (
	// for testing
	newLocalReplicatorFn  = NewLocalReplicator
	newRemoteReplicatorFn = NewRemoteReplicator
)

// Partition represents a partition of writeTask ahead log.
type Partition interface {
	io.Closer
	// BuildReplicaForLeader builds replica relation when handle writeTask connection.
	BuildReplicaForLeader(leader models.NodeID, replicas []models.NodeID) error
	// BuildReplicaForFollower builds replica relation when handle replica connection.
	BuildReplicaForFollower(leader models.NodeID, replica models.NodeID) error
	// ReplicaLog writes msg that leader send replica msg.
	// return appended index, if success.
	ReplicaLog(replicaIdx int64, msg []byte) (int64, error)
	// WriteLog writes msg that leader handle client writeTask request.
	WriteLog(msg []byte) error
	// ReplicaAckIndex returns the index which replica appended index.
	ReplicaAckIndex() int64
	ResetReplicaIndex(idx int64)
}

// partition implements Partition interface.
type partition struct {
	ctx           context.Context
	currentNodeID models.NodeID
	log           queue.FanOutQueue
	shardID       models.ShardID
	shard         tsdb.Shard
	family        tsdb.DataFamily

	peers    map[string]ReplicatorPeer
	cliFct   rpc.ClientStreamFactory
	stateMgr storage.StateManager

	mutex sync.Mutex

	logger *logger.Logger
}

// NewPartition creates a writeTask ahead log partition.
func NewPartition(
	ctx context.Context,
	shard tsdb.Shard,
	family tsdb.DataFamily,
	currentNodeID models.NodeID,
	log queue.FanOutQueue,
	cliFct rpc.ClientStreamFactory,
	stateMgr storage.StateManager,
) Partition {
	return &partition{
		ctx:           ctx,
		log:           log,
		shardID:       shard.ShardID(),
		shard:         shard,
		family:        family,
		currentNodeID: currentNodeID,
		cliFct:        cliFct,
		stateMgr:      stateMgr,
		peers:         make(map[string]ReplicatorPeer),
		logger:        logger.GetLogger("replica", "Partition"),
	}
}

// ReplicaLog writes msg that leader sends replica msg.
// return appended index, if success.
func (p *partition) ReplicaLog(replicaIdx int64, msg []byte) (int64, error) {
	appendIdx := p.log.HeadSeq()
	if replicaIdx != appendIdx {
		return appendIdx, nil
	}
	if err := p.log.Put(msg); err != nil {
		return -1, err
	}
	return appendIdx, nil
}

func (p *partition) ReplicaAckIndex() int64 {
	return p.log.HeadSeq() - 1
}

func (p *partition) ResetReplicaIndex(idx int64) {
	p.log.SetAppendSeq(idx)
}

// WriteLog writes msg that leader send replica msg.
func (p *partition) WriteLog(msg []byte) error {
	if len(msg) == 0 {
		return nil
	}
	return p.log.Put(msg)
}

// BuildReplicaForLeader builds replica relation when handle writeTask connection.
// local replicator: replica node == current node.
// remote replicator: replica node != current node.
func (p *partition) BuildReplicaForLeader(
	leader models.NodeID, replicas []models.NodeID,
) error {
	if leader != p.currentNodeID {
		return fmt.Errorf("leader not equals current node")
	}

	for _, replicaNodeID := range replicas {
		if err := p.buildReplica(leader, replicaNodeID); err != nil {
			p.logger.Error(
				"leader failed building replication channel to follower",
				logger.String("leader", leader.String()),
				logger.String("follower", replicaNodeID.String()),
				logger.Error(err),
			)
			return err
		}
	}
	return nil
}

// BuildReplicaForFollower builds replica relation when handle replica connection.
func (p *partition) BuildReplicaForFollower(leader models.NodeID, replica models.NodeID) error {
	if replica != p.currentNodeID {
		return fmt.Errorf("[BUG] replica not equals current node")
	}
	err := p.buildReplica(leader, replica)
	if err != nil {
		p.logger.Error("follower failed building replication channel from leader",
			logger.Int("leader", leader.Int()),
			logger.Int("follower", replica.Int()),
		)
	}
	return err
}

// Close shutdowns all replica workers.
func (p *partition) Close() error {
	var waiter sync.WaitGroup
	waiter.Add(len(p.peers))
	for k := range p.peers {
		r := p.peers[k]
		go func() {
			waiter.Done()
			r.Shutdown()
		}()
	}
	waiter.Wait()

	// close log
	p.log.Close()
	return nil
}

// buildReplica builds replica replication based on leader/follower node.
func (p *partition) buildReplica(leader models.NodeID, replica models.NodeID) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	key := fmt.Sprintf("[%d->%d]", leader, replica)
	_, ok := p.peers[key]
	if ok {
		// exist
		return nil
	}
	walConsumer, err := p.log.GetOrCreateFanOut(fmt.Sprintf("%d_%d", leader, replica))
	if err != nil {
		return err
	}
	var replicator Replicator
	channel := ReplicatorChannel{
		State: &models.ReplicaState{
			Database: p.shard.DatabaseName(),
			ShardID:  p.shardID,
			Leader:   leader,
			Follower: replica,
		},
		Queue: walConsumer,
	}
	if replica == p.currentNodeID {
		// local replicator
		replicator = newLocalReplicatorFn(&channel, p.shard, p.family)
	} else {
		// build remote replicator
		replicator = newRemoteReplicatorFn(p.ctx, &channel, p.stateMgr, p.cliFct)
	}

	// startup replicator peer
	peer := NewReplicatorPeer(replicator)
	p.peers[key] = peer
	peer.Startup()
	return nil
}
