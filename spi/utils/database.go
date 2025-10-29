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

package utils

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/store"
)

func FindSegments(db store.Database, shardIDs []int,
	interval timeutil.Interval, timeRange timeutil.TimeRange,
	fn func(shard store.Shard, partition store.Partition, segments []store.Segment),
) {
	storageInterval := db.FindMatchSmallestInterval(interval)
	for _, id := range shardIDs {
		shard, ok := db.GetShard(models.ShardID(id))
		if ok {
			pList := shard.GetPartitions(storageInterval, timeRange)
			if len(pList) > 0 {
				for _, partition := range pList {
					segments := partition.GetSegments(timeRange)
					if len(segments) > 0 {
						fn(shard, partition, segments)
					}
				}
			}
		}
	}
}
