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

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/streaming/cep/stream/input"
)

type OutputOperator struct {
	ctx          context.Context
	database     string
	outputStream string

	inputHandle input.InputHandler
	output      *plan.OutputNode

	child   operator.Operator
	inbound *operator.Queue
}

func NewOutputOperator(ctx context.Context, database, outputStream string,
	output *plan.OutputNode, child operator.Operator,
) operator.Operator {
	inputHandle := input.GetManager().GetInputHandler(database, outputStream)
	return &OutputOperator{
		ctx:          ctx,
		database:     database,
		outputStream: outputStream,
		inputHandle:  inputHandle,
		output:       output,
		child:        child,
		inbound:      operator.NewQueue(make(chan arrow.RecordBatch, 256)),
	}
}

func (op *OutputOperator) Run(ctx context.Context, output chan<- arrow.RecordBatch) {
	defer func() {
		fmt.Println("ouput....")
		input.GetManager().RemoveInputHandler(op.database, op.outputStream)
	}()

	// FIXME: get app from context
	rebuildRecord := false
	layout := op.output.GetOutputSymbols()

	sourceLayout := make(map[string]int)
	lo.ForEach(op.child.GetLayout(), func(symbol *plan.Symbol, index int) {
		sourceLayout[symbol.Name] = index
	})
	columnNames := op.output.ColumnNames

	fields := lo.Map(layout, func(symbol *plan.Symbol, index int) arrow.Field {
		name := symbol.Name
		if len(columnNames) > 0 {
			name = columnNames[index]
		}
		var metadata arrow.Metadata
		if symbol.AggType != types.ATUnknown {
			metadata = arrow.NewMetadata([]string{"agg"}, []string{symbol.AggType.String()})
		}
		return arrow.Field{Name: name, Type: symbol.DataType.ToArrowDataType(), Metadata: metadata}
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
			// 	panic(fmt.Errorf("output operator receive error page: %s", page.Error))
			// }
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
			op.inputHandle.Send(array.NewRecordBatch(schema, columns, record.NumRows()))
		} else {
			op.inputHandle.Send(record)
		}
	}
}

func (op *OutputOperator) GetLayout() []*plan.Symbol {
	panic("output operator should not get layout")
}

func (op *OutputOperator) Children() []operator.Operator {
	return []operator.Operator{op.child}
}

func (op *OutputOperator) GetInbounds() []chan arrow.RecordBatch {
	return []chan arrow.RecordBatch{op.inbound.GetInbound()}
}

func (op *OutputOperator) String() string {
	return "OutputOperator"
}
