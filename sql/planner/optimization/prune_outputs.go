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
		assigments := make(planpkg.Assignments)
		assigments.Add(restrictedOutputs)
		node.Assignments = assigments
	case *planpkg.TableScanNode:
		restrictedOutputs := restrictOutputs(node, permittedOutputs)
		node.OutputSymbols = restrictedOutputs
	}

	for _, child := range n.GetSources() {
		child.Accept(context, v)
	}
	return n
}

func restrictOutputs(node planpkg.PlanNode, permittedOutputs []*planpkg.Symbol) []*planpkg.Symbol {
	outputs := node.GetOutputSymbols()
	return lo.Filter(outputs, func(item *planpkg.Symbol, index int) bool {
		return lo.ContainsBy(permittedOutputs, func(other *planpkg.Symbol) bool {
			return other.Name == item.Name
		})
	})
}
