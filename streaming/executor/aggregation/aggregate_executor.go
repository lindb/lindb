package aggregation

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/streaming/executor"
	"github.com/lindb/lindb/streaming/query/context"
)

type aggregateExecutor struct {
	queryContext *context.QueryContext
	node         *plan.AggregationNode

	exprCtx expression.EvalContext
	exprs   []expression.Expression

	aggregators []Aggregator

	groupingKeys []tree.Expression

	source plan.PlanNode
}

func NewExecutor(ctx *context.QueryContext, node *plan.AggregationNode) executor.Executor {
	exec := &aggregateExecutor{
		queryContext: ctx,
		node:         node,
		groupingKeys: make([]tree.Expression, len(node.GroupingSets.GroupingKeys)),
	}
	exec.source = node.Source

	for i, key := range node.GroupingSets.GroupingKeys {
		exec.groupingKeys[i] = key.ToSymbolReference()
		if projection, ok := node.Source.(*plan.ProjectionNode); ok {
			exec.source = projection.Source

			assignment, exsit := lo.Find(projection.Assignments, func(item *plan.Assignment) bool {
				return item.Symbol.Name == key.Name
			})
			if exsit {
				exec.groupingKeys[i] = assignment.Expression
			}
		}
	}

	for _, agg := range node.Aggregations {
		aggregator, err := NewAggregator(agg.Aggregation.Function, agg.Aggregation.Arguments)
		if err != nil {
			panic(err)
		}
		exec.aggregators = append(exec.aggregators, aggregator)
	}

	exec.prepare()

	return exec
}

func (exec *aggregateExecutor) Process(event any) {
	for _, agg := range exec.aggregators {
		agg.Process(event)
	}

	values := make([]any, len(exec.exprs))
	v := types.String("rpc")
	row := types.NewArrayRow([]any{&v, map[string]string{"host": "1.1.1.1", "app": "order"}})

	for i, expr := range exec.exprs {
		fmt.Printf("do ..... projection op expr %T,%s ret type=%v\n", expr, expr.String(), expr.GetType().String())
		switch expr.GetType() {
		case types.DTString:
			val, _, _ := expr.EvalString(exec.exprCtx, row)
			values[i] = val
		case types.DTInt:
			val, _, _ := expr.EvalInt(exec.exprCtx, row)
			values[i] = val
		case types.DTFloat:
			val, _, _ := expr.EvalFloat(exec.exprCtx, row)
			values[i] = val
		case types.DTTimeSeries:
			val, _, _ := expr.EvalTimeSeries(exec.exprCtx, row)
			values[i] = val
		case types.DTTimestamp:
			val, _, _ := expr.EvalTime(exec.exprCtx, row)
			values[i] = val
		case types.DTDuration:
			val, _, _ := expr.EvalDuration(exec.exprCtx, row)
			values[i] = val
		case types.DTMap:
			val, _, _ := expr.EvalMap(exec.exprCtx, row)
			values[i] = val
		default:
			panic("projection operator error, unsupport data type:" + expr.GetType().String())
		}
	}
	fmt.Printf("result values=%v\n", values)
}

func (exec *aggregateExecutor) ResultSet(fn func(result any)) {
	values := make([]any, len(exec.aggregators))
	for i := range exec.aggregators {
		values[i] = exec.aggregators[i].GetValue()
	}

	fn(values)
}

func (exec *aggregateExecutor) prepare() {
	exec.exprCtx = expression.NewEvalContext(exec.queryContext.Context)
	exec.exprs = make([]expression.Expression, len(exec.groupingKeys))
	for i, assign := range exec.groupingKeys {
		exec.exprs[i] = expression.Rewrite(&expression.RewriteContext{
			SourceLayout: exec.source.GetOutputSymbols(),
		}, assign)
	}
}
