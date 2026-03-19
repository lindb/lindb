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

package base

import (
	"sync"
	"time"

	"github.com/lindb/common/pkg/fileutil"
	loggerpkg "github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/store"
)

type CalcPartitionTimeFn func(timestamp int64) (partitionTime int64)

type Shard struct {
	ID models.ShardID

	DB store.Database

	CreatePartitionFn   store.CreatePartitionFn
	CalcPartitionTimeFn CalcPartitionTimeFn
	Partitions          *store.Partitions // partition timestamp -> partition

	mutex sync.Mutex
}

func (s *Shard) GetOrCreatePartition(timestamp int64) (store.Partition, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	paritionTime := s.CalcPartitionTimeFn(timestamp)
	partition, ok, err := s.Partitions.GetPartition(paritionTime)
	if err != nil {
		return nil, err
	}
	if ok {
		return partition, nil
	}

	partition, err = s.CreatePartitionFn(paritionTime)
	if err != nil {
		return nil, err
	}

	s.Partitions.PutPartition(paritionTime, store.NewPartition(partition))

	return partition, nil
}

// ShardID implements store.Shard.
func (s *Shard) ShardID() models.ShardID {
	return s.ID
}

func (s *Shard) Replica() models.Replica {
	return s.Database().GetShardReplica(s.ID)
}

func (s *Shard) Database() store.Database {
	return s.DB
}

func (s *Shard) GetPartitions(interval timeutil.Interval, timeRange timeutil.TimeRange) (result []store.Partition) {
	var partitions []*store.LazyPartition
	s.mutex.Lock()
	partitions = s.Partitions.GetPartitions()
	s.mutex.Unlock()

	targetTimeRange := &timeutil.TimeRange{
		// need truncate start timestamp, e.g. 20190902 19:05:48 => 20190902 00:00:00
		Start: interval.Calculator().CalcSegmentTime(timeRange.Start),
		End:   timeRange.End,
	}

	for _, partition := range partitions {
		p, err := partition.Get()
		if err != nil {
			logger.Warn("load partition fail", loggerpkg.Error(err))
			continue
		}
		// TODO: maybe user not input time range of query(modify check condition)
		if timeRange.Start == timeRange.End || targetTimeRange.Contains(p.PartitionTime()) {
			result = append(result, p)
		}
	}

	return
}

func (s *Shard) Close() error {
	var partitions []*store.LazyPartition

	s.mutex.Lock()
	partitions = s.Partitions.GetPartitions()
	s.mutex.Unlock()

	for _, partition := range partitions {
		if partition.Loaded() {
			p, err := partition.Get()
			if err != nil {
				logger.Warn("load partition fail when close shard", loggerpkg.Error(err))
				continue
			}
			p.Close()
		}
	}
	return nil
}

func (s *Shard) TTL() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	database := s.Database()
	// normally, shard only has one retention policy
	policy := database.GetOption().Option.RetentionPolicies[0]
	expireTime := policy.CalcExpireTime(now.UnixMilli())

	partitions := s.Partitions.GetPartitions()
	for _, lp := range partitions {
		partition, err := lp.Get()
		if err != nil {
			logger.Warn("load partition fail when do ttl", loggerpkg.String("database", database.Name()), loggerpkg.Error(err))
			continue
		}
		// partition time is before expire time, need do ttl
		if partition.PartitionTime() > expireTime {
			continue
		}
		if err := partition.Close(); err != nil {
			logger.Warn("close partition fail when do ttl", loggerpkg.String("database", database.Name()), loggerpkg.Error(err))
			continue
		}
		// Remote partition from shard
		s.Partitions.RemovePartition(partition.PartitionTime())

		if err := fileutil.RemoveDir(partition.Path()); err != nil {
			logger.Warn("remove partition dir fail when do ttl", loggerpkg.String("database", database.Name()), loggerpkg.Error(err))
			continue
		}
		logger.Info("partition ttl completed",
			loggerpkg.String("database", s.Database().Name()),
			loggerpkg.String("path", partition.Path()))
	}
}
