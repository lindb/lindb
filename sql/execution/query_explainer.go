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
