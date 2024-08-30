package optimization

import (
	"fmt"

	"github.com/samber/lo"

	planpkg "github.com/lindb/lindb/sql/planner/plan"
)

type PruneColumns struct{}

func NewPruneColumns() PlanOptimizer {
	return &PruneColumns{}
}

// Optimize implements PlanOptimizer
func (opt *PruneColumns) Optimize(node planpkg.PlanNode, idAllocator *planpkg.PlanNodeIDAllocator) planpkg.PlanNode {
	outputs := node.GetOutputSymbols()
	result := node.Accept(outputs, &PruneColumnsVisitor{
		idAllocator: idAllocator,
	})
	if r, ok := result.(planpkg.PlanNode); ok {
		return r
	}
	// FIXME: need remove
	return node
}

type PruneColumnsVisitor struct {
	idAllocator *planpkg.PlanNodeIDAllocator
}

func (v *PruneColumnsVisitor) Visit(context any, n planpkg.PlanNode) (r any) {
	permittedOutputs := context.([]*planpkg.Symbol)
	fmt.Printf("table scan permitted outputs,%T=%v\n", n, permittedOutputs)
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
		child.Accept(n.GetOutputSymbols(), v)
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
