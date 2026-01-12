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

	"github.com/lindb/common/pkg/encoding"
	"github.com/samber/lo"

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
	fmt.Printf("create log source connector,table=%v,partitions=%v,predicate=%v\n", outputColumns, assignments, predicate)

	tableHandle := table.(*TableHandle)
	schema, err := GetManager().GetStreamManager(tableHandle.App).GetTableMetadata(tableHandle.App, "", tableHandle.Stream)
	if err != nil {
		panic(err)
	}

	connector := &sourceConnector{
		ctx: ctx,

		schema:        schema.Schema,
		predicate:     predicate,
		outputColumns: outputColumns,

		inbound: operator.NewQueue(make(chan *types.Page, 256)),
	}
	connector.initialize()

	inputHandle := input.GetManager().GetInputHandler(tableHandle.App, tableHandle.Stream)
	inputHandle.Subscribe(connector)

	fmt.Printf("create streaming connector .....%+v====>%+v\n", schema.Schema, connector.outputColumns)
	return connector
}

type sourceConnector struct {
	ctx context.Context

	schema        *types.TableSchema
	predicate     tree.Expression
	outputColumns []types.ColumnMetadata

	inbound *operator.Queue // TODO: queue need close
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
}

func (sc *sourceConnector) Receive(event any) {
	if page, ok := event.(*types.Page); ok {
		sc.inbound.Produce(page)
	}
}

func (sc *sourceConnector) Run(output chan<- *types.Page) {
	var v *visitor
	for {
		source, ok := sc.inbound.Consume(sc.ctx)
		if !ok {
			break
		}
		if sc.predicate == nil {
			// no filter, send page to next operator

			newPage := sc.createPage()

			it := source.Iterator()
			for row := it.Begin(); row != it.End(); row = it.Next() {
				sc.setPageValues(newPage, row)
			}
			output <- newPage
			continue
		}
		// do filter based on predicate
		if v == nil {
			columns := make(map[string]types.ColumnMetadata)
			lo.ForEach(source.Layout, func(item types.ColumnMetadata, index int) {
				item.Ref = index
				columns[item.Name] = item
			})
			// create predicate visitor if nil
			v = &visitor{
				sc:          sc,
				columns:     columns,
				evalContext: expression.NewEvalContext(sc.ctx),
			}

			v.expr = v.rewrite(sc.predicate)
		}

		page := v.filter(source)
		if page != nil {
			output <- page
		}
	}
}

func (sc *sourceConnector) createPage() *types.Page {
	newPage := types.NewPage()
	// FIXME: maybe page layout not match table schema
	newPage.Layout = sc.outputColumns
	newPage.Columns = lo.Map(sc.outputColumns, func(item types.ColumnMetadata, index int) *types.Column {
		return types.NewColumn()
	})
	return newPage
}

func (sc *sourceConnector) setPageValues(page *types.Page, row types.Row) {
	for i, column := range page.Columns {
		ref := sc.outputColumns[i].Ref
		if ref >= 0 {
			column.Append(row.Get(ref))
		}
	}
}

type visitor struct {
	sc          *sourceConnector
	columns     map[string]types.ColumnMetadata // column name -> column(include column ref(index))
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
		column, ok := v.columns[columName]
		if !ok {
			panic("column not exist")
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
			column: column,
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

func (v *visitor) filter(page *types.Page) *types.Page {
	newPage := v.sc.createPage()

	it := page.Iterator()
	for row := it.Begin(); row != it.End(); row = it.Next() {
		if v.check(row) {
			v.sc.setPageValues(newPage, row)
		}
	}

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
	fmt.Println("filtered page rows:", string(encoding.JSONMarshal(newPage)))
	return newPage
}

func (v *visitor) check(row types.Row) bool {
	switch node := v.expr.(type) {
	case *InExpr:
		return lo.Contains(node.values, string(*row.GetString(node.column.Ref)))
	default:
		panic(fmt.Errorf("not support,%T", v.expr))
	}
}
