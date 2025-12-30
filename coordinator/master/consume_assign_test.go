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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
)

func TestRangeConsumeAssign(t *testing.T) {
	tests := []struct {
		name      string
		consumers []models.NodeID
		shards    []models.ShardID
		expected  []models.ConsumeAssignment
	}{
		{
			name:      "3 consumers, 10 partitions",
			consumers: []models.NodeID{0, 1, 2},
			shards:    []models.ShardID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			expected: []models.ConsumeAssignment{
				{ConsumerID: 0, Shards: []models.ShardID{0, 1, 2, 3}},
				{ConsumerID: 1, Shards: []models.ShardID{4, 5, 6}},
				{ConsumerID: 2, Shards: []models.ShardID{7, 8, 9}},
			},
		},
		{
			name:      "4 consumers, 10 partitions",
			consumers: []models.NodeID{0, 1, 2, 3},
			shards:    []models.ShardID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			expected: []models.ConsumeAssignment{
				{ConsumerID: 0, Shards: []models.ShardID{0, 1, 2}},
				{ConsumerID: 1, Shards: []models.ShardID{3, 4, 5}},
				{ConsumerID: 2, Shards: []models.ShardID{6, 7}},
				{ConsumerID: 3, Shards: []models.ShardID{8, 9}},
			},
		},
		{
			name:      "2 consumers, 5 partitions",
			consumers: []models.NodeID{0, 1},
			shards:    []models.ShardID{0, 1, 2, 3, 4},
			expected: []models.ConsumeAssignment{
				{ConsumerID: 0, Shards: []models.ShardID{0, 1, 2}},
				{ConsumerID: 1, Shards: []models.ShardID{3, 4}},
			},
		},
		{
			name:      "5 consumers, 3 partitions",
			consumers: []models.NodeID{0, 1, 2, 3, 4},
			shards:    []models.ShardID{0, 1, 2},
			expected: []models.ConsumeAssignment{
				{ConsumerID: 0, Shards: []models.ShardID{0}},
				{ConsumerID: 1, Shards: []models.ShardID{1}},
				{ConsumerID: 2, Shards: []models.ShardID{2}},
				{ConsumerID: 3, Shards: []models.ShardID{}},
				{ConsumerID: 4, Shards: []models.ShardID{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RangeConsumeAssign(tt.consumers, tt.shards)
			if len(result) != len(tt.expected) {
				assert.Fail(t, "expected %d assignments, got %d", len(tt.expected), len(result))
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
