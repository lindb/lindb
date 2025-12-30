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

package master

import (
	"slices"

	"github.com/lindb/lindb/models"
)

// RangeConsumeAssign assigns shards to consumers using range assignment strategy.
// It distributes shards as evenly as possible among consumers in a round-robin fashion.
// Each consumer gets either floor(nShards/nConsumers) or ceil(nShards/nConsumers) shards.
// 1. Sorts both consumers and shardIDs to ensure consistent assignment order.
// 2. Calculates the base number of shards per consumer and the remainder.
// 3. Iterates over each consumer, assigning them their calculated number of shards.
// 4. Returns a slice of ConsumeAssignment, each containing a consumer ID and their assigned shards.
func RangeConsumeAssign(consumers []models.NodeID, shardIDs []models.ShardID) (assignments []models.ConsumeAssignment) {
	if len(consumers) == 0 || len(shardIDs) == 0 {
		// no consumers or no shards, return empty assignments
		return
	}

	slices.Sort(consumers)
	slices.Sort(shardIDs)

	nConsumers := len(consumers)
	nShards := len(shardIDs)
	perConsumer := nShards / nConsumers
	perConsumerRemainder := nShards % nConsumers

	offset := 0
	for i, c := range consumers {
		num := perConsumer
		if i < perConsumerRemainder {
			num++
		}

		start := offset
		end := offset + num
		offset = end

		assignments = append(assignments, models.ConsumeAssignment{
			ConsumerID: c,
			Shards:     shardIDs[start:end],
		})
	}

	return
}
