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

package wal

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/queue"
	"github.com/lindb/lindb/storage/store"
)

type writeAheadLog struct {
	segment store.Segment
	data    queue.FanOutQueue

	peers map[models.NodeID]store.ReplicatorPeer // follower nodeID => peer

	streamings map[string]store.ReplicatorPeer // streaming name => peer

	closed atomic.Bool
	mutex  sync.Mutex

	logger logger.Logger
}

func NewWriteAheadLog(path string, segment store.Segment) (store.WriteAheadLog, error) {
	// TODO: database level???
	pageSize := config.GlobalStorageConfig().WAL.PageSize
	fmt.Printf("wal page=%s\n", path)
	data, err := queue.NewFanOutQueue(path, int64(pageSize))
	if err != nil {
		return nil, err
	}
	return &writeAheadLog{
		segment:    segment,
		data:       data,
		peers:      make(map[models.NodeID]store.ReplicatorPeer),
		streamings: make(map[string]store.ReplicatorPeer),
		logger:     logger.GetLogger("WAL", "WriteAheadLog"),
	}, nil
}

func (w *writeAheadLog) Replica(replicaIndex int64, msg []byte) (int64, error) {
	if w.closed.Load() {
		return -1, errors.New("write ahead log is closed")
	}
	appendIdx := w.data.Queue().AppendedSeq() + 1
	if replicaIndex != appendIdx {
		return appendIdx, nil
	}
	if err := w.data.Queue().Put(msg); err != nil {
		return -1, err
	}
	return appendIdx, nil
}

func (w *writeAheadLog) ReplicaAckIndex() int64 {
	return w.data.Queue().AppendedSeq()
}

func (w *writeAheadLog) ResetReplicaIndex(idx int64) {
	w.data.SetAppendedSeq(idx - 1)
}

func (w *writeAheadLog) Get(index int64) ([]byte, error) {
	return w.data.Queue().Get(index)
}

func (w *writeAheadLog) Write(msg []byte) error {
	if len(msg) == 0 {
		return nil
	}
	if w.closed.Load() {
		return errors.New("write ahead log is closed")
	}
	// TODO: add metric
	return w.data.Queue().Put(msg)
}

func (w *writeAheadLog) Close() error {
	if w.closed.CompareAndSwap(false, true) {
		// close queue
		w.data.Close()
	}
	return nil
}

func (w *writeAheadLog) Consume(streaming string, consume models.NodeID) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if _, ok := w.streamings[streaming]; ok {
		// exist
		return
	}
	walConsumer, err := w.data.GetOrCreateConsumerGroup(streaming)
	if err != nil {
		w.logger.Error("failed to create wal consumer group",
			logger.String("streaming", streaming),
			logger.Int("consumer", consume.Int()),
			logger.Error(err))
		return
	}

	var replicator store.Replicator
	channel := store.ReplicatorChannel{
		State: &models.ReplicaState{
			Database:    w.segment.Partition().Shard().Database().Name(),
			ShardID:     w.segment.Partition().Shard().ShardID(),
			Streaming:   streaming,
			SegmentTime: w.segment.SegmentTimeRange().Start,
			Leader:      fmt.Sprintf("%d", meta.CurrentNode()),
			Follower:    fmt.Sprintf("%d", consume),
		},
		ConsumerGroup: walConsumer,
	}
	// build remote replicator
	// TODO: set context
	replicator = NewRemoteReplicator(context.TODO(), store.ReplicatorTypeObserve, &channel)

	// startup replicator peer
	peer := NewReplicatorPeer(replicator)
	w.streamings[streaming] = peer
	peer.Startup()
}

// BuildReplicaForLeader builds replica relation when handle writeTask connection.
// local replicator: replica node == current node.
// remote replicator: replica node != current node.
func (w *writeAheadLog) BuildReplicaForLeader(
	leader models.NodeID, replicas []models.NodeID,
) error {
	if leader != meta.CurrentNode() {
		return fmt.Errorf("leader not equals current node")
	}

	for _, replicaNodeID := range replicas {
		if err := w.buildReplica(leader, replicaNodeID); err != nil {
			w.logger.Error(
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
func (w *writeAheadLog) BuildReplicaForFollower(leader, replica models.NodeID) error {
	if replica != meta.CurrentNode() {
		return fmt.Errorf("replica not equals current node")
	}
	err := w.buildReplica(leader, replica)
	if err != nil {
		w.logger.Error("follower failed building replication channel from leader",
			logger.Int("leader", leader.Int()),
			logger.Int("follower", replica.Int()),
		)
	}
	return err
}

// buildReplica builds replica replication based on leader/follower node.
func (w *writeAheadLog) buildReplica(leader, replica models.NodeID) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if _, ok := w.peers[replica]; ok {
		// exist
		return nil
	}
	walConsumer, err := w.data.GetOrCreateConsumerGroup(fmt.Sprintf("%d", replica))
	if err != nil {
		return err
	}
	var replicator store.Replicator
	channel := store.ReplicatorChannel{
		State: &models.ReplicaState{
			Database:    w.segment.Partition().Shard().Database().Name(),
			ShardID:     w.segment.Partition().Shard().ShardID(),
			SegmentTime: w.segment.SegmentTimeRange().Start,
			Leader:      fmt.Sprintf("%d", leader),
			Follower:    fmt.Sprintf("%d", replica),
		},
		ConsumerGroup: walConsumer,
	}
	if replica == meta.CurrentNode() {
		// local replicator
		replicator = NewLocalReplicator(&channel, w.segment)
	} else {
		// build remote replicator
		// TODO: set context
		replicator = NewRemoteReplicator(context.TODO(), store.ReplicatorTypeRemote, &channel)
	}

	// startup replicator peer
	peer := NewReplicatorPeer(replicator)
	w.peers[replica] = peer
	peer.Startup()

	return nil
}
