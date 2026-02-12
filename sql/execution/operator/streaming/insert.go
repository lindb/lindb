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

package streaming

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/streaming/cep/stream/input"
)

type InsertOperator struct {
	ctx     context.Context
	exprCtx expression.EvalContext
	insert  *plan.InsertNode

	child   operator.Operator
	inbound *operator.Queue
}

func NewInsertOperator(ctx context.Context, insert *plan.InsertNode, child operator.Operator) operator.Operator {
	return &InsertOperator{
		ctx:     ctx,
		insert:  insert,
		child:   child,
		inbound: operator.NewQueue(make(chan arrow.RecordBatch, 256)),
	}
}

func (op *InsertOperator) Run(ctx context.Context, output chan<- arrow.RecordBatch) {
	streamName := op.insert.Table.Name.Name
	// FIXME: get app from context
	inputHandle := input.GetManager().GetInputHandler(op.insert.Database, streamName)
	for {
		page, ok := op.inbound.Consume(ctx)
		if !ok {
			break
		}
		inputHandle.Send(page)
	}
}

func (op *InsertOperator) GetLayout() []*plan.Symbol {
	return op.insert.GetOutputSymbols()
}

func (op *InsertOperator) Children() []operator.Operator {
	return []operator.Operator{op.child}
}

func (op *InsertOperator) GetInbounds() []chan arrow.RecordBatch {
	return []chan arrow.RecordBatch{op.inbound.GetInbound()}
}

func (op *InsertOperator) String() string {
	return "InsertOperator"
}
