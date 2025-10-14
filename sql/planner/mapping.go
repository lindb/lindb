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
	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/planner/plan"
	planpkg "github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

type NodeAndMappings struct {
	Node   plan.PlanNode
	Fields []*plan.Symbol
}

type PlanAndMappings struct {
	mappings map[tree.Expression]*plan.Symbol
}

func coerceExpressions(subPlan *PlanBuilder, expressions []tree.Expression, analysis *analyzer.Analysis,
	symbolAllocator *planpkg.SymbolAllocator, _ *planpkg.PlanNodeIDAllocator,
) *PlanAndMappings {
	mappings := make(map[tree.Expression]*planpkg.Symbol)

	for i := range expressions {
		expression := expressions[i]
		if _, ok := mappings[expression]; !ok {
			// TODO: need modify
			t := analysis.GetType(expression)
			symbol := symbolAllocator.FromExpression(subPlan.translations.Rewrite(expression), t)
			mappings[expression] = symbol
		}
	}
	return &PlanAndMappings{
		mappings: mappings,
	}
}
