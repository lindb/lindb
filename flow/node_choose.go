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

package flow

import (
	"math/rand"
	"time"

	"github.com/lindb/lindb/models"
)

//go:generate mockgen -source=./node_choose.go -destination=./node_choose_mock.go -package=flow

// NodeChoose represents node choose for data query.
type NodeChoose interface {
	// Choose chooses the compute nodes then builds physical plan.
	Choose(database string, numOfNodes int) ([]*models.PhysicalPlan, error)
}

// BuildPhysicalPlan returns physical plan based on live nodes and node number, need shuffle live node.
func BuildPhysicalPlan(database string, liveNodes []models.StatelessNode, numOfNodes int) *models.PhysicalPlan {
	physicalPlan := &models.PhysicalPlan{
		Database: database,
	}
	numOfLiveNodes := len(liveNodes)
	if numOfLiveNodes > 0 {
		// shuffle broker nodes
		rand.Seed(time.Now().Unix())
		rand.Shuffle(numOfLiveNodes, func(i, j int) {
			liveNodes[i], liveNodes[j] = liveNodes[j], liveNodes[i]
		})
		for i, node := range liveNodes {
			if i == numOfNodes {
				break
			}
			receiveOnly := true
			if i == 0 {
				receiveOnly = false
			}
			physicalPlan.AddTarget(&models.Target{
				Indicator:   node.Indicator(),
				ReceiveOnly: receiveOnly,
			})
		}
	}
	return physicalPlan
}
