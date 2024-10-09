package execution

import (
	"github.com/lindb/lindb/sql/analyzer"
	sqlContext "github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner"
	"github.com/lindb/lindb/sql/planner/plan"
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
	return logicalPlanner.Plan()
}

func (p *Planner) PlanDistribution(plan *plan.Plan) *plan.SubPlan {
	// fragment the plan
	fragmenter := planner.NewPlanFragmenter()
	return fragmenter.CreateSubPlans(plan)
}
