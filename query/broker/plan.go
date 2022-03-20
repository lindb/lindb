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

package brokerquery

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	querypkg "github.com/lindb/lindb/query"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// brokerPlan represents the broker execute plan
type brokerPlan struct {
	query             *stmtpkg.Query
	storageNodes      map[string][]models.ShardID
	currentBrokerNode models.StatelessNode
	brokerNodes       []models.StatelessNode
	intermediateNodes []models.StatelessNode
	databaseCfg       models.Database

	physicalPlan *models.PhysicalPlan
}

// newBrokerPlan creates broker execute plan
func newBrokerPlan(
	query *stmtpkg.Query,
	databaseCfg models.Database,
	storageNodes map[string][]models.ShardID,
	currentBrokerNode models.StatelessNode,
	brokerNodes []models.StatelessNode,
) *brokerPlan {
	return &brokerPlan{
		query:             query,
		databaseCfg:       databaseCfg,
		storageNodes:      storageNodes,
		currentBrokerNode: currentBrokerNode,
		brokerNodes:       brokerNodes,
	}
}

// Plan plans broker level query execute plan, there are some scenarios as below:
// 1) parse sql => stmt
// 2) build parallel exec tree
//    a) no group by => only need leafs
//    b) one storage node => only need leafs
//    c) no other active broker node => node need leafs
//    d) need intermediate computing nodes
func (p *brokerPlan) Plan() error {
	lenOfStorageNodes := len(p.storageNodes)
	if lenOfStorageNodes == 0 {
		return querypkg.ErrNoAvailableStorageNode
	}

	if p.query.Interval <= 0 {
		// TODO need get by time range
		p.query.Interval = p.databaseCfg.Option.Intervals[0].Interval
	}
	intervalVal := int64(p.query.Interval)
	p.query.TimeRange.Start = timeutil.Truncate(p.query.TimeRange.Start, intervalVal)
	p.query.TimeRange.End = timeutil.Truncate(p.query.TimeRange.End, intervalVal)

	root := p.currentBrokerNode

	p.buildIntermediateNodes()

	lenOfIntermediateNodes := len(p.intermediateNodes)

	if lenOfIntermediateNodes > 0 {
		// create parallel exec task
		p.physicalPlan = models.NewPhysicalPlan(models.Root{
			Indicator: root.Indicator(),
			NumOfTask: int32(lenOfIntermediateNodes)})

		p.buildIntermediates()
	} else {
		receivers := []models.StatelessNode{root}
		// create parallel exec task
		p.physicalPlan = models.NewPhysicalPlan(models.Root{
			Indicator: root.Indicator(),
			NumOfTask: int32(lenOfStorageNodes)})
		p.buildLeafs(root.Indicator(), p.getStorageNodeIDs(), receivers)
	}
	return nil
}

// buildIntermediateNodes builds intermediate nodes if need
func (p *brokerPlan) buildIntermediateNodes() {
	if len(p.query.GroupBy) == 0 {
		return
	}
	if len(p.brokerNodes) == 0 {
		return
	}
	if len(p.storageNodes) == 1 {
		return
	}

	for _, brokerNode := range p.brokerNodes {
		if brokerNode.Indicator() != p.currentBrokerNode.Indicator() {
			p.intermediateNodes = append(p.intermediateNodes, brokerNode)
		}
	}
}

// getStorageNodeIDs returns storage node ids
func (p *brokerPlan) getStorageNodeIDs() []string {
	var storageNodeIDs []string
	for nodeID := range p.storageNodes {
		storageNodeIDs = append(storageNodeIDs, nodeID)
	}
	return storageNodeIDs
}

// buildIntermediates builds the intermediates computing layer
func (p *brokerPlan) buildIntermediates() {
	lenOfIntermediateNodes := len(p.intermediateNodes)
	lenOfStorageNodes := len(p.storageNodes)
	// calc degree of parallelism
	parallel := lenOfStorageNodes / lenOfIntermediateNodes
	if lenOfStorageNodes%lenOfIntermediateNodes != 0 {
		parallel++
	}

	storageNodeIDs := p.getStorageNodeIDs()

	var pos, end, idx = 0, 0, 0
	for {
		if pos > lenOfStorageNodes {
			break
		}
		end += parallel

		if end > lenOfStorageNodes {
			end = lenOfStorageNodes
		}
		if idx >= lenOfIntermediateNodes {
			break
		}
		intermediateNodeID := p.intermediateNodes[idx].Indicator()

		// add intermediate task into parallel exec tree
		p.physicalPlan.AddIntermediate(models.Intermediate{
			BaseNode: models.BaseNode{
				Parent:    p.currentBrokerNode.Indicator(),
				Indicator: intermediateNodeID,
			},
			NumOfTask: int32(lenOfStorageNodes),
		})
		// add leaf tasks into parallel exec tree
		p.buildLeafs(intermediateNodeID, storageNodeIDs[pos:end], p.intermediateNodes)

		pos += parallel
		idx++
	}
}

// buildLeafs builds the leaf computing nodes based parent, nodes and result receivers
func (p *brokerPlan) buildLeafs(parentID string, nodeIDs []string, receivers []models.StatelessNode) {
	for _, nodeID := range nodeIDs {
		leaf := &models.Leaf{
			BaseNode: models.BaseNode{
				Parent:    parentID,
				Indicator: nodeID,
			},
			ShardIDs:  p.storageNodes[nodeID],
			Receivers: receivers,
		}
		p.physicalPlan.AddLeaf(leaf)
	}
}
