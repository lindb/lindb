package plan

import (
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/tree"
)

type Assignment struct {
	Symbol     *Symbol         `json:"symbol"`     // ouput
	Expression tree.Expression `json:"expression"` // input
}

type Assignments []*Assignment

func (a Assignments) Add(symbols []*Symbol) Assignments {
	// TODO: check dup
	for _, symbol := range symbols {
		a = append(a, &Assignment{
			Symbol:     symbol,
			Expression: symbol.ToSymbolReference(),
		})
	}
	return a
}

func (a Assignments) Put(symbol *Symbol, expression tree.Expression) Assignments {
	a = append(a, &Assignment{
		Symbol:     symbol,
		Expression: expression,
	})
	return a
}

func (a Assignments) GetExpressions() (r []tree.Expression) {
	for _, assignment := range a {
		r = append(r, assignment.Expression)
	}
	return
}

func (a Assignments) GetOutputs() (outputs []*Symbol) {
	for _, assignment := range a {
		outputs = append(outputs, assignment.Symbol)
	}
	return
}

func (a Assignments) IsIdentity() bool {
	for _, assignment := range a {
		if symbolRef, ok := assignment.Expression.(*tree.SymbolReference); ok && symbolRef.Name == assignment.Symbol.Name {
			continue
		} else {
			return false
		}
	}
	return true
}

func (a Assignments) Unique() Assignments {
	return lo.UniqBy(a, func(item *Assignment) string {
		return item.Symbol.Name
	})
}

type ProjectionNode struct {
	Source      PlanNode    `json:"source"`
	Assignments Assignments `json:"assignments"`

	BaseNode
}

func (n *ProjectionNode) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func (n *ProjectionNode) GetSources() []PlanNode {
	return []PlanNode{n.Source}
}

func (n *ProjectionNode) GetOutputSymbols() []*Symbol {
	return n.Assignments.GetOutputs()
}

func (n *ProjectionNode) ReplaceChildren(newChildren []PlanNode) PlanNode {
	return &ProjectionNode{
		BaseNode: BaseNode{
			ID: n.GetNodeID(),
		},
		Source:      newChildren[0],
		Assignments: n.Assignments,
	}
}
