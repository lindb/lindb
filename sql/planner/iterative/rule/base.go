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

// Base represents optimization base rule for common logic.
type Base[N plan.PlanNode] struct {
	// apply apply rule for specific node.
	apply func(context *iterative.Context, node N) plan.PlanNode
}

// Apply applies optimization rule for specific node.
func (rule *Base[N]) Apply(context *iterative.Context, node plan.PlanNode) plan.PlanNode {
	// check node if match target node type
	if targetNode, ok := node.(N); ok {
		return rule.apply(context, targetNode)
	}
	return nil
}
