package iterative

import (
	"fmt"

	"github.com/lindb/lindb/sql/planner/plan"
)

type Lookup interface {
	Resolve(node plan.PlanNode) plan.PlanNode
}

type lookup struct {
	resolver func(groupRef *plan.GroupReference) []plan.PlanNode
}

func NewLookup(resolver func(groupRef *plan.GroupReference) []plan.PlanNode) Lookup {
	return &lookup{
		resolver: resolver,
	}
}

func (l *lookup) Resolve(node plan.PlanNode) plan.PlanNode {
	if groupRef, ok := node.(*plan.GroupReference); ok {
		fmt.Printf("resolve lokkup88888%v\n", l.resolver(groupRef))
		return l.resolver(groupRef)[0] // FIXME: add check
	}
	return node
}
