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
	"fmt"
	"sync"

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
	fmt.Printf("query time range==%v,%v\n", timeRange, len(partitions))

	for _, partition := range partitions {
		p, err := partition.Get()
		if err != nil {
			logger.Warn("load partition fail", loggerpkg.Error(err))
			continue
		}
		fmt.Printf("query time range==%v,%v\n", timeRange, targetTimeRange.Contains(p.PartitionTime()))
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
