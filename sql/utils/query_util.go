package utils

import (
	"github.com/samber/lo"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
)

var nodeId = atomic.NewInt64(0)

func SingleValueQuery(columnName, value string) *tree.Query {
	page := types.NewPage()
	column := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{Name: columnName, DataType: types.DTString}, column)
	column.AppendString(value)
	body := &tree.QuerySpecification{
		Select: SelectAll(),
		From:   Values(page),
	}
	setNodeID(body)
	query := Query(body)
	return query
}

func SimpleQuery(selectItem *tree.Select, from tree.Relation, where tree.Expression) *tree.Query {
	body := &tree.QuerySpecification{
		Select: selectItem,
		From:   from,
		Where:  where,
	}
	setNodeID(body)
	query := Query(body)
	return query
}

func Query(body tree.QueryBody) *tree.Query {
	query := &tree.Query{
		QueryBody: body,
	}
	setNodeID(query)
	return query
}

func SelectItems(columns ...string) *tree.Select {
	return &tree.Select{
		SelectItems: lo.Map(columns, func(item string, index int) tree.SelectItem {
			return SingleSelect(item)
		}),
	}
}

func SingleSelect(columnName string) *tree.SingleColumn {
	singleColumn := &tree.SingleColumn{
		Expression: Identifier(columnName),
	}
	setNodeID(singleColumn)
	return singleColumn
}

func Table(database, namespace, tableName string) *tree.Table {
	table := &tree.Table{
		Name: tree.NewQualifiedName([]*tree.Identifier{
			Identifier(database),
			Identifier(namespace),
			Identifier(tableName),
		}),
	}
	setNodeID(table)
	return table
}

func SelectAll() *tree.Select {
	all := &tree.AllColumns{}
	setNodeID(all)
	return &tree.Select{
		SelectItems: []tree.SelectItem{all},
	}
}

func Values(rows *types.Page) *tree.Values {
	values := &tree.Values{
		Rows: rows,
	}
	setNodeID(values)
	return values
}

func LogicalAnd(term ...tree.Expression) *tree.LogicalExpression {
	logical := &tree.LogicalExpression{Terms: term, Operator: tree.LogicalAND}
	setNodeID(logical)
	return logical
}

func StringEqual(column, value string) tree.Expression {
	return Equal(Identifier(column), StringLiteral(value))
}

func Equal(left, right tree.Expression) tree.Expression {
	expr := &tree.ComparisonExpression{
		Left:     left,
		Right:    right,
		Operator: tree.ComparisonEQ,
	}
	setNodeID(expr)
	return expr
}

func Identifier(value string) *tree.Identifier {
	ident := &tree.Identifier{
		Value: value,
	}
	setNodeID(ident)
	return ident
}

func StringLiteral(value string) *tree.StringLiteral {
	ident := &tree.StringLiteral{
		Value: value,
	}
	setNodeID(ident)
	return ident
}

func setNodeID(node tree.Node) {
	node.SetID(tree.NodeID(nodeId.Inc()))
}
