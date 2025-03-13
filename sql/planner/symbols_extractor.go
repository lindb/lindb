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

package planner

import (
	"fmt"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

func ExtractSymbolsFromExpressions(expressions []tree.Expression) (symbols []*plan.Symbol) {
	visitor := &tree.DefaultTraversalVisitor{
		PostProcess: func(n tree.Node) {
			if ref, ok := n.(*tree.SymbolReference); ok {
				symbols = append(symbols, plan.SymbolFrom(ref))
			}
		},
	}
	for _, node := range expressions {
		visitor.Visit(nil, node)
	}
	return
}

func ExtractSymbolsFromAggreation(aggregation *plan.Aggregation) (symbols []*plan.Symbol) {
	visitor := &tree.DefaultTraversalVisitor{
		PostProcess: func(n tree.Node) {
			if ref, ok := n.(*tree.SymbolReference); ok {
				symbols = append(symbols, plan.SymbolFrom(ref))
			}
		},
	}
	for _, node := range aggregation.Arguments {
		fmt.Printf("extract symbols agg args.......=%T\n", node)
		visitor.Visit(nil, node)
	}
	// TODO: add agg other fields
	return
}

func ExtractSymbolsFromExpression(expression tree.Expression) []*plan.Symbol {
	panic("unimplemented implements extract symbols from expression")
}
