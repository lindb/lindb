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
	"fmt"
	"io"
	"sync"

	"github.com/lindb/lindb/models"
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

// Partition represents a partition of write ahead log.
type Partition interface {
	io.Closer
	// BuildReplicaForLeader builds replica relation when handle write connection.
	BuildReplicaForLeader(leader models.NodeID, replicas []models.NodeID) error
	// BuildReplicaForFollower builds replica relation when handle replica connection.
	BuildReplicaForFollower(leader models.NodeID, replica models.NodeID) error
	// ReplicaLog writes msg that leader send replica msg.
	// return appended index, if success.
	ReplicaLog(replicaIdx int64, msg []byte) (int64, error)
	// WriteLog writes msg that leader handle client write request.
	WriteLog(msg []byte) error
}

// partition implements  Partition interface.
type partition struct {
	database      string
	currentNodeID models.NodeID
	log           queue.FanOutQueue
	shardID       models.ShardID
	shard         tsdb.Shard
	peers         map[string]ReplicatorPeer
	cliFct        rpc.ClientStreamFactory

	mutex sync.Mutex
}

// NewPartition creates a write ahead log partition.
func NewPartition(shardID models.ShardID, shard tsdb.Shard,
	currentNodeID models.NodeID,
	log queue.FanOutQueue,
	cliFct rpc.ClientStreamFactory,
) Partition {
	return &partition{
		log:           log,
		shardID:       shardID,
		shard:         shard,
		currentNodeID: currentNodeID,
		cliFct:        cliFct,
		peers:         make(map[string]ReplicatorPeer),
	}
}

// ReplicaLog writes msg that leader send replica msg.
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

// ReplicaLog writes msg that leader send replica msg.
func (p *partition) WriteLog(msg []byte) error {
	if len(msg) == 0 {
		return nil
	}
	return p.log.Put(msg)
}

// BuildReplicaForLeader builds replica relation when handle write connection.
// local replicator: replica node == current node.
// remote replicator: replica node != current node.
func (p *partition) BuildReplicaForLeader(
	leader models.NodeID, replicas []models.NodeID,
) error {
	if leader != p.currentNodeID {
		return fmt.Errorf("leader not equals current node")
	}

	for _, replicaNodeID := range replicas {
		p.buildReplica(leader, replicaNodeID)
	}
	return nil
}

// BuildReplicaForFollower builds replica relation when handle replica connection.
func (p *partition) BuildReplicaForFollower(leader models.NodeID, replica models.NodeID) error {
	if replica != p.currentNodeID {
		return fmt.Errorf("replica not equals current node")
	}
	p.buildReplica(leader, replica)
	return nil
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
func (p *partition) buildReplica(leader models.NodeID, replica models.NodeID) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	key := fmt.Sprintf("[%d->%d]", leader, replica)
	_, ok := p.peers[key]
	if ok {
		// exist
		return
	}
	var replicator Replicator
	if replica == p.currentNodeID {
		// local replicator
		replicator = newLocalReplicatorFn(p.shard)
	} else {
		// build remote replicator
		replicator = newRemoteReplicatorFn(&ReplicatorChannel{
			Database: p.database,
			ShardID:  p.shardID,
			Queue:    nil, //TODO set queue
			From:     leader,
			To:       replica,
		}, p.cliFct)
	}

	// startup replicator peer
	peer := NewReplicatorPeer(replicator)
	p.peers[key] = peer
	peer.Startup()
}
