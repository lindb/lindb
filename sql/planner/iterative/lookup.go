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
