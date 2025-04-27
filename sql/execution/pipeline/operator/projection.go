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

// Projection means choosing which columns (or expressions) the query shall return.
type ProjectionOperator struct {
	ctx     context.Context
	exprCtx expression.EvalContext
	project *plan.ProjectionNode

	exprs []expression.Expression

	child Operator
}

func NewProjectionOperator(ctx context.Context, project *plan.ProjectionNode, child Operator) Operator {
	return &ProjectionOperator{
		ctx:     ctx,
		project: project,
		child:   child,
	}
}

// GetOutput implements Operator.
func (h *ProjectionOperator) Run(ctx context.Context, output chan<- *types.Page) {
	if len(h.exprs) == 0 {
		h.prepare()
	}
	fmt.Println(h.exprs)

	inbound := make(chan *types.Page)

	go func() {
		defer close(inbound)
		h.child.Run(ctx, inbound)
	}()

	for source := range inbound {
		newPage := types.NewPage()
		outputColumns := make([]*types.Column, len(h.project.Assignments))
		for i, assign := range h.project.Assignments {
			outputColumns[i] = types.NewColumn()
			newPage.AppendColumn(types.NewColumnInfo(assign.Symbol.Name, assign.Symbol.DataType), outputColumns[i])
		}
		it := source.Iterator()
		for row := it.Begin(); row != it.End(); row = it.Next() {
			fmt.Println("do projection op....")
			for i, expr := range h.exprs {
				fmt.Printf("do ..... projection op expr %T,%s ret type=%v\n", expr, expr.String(), expr.GetType().String())
				switch expr.GetType() {
				case types.DTString:
					val, _, _ := expr.EvalString(h.exprCtx, row)
					outputColumns[i].AppendString(val)
				case types.DTInt:
					val, _, _ := expr.EvalInt(h.exprCtx, row)
					outputColumns[i].AppendInt(val)
				case types.DTFloat:
					val, _, _ := expr.EvalFloat(h.exprCtx, row)
					outputColumns[i].AppendFloat(val)
				case types.DTTimeSeries:
					val, _, _ := expr.EvalTimeSeries(h.exprCtx, row)
					outputColumns[i].AppendTimeSeries(val)
				case types.DTTimestamp:
					val, _, _ := expr.EvalTime(h.exprCtx, row)
					outputColumns[i].AppendTimestamp(val)
				case types.DTDuration:
					val, _, _ := expr.EvalDuration(h.exprCtx, row)
					outputColumns[i].AppendDuration(val)
				default:
					panic("projection operator error, unsupport data type:" + expr.GetType().String())
				}
			}
		}

		output <- newPage
	}
}

func (h *ProjectionOperator) GetLayout() []*plan.Symbol {
	return h.project.GetOutputSymbols()
}

func (h *ProjectionOperator) Children() []Operator {
	return []Operator{h.child}
}

func (h *ProjectionOperator) prepare() {
	h.exprCtx = expression.NewEvalContext(h.ctx)
	h.exprs = make([]expression.Expression, len(h.project.Assignments))
	for i, assign := range h.project.Assignments {
		h.exprs[i] = expression.Rewrite(&expression.RewriteContext{
			SourceLayout: h.child.GetLayout(),
		}, assign.Expression)
	}
}
