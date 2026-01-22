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
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/streaming/cep/stream/input"
)

type OutputOperator struct {
	ctx         context.Context
	inputHandle input.InputHandler
	exprCtx     expression.EvalContext
	output      *plan.OutputNode

	child   operator.Operator
	inbound *operator.Queue
}

func NewOutputOperator(ctx context.Context, database, streamName string,
	output *plan.OutputNode, child operator.Operator,
) operator.Operator {
	inputHandle := input.GetManager().GetInputHandler(database, streamName)
	return &OutputOperator{
		ctx:         ctx,
		inputHandle: inputHandle,
		output:      output,
		child:       child,
		inbound:     operator.NewQueue(make(chan *types.Page, 256)),
	}
}

func (op *OutputOperator) Run(ctx context.Context, output chan<- *types.Page) {
	// FIXME: get app from context
	rebuildPage := false
	layout := op.output.GetOutputSymbols()

	sourceLayout := make(map[string]int)
	lo.ForEach(op.child.GetLayout(), func(symbol *plan.Symbol, index int) {
		sourceLayout[symbol.Name] = index
	})
	columnNames := op.output.ColumnNames

	for idx, symbol := range layout {
		sourceIdx, ok := sourceLayout[symbol.Name]
		if ok && (idx != sourceIdx || (len(columnNames) > 0 && columnNames[idx] != symbol.Name)) {
			rebuildPage = true
			break
		}
	}

	// process child output
	for {
		page, ok := op.inbound.Consume(ctx)
		if !ok {
			break
		}
		if page == nil || page.NumRows() == 0 {
			if page != nil && page.Error != "" {
				panic(fmt.Errorf("output operator receive error page: %s", page.Error))
			}
			break
		}
		if rebuildPage {
			targetPage := types.NewPage()
			for colIdx, col := range layout {
				if idx, ok := sourceLayout[col.Name]; ok {
					column := page.Layout[idx]
					if len(columnNames) > 0 {
						column.Name = columnNames[colIdx]
					}
					targetPage.AppendColumn(column, page.Columns[idx])
				}
			}
			op.inputHandle.Send(targetPage)
		} else {
			op.inputHandle.Send(page)
		}
	}
}

func (op *OutputOperator) GetLayout() []*plan.Symbol {
	panic("output operator should not get layout")
}

func (op *OutputOperator) Children() []operator.Operator {
	return []operator.Operator{op.child}
}

func (op *OutputOperator) GetInbounds() []chan *types.Page {
	return []chan *types.Page{op.inbound.GetInbound()}
}

func (op *OutputOperator) String() string {
	return "OutputOperator"
}
