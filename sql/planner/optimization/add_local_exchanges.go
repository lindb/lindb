package optimization

import (
	"fmt"
	"reflect"

	"github.com/lindb/lindb/sql/planner/plan"
)

type AddLocalExchanges struct {
}

func NewAddLocalExchanges() PlanOptimizer {
	return &AddLocalExchanges{}
}

// Optimize implements PlanOptimizer
func (opt *AddLocalExchanges) Optimize(plan plan.PlanNode, idAllocator *plan.PlanNodeIDAllocator) plan.PlanNode {
	result := plan.Accept(&StreamPreferredProps{}, &AddLocalExchangesRewrite{
		idAllocator: idAllocator,
	})
	if planProps, ok := result.(*PlanProps); ok {
		return planProps.node
	}
	//FIXME: need remove
	return plan
}

type AddLocalExchangesRewrite struct {
	idAllocator *plan.PlanNodeIDAllocator
}

func (v *AddLocalExchangesRewrite) Visit(context any, n plan.PlanNode) (r any) {
	fmt.Printf("logic exchange rewrite=%s\n", reflect.TypeOf(n))
	parentProps := context.(*StreamPreferredProps)
	switch node := n.(type) {
	case *plan.OutputNode:
		return v.visitOutput(parentProps, node)
	case *plan.ExchangeNode:
		return v.visitExchange(parentProps, node)
	case *plan.JoinNode:
		return v.visitJoin(parentProps, node)
	default:
		return v.planAndEnforceChildren(n,
			parentProps.withoutPreference().withDefaultParallelism(),
			parentProps.withDefaultParallelism(),
		)
	}
}

func (v *AddLocalExchangesRewrite) VisitPlan(context any, node plan.PlanNode) (r any) {
	fmt.Printf("logic exchange rewrite=%s\n", reflect.TypeOf(node))
	parentProps := context.(*StreamPreferredProps)
	return v.planAndEnforceChildren(node,
		parentProps.withoutPreference().withDefaultParallelism(),
		parentProps.withDefaultParallelism(),
	)
}

func (v *AddLocalExchangesRewrite) visitOutput(context any, node *plan.OutputNode) (r any) {
	return v.planAndEnforceChildren(node, empty().withOrderSensitivity(), empty().withOrderSensitivity())
}

func (v *AddLocalExchangesRewrite) visitJoin(context any, node *plan.JoinNode) (r any) {
	parentProps := context.(*StreamPreferredProps)
	probe := v.planAndEnforce(node.Left,
		defaultParallelism(),
		parentProps.constrainTo(node.GetOutputSymbols()).withDefaultParallelism())
	//FIXME:::
	buildProps := singleStream()
	build := v.planAndEnforce(node.Right, buildProps, buildProps)
	return v.rebaseAndDeriveProps(node, []*PlanProps{probe, build})
}

func (v *AddLocalExchangesRewrite) visitExchange(context any, node *plan.ExchangeNode) (r any) {
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
	//FIXME: verify properties are in terms of symbols produced by the node?

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
		exchangeNode := plan.GatheringExchange(v.idAllocator.Next(), plan.Local, planProps.node)
		return v.deriveProps(exchangeNode, []*StreamProps{planProps.props})
	}

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
		node: result,
	}
}
