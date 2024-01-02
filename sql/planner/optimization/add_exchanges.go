package optimization

import (
	"fmt"
	"reflect"

	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

type AddExchangesPlan struct {
	node  plan.PlanNode
	props *ActualProps
}

type AddExchanges struct {
}

func NewAddExchanges() PlanOptimizer {
	return &AddExchanges{}
}

// Optimize implements PlanOptimizer
func (opt *AddExchanges) Optimize(plan plan.PlanNode, idAllocator *plan.PlanNodeIDAllocator) plan.PlanNode {
	result := plan.Accept(&PreferredProps{}, &AddExchangesRewrite{
		idAllocator: idAllocator,
	})
	if planProps, ok := result.(*AddExchangesPlan); ok {
		return planProps.node
	}
	//FIXME: need remove
	return plan
}

type AddExchangesRewrite struct {
	idAllocator *plan.PlanNodeIDAllocator
}

func (v *AddExchangesRewrite) Visit(context any, n plan.PlanNode) (r any) {
	fmt.Printf("exchange rewrite=%s\n", reflect.TypeOf(n))
	parentProps := context.(*PreferredProps)
	switch node := n.(type) {
	case *plan.OutputNode:
		return v.VisitOutput(parentProps, node)
	case *plan.JoinNode:
		return v.VisitJoin(parentProps, node)
	case *plan.ProjectionNode:
		return v.VisitProjection(parentProps, node)
	case *plan.TableScanNode:
		return v.VisitTableScan(parentProps, node)
	default:
		return v.rebaseAndDeriveProps(n, v.planChild(n, parentProps))
	}
}

func (v *AddExchangesRewrite) VisitOutput(context any, node *plan.OutputNode) (r any) {
	child := v.planChild(node, undistributed())
	//FIXME:??? check sigle
	child = v.withDerivedProps(plan.GatheringExchange(v.idAllocator.Next(), plan.Remote, child.node), child.props)
	return v.rebaseAndDeriveProps(node, child)
}

func (v *AddExchangesRewrite) VisitJoin(context any, node *plan.JoinNode) (r any) {
	return v.planPartitionedJoin(node)
}

func (v *AddExchangesRewrite) VisitTableScan(context any, node *plan.TableScanNode) (r any) {
	return &AddExchangesPlan{
		node: node,
	}
}

func (v *AddExchangesRewrite) VisitProjection(context any, node *plan.ProjectionNode) (r any) {
	//FIXME: translate
	return v.rebaseAndDeriveProps(node, v.planChild(node, context.(*PreferredProps)))
}

func (v *AddExchangesRewrite) VisitFilter(node *plan.FilterNode, context any) (r any) {
	preferredProps := context.(*PreferredProps)
	if tableScan, ok := node.Source.(*plan.TableScanNode); ok {
		planNode := iterative.PushFilterIntoTableScan(node, tableScan)
		if planNode != nil {
			return &AddExchangesPlan{
				node: planNode,
			}
		}
	}

	return v.rebaseAndDeriveProps(node, v.planChild(node, preferredProps))
}

func (v *AddExchangesRewrite) planPartitionedJoin(node *plan.JoinNode) *AddExchangesPlan {
	left := node.Left.Accept(partitioned(), v).(*AddExchangesPlan)
	right := node.Right.Accept(partitioned(), v).(*AddExchangesPlan)

	left = v.withDerivedProps(plan.PartitionedExchange(v.idAllocator.Next(), plan.Remote, left.node), left.props)
	right = v.withDerivedProps(plan.PartitionedExchange(v.idAllocator.Next(), plan.Remote, right.node), right.props)

	return v.buildJoin(node, left, right, plan.Partitioned)
}

func (v *AddExchangesRewrite) planChild(node plan.PlanNode, preferredProps *PreferredProps) *AddExchangesPlan {
	return node.GetSources()[0].Accept(preferredProps, v).(*AddExchangesPlan)
}

func (v *AddExchangesRewrite) rebaseAndDeriveProps(node plan.PlanNode, child *AddExchangesPlan) *AddExchangesPlan {
	return v.withDerivedProps(plan.ReplaceChildren(node, []plan.PlanNode{child.node}), child.props)
}

func (v *AddExchangesRewrite) withDerivedProps(node plan.PlanNode, inputProps *ActualProps) *AddExchangesPlan {
	//FIXME::::
	return &AddExchangesPlan{
		node: node,
	}
}

func (v *AddExchangesRewrite) buildJoin(node *plan.JoinNode, newLeft, newRight *AddExchangesPlan, newDistributionType plan.DistributionType) *AddExchangesPlan {

	result := plan.JoinNode{
		BaseNode: plan.BaseNode{
			ID: node.GetNodeID(),
		},
		Type:             node.Type,
		DistributionType: newDistributionType,
		Left:             newLeft.node,
		Right:            newRight.node,
		Criteria:         node.Criteria,
	}
	return &AddExchangesPlan{
		node: &result,
	}
}
