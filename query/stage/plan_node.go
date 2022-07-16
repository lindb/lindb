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

package stage

import (
	"github.com/lindb/lindb/query/operator"
)

//go:generate mockgen -source=./plan_node.go -destination=./plan_node_mock.go -package=stage

// PlanNode represents the node of plan tree.
type PlanNode interface {
	// Execute executes the operator of current node.
	Execute() error
	// Children returns the children nodes of current node.
	Children() []PlanNode
	// AddChild adds child node.
	AddChild(node PlanNode)
}

// planNode implements PlanNode interface.
type planNode struct {
	op operator.Operator

	children []PlanNode
}

// NewEmptyPlanNode creates a PlanNode without operator.
func NewEmptyPlanNode() PlanNode {
	return &planNode{}
}

// NewPlanNode creates a PlanNode with operator.
func NewPlanNode(op operator.Operator) PlanNode {
	return &planNode{
		op: op,
	}
}

// Execute executes the operator of current node.
func (p *planNode) Execute() error {
	if p.op == nil {
		return nil
	}
	return p.op.Execute()
}

// Children returns the children nodes of current node.
func (p *planNode) Children() []PlanNode {
	return p.children
}

// AddChild adds child node.
func (p *planNode) AddChild(child PlanNode) {
	p.children = append(p.children, child)
}
