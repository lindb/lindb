package plan

type Plan struct {
	Root PlanNode
}

type SubPlan struct {
	Fragment *PlanFragment
	Children []*SubPlan
}

// GetAllFragments flattens and returns all plan fragment in the plan tree.
func (p *SubPlan) GetAllFragments() (result []*PlanFragment) {
	result = append(result, p.Fragment)
	for i := range p.Children {
		result = append(result, p.Children[i].GetAllFragments()...)
	}
	return
}
