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

import "fmt"

type DefaultTraversalVisitor struct {
	Process func(node PlanNode)
	Resolve func(node PlanNode) PlanNode
}

func (v *DefaultTraversalVisitor) Visit(context any, n PlanNode) (r any) {
	curNode := n
	switch node := n.(type) {
	case *OutputNode:
		_ = node.Source.Accept(context, v)
	case *AggregationNode:
		_ = node.Source.Accept(context, v)
	case *FilterNode:
		_ = node.Source.Accept(context, v)
	case *ProjectionNode:
		_ = node.Source.Accept(context, v)
	case *GroupReference:
		// need resolve group reference(raw plan node)
		if v.Resolve != nil {
			rawNode := v.Resolve(node)
			// set current node using raw node
			curNode = rawNode
			sources := rawNode.GetSources()
			for _, source := range sources {
				_ = source.Accept(context, v)
			}
		}
	default:
		// TODO: remove
		fmt.Printf("plan node default traversal visitor not support..................=%T\n", n)
	}
	if v.Process != nil {
		v.Process(curNode)
	}
	return
}
