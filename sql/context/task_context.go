package context

import (
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/planner/plan"
)

type TaskContext struct {
	TaskID     model.TaskID
	Fragment   *plan.PlanFragment
	Partitions []int
}
