package plan

func ReplaceChildren(node PlanNode, children []PlanNode) PlanNode {
	sources := node.GetSources()
	for i := range sources {
		child := sources[0]
		if children[i] != child {
			return node.ReplaceChildren(children)
		}
	}
	return node
}
