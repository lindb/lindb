package context

import (
	"context"

	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/planner/plan"
)

type TaskContext struct {
	Context    context.Context
	TaskID     model.TaskID
	Fragment   *plan.PlanFragment
	Partitions []int
}
