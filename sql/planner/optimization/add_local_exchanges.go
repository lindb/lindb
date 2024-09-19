package optimization

import (
	"fmt"
	"reflect"

	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
)

type AddLocalExchanges struct{}

func NewAddLocalExchanges() PlanOptimizer {
	return &AddLocalExchanges{}
}

// Optimize implements PlanOptimizer
func (opt *AddLocalExchanges) Optimize(ctx *context.PlannerContext, plan plan.PlanNode) plan.PlanNode {
	result := plan.Accept(&StreamPreferredProps{}, &AddLocalExchangesRewrite{
		idAllocator: ctx.PlanNodeIDAllocator,
	})
	if planProps, ok := result.(*PlanProps); ok {
		return planProps.node
	}
	// FIXME: need remove
	return plan
}

type AddLocalExchangesRewrite struct {
	idAllocator *plan.PlanNodeIDAllocator
}

func (v *AddLocalExchangesRewrite) Visit(context any, n plan.PlanNode) (r any) {
	fmt.Printf("add local exchange rewrite=%s,%v\n", reflect.TypeOf(n), n)
	parentProps := context.(*StreamPreferredProps)
	switch node := n.(type) {
	case *plan.OutputNode:
		return v.visitOutput(parentProps, node)
	case *plan.ExchangeNode:
		return v.visitExchange(parentProps, node)
	case *plan.AggregationNode:
		return v.visitAggregation(parentProps, node)
	case *plan.JoinNode:
		return v.visitJoin(parentProps, node)
	default:
		return v.planAndEnforceChildren(n,
			parentProps.withoutPreference().withDefaultParallelism(),
			parentProps.withDefaultParallelism(),
		)
	}
}

//	func (v *AddLocalExchangesRewrite) VisitPlan(context any, node plan.PlanNode) (r any) {
//		panic("kkkkkk....")
//		fmt.Printf("logic exchange rewrite=%s\n", reflect.TypeOf(node))
//		parentProps := context.(*StreamPreferredProps)
//		return v.planAndEnforceChildren(node,
//			parentProps.withoutPreference().withDefaultParallelism(),
//			parentProps.withDefaultParallelism(),
//		)
//	}
func (v *AddLocalExchangesRewrite) visitAggregation(parentProps *StreamPreferredProps, node *plan.AggregationNode) *PlanProps {
	if node.IsSingleNodeExecutionPreference() {
		return v.planAndEnforceChildren(node, singleStream(), defaultParallelism())
	}
	groupingKeys := node.GetGroupingKeys()

	childRequirements := parentProps.constrainTo(node.Source.GetOutputSymbols()).withDefaultParallelism().withPartitioning(groupingKeys)
	child := v.planAndEnforce(node.Source, childRequirements, childRequirements)
	fmt.Printf("agg child:=%v\n", child.node)
	result := plan.NewAggregationNode(node.GetNodeID(), child.node, node.Aggregations, node.GroupingSets, node.Step)
	return v.deriveProps(result, []*StreamProps{child.props})
}

func (v *AddLocalExchangesRewrite) visitOutput(context any, node *plan.OutputNode) *PlanProps {
	return v.planAndEnforceChildren(node, empty().withOrderSensitivity(), empty().withOrderSensitivity())
}

func (v *AddLocalExchangesRewrite) visitJoin(parentProps *StreamPreferredProps, node *plan.JoinNode) *PlanProps {
	probe := v.planAndEnforce(node.Left,
		defaultParallelism(),
		parentProps.constrainTo(node.GetOutputSymbols()).withDefaultParallelism())
	// FIXME:::
	buildProps := singleStream()
	build := v.planAndEnforce(node.Right, buildProps, buildProps)
	return v.rebaseAndDeriveProps(node, []*PlanProps{probe, build})
}

func (v *AddLocalExchangesRewrite) visitExchange(context any, node *plan.ExchangeNode) *PlanProps {
	if node.Scope == plan.Local {
		panic("add local exchange cannot process a plan containing a local exchange")
	}
	return v.planAndEnforceChildren(node, empty(), defaultParallelism())
}

func (v *AddLocalExchangesRewrite) planAndEnforceChildren(node plan.PlanNode, requiredProps, preferredProps *StreamPreferredProps) *PlanProps {
	sources := node.GetSources()
	var children []*PlanProps
	for i := range sources {
		source := sources[i]
		child := v.planAndEnforce(source,
			requiredProps.constrainTo(source.GetOutputSymbols()),
			preferredProps.constrainTo(source.GetOutputSymbols()))
		children = append(children, child)
	}
	return v.rebaseAndDeriveProps(node, children)
}

func (v *AddLocalExchangesRewrite) planAndEnforce(node plan.PlanNode, requiredProps, preferredProps *StreamPreferredProps) *PlanProps {
	// FIXME: verify properties are in terms of symbols produced by the node?

	// plan the node using the preferred props
	result := node.Accept(preferredProps, v).(*PlanProps)
	// enforce the required props
	result = v.enforce(result, requiredProps)
	return result
}

func (v *AddLocalExchangesRewrite) enforce(planProps *PlanProps, requiredProps *StreamPreferredProps) *PlanProps {
	if requiredProps.isSatisfiedBy(planProps.props) {
		return planProps
	}

	if requiredProps.isSingleStreamPreferred() {
		fmt.Println("single stream.......")
		exchangeNode := plan.GatheringExchange(v.idAllocator.Next(), plan.Local, planProps.node)
		return v.deriveProps(exchangeNode, []*StreamProps{planProps.props})
	}

	fmt.Println("other stream.......")
	// no explicit parallel requirement, so gather to a single stream
	exchangeNode := plan.GatheringExchange(v.idAllocator.Next(), plan.Local, planProps.node)
	return v.deriveProps(exchangeNode, []*StreamProps{planProps.props})
}

func (v *AddLocalExchangesRewrite) rebaseAndDeriveProps(node plan.PlanNode, children []*PlanProps) *PlanProps {
	var childrenNode []plan.PlanNode
	var inputProps []*StreamProps
	for i := range children {
		childrenNode = append(childrenNode, children[i].node)
		inputProps = append(inputProps, children[i].props)
	}

	result := plan.ReplaceChildren(node, childrenNode)

	return v.deriveProps(result, inputProps)
}

func (v *AddLocalExchangesRewrite) deriveProps(result plan.PlanNode, inputProps []*StreamProps) *PlanProps {
	return &PlanProps{
		node:  result,
		props: deriveStreamProps(result, inputProps),
	}
}
