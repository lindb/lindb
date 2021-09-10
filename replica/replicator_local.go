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

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb"
)

var (
	localReplicaScope       = linmetric.NewScope("lindb.replica.local")
	localMaxDecodedBlockVec = localReplicaScope.NewMaxVec("max_decoded_block", "db", "shard")
	localReplicaCountsVec   = localReplicaScope.NewCounterVec("replica_count", "db", "shard")
	localReplicaBytesVec    = localReplicaScope.NewCounterVec("replica_bytes", "db", "shard")
	localReplicaRowsVec     = localReplicaScope.NewCounterVec("replica_rows", "db", "shard")
	localReplicaSequenceVec = localReplicaScope.NewGaugeVec("replica_sequence", "db", "shard")
)

type localReplicator struct {
	replicator

	shard     tsdb.Shard
	logger    *logger.Logger
	batchRows *metric.StorageBatchRows

	block []byte

	statistics struct {
		localMaxDecodedBlock *linmetric.BoundMax
		localReplicaCounts   *linmetric.BoundCounter
		localReplicaBytes    *linmetric.BoundCounter
		localReplicaRows     *linmetric.BoundCounter
		localReplicaSequence *linmetric.BoundGauge
	}
}

func NewLocalReplicator(channel *ReplicatorChannel, shard tsdb.Shard) Replicator {
	lr := &localReplicator{
		replicator: replicator{
			channel: channel,
		},
		shard:     shard,
		batchRows: metric.NewStorageBatchRows(),
		logger:    logger.GetLogger("replica", "LocalReplicator"),
		block:     make([]byte, 256*1024),
	}

	shardStr := shard.ShardID().String()
	lr.statistics.localMaxDecodedBlock = localMaxDecodedBlockVec.WithTagValues(shard.DatabaseName(), shardStr)
	lr.statistics.localReplicaCounts = localReplicaCountsVec.WithTagValues(shard.DatabaseName(), shardStr)
	lr.statistics.localReplicaBytes = localReplicaBytesVec.WithTagValues(shard.DatabaseName(), shardStr)
	lr.statistics.localReplicaRows = localReplicaRowsVec.WithTagValues(shard.DatabaseName(), shardStr)
	lr.statistics.localReplicaSequence = localReplicaSequenceVec.WithTagValues(shard.DatabaseName(), shardStr)
	return lr
}

func (r *localReplicator) Replica(sequence int64, msg []byte) {
	//TODO add util
	var err error
	r.block, err = snappy.Decode(r.block, msg)
	if err != nil {
		r.logger.Error("decompress replica data error", logger.Error(err))
		return
	}

	r.statistics.localMaxDecodedBlock.Update(float64(len(r.block)))
	r.statistics.localReplicaBytes.Add(float64(len(r.block)))
	r.statistics.localReplicaSequence.Update(float64(sequence))
	r.statistics.localReplicaCounts.Incr()

	// flat will always panic when data are corrupted,
	// or data are not serialized correctly
	defer func() {
		if recovered := recover(); recovered != nil {
			r.logger.Error("corrupted flat block",
				logger.Int("message-length", len(msg)),
				logger.Int("decoded-length", len(r.block)),
				logger.Stack(),
			)
		}
		r.block = r.block[:0]
	}()

	r.batchRows.UnmarshalRows(r.block)
	r.statistics.localReplicaRows.Add(float64(r.batchRows.Len()))

	familyIterator := r.batchRows.NewFamilyIterator(r.shard.CurrentInterval())
	for familyIterator.HasNextFamily() {
		familyTime, rows := familyIterator.NextFamily()
		if err := r.shard.WriteRows(familyTime, rows); err != nil {
			r.logger.Error("failed writing family rows",
				logger.Int64("family", familyTime),
				logger.Int("rows", len(rows)),
				logger.String("database", r.shard.DatabaseName()),
				logger.Int("shardID", int(r.shard.ShardID())),
				logger.Error(err))
		}
	}
}
