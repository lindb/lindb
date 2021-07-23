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
)

func TestPhysicalPlan(t *testing.T) {
	physicalPlan := NewPhysicalPlan(Root{Indicator: "1.1.1.3:8000", NumOfTask: 1})
	physicalPlan.AddLeaf(Leaf{
		BaseNode: BaseNode{
			Parent:    "1.1.1.2:8000",
			Indicator: "1.1.1.1:9000",
		},
		Receivers: []StatelessNode{{HostIP: "1.1.1.5", GRPCPort: 8000}},
		ShardIDs:  []ShardID{1, 2, 4},
	})
	physicalPlan.AddIntermediate(Intermediate{
		BaseNode: BaseNode{
			Parent:    "1.1.1.3:8000",
			Indicator: "1.1.1.2:8000",
		},
		NumOfTask: 1,
	})
	physicalPlan.Database = "test_db"

	assert.Equal(t, PhysicalPlan{
		Database: "test_db",
		Root:     Root{Indicator: "1.1.1.3:8000", NumOfTask: 1},
		Intermediates: []Intermediate{{
			BaseNode: BaseNode{
				Parent:    "1.1.1.3:8000",
				Indicator: "1.1.1.2:8000",
			},
			NumOfTask: 1}},
		Leafs: []Leaf{{
			BaseNode: BaseNode{
				Parent:    "1.1.1.2:8000",
				Indicator: "1.1.1.1:9000",
			},
			Receivers: []StatelessNode{{HostIP: "1.1.1.5", GRPCPort: 8000}},
			ShardIDs:  []ShardID{1, 2, 4},
		}},
	}, *physicalPlan)
}
