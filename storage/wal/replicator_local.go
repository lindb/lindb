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
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/storage/store"
)

// localReplicator represents local replicator which writes data into local tsdb storage.
type localReplicator struct {
	replicator

	leader  models.NodeID
	segment store.Segment

	statistics *metrics.StorageLocalReplicatorStatistics

	logger logger.Logger
}

func NewLocalReplicator(channel *store.ReplicatorChannel, segment store.Segment) store.Replicator {
	lr := &localReplicator{
		leader: channel.State.Leader,
		replicator: replicator{
			channel: channel,
		},
		segment:    segment,
		statistics: metrics.NewStorageLocalReplicatorStatistics(channel.State.Database, channel.State.ShardID.String()),
		logger:     logger.GetLogger("Replica", "LocalReplicator"),
	}

	// add ack sequence callback
	segment.AckSequence(lr.leader, func(seq int64) {
		lr.SetAckIndex(seq)
		lr.statistics.AckSequence.Incr()
		lr.logger.Info("ack local replica index",
			logger.String("replicator", lr.String()),
			logger.Int64("ackIdx", seq))
	})

	// reset replica index = ack index + 1, replay wal log
	lr.ResetReplicaIndex(lr.AckIndex() + 1)
	segment.Retain() // mark segment will write data

	lr.logger.Info("start local replicator", logger.String("replica", lr.String()),
		logger.Int64("replicaIndex", lr.channel.ConsumerGroup.ConsumedSeq()),
		logger.Int64("ackIndex", lr.AckIndex()))
	return lr
}

// State returns the state of local replicator, it's always ready.
func (r *localReplicator) State() *store.ReplicatorState {
	return &store.ReplicatorState{State: models.ReplicatorReadyState}
}

// Replica replicas local data into local storage,
// 1. check replica replica if valid
// 2. un-compress/unmarshal msg
// 3. write metric data
// 4. commit sequence in data segment
func (r *localReplicator) Replica(sequence int64, msg []byte) {
	var err error

	if !r.segment.ValidateSequence(r.leader, sequence) {
		r.statistics.InvalidSequence.Incr()
		return
	}

	// flat will always panic when data are corrupted,
	// or data are not serialized correctly
	defer func() {
		if err != nil {
			r.IgnoreMessage(sequence)
			r.logger.Warn("ack sequence when replica message failure, will ignore message",
				logger.Int64("sequence", sequence),
				logger.String("replicator", r.String()),
				logger.Error(err))
		}

		// after write need commit sequence, drop write failure data.
		r.segment.CommitSequence(r.leader, sequence)
	}()

	// write data
	rows, err := r.segment.Write(r.leader, sequence, msg)
	if err != nil {
		r.statistics.ReplicaFailures.Incr()
		r.logger.Error("failed writing segment rows",
			logger.Int64("sequence", sequence),
			logger.String("replicator", r.String()),
			logger.Error(err))
		return
	}
	r.statistics.ReplicaRows.Add(float64(rows))
}

// Close closes local replicator.
func (r *localReplicator) Close() {
	// mark write data completed.
	r.segment.Release()
}
