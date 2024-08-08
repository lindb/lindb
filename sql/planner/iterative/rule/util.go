package rule

import (
	"fmt"

	lo "github.com/samber/lo"

	"github.com/lindb/lindb/sql/planner/plan"
)

func restrictOutputs(idAllcator *plan.PlanNodeIDAllocator, node plan.PlanNode, permittedOutputs []*plan.Symbol) plan.PlanNode {
	outputs := node.GetOutputSymbols()
	restrictedOutputs := lo.Filter(outputs, func(item *plan.Symbol, index int) bool {
		return lo.ContainsBy(permittedOutputs, func(other *plan.Symbol) bool {
			return other.Name == item.Name
		})
	})
	if len(outputs) == len(restrictedOutputs) {
		fmt.Println("outputs same.....")
		return nil
	}
	fmt.Printf("restrictedOutputs, a=%v,b=%v,c=%v\n", outputs, restrictedOutputs, permittedOutputs)

	var assigments plan.Assignments
	assigments = assigments.Add(restrictedOutputs)

	return &plan.ProjectionNode{
		BaseNode: plan.BaseNode{
			ID: idAllcator.Next(),
		},
		Source:      node,
		Assignments: assigments,
	}
}

func restrictChildOutputs(idAllcator *plan.PlanNodeIDAllocator, node plan.PlanNode, permittedChildOutputs ...[]*plan.Symbol) plan.PlanNode {
	if len(node.GetSources()) != len(permittedChildOutputs) {
		panic(fmt.Sprintf("mismatched child (%d) and permitted outputs (%d) sizes",
			len(node.GetSources()), len(permittedChildOutputs)))
	}

	var newChildren []plan.PlanNode
	rewriteChildren := false

	for i, oldChild := range node.GetSources() {
		newChild := restrictOutputs(idAllcator, oldChild, permittedChildOutputs[i])
		if newChild != nil {
			rewriteChildren = true
			newChildren = append(newChildren, newChild)
		} else {
			newChildren = append(newChildren, oldChild)
		}
	}

	if !rewriteChildren {
		return nil
	}

	return node.ReplaceChildren(newChildren)
}
