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

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/lindb/common/pkg/logger"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/discovery"
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
	data, err := queue.NewFanOutQueue(path, int64(pageSize))
	if err != nil {
		return nil, err
	}

	wal := &writeAheadLog{
		segment:    segment,
		data:       data,
		peers:      make(map[models.NodeID]store.ReplicatorPeer),
		streamings: make(map[string]store.ReplicatorPeer),
		logger:     logger.GetLogger("WAL", "WriteAheadLog"),
	}

	// TODO: only writable write ahead log need build replica for leader
	// build wal replica relation for leader
	wal.BuildReplicaForLeader(meta.CurrentNode(), segment.Partition().Shard().Replica().Replicas)
	return wal, nil
}

// Database implements [store.WriteAheadLog].
func (w *writeAheadLog) Database() string {
	return w.segment.Partition().Shard().Database().Name()
}

// Segment implements [store.WriteAheadLog].
func (w *writeAheadLog) SegmentTime() int64 {
	return w.segment.SegmentTimeRange().Start
}

func (w *writeAheadLog) SwitchSequence(leader models.NodeID) int64 {
	local := w.getLocalReplicator(leader)
	if local == nil {
		return -1
	}
	return local.SwitchSequence()
}

func (w *writeAheadLog) PersistSequence(leader models.NodeID) {
	local := w.getLocalReplicator(leader)
	if local == nil {
		return
	}
	local.PersistSequence()
}

func (w *writeAheadLog) GetSequence(leader models.NodeID) int64 {
	local := w.getLocalReplicator(leader)
	if local == nil {
		return -1
	}
	return local.sequence.Load()
}

func (w *writeAheadLog) GetImmutableSequence(leader models.NodeID) int64 {
	local := w.getLocalReplicator(leader)
	if local == nil {
		return -1
	}
	return local.sequence.Load()
}

func (w *writeAheadLog) AckSequence(leader models.NodeID, ack int64) {
	local := w.getLocalReplicator(leader)
	if local == nil {
		return
	}
	local.AckSequence(ack)
}

func (w *writeAheadLog) getLocalReplicator(leader models.NodeID) *localReplicator {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	peer, ok := w.peers[leader]
	if !ok {
		return nil
	}
	replicator := peer.Replicator()
	if local, ok := replicator.(*localReplicator); ok {
		return local
	}
	return nil
}

// Shard implements [store.WriteAheadLog].
func (w *writeAheadLog) Shard() models.ShardID {
	return w.segment.Partition().Shard().ShardID()
}

func (w *writeAheadLog) Receive(event discovery.MetaEvent) {
	switch stateEvent := event.(type) {
	case *store.ConsumerStateChange:
		if stateEvent.IsDelete {
			// delete consumer
			w.mutex.Lock()
			defer w.mutex.Unlock()

			peer, ok := w.streamings[stateEvent.Streaming]
			if ok {
				// shutdown peer
				peer.Shutdown()

				w.data.DeleteConsumerGroup(stateEvent.Streaming)
				delete(w.streamings, stateEvent.Streaming)
			}
		} else {
			w.BuildConsumer(stateEvent.Streaming, stateEvent.ConsumerID)
		}
	}
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

func (w *writeAheadLog) WriteArrow(records []arrow.RecordBatch) error {
	if w.closed.Load() {
		return errors.New("write ahead log is closed")
	}
	for _, record := range records {
		record.Release()
		fmt.Println(record)
	}
	return nil
}

func (w *writeAheadLog) Close() error {
	if w.closed.CompareAndSwap(false, true) {
		// Collect peers while holding the lock, then shut them down without the lock
		// so Shutdown() can wait for goroutines to exit without risk of deadlock.
		w.mutex.Lock()
		peers := make([]store.ReplicatorPeer, 0, len(w.peers)+len(w.streamings))
		for _, peer := range w.peers {
			peers = append(peers, peer)
		}
		for _, peer := range w.streamings {
			peers = append(peers, peer)
		}
		w.mutex.Unlock()

		for _, peer := range peers {
			peer.Shutdown()
		}

		// close queue
		w.data.Close()
		// unsubscribe wal consumer events
		meta.GetStorageMetaManager().Unsubscribe(w)
	}
	return nil
}

func (w *writeAheadLog) BuildConsumer(streaming string, consume models.NodeID) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	peer, ok := w.streamings[streaming]
	if ok {
		peer.Replicator().Resume()
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
			Type:        models.ReplicatorTypeObserve,
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
	replicator = NewRemoteReplicator(context.TODO(), &channel)

	// startup replicator peer
	peer = NewReplicatorPeer(replicator)
	w.streamings[streaming] = peer
	peer.Startup()
}

// BuildReplicaForLeader builds replica relation when handle writeTask connection.
// local replicator: replica node == current node.
// remote replicator: replica node != current node.
func (w *writeAheadLog) BuildReplicaForLeader(
	leader models.NodeID, replicas []models.NodeID,
) error {
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
		channel.State.Type = models.ReplicatorTypeLocal
		// local replicator
		replicator = NewLocalReplicator(&channel, w.segment)
	} else {
		// build remote replicator
		channel.State.Type = models.ReplicatorTypeRemote
		// TODO: set context
		replicator = NewRemoteReplicator(context.TODO(), &channel)
	}

	// startup replicator peer
	peer := NewReplicatorPeer(replicator)
	w.peers[replica] = peer
	peer.Startup()

	return nil
}
