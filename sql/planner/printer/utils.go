package printer

import (
	"fmt"
	"strings"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
)

func formatSymbols(symbols []*plan.Symbol) string {
	var columns []string
	for i := range symbols {
		columns = append(columns, symbols[i].String())
	}
	return "[" + strings.Join(columns, ", ") + "]"
}

func formatAggregation(aggregation *plan.Aggregation) string {
	var args []string
	for _, arg := range aggregation.Arguments {
		args = append(args, fmt.Sprintf("\"%s\"", tree.FormatExpression(arg)))
	}
	return fmt.Sprintf("%s(%s)", aggregation.ResolvedFunction.Signature.Name, strings.Join(args, ", "))
}
