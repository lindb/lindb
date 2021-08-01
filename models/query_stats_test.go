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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/ltoml"
)

func TestNewQueryStats(t *testing.T) {
	stats := NewQueryStats()
	stats.MergeStorageTaskStats("task-1", NewStorageStats())
	storageStats, ok := stats.StorageNodes["task-1"]
	assert.NotNil(t, storageStats)
	assert.True(t, ok)
}

func TestStorageStats(t *testing.T) {
	stats := NewStorageStats()
	stats.SetPlanCost(10)
	stats.SetTagFilterCost(10)
	stats.SetCollectTagValuesStats("test-1", 10)
	stats.SetShardGroupBuildStats(10, 10)
	stats.SetShardScanStats(10, "id", 10)
	stats.SetShardGroupingCost(10, 10)
	stats.SetShardKVDataFilterCost(10, 10)
	stats.SetShardMemoryDataFilterCost(10, 10)
	shard, ok := stats.Shards[10]
	assert.False(t, ok)
	assert.Nil(t, shard)
	assert.Equal(t, ltoml.Duration(10), stats.PlanCost)
	assert.Equal(t, ltoml.Duration(10), stats.TagFilterCost)

	stats.SetShardSeriesIDsSearchStats(10, 10, 10)
	stats.SetCollectTagValuesStats("test-1", 10)
	stats.SetShardGroupBuildStats(10, 10)
	stats.SetShardScanStats(10, "id", 10)
	stats.SetShardGroupBuildStats(10, 10)
	stats.SetShardScanStats(10, "id", 10)
	stats.SetShardGroupingCost(10, 10)
	stats.SetShardKVDataFilterCost(10, 10)
	stats.SetShardMemoryDataFilterCost(10, 10)
	stats.Complete()
	assert.True(t, stats.TotalCost > 0)
	shard, ok = stats.Shards[10]
	assert.True(t, ok)
	assert.NotNil(t, shard)

	assert.Equal(t, ltoml.Duration(10), shard.SeriesFilterCost)
	assert.Equal(t, ltoml.Duration(10), shard.MemFilterCost)
	assert.Equal(t, ltoml.Duration(10), shard.KVFilterCost)
	assert.Equal(t, ltoml.Duration(10), shard.GroupBuildStats.Max)
	assert.Equal(t, ltoml.Duration(10), shard.GroupBuildStats.Min)
	assert.Equal(t, 2, shard.GroupBuildStats.Count)
	scan, ok := shard.ScanStats["id"]
	assert.True(t, ok)
	assert.Equal(t, ltoml.Duration(10), scan.Max)
	assert.Equal(t, ltoml.Duration(10), scan.Min)
	assert.Equal(t, 2, scan.Count)
}

func TestShardStats(t *testing.T) {
	stats := newShardStats()
	stats.SetGroupBuildStats(10)
	stats.SetGroupBuildStats(20)
	stats.SetGroupBuildStats(5)
	assert.Equal(t, 3, stats.GroupBuildStats.Count)
	assert.Equal(t, ltoml.Duration(5), stats.GroupBuildStats.Min)
	assert.Equal(t, ltoml.Duration(20), stats.GroupBuildStats.Max)

	stats.SetScanStats("id", 10)
	stats.SetScanStats("id", 20)
	stats.SetScanStats("id", 5)
	s := stats.ScanStats["id"]
	assert.Equal(t, 3, s.Count)
	assert.Equal(t, ltoml.Duration(5), s.Min)
	assert.Equal(t, ltoml.Duration(20), s.Max)
}
