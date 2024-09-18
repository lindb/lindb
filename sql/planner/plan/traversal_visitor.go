package plan

import "fmt"

type DefaultTraversalVisitor struct {
	Process func(node PlanNode)
	Resolve func(node PlanNode) PlanNode
}

func (v *DefaultTraversalVisitor) Visit(context any, n PlanNode) (r any) {
	curNode := n
	switch node := n.(type) {
	case *OutputNode:
		_ = node.Source.Accept(context, v)
	case *AggregationNode:
		_ = node.Source.Accept(context, v)
	case *FilterNode:
		_ = node.Source.Accept(context, v)
	case *ProjectionNode:
		_ = node.Source.Accept(context, v)
	case *GroupReference:
		// need resolve group reference(raw plan node)
		if v.Resolve != nil {
			rawNode := v.Resolve(node)
			// set current node using raw node
			curNode = rawNode
			sources := rawNode.GetSources()
			for _, source := range sources {
				_ = source.Accept(context, v)
			}
		}
	default:
		// TODO: remove
		fmt.Printf("plan node default traversal visitor not support..................=%T\n", n)
	}
	if v.Process != nil {
		v.Process(curNode)
	}
	return
}
