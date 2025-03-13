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

package execution

import (
	"github.com/lindb/lindb/sql/planner/printer"
	"github.com/lindb/lindb/sql/tree"
)

type QueryExplainer struct {
	planner *Planner
}

func NewQueryExplainer(planner *Planner) *QueryExplainer {
	return &QueryExplainer{planner: planner}
}

func (qe *QueryExplainer) ExplainPlan(session *Session,
	statement tree.Statement, explainType string,
) string {
	plan := qe.planner.Plan(session, statement)
	printer := printer.NewPlanPrinter(printer.NewTextRender(0))
	if explainType == tree.DistributedExplain {
		fragmentedPlan := qe.planner.PlanDistribution(plan)
		return printer.PrintDistributedPlan(fragmentedPlan)
	}
	return printer.PrintLogicPlan(plan.Root)
}
