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

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"

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
		inbound: NewQueue(make(chan arrow.RecordBatch, 1024)),
	}
}

func (op *ProjectionOperator) Run(ctx context.Context, output chan<- arrow.RecordBatch) {
	if len(op.exprs) == 0 {
		op.prepare()
	}

	for {
		record, ok := op.inbound.Consume(ctx)
		if !ok {
			break
		}
		output <- op.process(record)
	}
}

func (op *ProjectionOperator) process(record arrow.RecordBatch) arrow.RecordBatch {
	result := make([]arrow.Array, len(op.exprs))
	success := false

	defer func() {
		// record.Release()
		if !success {
			for _, array := range result {
				if array != nil {
					array.Release()
				}
			}
		}
	}()

	for i, expr := range op.exprs {
		array, err := expr.Eval(record)
		if err != nil {
			panic(err)
		}
		result[i] = array
	}

	fields := make([]arrow.Field, len(op.project.Assignments))
	for i, assign := range op.project.Assignments {
		fields[i] = arrow.Field{
			Name: assign.Symbol.Name,
			Type: result[i].DataType(),
		}
	}

	// fmt.Println(record)
	// fmt.Println(result)
	rs := array.NewRecordBatch(arrow.NewSchema(fields, nil), result, record.NumRows())
	// fmt.Println(rs)
	success = true
	return rs
}

func (op *ProjectionOperator) GetLayout() []*plan.Symbol {
	return op.project.GetOutputSymbols()
}

func (op *ProjectionOperator) Children() []Operator {
	return []Operator{op.child}
}

func (op *ProjectionOperator) GetInbounds() []chan arrow.RecordBatch {
	return []chan arrow.RecordBatch{op.inbound.GetInbound()}
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
