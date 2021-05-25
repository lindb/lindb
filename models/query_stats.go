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

package models

import (
	"sync"

	"github.com/lindb/lindb/pkg/timeutil"
)

// QueryStats represents the query stats when need explain query flow stat
type QueryStats struct {
	StorageNodes map[string]*StorageStats `json:"storageNodes,omitempty"`
	Cost         int64                    `json:"cost"` // total query cost
	ExpressCost  int64                    `json:"expressCost"`
}

// NewQueryStats creates the query stats
func NewQueryStats() *QueryStats {
	return &QueryStats{
		StorageNodes: make(map[string]*StorageStats),
	}
}

// MergeStorageTaskStats merges storage task execution stats
func (s *QueryStats) MergeStorageTaskStats(taskID string, stats *StorageStats) {
	s.StorageNodes[taskID] = stats
}

// StorageStats represents query stats in storage side
type StorageStats struct {
	NetPayload            int                   `json:"netPayload"`
	NetCost               int64                 `json:"netCost"`
	TotalCost             int64                 `json:"totalCost"`
	PlanCost              int64                 `json:"planCost"`
	TagFilterCost         int64                 `json:"tagFilterCost"`
	Shards                map[int32]*ShardStats `json:"shards,omitempty"`
	CollectTagValuesStats map[string]int64      `json:"collectTagValuesStats,omitempty"`

	start int64      // track search start time in storage side
	mutex sync.Mutex // need add lock for goroutine update stats data
}

// NewStorageStats creates the query stats in storage side
func NewStorageStats() *StorageStats {
	return &StorageStats{
		Shards:                make(map[int32]*ShardStats),
		CollectTagValuesStats: make(map[string]int64),
		start:                 timeutil.NowNano(),
	}
}

// Complete completes the query stats when query completed
func (s *StorageStats) Complete() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TotalCost = timeutil.NowNano() - s.start
}

// SetPlanCost sets plan cost
func (s *StorageStats) SetPlanCost(cost int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.PlanCost = cost
}

// SetTagFilterCost sets tag filter cost
func (s *StorageStats) SetTagFilterCost(cost int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TagFilterCost = cost
}

// SetShardSeriesIDsSearchStats sets shard series ids search stats
func (s *StorageStats) SetShardSeriesIDsSearchStats(shardID int32, numOfSeries uint64, seriesFilterCost int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	stats := newShardStats()
	stats.NumOfSeries = numOfSeries
	stats.SeriesFilterCost = seriesFilterCost
	s.Shards[shardID] = stats
}

// SetShardMemoryDataFilterCost sets shard memory data filter cost
func (s *StorageStats) SetShardMemoryDataFilterCost(shardID int32, cost int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.MemFilterCost = cost
	}
}

// SetShardKVDataFilterCost sets shard data filter cost in kv store
func (s *StorageStats) SetShardKVDataFilterCost(shardID int32, cost int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.KVFilterCost = cost
	}
}

// SetShardGroupingCost sets get shard grouping context cost
func (s *StorageStats) SetShardGroupingCost(shardID int32, cost int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.GroupingCost = cost
	}
}

// SetShardScanStats sets data scan cost in shard level
func (s *StorageStats) SetShardScanStats(shardID int32, identifier string, cost int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.SetScanStats(identifier, cost)
	}
}

// SetShardGroupBuildStats sets grouping build stats in shard level
func (s *StorageStats) SetShardGroupBuildStats(shardID int32, cost int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.SetGroupBuildStats(cost)
	}
}

// SetCollectTagValuesStats sets collect tag values stats after search for group by query
func (s *StorageStats) SetCollectTagValuesStats(tagKey string, cost int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.CollectTagValuesStats[tagKey] = cost
}

// ShardStats represents the shard level stats
type ShardStats struct {
	SeriesFilterCost int64             `json:"seriesFilterCost"`
	NumOfSeries      uint64            `json:"numOfSeries"`
	MemFilterCost    int64             `json:"memFilterCost"`
	KVFilterCost     int64             `json:"kvFilterCost"`
	GroupingCost     int64             `json:"groupingCost"`
	ScanStats        map[string]*Stats `json:"scanStats,omitempty"`
	GroupBuildStats  *Stats            `json:"groupBuildStats,omitempty"`
}

// newShardStats creates the shard level stats
func newShardStats() *ShardStats {
	return &ShardStats{
		ScanStats: make(map[string]*Stats),
	}
}

// SetScanStats sets the data scan stats in shard level
func (s *ShardStats) SetScanStats(identifier string, cost int64) {
	stats, ok := s.ScanStats[identifier]
	if ok {
		stats.Count++
		if stats.Max < cost {
			stats.Max = cost
		} else if stats.Min > cost {
			stats.Min = cost
		}
	} else {
		s.ScanStats[identifier] = &Stats{
			Min:   cost,
			Max:   cost,
			Count: 1,
		}
	}
}

// SetGroupBuildStats sets the group build stats in shard level
func (s *ShardStats) SetGroupBuildStats(cost int64) {
	if s.GroupBuildStats == nil {
		s.GroupBuildStats = &Stats{
			Min:   cost,
			Max:   cost,
			Count: 1,
		}
	} else {
		s.GroupBuildStats.Count++
		if s.GroupBuildStats.Max < cost {
			s.GroupBuildStats.Max = cost
		} else if s.GroupBuildStats.Min > cost {
			s.GroupBuildStats.Min = cost
		}
	}
}

// Stats represents the time stats
type Stats struct {
	Min   int64 `json:"min"`
	Max   int64 `json:"max"`
	Count int   `json:"count"`
}
