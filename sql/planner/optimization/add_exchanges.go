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
	"reflect"

	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
)

type AddExchangesPlan struct {
	node  plan.PlanNode
	props *ActualProps
}

type AddExchanges struct{}

func NewAddExchanges() PlanOptimizer {
	return &AddExchanges{}
}

// Optimize implements PlanOptimizer
func (opt *AddExchanges) Optimize(ctx *context.PlannerContext, plan plan.PlanNode) plan.PlanNode {
	result := plan.Accept(&PreferredProps{}, &AddExchangesRewrite{
		idAllocator: ctx.PlanNodeIDAllocator,
	})
	if planProps, ok := result.(*AddExchangesPlan); ok {
		fmt.Printf("after exchange rewrite=%v\n", planProps.node)
		return planProps.node
	}
	// FIXME: need remove
	fmt.Printf("after exchange rewrite=%v\n", plan)
	return plan
}

type AddExchangesRewrite struct {
	idAllocator *plan.PlanNodeIDAllocator
}

func (v *AddExchangesRewrite) Visit(context any, n plan.PlanNode) (r any) {
	fmt.Printf("222exchange rewrite=%s\n", reflect.TypeOf(n))
	parentProps := context.(*PreferredProps)
	switch node := n.(type) {
	case *plan.OutputNode:
		return v.visitOutput(parentProps, node)
	case *plan.JoinNode:
		return v.visitJoin(parentProps, node)
	case *plan.ProjectionNode:
		return v.visitProjection(parentProps, node)
	case *plan.FilterNode:
		return v.visitFilter(parentProps, node)
	case *plan.TableScanNode:
		return v.visitTableScan(parentProps, node)
	case *plan.AggregationNode:
		return v.visitAggregation(parentProps, node)
	case *plan.ValuesNode:
		return v.visitValues(parentProps, node)
	default:
		return v.rebaseAndDeriveProps(n, v.planChild(n, parentProps))
	}
}

func (v *AddExchangesRewrite) visitValues(_ any, node *plan.ValuesNode) (r any) {
	return &AddExchangesPlan{
		node:  node,
		props: NewActualPropsBuilder(singlePartition()).Build(),
	}
}

func (v *AddExchangesRewrite) visitOutput(_ any, node *plan.OutputNode) (r any) {
	child := v.planChild(node, Undistributed())
	return v.rebaseAndDeriveProps(node, child)
}

func (v *AddExchangesRewrite) visitJoin(_ any, node *plan.JoinNode) (r any) {
	return v.planPartitionedJoin(node)
}

func (v *AddExchangesRewrite) visitFilter(_ any, node *plan.FilterNode) (r any) {
	fmt.Printf("node=%v,child=%v\n", node, node.GetSources())
	if _, ok := node.GetSources()[0].(*plan.TableScanNode); ok {
		fmt.Println("1231232...")
		child := &AddExchangesPlan{
			node:  node,
			props: v.dervieProps(node, nil),
		}
		return v.rebaseAndDeriveProps(plan.GatheringExchange(v.idAllocator.Next(), plan.Remote, node), child)
	}
	return v.rebaseAndDeriveProps(node, v.planChild(node, Any()))
}

func (v *AddExchangesRewrite) visitTableScan(_ any, node *plan.TableScanNode) (r any) {
	child := &AddExchangesPlan{
		node:  node,
		props: v.dervieProps(node, nil),
	}
	return v.rebaseAndDeriveProps(plan.GatheringExchange(v.idAllocator.Next(), plan.Remote, node), child)
}

func (v *AddExchangesRewrite) visitProjection(context any, node *plan.ProjectionNode) (r any) {
	// FIXME: translate
	return v.rebaseAndDeriveProps(node, v.planChild(node, context.(*PreferredProps)))
}

func (v *AddExchangesRewrite) visitAggregation(_ any, node *plan.AggregationNode) (r any) {
	preferSingleNode := node.IsSingleNodeExecutionPreference()
	var preferredProps *PreferredProps
	if preferSingleNode {
		preferredProps = Undistributed()
	} else {
		preferredProps = Any()
	}

	child := v.planChild(node, preferredProps)
	return v.rebaseAndDeriveProps(node, child)
}

func (v *AddExchangesRewrite) planPartitionedJoin(node *plan.JoinNode) *AddExchangesPlan {
	left := node.Left.Accept(Partitioned(), v).(*AddExchangesPlan)
	right := node.Right.Accept(Partitioned(), v).(*AddExchangesPlan)

	// TODO: set partitioning scheme
	// left = v.withDerivedProps(plan.PartitionedExchange(v.idAllocator.Next(), plan.Remote, left.node, &plan.PartitioningScheme{OutputLayout: left.node.GetOutputSymbols()}), left.props)
	// right = v.withDerivedProps(plan.PartitionedExchange(v.idAllocator.Next(), plan.Remote, right.node, &plan.PartitioningScheme{OutputLayout: right.node.GetOutputSymbols()}), right.props)

	return v.buildJoin(node, left, right, plan.Partitioned)
	// return &AddExchangesPlan{
	// 	node:  node,
	// 	props: v.dervieProps(node, []*ActualProps{left.props, right.props}),
	// }
}

func (v *AddExchangesRewrite) planChild(node plan.PlanNode, preferredProps *PreferredProps) *AddExchangesPlan {
	child := node.GetSources()[0].Accept(preferredProps, v).(*AddExchangesPlan)
	fmt.Printf("add exchange child====%T\n", child.node)
	return child
}

func (v *AddExchangesRewrite) rebaseAndDeriveProps(node plan.PlanNode, child *AddExchangesPlan) *AddExchangesPlan {
	return v.withDerivedProps(plan.ReplaceChildren(node, []plan.PlanNode{child.node}), child.props)
}

func (v *AddExchangesRewrite) withDerivedProps(node plan.PlanNode, inputProps *ActualProps) *AddExchangesPlan {
	// FIXME::::
	return &AddExchangesPlan{
		node:  node,
		props: v.dervieProps(node, []*ActualProps{inputProps}),
	}
}

func (v *AddExchangesRewrite) dervieProps(node plan.PlanNode, inputProperties []*ActualProps) *ActualProps {
	return deriveProps(node, inputProperties)
}

func (v *AddExchangesRewrite) buildJoin(node *plan.JoinNode,
	newLeft, newRight *AddExchangesPlan, newDistributionType plan.DistributionType,
) *AddExchangesPlan {
	result := &plan.JoinNode{
		BaseNode: plan.BaseNode{
			ID: node.GetNodeID(),
		},
		Type:               node.Type,
		DistributionType:   newDistributionType,
		Left:               newLeft.node,
		Right:              newRight.node,
		LeftOutputSymbols:  newLeft.node.GetOutputSymbols(),
		RightOutputSymbols: newRight.node.GetOutputSymbols(),
		Criteria:           node.Criteria,
	}
	return &AddExchangesPlan{
		node:  result,
		props: v.dervieProps(result, []*ActualProps{newLeft.props, newRight.props}),
	}
}
