package validate

import (
	"errors"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/context"
	"github.com/lindb/lindb/sql/planner/plan"
)

type OutputValidator struct {
	Base[*plan.OutputNode]
}

func NewOutputValidator() Validator {
	v := &OutputValidator{}
	v.validate = func(ctx *context.PlannerContext, node *plan.OutputNode) error {
		_, ok := lo.Find(node.GetOutputSymbols(), func(item *plan.Symbol) bool {
			return item.DataType == types.DTTimestamp
		})
		if !ok {
			// output node has no timestamp column
			return nil
		}
		if _, ok = lo.Find(node.GetOutputSymbols(), func(item *plan.Symbol) bool {
			return item.DataType == types.DTTimeSeries
		}); !ok {
			return errors.New("push timestamp column failed, output node has no time series column")
		}
		return nil
	}
	return v
}
