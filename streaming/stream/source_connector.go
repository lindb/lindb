package stream

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/streaming/stream/input"
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
	inputHandle := input.GetManager().GetInputHandler(tableHandle.App, tableHandle.Stream)

	connector := &sourceConnector{
		ctx: ctx,

		predicate:     predicate,
		outputColumns: outputColumns,

		inbound: operator.NewQueue(make(chan *types.Page, 256)),
	}

	inputHandle.Subscribe(connector)
	fmt.Println("create streaming connector .....")
	return connector
}

type sourceConnector struct {
	ctx context.Context

	predicate     tree.Expression
	outputColumns []types.ColumnMetadata

	inbound *operator.Queue // TODO: queue need close
}

func (sc *sourceConnector) Receive(event any) {
	// fmt.Printf("stream h connector receiver, receive event:%v\n", event)
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
			output <- source
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

type visitor struct {
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
	newPage := types.NewPage()
	newPage.Layout = page.Layout
	newPage.Columns = lo.Map(page.Columns, func(item *types.Column, index int) *types.Column {
		return types.NewColumn()
	})

	it := page.Iterator()
	for row := it.Begin(); row != it.End(); row = it.Next() {
		if v.check(row) {
			for i, column := range newPage.Columns {
				column.Append(row.Get(i))
			}
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
