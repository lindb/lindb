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

package log

import (
	"github.com/lindb/common/pkg/fileutil"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/storage/base"
	"github.com/lindb/lindb/storage/store"
)

type shard struct {
	base.Shard

	dir string
}

func NewShard(id models.ShardID, database *Database) (store.Shard, error) {
	shardPath := store.ShardPath(database.Name(), id)
	if err := fileutil.MkDirIfNotExist(shardPath); err != nil {
		return nil, err
	}
	s := &shard{
		Shard: base.Shard{
			ID:                  id,
			CalcPartitionTimeFn: store.MinuteIntervalCalc.CalcSegmentTime,
			Partitions:          store.NewPartitions(store.MinuteInterval),
			DB:                  database,
		},
		dir: shardPath,
	}
	s.CreatePartitionFn = s.createPartition

	if err := s.Partitions.Load(store.PartitionsPath(database.Name(), id, store.MinuteInterval), func(timestamp int64) (*store.LazyPartition, error) {
		return store.NewLazyPartition(timestamp, s.createPartition), nil
	}); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *shard) createPartition(timestamp int64) (store.Partition, error) {
	return NewPartition(timestamp, s)
}
