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
	"fmt"

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
		inbound: NewQueue(make(chan *types.Page)),
	}
}

func (op *ProjectionOperator) Run(ctx context.Context, output chan<- *types.Page) {
	if len(op.exprs) == 0 {
		op.prepare()
	}
	fmt.Println(op.exprs)

	for {
		source, ok := op.inbound.Consume(ctx)
		if !ok {
			break
		}

		newPage := types.NewPage()
		outputColumns := make([]*types.Column, len(op.project.Assignments))
		for i, assign := range op.project.Assignments {
			outputColumns[i] = types.NewColumn()
			newPage.AppendColumn(types.NewColumnInfo(assign.Symbol.Name, assign.Symbol.DataType), outputColumns[i])
		}
		it := source.Iterator()
		for row := it.Begin(); row != it.End(); row = it.Next() {
			fmt.Println("do projection op....")
			for i, expr := range op.exprs {
				fmt.Printf("do ..... projection op expr %T,%s ret type=%v\n", expr, expr.String(), expr.GetType().String())
				switch expr.GetType() {
				case types.DTString:
					val, _, _ := expr.EvalString(op.exprCtx, row)
					outputColumns[i].AppendString(val)
				case types.DTInt:
					val, _, _ := expr.EvalInt(op.exprCtx, row)
					outputColumns[i].AppendInt(val)
				case types.DTFloat:
					val, _, _ := expr.EvalFloat(op.exprCtx, row)
					outputColumns[i].AppendFloat(val)
				case types.DTTimeSeries:
					val, _, _ := expr.EvalTimeSeries(op.exprCtx, row)
					outputColumns[i].AppendTimeSeries(val)
				case types.DTTimestamp:
					val, _, _ := expr.EvalTime(op.exprCtx, row)
					outputColumns[i].AppendTimestamp(val)
				case types.DTDuration:
					val, _, _ := expr.EvalDuration(op.exprCtx, row)
					outputColumns[i].AppendDuration(val)
				default:
					panic("projection operator error, unsupport data type:" + expr.GetType().String())
				}
			}
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

func (op *ProjectionOperator) prepare() {
	op.exprCtx = expression.NewEvalContext(op.ctx)
	op.exprs = make([]expression.Expression, len(op.project.Assignments))
	for i, assign := range op.project.Assignments {
		op.exprs[i] = expression.Rewrite(&expression.RewriteContext{
			SourceLayout: op.child.GetLayout(),
		}, assign.Expression)
	}
}
