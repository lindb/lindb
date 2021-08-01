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
	PlanCost     ltoml.Duration           `json:"planCost,omitempty"`
	WaitCost     ltoml.Duration           `json:"waitCost,omitempty"` // wait intermediate or leaf response duration
	ExpressCost  ltoml.Duration           `json:"expressCost,omitempty"`
	TotalCost    ltoml.Duration           `json:"totalCost,omitempty"` // total query cost
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
	NetPayload            ltoml.Size                `json:"netPayload"`
	TotalCost             ltoml.Duration            `json:"totalCost"`
	PlanCost              ltoml.Duration            `json:"planCost"`
	TagFilterCost         ltoml.Duration            `json:"tagFilterCost"`
	Shards                map[int32]*ShardStats     `json:"shards,omitempty"`
	CollectTagValuesStats map[string]ltoml.Duration `json:"collectTagValuesStats,omitempty"`

	start time.Time  // track search start time in storage side
	mutex sync.Mutex // need add lock for goroutine update stats data
}

// NewStorageStats creates the query stats in storage side
func NewStorageStats() *StorageStats {
	return &StorageStats{
		Shards:                make(map[int32]*ShardStats),
		CollectTagValuesStats: make(map[string]ltoml.Duration),
		start:                 time.Now(),
	}
}

// Complete completes the query stats when query completed
func (s *StorageStats) Complete() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TotalCost = ltoml.Duration(time.Since(s.start))
}

// SetPlanCost sets plan cost
func (s *StorageStats) SetPlanCost(cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.PlanCost = ltoml.Duration(cost)
}

// SetTagFilterCost sets tag filter cost
func (s *StorageStats) SetTagFilterCost(cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TagFilterCost = ltoml.Duration(cost)
}

// SetShardSeriesIDsSearchStats sets shard series ids search stats
func (s *StorageStats) SetShardSeriesIDsSearchStats(shardID int32, numOfSeries uint64, seriesFilterCost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	stats := newShardStats()
	stats.NumOfSeries = numOfSeries
	stats.SeriesFilterCost = ltoml.Duration(seriesFilterCost)
	s.Shards[shardID] = stats
}

// SetShardMemoryDataFilterCost sets shard memory data filter cost
func (s *StorageStats) SetShardMemoryDataFilterCost(shardID int32, cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.MemFilterCost = ltoml.Duration(cost)
	}
}

// SetShardKVDataFilterCost sets shard data filter cost in kv store
func (s *StorageStats) SetShardKVDataFilterCost(shardID int32, cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.KVFilterCost = ltoml.Duration(cost)
	}
}

// SetShardGroupingCost sets get shard grouping context cost
func (s *StorageStats) SetShardGroupingCost(shardID int32, cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.GroupingCost = ltoml.Duration(cost)
	}
}

// SetShardScanStats sets data scan cost in shard level
func (s *StorageStats) SetShardScanStats(shardID int32, identifier string, cost time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stats, ok := s.Shards[shardID]
	if ok {
		stats.SetScanStats(identifier, cost)
	}
}

// SetShardGroupBuildStats sets grouping build stats in shard level
func (s *StorageStats) SetShardGroupBuildStats(shardID int32, cost time.Duration) {
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
	s.CollectTagValuesStats[tagKey] = ltoml.Duration(cost)
}

// ShardStats represents the shard level stats
type ShardStats struct {
	SeriesFilterCost ltoml.Duration    `json:"seriesFilterCost"`
	NumOfSeries      uint64            `json:"numOfSeries"`
	MemFilterCost    ltoml.Duration    `json:"memFilterCost"`
	KVFilterCost     ltoml.Duration    `json:"kvFilterCost"`
	GroupingCost     ltoml.Duration    `json:"groupingCost"`
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
func (s *ShardStats) SetScanStats(identifier string, cost time.Duration) {
	stats, ok := s.ScanStats[identifier]
	if ok {
		stats.Count++
		if stats.Max < ltoml.Duration(cost) {
			stats.Max = ltoml.Duration(cost)
		} else if stats.Min > ltoml.Duration(cost) {
			stats.Min = ltoml.Duration(cost)
		}
	} else {
		s.ScanStats[identifier] = &Stats{
			Min:   ltoml.Duration(cost),
			Max:   ltoml.Duration(cost),
			Count: 1,
		}
	}
}

// SetGroupBuildStats sets the group build stats in shard level
func (s *ShardStats) SetGroupBuildStats(cost time.Duration) {
	if s.GroupBuildStats == nil {
		s.GroupBuildStats = &Stats{
			Min:   ltoml.Duration(cost),
			Max:   ltoml.Duration(cost),
			Count: 1,
		}
	} else {
		s.GroupBuildStats.Count++
		if s.GroupBuildStats.Max < ltoml.Duration(cost) {
			s.GroupBuildStats.Max = ltoml.Duration(cost)
		} else if s.GroupBuildStats.Min > ltoml.Duration(cost) {
			s.GroupBuildStats.Min = ltoml.Duration(cost)
		}
	}
}

// Stats represents the time stats
type Stats struct {
	Min   ltoml.Duration `json:"min"`
	Max   ltoml.Duration `json:"max"`
	Count int            `json:"count"`
}
