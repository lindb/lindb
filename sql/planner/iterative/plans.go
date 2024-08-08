package iterative

import (
	lo "github.com/samber/lo"

	"github.com/lindb/lindb/sql/planner/plan"
)

type ResolvingVisitor struct {
	lookup func(groupRef *plan.GroupReference) plan.PlanNode
}

func NewResolvingVisitor(lookup func(groupRef *plan.GroupReference) plan.PlanNode) plan.Visitor {
	return ResolvingVisitor{
		lookup: lookup,
	}
}

func (v ResolvingVisitor) Visit(context any, n plan.PlanNode) (r any) {
	switch node := n.(type) {
	case *plan.GroupReference:
		pNode := v.lookup(node)
		return pNode.Accept(context, v)
	default:
		newChildren := lo.Map(node.GetSources(), func(child plan.PlanNode, index int) plan.PlanNode {
			return child.Accept(context, v).(plan.PlanNode)
		})
		return node.ReplaceChildren(newChildren)
	}
}

func resolveGroupReferences(node plan.PlanNode, lookup func(groupRef *plan.GroupReference) plan.PlanNode) plan.PlanNode {
	return node.Accept(nil, NewResolvingVisitor(lookup)).(plan.PlanNode)
}
