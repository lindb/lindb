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

package output

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/planner/plan"
)

type ResultSetOutputOperator struct {
	node    *plan.OutputNode
	child   operator.Operator
	inbound *operator.Queue
}

func NewRSOutputOperator(node *plan.OutputNode, child operator.Operator) operator.Operator {
	return &ResultSetOutputOperator{
		node:    node,
		child:   child,
		inbound: operator.NewQueue(make(chan arrow.RecordBatch)),
	}
}

func (op *ResultSetOutputOperator) Run(ctx context.Context, output chan<- arrow.RecordBatch) {
	rebuildRecord := false
	layout := op.node.GetOutputSymbols()

	sourceLayout := make(map[string]int)
	lo.ForEach(op.child.GetLayout(), func(symbol *plan.Symbol, index int) {
		sourceLayout[symbol.Name] = index
	})
	columnNames := op.node.ColumnNames

	fields := lo.Map(layout, func(symbol *plan.Symbol, _ int) arrow.Field {
		return arrow.Field{Name: symbol.Name, Type: symbol.DataType.ToArrowDataType()}
	})
	schema := arrow.NewSchema(fields, nil)

	for idx, symbol := range layout {
		sourceIdx, ok := sourceLayout[symbol.Name]
		if ok && (idx != sourceIdx || (len(columnNames) > 0 && columnNames[idx] != symbol.Name)) {
			rebuildRecord = true
			break
		}
	}

	// process child output
	for {
		record, ok := op.inbound.Consume(ctx)
		if !ok {
			break
		}
		if record == nil || record.NumRows() == 0 {
			// FIXME: if page != nil && page.Error != "" {
			// 	output <- page
			// }
			// TODO: if has error or no datareturn directly?
			continue
		}
		if rebuildRecord {
			defer record.Release()

			columns := make([]arrow.Array, len(layout))
			for colIdx, col := range layout {
				if idx, ok := sourceLayout[col.Name]; ok {
					columns[colIdx] = record.Column(idx)
				}
			}
			output <- array.NewRecordBatch(schema, columns, record.NumRows())
		} else {
			output <- record
		}
	}
}

func (op *ResultSetOutputOperator) GetLayout() []*plan.Symbol {
	panic("result set output operator should not get layout")
}

func (op *ResultSetOutputOperator) Children() []operator.Operator {
	return []operator.Operator{op.child}
}

func (op *ResultSetOutputOperator) GetInbounds() []chan arrow.RecordBatch {
	return []chan arrow.RecordBatch{op.inbound.GetInbound()}
}

func (op *ResultSetOutputOperator) String() string {
	return "ResultSetOutputOperator"
}
