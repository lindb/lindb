// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package plan

import (
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/tree"
)

type Assignment struct {
	Symbol     *Symbol         `json:"symbol"`     // output
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
