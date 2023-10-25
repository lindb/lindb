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

package operator

import (
	"github.com/lindb/lindb/query/context"
)

// physicalPlan represents physical plan operator.
type physicalPlan struct {
	ctx context.TaskContext
}

// NewPhysicalPlan creates a physicalPlan instance.
func NewPhysicalPlan(ctx context.TaskContext) Operator {
	return &physicalPlan{ctx: ctx}
}

// Execute returns physical plan by given task context.
func (op *physicalPlan) Execute() error {
	return op.ctx.MakePlan()
}

// Identifier returns identifier string value of physical plan operator.
func (op *physicalPlan) Identifier() string {
	return "Physical Plan"
}
