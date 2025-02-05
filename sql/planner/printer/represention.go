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

package printer

import (
	"github.com/lindb/lindb/sql/planner/plan"
)

type PlanRepresentation struct {
	root     plan.PlanNode
	nodeInfo map[plan.PlanNodeID]*NodeRepresentation
}

func NewPlanRepresentation(root plan.PlanNode) *PlanRepresentation {
	return &PlanRepresentation{
		root:     root,
		nodeInfo: make(map[plan.PlanNodeID]*NodeRepresentation),
	}
}

func (pr *PlanRepresentation) getRoot() (nr *NodeRepresentation) {
	nr = pr.nodeInfo[pr.root.GetNodeID()]
	return nr
}

func (pr *PlanRepresentation) getNode(nodeID plan.PlanNodeID) (node *NodeRepresentation) {
	node = pr.nodeInfo[nodeID]
	return
}

func (pr *PlanRepresentation) addNode(node *NodeRepresentation) {
	pr.nodeInfo[node.getID()] = node
}

type NodeRepresentation struct {
	descriptor map[string]string
	name       string
	children   []plan.PlanNodeID
	details    []string
	outputs    []*plan.Symbol
	id         plan.PlanNodeID
}

func (nr *NodeRepresentation) appendDetails(detail string) {
	nr.details = append(nr.details, detail)
}

func (nr *NodeRepresentation) getID() plan.PlanNodeID {
	return nr.id
}

func (nr *NodeRepresentation) getName() string {
	return nr.name
}
