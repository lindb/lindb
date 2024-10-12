package operator

import (
	"context"
	"fmt"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
)

type ProjectionOperatorFactory struct {
	project      *plan.ProjectionNode
	sourceLayout []*plan.Symbol
}

func NewProjectionOperatorFactory(project *plan.ProjectionNode, sourceLayout []*plan.Symbol) OperatorFactory {
	return &ProjectionOperatorFactory{
		project:      project,
		sourceLayout: sourceLayout,
	}
}

// CreateOperator implements OperatorFactory.
func (fct *ProjectionOperatorFactory) CreateOperator(ctx context.Context) Operator {
	return NewProjectionOperator(ctx, fct.project, fct.sourceLayout)
}

type ProjectionOperator struct {
	ctx     context.Context
	exprCtx expression.EvalContext
	project *plan.ProjectionNode
	source  *types.Page // TODO: refact
	ouput   *types.Page

	sourceLayout  []*plan.Symbol
	outputColumns []*types.Column
	exprs         []expression.Expression
}

func NewProjectionOperator(ctx context.Context, project *plan.ProjectionNode, sourceLayout []*plan.Symbol) Operator {
	return &ProjectionOperator{ctx: ctx, project: project, sourceLayout: sourceLayout}
}

// AddInput implements Operator.
func (h *ProjectionOperator) AddInput(page *types.Page) {
	h.source = page
}

// Finish implements Operator.
func (h *ProjectionOperator) Finish() {
}

// GetOutput implements Operator.
func (h *ProjectionOperator) GetOutput() *types.Page {
	if len(h.exprs) == 0 {
		h.prepare()
	}
	fmt.Println(h.exprs)
	fmt.Println(h.source)
	it := h.source.Iterator()
	for row := it.Begin(); row != it.End(); row = it.Next() {
		fmt.Println("do projection op....")
		for i, expr := range h.exprs {
			fmt.Printf("do ..... projection op expr %T,%s ret type=%v\n", expr, expr.String(), expr.GetType().String())
			switch expr.GetType() {
			case types.DTString:
				val, _, _ := expr.EvalString(h.exprCtx, row)
				h.outputColumns[i].AppendString(val)
			case types.DTInt:
				val, _, _ := expr.EvalInt(h.exprCtx, row)
				h.outputColumns[i].AppendInt(val)
			case types.DTFloat:
				val, _, _ := expr.EvalFloat(h.exprCtx, row)
				h.outputColumns[i].AppendFloat(val)
			case types.DTTimeSeries:
				val, _, _ := expr.EvalTimeSeries(h.exprCtx, row)
				h.outputColumns[i].AppendTimeSeries(val)
			case types.DTTimestamp:
				val, _, _ := expr.EvalTime(h.exprCtx, row)
				h.outputColumns[i].AppendTimestamp(val)
			case types.DTDuration:
				val, _, _ := expr.EvalDuration(h.exprCtx, row)
				h.outputColumns[i].AppendDuration(val)
			default:
				panic("projection operator error, unsupport data type:" + expr.GetType().String())
			}
		}
	}
	return h.ouput
}

// IsFinished implements Operator.
func (h *ProjectionOperator) IsFinished() bool {
	return true
}

func (h *ProjectionOperator) prepare() {
	h.exprCtx = expression.NewEvalContext(h.ctx)
	h.exprs = make([]expression.Expression, len(h.project.Assignments))
	h.outputColumns = make([]*types.Column, len(h.project.Assignments))
	h.ouput = types.NewPage()
	for i, assign := range h.project.Assignments {
		h.exprs[i] = expression.Rewrite(&expression.RewriteContext{
			SourceLayout: h.sourceLayout,
		}, assign.Expression)
		h.outputColumns[i] = types.NewColumn()
		h.ouput.AppendColumn(types.NewColumnInfo(assign.Symbol.Name, assign.Symbol.DataType), h.outputColumns[i])
	}
}
