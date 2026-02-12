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

package stream

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/lindb/common/pkg/logger"
	"github.com/samber/lo"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/streaming/cep/stream/input"
)

type sourceConnectorProvider struct{}

func NewSourceConnectorProvider() spi.SourceConnectorProvider {
	return &sourceConnectorProvider{}
}

func (s *sourceConnectorProvider) CreateSourceConnector(ctx context.Context,
	table spi.TableHandle, partitions []int, columnMapping map[string]string,
	predicate tree.Expression,
	outputColumns []types.ColumnMetadata, assignments []*spi.ColumnAssignment,
) spi.SourceConnector {
	tableHandle := table.(*TableHandle)
	schema, err := GetManager().GetStreamManager(tableHandle.Database).GetTableMetadata(tableHandle.Database, "", tableHandle.Stream)
	if err != nil {
		panic(err)
	}

	inputHandle := input.GetManager().GetInputHandler(tableHandle.Database, tableHandle.Stream)
	connector := &sourceConnector{
		ctx: ctx,

		table:   tableHandle,
		input:   inputHandle,
		running: atomic.NewBool(true),

		schema:        schema.Schema,
		predicate:     predicate,
		outputColumns: outputColumns,

		inbound: operator.NewQueue(make(chan arrow.RecordBatch, 256)),

		logger: logger.GetLogger("CEP", "SourceConnector"),
	}
	connector.initialize()

	inputHandle.Subscribe(connector)

	return connector
}

type sourceConnector struct {
	ctx context.Context

	table   *TableHandle
	input   input.InputHandler
	running *atomic.Bool

	schema        *types.TableSchema
	outputColumns []types.ColumnMetadata

	recordSchema *arrow.Schema

	inbound *operator.Queue

	predicate tree.Expression
	visitor   *visitor

	logger logger.Logger
}

func (sc *sourceConnector) initialize() {
	index := 0
	columnMap := lo.Associate(sc.schema.Columns, func(item types.ColumnMetadata) (string, int) {
		i := index
		index++
		return item.Name, i
	})
	for i := range sc.outputColumns {
		colMeta := &sc.outputColumns[i]
		colIndex, ok := columnMap[colMeta.Name]
		if !ok {
			colMeta.Ref = -1
			continue
		}
		colMeta.Ref = colIndex
	}

	sc.recordSchema = arrow.NewSchema(lo.Map(sc.outputColumns, func(item types.ColumnMetadata, index int) arrow.Field {
		return arrow.Field{Name: item.Name, Type: item.DataType.ToArrowDataType()}
	}), nil)
}

func (sc *sourceConnector) Receive(event models.Event) {
	if !sc.running.Load() {
		sc.logger.Warn("source connector has been stopped, drop received event",
			logger.String("database", sc.table.Database),
			logger.String("table", sc.table.Stream))
		return
	}
	if record, ok := event.(arrow.RecordBatch); ok {
		sc.inbound.Produce(record)
	}
}

func (sc *sourceConnector) Run(output chan<- arrow.RecordBatch) {
	defer func() {
		if err := recover(); err != nil {
			sc.logger.Error("source connector panicked", logger.Any("error", err), logger.Stack())
		}
		if sc.running.CompareAndSwap(true, false) {
			// unsubscribe input handler
			sc.input.Unsubscribe(sc)
			sc.inbound.Close()

			sc.logger.Info("source connector stopped",
				logger.String("database", sc.table.Database),
				logger.String("table", sc.table.Stream))
		}
	}()
	fmt.Println("run source connector")

	for {
		record, ok := sc.inbound.Consume(sc.ctx)
		// fmt.Printf("consume record: %v, ok: %v\n", record, ok)
		if !ok {
			break
		}
		sc.process(record, output)
	}
}

