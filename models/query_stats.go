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
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

// QueryStats represents the query stats when need explain query flow stat
type QueryStats struct {
	BrokerNodes  map[string]*QueryStats   `json:"brokerNodes,omitempty"`
	StorageNodes map[string]*StorageStats `json:"storageNodes,omitempty"`
	PlanCost     int64                    `json:"planCost,omitempty"`
	WaitCost     int64                    `json:"waitCost,omitempty"` // wait intermediate or leaf response duration
	ExpressCost  int64                    `json:"expressCost,omitempty"`
	TotalCost    int64                    `json:"totalCost,omitempty"` // total query cost
}

// NewQueryStats creates the query stats
func NewQueryStats() *QueryStats {
	return &QueryStats{
		BrokerNodes:  make(map[string]*QueryStats),
		StorageNodes: make(map[string]*StorageStats),
	}
}

// MergeBrokerTaskStats merges intermediate task execution stats
func (s *QueryStats) MergeBrokerTaskStats(nodeID string, stats *QueryStats) {
	s.BrokerNodes[nodeID] = stats
}

// MergeStorageTaskStats merges storage task execution stats
func (s *QueryStats) MergeStorageTaskStats(nodeID string, stats *StorageStats) {
	s.StorageNodes[nodeID] = stats
}

// StorageStats represents query stats in storage side
type StorageStats struct {
	NetPayload            ltoml.Size              `json:"netPayload"`
	TotalCost             int64                   `json:"totalCost"`
	PlanCost              int64                   `json:"planCost"`
	TagFilterCost         int64                   `json:"tagFilterCost"`
	Shards                map[ShardID]*ShardStats `json:"shards,omitempty"`
	CollectTagValuesStats map[string]int64        `json:"collectTagValuesStats,omitempty"`

	start time.Time  // track search start time in storage side
	mutex sync.Mutex // need add lock for goroutine update stats data
}

// NewStorageStats creates the query stats in storage side
func NewStorageStats() *StorageStats {
	return &StorageStats{
		Shards:                make(map[ShardID]*ShardStats),
		CollectTagValuesStats: make(map[string]int64),
		start:                 time.Now(),
	}
}

// Complete completes the query stats when query completed
func (s *StorageStats) Complete() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TotalCost = time.Since(s.start).Nanoseconds()
}

// SetPlanCost sets plan cost
func (s *StorageStats) SetPlanCost(cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.PlanCost = cost.Nanoseconds()
}

// SetTagFilterCost sets tag filter cost
func (s *StorageStats) SetTagFilterCost(cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TagFilterCost = cost.Nanoseconds()
}

// SetShardSeriesIDsSearchStats sets shard series ids search stats
func (s *StorageStats) SetShardSeriesIDsSearchStats(shardID ShardID, numOfSeries uint64, seriesFilterCost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	stats := newShardStats()
	stats.NumOfSeries = numOfSeries
	stats.SeriesFilterCost = seriesFilterCost.Nanoseconds()
	s.Shards[shardID] = stats
}

// SetShardMemoryDataFilterCost sets shard memory data filter cost
func (s *StorageStats) SetShardMemoryDataFilterCost(shardID ShardID, cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.MemFilterCost = cost.Nanoseconds()
	}
}

// SetShardKVDataFilterCost sets shard data filter cost in kv store
func (s *StorageStats) SetShardKVDataFilterCost(shardID ShardID, cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.KVFilterCost = cost.Nanoseconds()
	}
}

// SetShardGroupingCost sets get shard grouping context cost
func (s *StorageStats) SetShardGroupingCost(shardID ShardID, cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.GroupingCost = cost.Nanoseconds()
	}
}

// SetShardScanStats sets data scan cost in shard level
func (s *StorageStats) SetShardScanStats(shardID ShardID, identifier string, cost time.Duration, foundSeries int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.SetScanStats(identifier, cost, foundSeries)
	}
}

// SetShardGroupBuildStats sets grouping build stats in shard level
func (s *StorageStats) SetShardGroupBuildStats(shardID ShardID, cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.SetGroupBuildStats(cost)
	}
}

// SetCollectTagValuesStats sets collect tag values stats after search for group by query
func (s *StorageStats) SetCollectTagValuesStats(tagKey string, cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.CollectTagValuesStats[tagKey] = cost.Nanoseconds()
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
func (s *ShardStats) SetScanStats(identifier string, cost time.Duration, foundSeries int) {
	costVal := cost.Nanoseconds()
	stats, ok := s.ScanStats[identifier]
	if ok {
		stats.Count++
		stats.Series += foundSeries
		if stats.Max < costVal {
			stats.Max = costVal
		} else if stats.Min > costVal {
			stats.Min = costVal
		}
	} else {
		s.ScanStats[identifier] = &Stats{
			Min:    costVal,
			Max:    costVal,
			Count:  1,
			Series: foundSeries,
		}
	}
}

// SetGroupBuildStats sets the group build stats in shard level
func (s *ShardStats) SetGroupBuildStats(cost time.Duration) {
	costVal := cost.Nanoseconds()
	if s.GroupBuildStats == nil {
		s.GroupBuildStats = &Stats{
			Min:   costVal,
			Max:   costVal,
			Count: 1,
		}
	} else {
		s.GroupBuildStats.Count++
		if s.GroupBuildStats.Max < costVal {
			s.GroupBuildStats.Max = costVal
		} else if s.GroupBuildStats.Min > costVal {
			s.GroupBuildStats.Min = costVal
		}
	}
}

// Stats represents the time stats
type Stats struct {
	Min    int64 `json:"min"`
	Max    int64 `json:"max"`
	Count  int   `json:"count"`
	Series int   `json:"series,omitempty"`
}
