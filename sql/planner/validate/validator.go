package validate

import (
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
)

type Validator interface {
	Validate(ctx *context.PlannerContext, node plan.PlanNode) error
}

type Base[N plan.PlanNode] struct {
	validate func(ctx *context.PlannerContext, node N) error
}

func (v *Base[N]) Validate(ctx *context.PlannerContext, node plan.PlanNode) error {
	if targetNode, ok := node.(N); ok {
		return v.validate(ctx, targetNode)
	}
	return nil
}
