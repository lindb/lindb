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

package optimization

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/planner/plan"
)

func deriveStreamProps(node plan.PlanNode, inputProps []*StreamProps) (result *StreamProps) {
	result = node.Accept(inputProps, &StreamPropsDerivationVisitor{}).(*StreamProps)
	// TODO: add check
	return
}

type StreamPropsDerivationVisitor struct{}

func (v *StreamPropsDerivationVisitor) Visit(context any, n plan.PlanNode) (r any) {
	inputProps := context.([]*StreamProps)
	switch node := n.(type) {
	case *plan.OutputNode:
		return v.visitOutput(inputProps, node)
	case *plan.AggregationNode:
		return v.visitAggregation(inputProps, node)
	case *plan.ExchangeNode:
		return v.visitExchange(inputProps, node)
	case *plan.ProjectionNode:
		return v.visitProjection(inputProps, node)
	case *plan.FilterNode:
		return v.visitFilter(inputProps, node)
	case *plan.TableScanNode:
		return v.visitTableScan(inputProps, node)
	case *plan.ValuesNode:
		return v.visitValues(inputProps, node)
	case *plan.JoinNode:
		return v.visitJoin(inputProps, node)
	default:
		panic(fmt.Sprintf("impl stream props derivation visitor:%T", n))
	}
}

func (v *StreamPropsDerivationVisitor) visitJoin(inputProps []*StreamProps, _ *plan.JoinNode) *StreamProps {
	return inputProps[0]
}

func (v *StreamPropsDerivationVisitor) visitFilter(inputProps []*StreamProps, _ *plan.FilterNode) *StreamProps {
	return inputProps[0]
}

func (v *StreamPropsDerivationVisitor) visitProjection(inputProps []*StreamProps, node *plan.ProjectionNode) *StreamProps {
	props := inputProps[0]
	identities := computeIdentityTranslations(node.Assignments)
	return props.translate(func(column *plan.Symbol) *plan.Symbol {
		return identities[column.Name]
	})
}

func (v *StreamPropsDerivationVisitor) visitOutput(inputProps []*StreamProps, node *plan.OutputNode) *StreamProps {
	return inputProps[0].translate(func(column *plan.Symbol) *plan.Symbol {
		return filterIfMissing(node.Outputs, column)
	})
}

func (v *StreamPropsDerivationVisitor) visitAggregation(inputProps []*StreamProps, node *plan.AggregationNode) *StreamProps {
	props := inputProps[0]
	groupingKeys := node.GetGroupingKeys()
	return props.translate(func(column *plan.Symbol) *plan.Symbol {
		if lo.ContainsBy(groupingKeys, func(item *plan.Symbol) bool {
			return column.Name == item.Name
		}) {
			return column
		}
		return nil
	})
}

func (v *StreamPropsDerivationVisitor) visitExchange(_ []*StreamProps, node *plan.ExchangeNode) *StreamProps {
	if node.Scope == plan.Remote {
		return FixedStreams()
	}
	switch node.Type {
	case plan.Gather:
		return SingleStream()
	case plan.Repartition:
		// FIXME: add partition handle check
		return FixedStreams()
	}
	return nil
}

func (v *StreamPropsDerivationVisitor) visitTableScan(_ []*StreamProps, _ *plan.TableScanNode) *StreamProps {
	// FIXME: add partition
	return &StreamProps{
		distribution: Multiple,
	}
}

func (v *StreamPropsDerivationVisitor) visitValues(_ []*StreamProps, _ *plan.ValuesNode) *StreamProps {
	return &StreamProps{
		distribution: Single,
	}
}