func (sc *sourceConnector) process(record arrow.RecordBatch, output chan<- arrow.RecordBatch) {
	// defer record.Release()

	if sc.predicate == nil {
		// no filter, send page to next operator
		columns := make([]arrow.Array, len(sc.outputColumns))
		for i, column := range sc.outputColumns {
			if column.Ref >= 0 {
				columns[i] = record.Column(column.Ref)
			}
		}
		rs := array.NewRecordBatch(sc.recordSchema, columns, record.NumRows())

		output <- rs
		return
	}
	// do filter based on predicate
	if sc.visitor == nil {
		// create predicate visitor if nil
		sc.visitor = &visitor{
			sc:          sc,
			schema:      record.Schema(),
			evalContext: expression.NewEvalContext(sc.ctx),
		}

		sc.visitor.expr = sc.visitor.rewrite(sc.predicate)
	}

	result, err := sc.visitor.expr.Eval(record)
	if err != nil {
		fmt.Println("failed to evaluate predicate:", err)
		return
	}
	if result.IsEmpty() {
		fmt.Println("filter result is empty, skip this record")
		return
	}
	// inputDatum := compute.NewDatum(record)
	// defer inputDatum.Release()
	// indices := result.ToArray()
	//
	// it32Builder := array.NewInt32Builder(memory.DefaultAllocator)
	// defer it32Builder.Release()
	//
	// for _, idx := range indices {
	// 	it32Builder.Append(int32(idx))
	// }
	// indexArray := it32Builder.NewArray()
	// defer indexArray.Release()
	// indexDatum := compute.NewDatum(indexArray)
	// defer indexDatum.Release()
	//
	// r, err := compute.Take(context.TODO(), *compute.DefaultTakeOptions(), inputDatum, indexDatum)
	// if err != nil {
	// 	fmt.Println("failed to take record based on filter result:", err)
	// 	return
	// }
	// defer r.Release()
	// _ = r.(*compute.RecordDatum).Value
	// fmt.Println(rr)
	// fmt.Printf("filter result: %v\n", result)

	// if page != nil {
	columns := make([]arrow.Array, len(sc.outputColumns))
	fields := make([]arrow.Field, len(sc.outputColumns))
	for i, column := range sc.outputColumns {
		if column.Ref >= 0 {
			columns[i] = record.Column(column.Ref)
		}
		fields[i] = arrow.Field{Name: column.Name, Type: columns[i].DataType()}
	}
	// fmt.Printf("output record: %v\n", columns)
	rs := array.NewRecordBatch(arrow.NewSchema(fields, nil), columns, record.NumRows())

	output <- rs
	// output <- record
	// }
}

type visitor struct {
	sc          *sourceConnector
	schema      *arrow.Schema
	evalContext expression.EvalContext

	expr Expr
}

func (v *visitor) rewrite(n tree.Expression) Expr {
	switch node := n.(type) {
	// case *tree.ComparisonExpression:
	// v, err := GetFieldValue(event, getValue(node.Left))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	case *tree.InPredicate:
		columName, err := expression.EvalString(v.evalContext, node.Value)
		if err != nil {
			panic(err)
		}
		indexes := v.schema.FieldIndices(columName)
		if len(indexes) != 1 {
			panic(fmt.Sprintf("invalid column %s, found %d fields", columName, len(indexes)))
		}
		var values []string
		if inListExpression, ok := node.ValueList.(*tree.InListExpression); ok {
			values = lo.Map(inListExpression.Values, func(item tree.Expression, index int) string {
				value, err := expression.EvalString(v.evalContext, item)
				if err != nil {
					panic(err)
				}
				return value
			})
		}

		return &InExpr{
			column: column{
				name:  columName,
				index: indexes[0],
			},
			values: values,
		}
	// case *tree.LogicalExpression:
	// for _, term := range node.Terms {
	// 	val, ok := term.Accept(event, v).(bool)
	// 	if !ok {
	// 		return false
	// 	}
	// 	if node.Operator == tree.LogicalOR && val {
	// 		return true
	// 	} else if !val {
	// 		return false
	// 	}
	// }
	// return true
	default:
		panic(fmt.Errorf("not support,%T", n))
	}
}

// func (v *visitor) filter(record arrow.RecordBatch) (*roaring.Bitmap, error) {
// 	return v.expr.Eval(record)
// 	switch node := n.(type) {
// 	case *tree.ComparisonExpression:
// 		v, err := GetFieldValue(event, getValue(node.Left))
// 		if err != nil {
// 			fmt.Println(err)
// 			return false
// 		}
// 		return v == getValue(node.Right)
// 	case *tree.InPredicate:
// 		var values []string
// 		if inListExpression, ok := node.ValueList.(*tree.InListExpression); ok {
// 			values = lo.Map(inListExpression.Values, func(item tree.Expression, index int) string {
// 				return getValue(item)
// 			})
// 		}
// 		v, err := GetFieldValue(event, getValue(node.Value))
// 		if err != nil {
// 			fmt.Println(err)
// 			return false
// 		}
// 		return lo.Contains(values, v.(string))
// 	case *tree.LogicalExpression:
// 		for _, term := range node.Terms {
// 			val, ok := term.Accept(event, v).(bool)
// 			if !ok {
// 				return false
// 			}
// 			if node.Operator == tree.LogicalOR && val {
// 				return true
// 			} else if !val {
// 				return false
// 			}
// 		}
// 		return true
// 	default:
// 		panic(fmt.Errorf("not support,%T", n))
// 	}
// fmt.Println("filtered page rows:", string(encoding.JSONMarshal(newPage)))
// 	return newPage
// }

// func (v *visitor) eval(record arrow.RecordBatch, row int) bool {
// 	switch node := v.expr.(type) {
// 	case *InExpr:
// 		return lo.Contains(node.values, row.GetString(node.column.Ref))
// 	default:
// 		panic(fmt.Errorf("not support,%T", v.expr))
// 	}
// }
