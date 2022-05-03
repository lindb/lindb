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
	"github.com/golang/snappy"

	"github.com/lindb/lindb/metrics"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb"
)

// localReplicator represents local replicator which writes data into local tsdb storage.
type localReplicator struct {
	replicator

	leader    int32
	shard     tsdb.Shard
	family    tsdb.DataFamily
	logger    *logger.Logger
	batchRows *metric.StorageBatchRows

	block []byte

	statistics *metrics.StorageLocalReplicatorStatistics
}

func NewLocalReplicator(channel *ReplicatorChannel, shard tsdb.Shard, family tsdb.DataFamily) Replicator {
	lr := &localReplicator{
		leader: int32(channel.State.Leader),
		replicator: replicator{
			channel: channel,
		},
		shard:      shard,
		family:     family,
		batchRows:  metric.NewStorageBatchRows(),
		statistics: metrics.NewStorageLocalReplicatorStatistics(channel.State.Database, channel.State.ShardID.String()),
		logger:     logger.GetLogger("replica", "LocalReplicator"),
		block:      make([]byte, 256*1024),
	}

	// add ack sequence callback
	family.AckSequence(lr.leader, func(seq int64) {
		lr.SetAckIndex(seq)
		lr.statistics.AckSequence.Incr()
		lr.logger.Info("ack local replica index",
			logger.String("replica", lr.String()),
			logger.Int64("ackIdx", seq))
	})

	lr.logger.Info("start local replicator", logger.String("replica", lr.String()))
	return lr
}

// Replica replicas local data,
// 1. check replica replica if valid
// 2. un-compress/unmarshal msg
// 3. lookup metadata
// 4. write metric data
// 5. commit sequence in data family
func (r *localReplicator) Replica(sequence int64, msg []byte) {
	if !r.family.ValidateSequence(r.leader, sequence) {
		r.statistics.InvalidSequence.Incr()
		return
	}

	// TODO add util
	var err error
	r.block, err = snappy.Decode(r.block, msg)
	if err != nil {
		r.logger.Error("decompress replica data error", logger.Error(err))
		return
	}

	// flat will always panic when data are corrupted,
	// or data are not serialized correctly
	defer func() {
		r.block = r.block[:0]

		// after write need commit sequence, drop write failure data.
		r.family.CommitSequence(r.leader, sequence)
	}()

	r.batchRows.UnmarshalRows(r.block)
	rowsLen := r.batchRows.Len()
	if rowsLen == 0 {
		return
	}
	rows := r.batchRows.Rows()

	// lookup metric metadata
	if err := r.shard.LookupRowMetricMeta(rows); err != nil {
		r.statistics.ReplicaFailures.Incr()
		r.logger.Error("failed writing family rows",
			logger.Int("rows", r.batchRows.Len()),
			logger.String("database", r.shard.Database().Name()),
			logger.Int("shardID", int(r.shard.ShardID())),
			logger.Error(err))
		return
	}
	// write metric data
	if err := r.family.WriteRows(rows); err != nil {
		r.statistics.ReplicaFailures.Incr()
		r.logger.Error("failed writing family rows",
			logger.Int("rows", r.batchRows.Len()),
			logger.String("database", r.shard.Database().Name()),
			logger.Int("shardID", int(r.shard.ShardID())),
			logger.Error(err))
		return
	}
	r.statistics.ReplicaRows.Add(float64(rowsLen))
}
