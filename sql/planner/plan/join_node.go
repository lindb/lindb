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

package plan

import (
	"github.com/lindb/lindb/sql/tree"
)

type JoinType string

type DistributionType string

var (
	Inner JoinType = "InnerJoin"
	Left  JoinType = "LeftJoin"
	Right JoinType = "RightJoin"
	Full  JoinType = "FullJoin"

	Partitioned DistributionType = "Partitioned"
)

type EqualJoinCriteria struct {
	Left  *Symbol
	Right *Symbol
}

func (n *EqualJoinCriteria) ToExpression() *tree.ComparisonExpression {
	return &tree.ComparisonExpression{
		Left:     n.Left.ToSymbolReference(),
		Operator: tree.ComparisonEQ,
		Right:    n.Right.ToSymbolReference(),
	}
}

type JoinNode struct {
	Left               PlanNode
	Right              PlanNode
	DistributionType   DistributionType
	Type               JoinType
	Criteria           []*EqualJoinCriteria
	LeftOutputSymbols  []*Symbol // FIXME: remove it
	RightOutputSymbols []*Symbol
	Filter             tree.Expression // FIXME: set expression

	BaseNode
}

func (n *JoinNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *JoinNode) GetSources() []PlanNode {
	return []PlanNode{n.Left, n.Right}
}

func (n *JoinNode) GetOutputSymbols() (output []*Symbol) {
	output = append(output, n.Left.GetOutputSymbols()...)
	output = append(output, n.Right.GetOutputSymbols()...)
	return output
}

func (n *JoinNode) IsCrossJoin() bool {
	return len(n.Criteria) == 0 && n.Filter == nil && n.Type == Inner
}

func (n *JoinNode) Clone() *JoinNode {
	return n.ReplaceChildren(n.GetSources()).(*JoinNode)
}

func (n *JoinNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return &JoinNode{
		BaseNode: BaseNode{
			ID: n.GetNodeID(),
		},
		Left:               newChildren[0],
		Right:              newChildren[1],
		LeftOutputSymbols:  n.LeftOutputSymbols,
		RightOutputSymbols: n.RightOutputSymbols,
		DistributionType:   n.DistributionType,
		Type:               n.Type,
		Criteria:           n.Criteria,
		Filter:             n.Filter,
	}
}

func JoinTypeConvert(joinType tree.JoinType) JoinType {
	switch joinType {
	case tree.CROSS, tree.IMPLICIT, tree.INNER:
		return Inner
	case tree.LEFT:
		return Left
	case tree.RIGHT:
		return Right
	case tree.FULL:
		return Full
	default:
		panic("unsupport join type")
	}
}
