package execution

import "github.com/lindb/lindb/sql/planner/plan"

type PlanRoot struct {
	root *plan.SubPlan
}
