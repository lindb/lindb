package rule

import (
	"fmt"
	"sort"

	"github.com/lindb/lindb/sql/matching"
	"github.com/lindb/lindb/sql/planner/iterative"
	"github.com/lindb/lindb/sql/planner/plan"
)

// RemoveRedundantIdentityProjections removes projection nodes that only perform non-renaming identity projections.
type RemoveRedundantIdentityProjections struct{}

func NewRemoveRedundantIdentityProjections() iterative.Rule {
	return &RemoveRedundantIdentityProjections{}
}

func (rule *RemoveRedundantIdentityProjections) GetPattern() *matching.Pattern {
	return project()
}

func (rule *RemoveRedundantIdentityProjections) Apply(context *iterative.Context, captures *matching.Captures, node plan.PlanNode) plan.PlanNode {
	if project, ok := node.(*plan.ProjectionNode); ok {
		fmt.Printf("remove identity project...............................................%v,%v\n", project.Assignments.IsIdentity(),
			symbolsEquals(project.GetOutputSymbols(), project.Source.GetOutputSymbols()),
		)
		if project.Assignments.IsIdentity() &&
			symbolsEquals(project.GetOutputSymbols(), project.Source.GetOutputSymbols()) {
			return project.Source
		}
	}
	return nil
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
