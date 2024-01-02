package iterative

import (
	"github.com/lindb/lindb/sql/planner/plan"
	lo "github.com/samber/lo"
)

type ResolvingVisitor struct {
	lookup func(groupRef *GroupReference) plan.PlanNode
}

func NewResolvingVisitor(lookup func(groupRef *GroupReference) plan.PlanNode) plan.Visitor {
	return ResolvingVisitor{
		lookup: lookup,
	}
}

func (v ResolvingVisitor) Visit(context any, n plan.PlanNode) (r any) {
	switch node := n.(type) {
	case *GroupReference:
		return v.lookup(node).Accept(context, v)
	default:
		newChildren := lo.Map(node.GetSources(), func(item plan.PlanNode, index int) plan.PlanNode {
			return item.Accept(context, v).(plan.PlanNode)
		})
		return node.ReplaceChildren(newChildren)
	}
}

func resolveGroupReferences(node plan.PlanNode, lookup func(groupRef *GroupReference) plan.PlanNode) plan.PlanNode {
	return node.Accept(nil, NewResolvingVisitor(lookup)).(plan.PlanNode)
}
