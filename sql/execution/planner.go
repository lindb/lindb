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
	"github.com/lindb/lindb/sql/analyzer"
	sqlContext "github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/planner/validate"
	"github.com/lindb/lindb/sql/tree"
)

type Planner struct {
	analyzerFct *analyzer.AnalyzerFactory
}

func NewPlanner(analyzerFct *analyzer.AnalyzerFactory) *Planner {
	return &Planner{analyzerFct: analyzerFct}
}

func (p *Planner) Plan(session *Session,
	statement tree.Statement,
) *plan.Plan {
	analyzerContext := analyzer.NewAnalyzerContext(session.Database, statement, session.NodeIDAllocator)
	plannerContext := sqlContext.NewPlannerContext(
		session.Context,
		session.Database,
		session.NodeIDAllocator,
		statement,
	)
	plannerContext.AnalyzerContext = analyzerContext
	plannerContext.SymbolAllocator = plan.NewSymbolAllocator(analyzerContext)
	analyzer := p.analyzerFct.CreateAnalyzer(analyzerContext)
	// do analyze
	analyzer.Analyze(statement)

	// plan query
	logicalPlanner := planner.NewLogicalPlanner(plannerContext, planOptimizers())
	plan := logicalPlanner.Plan()
	v := validate.NewValidators()
	if err := v.Validate(plannerContext, plan.Root); err != nil {
		panic(err)
	}
	return plan
}

func (p *Planner) PlanDistribution(plan *plan.Plan) *plan.SubPlan {
	// fragment the plan
	fragmenter := planner.NewPlanFragmenter()
	return fragmenter.CreateSubPlans(plan)
}
