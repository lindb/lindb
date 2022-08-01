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
	"time"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/query/operator"
)

//go:generate mockgen -source=./plan_node.go -destination=./plan_node_mock.go -package=stage

// PlanNode represents the node of plan tree.
type PlanNode interface {
	// Execute executes the operator of current node.
	Execute() error
	// ExecuteWithStats executes the operator of current node with stats.
	ExecuteWithStats() (stats *models.OperatorStats, err error)
	// Children returns the children nodes of current node.
	Children() []PlanNode
	// AddChild adds child node.
	AddChild(node PlanNode)
	// IgnoreNotFound returns if the stage ignore not found error.
	IgnoreNotFound() bool
}

// planNode implements PlanNode interface.
type planNode struct {
	op       operator.Operator
	children []PlanNode

	ignore bool
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

// NewPlanNodeWithIgnore creates a PlanNode with operator, need ignore not found error.
func NewPlanNodeWithIgnore(op operator.Operator) PlanNode {
	return &planNode{
		op:     op,
		ignore: true,
	}
}

// Execute executes the operator of current node.
func (p *planNode) Execute() error {
	if p.op == nil {
		return nil
	}
	return p.op.Execute()
}

// ExecuteWithStats executes the operator of current node with stats.
func (p *planNode) ExecuteWithStats() (stats *models.OperatorStats, err error) {
	if p.op == nil {
		return nil, nil
	}
	start := time.Now()
	defer func() {
		end := time.Now()
		stats = &models.OperatorStats{
			Identifier: p.op.Identifier(),
			Start:      start.UnixMilli(),
			End:        end.UnixMilli(),
			Cost:       end.Sub(start).Nanoseconds(),
		}
		if track, ok := p.op.(operator.TrackableOperator); ok {
			stats.Stats = track.Stats()
		}
	}()
	err = p.op.Execute()
	return
}

// Children returns the children nodes of current node.
func (p *planNode) Children() []PlanNode {
	return p.children
}

// AddChild adds child node.
func (p *planNode) AddChild(child PlanNode) {
	p.children = append(p.children, child)
}

// IgnoreNotFound returns if the stage ignore not found error.
func (p *planNode) IgnoreNotFound() bool {
	return p.ignore
}
