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

type AddExchanges struct{}

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
	// FIXME: need remove
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
		return v.visitOutput(parentProps, node)
	case *plan.JoinNode:
		return v.visitJoin(parentProps, node)
	case *plan.ProjectionNode:
		return v.visitProjection(parentProps, node)
	case *plan.TableScanNode:
		return v.visitTableScan(parentProps, node)
	case *plan.AggregationNode:
		return v.visitAggregation(parentProps, node)
	default:
		return v.rebaseAndDeriveProps(n, v.planChild(n, parentProps))
	}
}

func (v *AddExchangesRewrite) visitOutput(context any, node *plan.OutputNode) (r any) {
	child := v.planChild(node, Undistributed())
	// 	// FIXME:??? check sigle/force single node output
	// if !child.props.isSingleNode() {
	// 	child = v.withDerivedProps(plan.GatheringExchange(v.idAllocator.Next(), plan.Remote, child.node), child.props)
	// }
	return v.rebaseAndDeriveProps(node, child)
}

func (v *AddExchangesRewrite) visitJoin(context any, node *plan.JoinNode) (r any) {
	return v.planPartitionedJoin(node)
}

func (v *AddExchangesRewrite) visitTableScan(context any, node *plan.TableScanNode) (r any) {
	return &AddExchangesPlan{
		node:  node,
		props: v.dervieProps(node, nil),
	}
}

func (v *AddExchangesRewrite) visitProjection(context any, node *plan.ProjectionNode) (r any) {
	// FIXME: translate
	return v.rebaseAndDeriveProps(node, v.planChild(node, context.(*PreferredProps)))
}

func (v *AddExchangesRewrite) visitFilter(node *plan.FilterNode, context any) (r any) {
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

func (v *AddExchangesRewrite) visitAggregation(context any, node *plan.AggregationNode) (r any) {
	parentPreferredProps := context.(*PreferredProps)
	partitioningRequirement := node.GetGroupingKeys() // TODO: cope it?
	preferSingleNode := node.IsSingleNodeExecutionPreference()
	var preferredProps *PreferredProps
	if preferSingleNode {
		preferredProps = Undistributed()
	} else {
		preferredProps = Any()
	}

	if len(node.GetGroupingKeys()) > 0 {
		preferredProps = v.computePreference(PartitionedWithLocal(partitioningRequirement), parentPreferredProps)
	}

	child := v.planChild(node, preferredProps)
	if child.props.isSingleNode() {
		return v.rebaseAndDeriveProps(node, child)
	}
	if preferSingleNode {
		child = v.withDerivedProps(plan.GatheringExchange(v.idAllocator.Next(), plan.Remote, child.node), child.props)
	} else {
		// TODO: partition keys
		child = v.withDerivedProps(plan.PartitionedExchange(v.idAllocator.Next(), plan.Remote, child.node), child.props)
	}
	return v.rebaseAndDeriveProps(node, child)
}

func (v *AddExchangesRewrite) planPartitionedJoin(node *plan.JoinNode) *AddExchangesPlan {
	left := node.Left.Accept(Partitioned(), v).(*AddExchangesPlan)
	right := node.Right.Accept(Partitioned(), v).(*AddExchangesPlan)

	left = v.withDerivedProps(plan.PartitionedExchange(v.idAllocator.Next(), plan.Remote, left.node), left.props)
	right = v.withDerivedProps(plan.PartitionedExchange(v.idAllocator.Next(), plan.Remote, right.node), right.props)

	return v.buildJoin(node, left, right, plan.Partitioned)
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

func (a *AddExchangesRewrite) computePreference(preferredProps, parentPreferredProperties *PreferredProps) *PreferredProps {
	// TODO: check ignore down stream preferences

	return preferredProps
}
