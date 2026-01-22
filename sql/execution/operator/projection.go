// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package operator

import (
	"context"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
)

// ProjectionOperator means choosing which columns (or expressions) the query shall return.
type ProjectionOperator struct {
	ctx     context.Context
	exprCtx expression.EvalContext
	project *plan.ProjectionNode

	exprs []expression.Expression

	child   Operator
	inbound *Queue
}

func NewProjectionOperator(ctx context.Context, project *plan.ProjectionNode, child Operator) Operator {
	return &ProjectionOperator{
		ctx:     ctx,
		project: project,
		child:   child,
		inbound: NewQueue(make(chan *types.Page, 1024)),
	}
}

func (op *ProjectionOperator) Run(ctx context.Context, output chan<- *types.Page) {
	if len(op.exprs) == 0 {
		op.prepare()
	}

	for {
		source, ok := op.inbound.Consume(ctx)
		if !ok {
			break
		}
		newPage := types.NewPage()
		outputColumns := make([]*types.Column, len(op.project.Assignments))
		for i, assign := range op.project.Assignments {
			outputColumns[i] = types.NewColumn()
			newPage.AppendColumn(
				types.NewColumnInfo(assign.Symbol.Name, assign.Symbol.DataType, assign.Symbol.Hidden, assign.Symbol.AggType),
				outputColumns[i])
		}
		rowNum := 0
		// logger.GetLogger("sql", "projection").Info("do projection op start", logger.Any("page", source), logger.Any("exp", op.project.Assignments))

		it := source.Iterator()
		for row := it.Begin(); row != it.End(); row = it.Next() {
			for i, expr := range op.exprs {
				switch expr.GetType() {
				case types.DTString:
					val, _, _ := expr.EvalString(row)
					outputColumns[i].Append(val)
				case types.DTInt:
					val, _, _ := expr.EvalInt(row)
					outputColumns[i].Append(val)
				case types.DTFloat:
					val, _, _ := expr.EvalFloat(row)
					outputColumns[i].Append(val)
				case types.DTTimeSeries:
					val, _, _ := expr.EvalTimeSeries(row)
					outputColumns[i].Append(val)
				case types.DTTimestamp:
					val, _, _ := expr.EvalTime(row)
					outputColumns[i].Append(val)
				case types.DTDuration:
					val, _, _ := expr.EvalDuration(row)
					outputColumns[i].Append(val)
				case types.DTMap:
					val, _, _ := expr.EvalMap(row)
					outputColumns[i].Append(val)
				case types.DTExemplar:
					val, _, _ := expr.EvalExemplar(row)
					outputColumns[i].Append(val)
				default:
					panic("projection operator error, unsupport data type:" + expr.GetType().String())
				}
			}

			rowNum++
		}

		output <- newPage
	}
}

func (op *ProjectionOperator) GetLayout() []*plan.Symbol {
	return op.project.GetOutputSymbols()
}

func (op *ProjectionOperator) Children() []Operator {
	return []Operator{op.child}
}

func (op *ProjectionOperator) GetInbounds() []chan *types.Page {
	return []chan *types.Page{op.inbound.GetInbound()}
}

func (op *ProjectionOperator) String() string {
	return "ProjectionOperator"
}

func (op *ProjectionOperator) prepare() {
	op.exprCtx = expression.NewEvalContext(op.ctx)
	op.exprs = make([]expression.Expression, len(op.project.Assignments))
	for i, assign := range op.project.Assignments {
		op.exprs[i] = expression.Rewrite(&expression.RewriteContext{
			SourceLayout: op.child.GetLayout(),
			EvalContext:  op.exprCtx,
		}, assign.Expression)
	}
}
