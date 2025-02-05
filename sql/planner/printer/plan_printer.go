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

package printer

import (
	"fmt"
	"strings"

	"github.com/lindb/lindb/sql/planner/plan"
)

type PlanPrinter struct {
	render Render
}

func NewPlanPrinter(render Render) *PlanPrinter {
	return &PlanPrinter{
		render: render,
	}
}

func (p *PlanPrinter) PrintLogicPlan(plan plan.PlanNode) string {
	representation := NewPlanRepresentation(plan)

	visitor := NewVisitor(representation)
	plan.Accept(nil, visitor)

	return p.render.Render(representation)
}

func (p *PlanPrinter) PrintDistributedPlan(plan *plan.SubPlan) string {
	allFragments := plan.GetAllFragments()
	sb := &strings.Builder{}
	for _, fragment := range allFragments {
		sb.WriteString(p.formatFragment(fragment))
	}
	return sb.String()
}

func (p *PlanPrinter) formatFragment(fragment *plan.PlanFragment) string {
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "Plan Fragemnt %v\n", fragment.ID)
	printer := NewPlanPrinter(NewTextRender(1))
	sb.WriteString(printer.PrintLogicPlan(fragment.Root))
	sb.WriteString("\n")
	return sb.String()
}

func TextLogicalPlan(plan plan.PlanNode) string {
	printer := NewPlanPrinter(NewTextRender(0))
	return printer.PrintLogicPlan(plan)
}
