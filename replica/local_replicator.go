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

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb"
)

type localReplicator struct {
	replicator

	shard     tsdb.Shard
	logger    *logger.Logger
	batchRows *metric.BatchRows

	block []byte
}

func NewLocalReplicator(channel *ReplicatorChannel, shard tsdb.Shard) Replicator {
	lr := &localReplicator{
		replicator: replicator{
			channel: channel,
		},
		shard:     shard,
		batchRows: metric.NewBatchRows(),
		logger:    logger.GetLogger("replica", "LocalReplicator"),
		block:     make([]byte, 256*1024),
	}
	return lr
}

func (r *localReplicator) Replica(_ int64, msg []byte) {
	//TODO add util
	var err error
	r.block, err = snappy.Decode(r.block, msg)
	if err != nil {
		r.logger.Error("decompress replica data error", logger.Error(err))
		return
	}

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

	familyIterator := r.batchRows.NewFamilyIterator(r.shard.CurrentInterval())
	for familyIterator.HasNextFamily() {
		familyTime, rows := familyIterator.NextFamily()
		if err := r.shard.WriteBatchRows(familyTime, rows); err != nil {
			r.logger.Error("failed writing family rows",
				logger.Int64("family", familyTime),
				logger.Int("rows", len(rows)),
				logger.String("database", r.shard.DatabaseName()),
				logger.Int("shardID", int(r.shard.ShardID())),
				logger.Error(err))
		}
	}
}
