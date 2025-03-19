package validate

import (
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
)

type Validators struct {
	validators []Validator
}

func NewValidators() Validator {
	return &Validators{
		validators: []Validator{
			NewOutputValidator(),
		},
	}
}

func (v *Validators) Validate(ctx *context.PlannerContext, node plan.PlanNode) error {
	if err := v.validate(ctx, node); err != nil {
		return err
	}

	for _, child := range node.GetSources() {
		if err := v.Validate(ctx, child); err != nil {
			return err
		}
	}
	return nil
}

func (v *Validators) validate(ctx *context.PlannerContext, node plan.PlanNode) error {
	for _, validator := range v.validators {
		if err := validator.Validate(ctx, node); err != nil {
			return err
		}
	}
	return nil
}
