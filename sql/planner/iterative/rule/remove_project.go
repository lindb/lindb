package rule

import (
	"fmt"
	"sort"

	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

// RemoveRedundantIdentityProjections removes projection nodes that only perform non-renaming identity projections.
type RemoveRedundantIdentityProjections struct {
	Base[*plan.ProjectionNode]
}

func NewRemoveRedundantIdentityProjections() iterative.Rule {
	rule := &RemoveRedundantIdentityProjections{}
	rule.apply = func(context *iterative.Context, node *plan.ProjectionNode) plan.PlanNode {
		fmt.Printf("remove projection node: %v,%v,source=%T\n", node.Assignments, node.Assignments.IsIdentity(), node.Source)
		if node.Assignments.IsIdentity() &&
			symbolsEquals(node.GetOutputSymbols(), node.Source.GetOutputSymbols()) {
			return node.Source
		}
		return nil
	}
	return rule
}

func symbolsEquals(a, b []*plan.Symbol) bool {
	fmt.Printf("check symbols========a=%v,b=%v\n", a, b)
	if len(a) != len(b) {
		return false
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i].Name > a[j].Name
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i].Name > b[j].Name
	})

	for i := range a {
		if a[i].Name != b[i].Name {
			return false
		}
	}

	return true
}
