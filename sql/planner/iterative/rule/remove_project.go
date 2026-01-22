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

package rule

import (
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
		if node.Assignments.IsIdentity() &&
			symbolsEquals(node.GetOutputSymbols(), node.Source.GetOutputSymbols()) {
			return node.Source
		}
		return nil
	}
	return rule
}

func symbolsEquals(a, b []*plan.Symbol) bool {
	if len(a) != len(b) {
		return false
	}
	names := make(map[string]struct{})
	for _, aName := range a {
		names[aName.Name] = struct{}{}
	}

	for _, bName := range b {
		if _, ok := names[bName.Name]; !ok {
			return false
		}
	}

	return true
}
