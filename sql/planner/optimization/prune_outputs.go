package optimization

import (
	"github.com/samber/lo"

	planpkg "github.com/lindb/lindb/sql/planner/plan"
)

type PruneOutputs struct{}

func NewPruneOutputs() PlanOptimizer {
	return &PruneOutputs{}
}

// Optimize implements PlanOptimizer
func (opt *PruneOutputs) Optimize(node planpkg.PlanNode, idAllocator *planpkg.PlanNodeIDAllocator) planpkg.PlanNode {
	outputs := node.GetOutputSymbols()
	result := node.Accept(outputs, &PruneOutputsVisitor{
		idAllocator: idAllocator,
	})
	if r, ok := result.(planpkg.PlanNode); ok {
		return r
	}
	// FIXME: need remove
	return node
}

type PruneOutputsVisitor struct {
	idAllocator *planpkg.PlanNodeIDAllocator
}

func (v *PruneOutputsVisitor) Visit(context any, n planpkg.PlanNode) (r any) {
	permittedOutputs := context.([]*planpkg.Symbol)
	switch node := n.(type) {
	case *planpkg.ProjectionNode:
		restrictedOutputs := restrictOutputs(node, permittedOutputs)
		var assigments planpkg.Assignments
		node.Assignments = assigments.Add(restrictedOutputs)
	case *planpkg.TableScanNode:
		restrictedOutputs := restrictOutputs(node, permittedOutputs)
		node.OutputSymbols = restrictedOutputs
	}

	for _, child := range n.GetSources() {
		child.Accept(context, v)
	}
	return n
}

func restrictOutputs(node planpkg.PlanNode, permittedOutputs []*planpkg.Symbol) (newOutputs []*planpkg.Symbol) {
	outputs := node.GetOutputSymbols()
	// restrict outputs based on permitted outputs
	for _, output := range permittedOutputs {
		o, ok := lo.Find(outputs, func(item *planpkg.Symbol) bool {
			return item.Name == output.Name
		})
		if ok {
			newOutputs = append(newOutputs, o)
		}
	}
	return
}
